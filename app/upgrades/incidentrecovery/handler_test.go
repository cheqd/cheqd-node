package incidentrecovery

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	locked          map[string]sdk.Coins // keyed by account address string; subtracted from balances for SpendableCoins
	sendErr         map[string]error     // keyed by "from" address string
	undelegateErr   map[string]error     // keyed by sourceModule+"|"+recipient address string
	// undelegateErrSeq overrides undelegateErr for a key, if present:
	// each call for that key pops the next entry (nil = succeed that
	// call), letting a test make e.g. the first of two credits for the
	// same key succeed and the second fail.
	undelegateErrSeq map[string][]error
	sendCoinsCalled  []string // records every "from" address SendCoins was called for
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		balances:         map[string]sdk.Coins{},
		locked:           map[string]sdk.Coins{},
		sendErr:          map[string]error{},
		undelegateErr:    map[string]error{},
		undelegateErrSeq: map[string][]error{},
	}
}

func (m *mockBankKeeper) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()].Sub(m.locked[addr.String()]...)
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
	if seq, ok := m.undelegateErrSeq[key]; ok && len(seq) > 0 {
		err := seq[0]
		m.undelegateErrSeq[key] = seq[1:]
		if err != nil {
			return err
		}
		m.balances[recipientAddr.String()] = m.balances[recipientAddr.String()].Add(amt...)
		return nil
	}
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
	setUBDs      []stakingtypes.UnbondingDelegation // records every record SetUnbondingDelegation was called with
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

func (m *mockStakingKeeper) GetAllDelegatorDelegations(_ context.Context, delegator sdk.AccAddress) ([]stakingtypes.Delegation, error) {
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

func (m *mockStakingKeeper) SetUnbondingDelegation(_ context.Context, ubd stakingtypes.UnbondingDelegation) error {
	m.setUBDs = append(m.setUBDs, ubd)
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

// mockAccountKeeper models whether a given address is a module account.
// Defaults to "no account at all" (GetAccount returns nil), same as a real
// AccountKeeper for any ordinary, never-module address.
type mockAccountKeeper struct {
	moduleAccounts map[string]bool // keyed by address string
}

func newMockAccountKeeper() *mockAccountKeeper {
	return &mockAccountKeeper{moduleAccounts: map[string]bool{}}
}

func (m *mockAccountKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	if m.moduleAccounts[addr.String()] {
		return authtypes.NewEmptyModuleAccount("mock-module")
	}
	return nil
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
	accounts := newMockAccountKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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
	accounts := newMockAccountKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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
	accounts := newMockAccountKeeper()

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
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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

func TestRunIncidentRecovery_PartialUnbondingFailureKeepsFailedEntryOnRecord(t *testing.T) {
	// Review finding #5: if one unbonding entry in a record succeeds and a
	// sibling entry in the SAME record fails to credit, only the
	// successful entry's tokens have actually moved. Deleting the whole
	// record (the original bug) would strand the failed entry's tokens in
	// the NotBonded pool with nothing pointing to them ever again. The
	// fix must keep the record with only the failed entry remaining.
	bank := newMockBankKeeper()
	key := stakingtypes.NotBondedPoolName + "|" + testDelAddr1
	// first entry's credit succeeds (nil), second entry's credit fails.
	bank.undelegateErrSeq[key] = []error{nil, errors.New("simulated credit failure")}
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	staking.unbondings[testDelAddr1] = []stakingtypes.UnbondingDelegation{
		{
			DelegatorAddress: testDelAddr1,
			ValidatorAddress: testValAddr1,
			Entries: []stakingtypes.UnbondingDelegationEntry{
				{CreationHeight: 12345, Balance: math.NewInt(5000)},
				{CreationHeight: 12400, Balance: math.NewInt(3000)},
			},
		},
	}

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if !r.UnbondingMoved.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(5000)))) {
		t.Fatalf("expected only the successful 5000%s entry to be recorded as moved, got %s", testBondDenom, r.UnbondingMoved)
	}
	if len(staking.removedUBDs) != 0 {
		t.Fatalf("expected RemoveUnbondingDelegation to NOT be called while one entry is still outstanding, got %d calls", len(staking.removedUBDs))
	}
	if len(staking.setUBDs) != 1 {
		t.Fatalf("expected SetUnbondingDelegation to persist the trimmed record once, got %d calls", len(staking.setUBDs))
	}
	remaining := staking.setUBDs[0]
	if len(remaining.Entries) != 1 || !remaining.Entries[0].Balance.Equal(math.NewInt(3000)) {
		t.Fatalf("expected the trimmed record to retain only the failed 3000%s entry, got %+v", testBondDenom, remaining.Entries)
	}
}

