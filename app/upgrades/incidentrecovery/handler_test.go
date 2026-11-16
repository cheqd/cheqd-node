package incidentrecovery

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// --- mocks ---
//
// Minimal, hand-written mocks against the small interfaces this package
// defines itself — no mockgen/testify needed. Each mock is keyed by
// delegator address string so a single mock instance can represent
// multiple affected addresses' state in one test.

type mockBankKeeper struct {
	balances        map[string]sdk.Coins // keyed by account address string
	sendErr         map[string]error     // keyed by "from" address string
	undelegateErr   map[string]error     // keyed by sourceModule+"|"+recipient address string
	sendCoinsCalled []string             // records every "from" address SendCoins was called for
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		balances:      map[string]sdk.Coins{},
		sendErr:       map[string]error{},
		undelegateErr: map[string]error{},
	}
}

func (m *mockBankKeeper) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

func (m *mockBankKeeper) SendCoins(_ context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	m.sendCoinsCalled = append(m.sendCoinsCalled, fromAddr.String())
	if err, ok := m.sendErr[fromAddr.String()]; ok {
		return err
	}
	m.balances[fromAddr.String()] = m.balances[fromAddr.String()].Sub(amt...)
	m.balances[toAddr.String()] = m.balances[toAddr.String()].Add(amt...)
	return nil
}

func (m *mockBankKeeper) UndelegateCoinsFromModuleToAccount(_ context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	key := senderModule + "|" + recipientAddr.String()
	if err, ok := m.undelegateErr[key]; ok {
		return err
	}
	m.balances[recipientAddr.String()] = m.balances[recipientAddr.String()].Add(amt...)
	return nil
}

type mockStakingKeeper struct {
	bondDenom    string
	delegations  map[string][]stakingtypes.Delegation          // keyed by delegator address string
	validators   map[string]stakingtypes.Validator             // keyed by validator operator address string
	unbondings   map[string][]stakingtypes.UnbondingDelegation // keyed by delegator address string
	unbondErr    map[string]error                              // keyed by delegator+"|"+validator address string
	removedUBDs  []string                                      // records every delegator|validator RemoveUnbondingDelegation was called for
	removeUBDErr map[string]error
}

func newMockStakingKeeper(bondDenom string) *mockStakingKeeper {
	return &mockStakingKeeper{
		bondDenom:    bondDenom,
		delegations:  map[string][]stakingtypes.Delegation{},
		validators:   map[string]stakingtypes.Validator{},
		unbondings:   map[string][]stakingtypes.UnbondingDelegation{},
		unbondErr:    map[string]error{},
		removeUBDErr: map[string]error{},
	}
}

func (m *mockStakingKeeper) BondDenom(context.Context) (string, error) { return m.bondDenom, nil }

func (m *mockStakingKeeper) GetDelegatorDelegations(_ context.Context, delegator sdk.AccAddress, _ uint16) ([]stakingtypes.Delegation, error) {
	return m.delegations[delegator.String()], nil
}

func (m *mockStakingKeeper) GetValidator(_ context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
	v, ok := m.validators[addr.String()]
	if !ok {
		return stakingtypes.Validator{}, errors.New("validator not found")
	}
	return v, nil
}

func (m *mockStakingKeeper) Unbond(_ context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
	key := delAddr.String() + "|" + valAddr.String()
	if err, ok := m.unbondErr[key]; ok {
		return math.ZeroInt(), err
	}
	// 1 share = 1 token for test simplicity. Real Unbond returns the final
	// truncated token amount directly (math.Int), not a Dec.
	return shares.TruncateInt(), nil
}

func (m *mockStakingKeeper) GetUnbondingDelegation(_ context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.UnbondingDelegation, error) {
	for _, u := range m.unbondings[delAddr.String()] {
		if u.ValidatorAddress == valAddr.String() {
			return u, nil
		}
	}
	return stakingtypes.UnbondingDelegation{}, errors.New("not found")
}

func (m *mockStakingKeeper) GetAllUnbondingDelegations(_ context.Context, delegator sdk.AccAddress) ([]stakingtypes.UnbondingDelegation, error) {
	return m.unbondings[delegator.String()], nil
}

