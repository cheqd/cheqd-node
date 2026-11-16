package incidentrecovery

import (
	"context"
	"fmt"
	"os"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// ConfigPath resolves the config file path at upgrade-execution time: the
// env var if set, otherwise DefaultConfigPath. A plain function rather than
// an init-time global so it's re-evaluated lazily, inside the handler
// closure, every time — not read once at app startup before the env might
// even be set.
func ConfigPath() string {
	if p := os.Getenv(ConfigPathEnvVar); p != "" {
		return p
	}
	return DefaultConfigPath
}

// NewUpgradeHandler is the one entry point app.go should wire into
// app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, ...). Unlike
// RegisterUpgradeHandler below (which takes an already-loaded *Config), this
// loads and validates the config from disk lazily, inside the returned
// closure — not at app construction/RegisterUpgradeHandlers time. That
// matters because RegisterUpgradeHandlers runs on every single node on
// every startup, forever, including the overwhelming majority of nodes
// that will never execute this specific upgrade plan. A missing or
// malformed config file must never prevent the binary from starting, or
// from running any other, unrelated upgrade — it should only ever error
// out at the moment this specific upgrade plan actually fires.
func NewUpgradeHandler(
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	distrKeeper DistrKeeper,
	accountKeeper AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// plan.Info is part of the governance proposal itself — every node
		// evaluates the exact same value, since it came through consensus,
		// unlike the config file on disk. Binding the local file to it here
		// is what makes "the config actually run" something voters approved,
		// not just the plan's name. See LoadConfigVerified's doc comment.
		cfg, err := LoadConfigVerified(ConfigPath(), plan.Info)
		if err != nil {
			return nil, fmt.Errorf("incident recovery upgrade fired but config at %q (env var %s) could not be verified against the governance-approved plan.Info: %w", ConfigPath(), ConfigPathEnvVar, err)
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
		}

		// Returning the version map unchanged for modules this upgrade
		// doesn't touch is standard — whoever wires this in should run it
		// through the real module manager's RunMigrations as usual.
		return vm, nil
	}
}