func TestRunIncidentRecovery_ConvertsActiveDelegationToLiquid(t *testing.T) {
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	staking.delegations[testDelAddr1] = []stakingtypes.Delegation{
		{DelegatorAddress: testDelAddr1, ValidatorAddress: testValAddr1, Shares: math.LegacyNewDec(2000)},
	}
	staking.validators[testValAddr1] = stakingtypes.Validator{
		OperatorAddress: testValAddr1,
		Status:          stakingtypes.Bonded,
	}

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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

func TestRunIncidentRecovery_AbortsEntireRunWhenTokensWouldBeUnaccountedFor(t *testing.T) {
	// The specific dangerous case this test proves is caught: Unbond
	// succeeds (shares already removed) but crediting the account
	// afterwards fails. An earlier draft treated this as a recoverable
	// per-address warning and kept going — but by the time Unbond has
	// returned successfully, the shares are already gone with no way
	// back, so a credit failure here must abort the WHOLE run with an
	// error (the chain halts before this commits), not just warn.
	bank := newMockBankKeeper()
	bank.undelegateErr[stakingtypes.BondedPoolName+"|"+testDelAddr1] = errors.New("simulated credit failure")
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	staking.delegations[testDelAddr1] = []stakingtypes.Delegation{
		{DelegatorAddress: testDelAddr1, ValidatorAddress: testValAddr1, Shares: math.LegacyNewDec(2000)},
	}
	staking.validators[testValAddr1] = stakingtypes.Validator{
		OperatorAddress: testValAddr1,
		Status:          stakingtypes.Bonded,
	}

	cfg := testConfig(testDelAddr1)
	_, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
	if err == nil {
		t.Fatal("expected RunIncidentRecovery to return an error aborting the whole run, got nil")
	}
}

func TestRunIncidentRecovery_AbortsEntireRunOnFinalSendFailure(t *testing.T) {
	// By the time the final SendCoins runs, steps 1-3 may have already
	// liquidated rewards/delegations/unbonding entries into this
	// address's spendable balance. Skipping on a send failure here (the
	// original behaviour) would leave that freshly-liquidated balance
	// sitting immediately spendable in the compromised account — worse
	// than doing nothing. It must abort the whole run instead.
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	bank.sendErr[testDelAddr1] = errors.New("simulated send failure")
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	cfg := testConfig(testDelAddr1)
	_, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
	if err == nil {
		t.Fatal("expected RunIncidentRecovery to return an error aborting the whole run, got nil")
	}
}

func TestRunIncidentRecovery_DoesNotSweepLockedVestingBalance(t *testing.T) {
	// Review finding #4: GetAllBalances includes locked/vesting coins,
	// which SendCoins refuses to move — sweeping the full balance would
	// make the final step fail outright for a vesting account. The fix
	// sweeps only SpendableCoins and flags any locked remainder as a
	// warning instead of silently losing track of it.
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	bank.locked[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(400)))
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	cfg := testConfig(testDelAddr1)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if !r.TotalSwept.Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(600)))) {
		t.Fatalf("expected only the 600%s spendable amount to be swept, got %s", testBondDenom, r.TotalSwept)
	}
	if len(r.Warnings) == 0 {
		t.Fatal("expected a warning flagging the unswept locked/vesting balance, got none")
	}
	if !bank.balances[testDelAddr1].Equal(sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(400)))) {
		t.Fatalf("expected the locked 400%s to remain with the affected address, got %s", testBondDenom, bank.balances[testDelAddr1])
	}
}

func TestRunIncidentRecovery_RefusesModuleAccountDestinationAtExecutionTime(t *testing.T) {
	// Review finding #7: config.go's static deny-list only covers the
	// standard SDK module names. This proves the dynamic, execution-time
	// check (via the real AccountKeeper) catches a module account the
	// static list doesn't know about at all.
	bank := newMockBankKeeper()
	bank.balances[testDelAddr1] = sdk.NewCoins(sdk.NewCoin(testBondDenom, math.NewInt(1000)))
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()
	accounts.moduleAccounts[wellFormedDest] = true

	cfg := testConfig(testDelAddr1)
	_, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
	if err == nil {
		t.Fatal("expected RunIncidentRecovery to refuse a destination that is a module account on-chain, got nil error")
	}
}

func TestRunIncidentRecovery_RefusesToRunWithInvalidConfig(t *testing.T) {
	bank := newMockBankKeeper()
	staking := newMockStakingKeeper(testBondDenom)
	distr := newMockDistrKeeper()
	accounts := newMockAccountKeeper()

	cfg := &Config{
		UpgradeName:        "v4.3.0-incident-recovery",
		DestinationAddress: placeholderDestination,
		AffectedAddresses:  []string{testDelAddr1},
	}
	_, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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
	accounts := newMockAccountKeeper()

	cfg := testConfig(testDelAddr1, testDelAddr2)
	results, err := RunIncidentRecovery(context.Background(), cfg, bank, staking, distr, accounts)
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