func (m *mockStakingKeeper) RemoveUnbondingDelegation(_ context.Context, ubd stakingtypes.UnbondingDelegation) error {
	key := ubd.DelegatorAddress + "|" + ubd.ValidatorAddress
	m.removedUBDs = append(m.removedUBDs, key)
	if err, ok := m.removeUBDErr[key]; ok {
		return err
	}
	return nil
}

func (m *mockStakingKeeper) SetUnbondingDelegation(context.Context, stakingtypes.UnbondingDelegation) error {
	return nil
}

type mockDistrKeeper struct {
	rewards     map[string]sdk.Coins // keyed by delegator+"|"+validator address string
	withdrawErr map[string]error
}

func newMockDistrKeeper() *mockDistrKeeper {
	return &mockDistrKeeper{rewards: map[string]sdk.Coins{}, withdrawErr: map[string]error{}}
}

func (m *mockDistrKeeper) WithdrawDelegationRewards(_ context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	key := delAddr.String() + "|" + valAddr.String()
	if err, ok := m.withdrawErr[key]; ok {
		return nil, err
	}
	return m.rewards[key], nil
}

// --- tests ---

const (
	testBondDenom = "ncheq"
	testDelAddr1  = "cheqd1pav86kd046tv7em6nc2dgcycxdppdnyzu2xycv"
	testValAddr1  = "cheqdvaloper1qsp3a2qd6km9g0hczsac8279wcwmzmvzgvre8w"
)

func testConfig(affected ...string) *Config {
	return &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         affected,
	}
}

func TestRunIncidentRecovery_SweepsLiquidBalance(t *testing.T) {
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Skipped {
		t.Fatalf("expected not skipped, got skipped: %s", r.SkipReason)
	}
	if !r.TotalSwept.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))) {
		t.Fatalf("expected 1000%s swept, got %s", testBondDenom, r.TotalSwept)
	}
	destBalance := bank.balances[wellFormedDest]
	if !destBalance.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))) {
		t.Fatalf("expected destination to receive 1000%s, got %s", testBondDenom, destBalance)
	}
	if bank.balances[testDelAddr1].IsZero() == false {
		t.Fatalf("expected source address balance to be zero after sweep, got %s", bank.balances[testDelAddr1])
	}
}

func TestRunIncidentRecovery_SkipsZeroBalanceAddress(t *testing.T) {
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Fatal("expected address with zero balance to be skipped")
	}
	if len(bank.sendCoinsCalled) != 0 {
		t.Fatalf("expected SendCoins to never be called for a zero-balance address, got %d calls", len(bank.sendCoinsCalled))
	}
}

func TestRunIncidentRecovery_CompletesUnbondingBypassingJailedCheck(t *testing.T) {
	// This is the core claim of the whole handler: force-completing an
	// unbonding entry works regardless of jailed status, because it's a
	// direct state operation, not MsgCancelUnbondingDelegation. The mock
	// staking keeper has no jailed-validator check at all in
	// GetDelegatorUnbondingDelegations/RemoveUnbondingDelegation — which is
	// the point: this path never consults validator.IsJailed() anywhere.
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	staking.unbondings[testDelAddr1] = []stakingtypes.UnbondingDelegation{
		{
			DelegatorAddress: testDelAddr1,
			ValidatorAddress: testValAddr1,
			Entries: []stakingtypes.UnbondingDelegationEntry{
				{CreationHeight: 12345, Balance: math.NewInt(5000)},
			},
		},
	}

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if r.Skipped {
		t.Fatalf("expected not skipped, got: %s", r.SkipReason)
	}
	if !r.UnbondingMoved.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(5000)))) {
		t.Fatalf("expected 5000%s moved from unbonding, got %s", testBondDenom, r.UnbondingMoved)
	}
	if !r.TotalSwept.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(5000)))) {
		t.Fatalf("expected 5000%s total swept, got %s", testBondDenom, r.TotalSwept)
	}
	if len(staking.removedUBDs) != 1 {
		t.Fatalf("expected RemoveUnbondingDelegation to be called once, got %d", len(staking.removedUBDs))
	}
}

