package incidentrecovery

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// NewUpgradeHandler is the one entry point app.go should wire into
// app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, ...). It reads the
// recovery config from the governance-approved plan's Info field (see
// ParsePlanInfo), lazily inside the returned closure, so nothing about this
// upgrade is evaluated at node startup or affects any other upgrade. Every
// node evaluates the same plan.Info because it came through consensus, so
// there is no per-node file that could differ or be missing.
func NewUpgradeHandler(
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	distrKeeper DistrKeeper,
	accountKeeper AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		cfg, err := ParsePlanInfo(plan.Info)
		if err != nil {
			return nil, fmt.Errorf("incident recovery upgrade fired but its config could not be read from the governance-approved plan.Info: %w", err)
		}
		return RegisterUpgradeHandler(cfg, bankKeeper, stakingKeeper, distrKeeper, accountKeeper)(ctx, plan, vm)
	}
}

// RegisterUpgradeHandler builds the cosmos-sdk UpgradeHandler for
// cfg.UpgradeName. This is the piece that actually needs wiring into the
// real cheqd-node app.go, e.g.:
//
//	app.UpgradeKeeper.SetUpgradeHandler(
//	    cfg.UpgradeName,
//	    incidentrecovery.RegisterUpgradeHandler(cfg, app.BankKeeper, app.StakingKeeper, app.DistrKeeper, app.AccountKeeper),
//	)
//
// NOT done here, since this package isn't compiled against the real app.
func RegisterUpgradeHandler(
	cfg *Config,
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	distrKeeper DistrKeeper,
	accountKeeper AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		results, err := RunIncidentRecovery(ctx, cfg, bankKeeper, stakingKeeper, distrKeeper, accountKeeper)
		if err != nil {
			return vm, fmt.Errorf("incident recovery upgrade handler failed: %w", err)
		}

		for _, r := range results {
			attrs := []sdk.Attribute{
				sdk.NewAttribute("address", r.Address),
			}
			if r.Skipped {
				attrs = append(attrs,
					sdk.NewAttribute("skipped", "true"),
					sdk.NewAttribute("skip_reason", r.SkipReason),
				)
			} else {
				attrs = append(attrs,
					sdk.NewAttribute("rewards_withdrawn", r.RewardsWithdrawn.String()),
					sdk.NewAttribute("delegations_moved", r.DelegationsMoved.String()),
					sdk.NewAttribute("unbonding_moved", r.UnbondingMoved.String()),
					sdk.NewAttribute("total_swept", r.TotalSwept.String()),
					sdk.NewAttribute("destination", cfg.DestinationAddress),
				)
			}
			// Warnings go on-chain too, even for an otherwise-successful
			// address — a partial per-validator failure inside a
			// successful overall sweep is exactly the kind of thing that
			// must be visible to anyone auditing this afterwards, not just
			// available in an in-memory struct that nothing ever reads.
			for i, w := range r.Warnings {
				attrs = append(attrs, sdk.NewAttribute(fmt.Sprintf("warning_%d", i), w))
			}
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent("incident_recovery_sweep", attrs...),
			)
			// Also logged: if the upgrade block is applied during a restart's
			// handshake replay, its events never reach block_results, so the
			// node log is the only place operators can see the outcome.
			logArgs := make([]any, 0, len(attrs)*2)
			for _, a := range attrs {
				logArgs = append(logArgs, a.Key, a.Value)
			}
			sdkCtx.Logger().Info("incident_recovery_sweep", logArgs...)
		}

		// Returning the version map unchanged for modules this upgrade
		// doesn't touch is standard — whoever wires this in should run it
		// through the real module manager's RunMigrations as usual.
		return vm, nil
	}
}