func TestRunIncidentRecovery_ConvertsActiveDelegationToLiquid(t *testing.T) {
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	staking.delegations[testDelAddr1] = []stakingtypes.Delegation{
		{DelegatorAddress: testDelAddr1, ValidatorAddress: testValAddr1, Shares: math.LegacyNewDec(2000)},
	}
	staking.validators[testValAddr1] = stakingtypes.Validator{
		OperatorAddress: testValAddr1,
		Status:          stakingtypes.Bonded,
	}

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if r.Skipped {
		t.Fatalf("expected not skipped, got: %s", r.SkipReason)
	}
	if !r.DelegationsMoved.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(2000)))) {
		t.Fatalf("expected 2000%s moved from active delegation, got %s", testBondDenom, r.DelegationsMoved)
	}
}

func TestRunIncidentRecovery_WarnsLoudlyWhenTokensAreUnaccountedFor(t *testing.T) {
	// The specific dangerous case this test proves is caught: Unbond
	// succeeds (shares removed) but crediting the account afterwards
	// fails. Before the review fix, this was a silent `continue` with no
	// trace anywhere. Now it must show up as a Warning.
	bank := newMockBankKeeper()
	bank.undelegateErr[stakingtypes.BondedPoolName+"|"+testDelAddr1] = errors.New("simulated credit failure")
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	staking.delegations[testDelAddr1] = []stakingtypes.Delegation{
		{DelegatorAddress: testDelAddr1, ValidatorAddress: testValAddr1, Shares: math.LegacyNewDec(2000)},
	}
	staking.validators[testValAddr1] = stakingtypes.Validator{
		OperatorAddress: testValAddr1,
		Status:          stakingtypes.Bonded,
	}

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if len(r.Warnings) == 0 {
		t.Fatal("expected at least one warning recording the unaccounted-for tokens, got none")
	}
	found := false
	for _, w := range r.Warnings {
		if containsUnaccountedForWarning(w) {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a warning flagging tokens as unaccounted for, got warnings: %v", r.Warnings)
	}
	if !r.DelegationsMoved.IsZero() {
		t.Fatalf("expected DelegationsMoved to be zero since crediting failed, got %s", r.DelegationsMoved)
	}
}

func containsUnaccountedForWarning(s string) bool {
	return len(s) > 0 && (contains(s, "UNACCOUNTED FOR"))
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestRunIncidentRecovery_SkipsOnFinalSendFailure(t *testing.T) {
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	bank.sendErr[testDelAddr1] = errors.New("simulated send failure")
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if !r.Skipped {
		t.Fatal("expected address to be skipped when the final SendCoins fails")
	}
	if r.TotalSwept != nil && !r.TotalSwept.IsZero() {
		t.Fatalf("expected TotalSwept to be empty/zero on failure, got %s", r.TotalSwept)
	}
}

func TestRunIncidentRecovery_RefusesToRunWithInvalidConfig(t *testing.T) {
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	cfg := &Config{
		UpgradeName:        "v4.3.0-incident-recovery",
		DestinationAddress: placeholderDestination,
		AffectedAddresses:  []string{testDelAddr1},
	}
	_, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err == nil {
		t.Fatal("expected RunIncidentRecovery to refuse an invalid (placeholder-destination) config, got nil error")
	}
}

func TestRunIncidentRecovery_SweepsMultipleAddressesIndependently(t *testing.T) {
	const testDelAddr2 = "cheqd18wm3f6s7yw0f0la09rnczvu9xp8wkkxfue8j5m"
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	bank.balances[testDelAddr2] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(2000)))
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()

	cfg := testConfig(testDelAddr1, testDelAddr2)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	destBalance := bank.balances[wellFormedDest]
	expected := sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(3000)))
	if !destBalance.Equal(expected) {
		t.Fatalf("expected destination to receive combined 3000%s, got %s", testBondDenom, destBalance)
	}
}
