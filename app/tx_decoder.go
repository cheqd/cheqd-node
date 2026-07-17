package app

import (
	"github.com/cosmos/gogoproto/proto"
	protov2 "google.golang.org/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	coretypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// CustomTxDecoder wraps the standard SDK TxDecoder to handle the custom AuthInfo
// format produced by the fee-abstraction module's reward distribution transactions.
//
// The standard decoder calls unknownproto.RejectUnknownFieldsStrict on AuthInfo
// bytes, which rejects the custom proto fields used by fee-abstraction's
// CustomRewardEntry messages. This wrapper falls back to lenient decoding when
// strict rejection fails, allowing the transaction to be returned for querying
// via RPC or block explorers.
//
// Note: the fallback is intentionally minimal — just enough to decode and return
// the tx. Fee fields may be zeroed/empty since the custom AuthInfo format cannot
// be fully unmarshalled into the standard Fee structure. Ante handlers that
// require fee checking will reject these txs (which is correct: reward-distribution
// txs are internal accounting entries, not user-submitted transactions).
func CustomTxDecoder(strictDecoder sdk.TxDecoder, cdc codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		tx, err := strictDecoder(txBytes)
		if err == nil {
			return tx, nil
		}

		// Strict decoding failed — likely due to custom AuthInfo fields from
		// the fee-abstraction module. Try lenient decoding as a fallback.
		fallbackTx, fbErr := fallbackDecode(txBytes, cdc)
		if fbErr != nil {
			// Return the original strict-decoder error instead of the fallback
			// error, since it's more informative about what went wrong.
			return nil, err
		}

		return fallbackTx, nil
	}
}

// fallbackTx wraps a *coretypes.Tx to implement sdk.Tx and related interfaces,
// using the supplied codec for lazy Any type resolution.
//
// It implements:
//   - sdk.Tx (GetMsgs, GetMsgsV2)
//   - sdk.HasValidateBasic (ValidateBasic)
//   - sdk.FeeTx (GetGas, GetFee, FeePayer, FeeGranter)
//   - sdk.TxWithMemo (GetMemo)
//   - sdk.TxWithTimeoutHeight (GetTimeoutHeight)
type fallbackTx struct {
	tx  *coretypes.Tx
	cdc codec.Codec

	// cached results from GetSigners/GetMsgsV2
	msgsV2  []protov2.Message
	signers [][]byte
}

func (f *fallbackTx) GetMsgs() []sdk.Msg { return f.tx.GetMsgs() }

func (f *fallbackTx) GetMsgsV2() ([]protov2.Message, error) {
	_, err := f.ensureInit()
	if err != nil {
		return nil, err
	}
	return f.msgsV2, nil
}

func (f *fallbackTx) ValidateBasic() error { return f.tx.ValidateBasic() }

func (f *fallbackTx) GetGas() uint64 { return f.tx.GetGas() }

func (f *fallbackTx) GetFee() sdk.Coins { return f.tx.GetFee() }

func (f *fallbackTx) FeePayer() []byte {
	feePayer := f.tx.AuthInfo.Fee.Payer
	if feePayer != "" {
		feePayerAddr, err := f.cdc.InterfaceRegistry().SigningContext().AddressCodec().StringToBytes(feePayer)
		if err != nil {
			panic(err)
		}
		return feePayerAddr
	}
	signers, err := f.GetSigners()
	if err != nil || len(signers) == 0 {
		return nil
	}
	return signers[0]
}

func (f *fallbackTx) FeeGranter() []byte {
	return f.tx.FeeGranter(f.cdc)
}

func (f *fallbackTx) GetMemo() string {
	if f.tx.Body != nil {
		return f.tx.Body.Memo
	}
	return ""
}

func (f *fallbackTx) GetTimeoutHeight() uint64 {
	if f.tx.Body != nil {
		return f.tx.Body.TimeoutHeight
	}
	return 0
}

// GetSigners lazily resolves message signers via the codec.
func (f *fallbackTx) GetSigners() ([][]byte, error) {
	_, err := f.ensureInit()
	if err != nil {
		return nil, err
	}
	return f.signers, nil
}

// ensureInit lazily resolves signers and v2 messages using the codec.
func (f *fallbackTx) ensureInit() ([]protov2.Message, error) {
	if f.msgsV2 != nil && f.signers != nil {
		return f.msgsV2, nil
	}
	signers, msgsV2, err := f.tx.GetSigners(f.cdc)
	if err != nil {
		return nil, err
	}
	f.signers = signers
	f.msgsV2 = msgsV2
	return msgsV2, nil
}

// fallbackDecode unmarshals a transaction from raw bytes with lenient handling
// of the AuthInfo section. It uses gogo/protobuf's proto.Unmarshal directly
// (which does not reject unknown fields), allowing custom AuthInfo formats
// (such as fee-abstraction's CustomRewardEntry) to decode without error.
func fallbackDecode(txBytes []byte, cdc codec.Codec) (*fallbackTx, error) {
	var raw coretypes.TxRaw
	if err := proto.Unmarshal(txBytes, &raw); err != nil {
		return nil, sdkerrors.ErrTxDecode.Wrapf("cannot unmarshal TxRaw: %v", err)
	}

	var body coretypes.TxBody
	if err := proto.Unmarshal(raw.BodyBytes, &body); err != nil {
		return nil, sdkerrors.ErrTxDecode.Wrapf("cannot unmarshal TxBody: %v", err)
	}

	var authInfo coretypes.AuthInfo
	if err := proto.Unmarshal(raw.AuthInfoBytes, &authInfo); err != nil {
		// If AuthInfo bytes contain a completely incompatible format, create a
		// minimal valid AuthInfo so the tx can still be returned for inspection.
		authInfo = coretypes.AuthInfo{
			Fee: &coretypes.Fee{},
		}
	}

	// Ensure nil fields that would cause panics downstream.
	if authInfo.Fee == nil {
		authInfo.Fee = &coretypes.Fee{}
	}
	if authInfo.SignerInfos == nil {
		authInfo.SignerInfos = []*coretypes.SignerInfo{}
	}

	theTx := &coretypes.Tx{
		Body:       &body,
		AuthInfo:   &authInfo,
		Signatures: raw.Signatures,
	}

	// Unpack Any interfaces in the body so GetMsgs works.
	if err := theTx.UnpackInterfaces(cdc.InterfaceRegistry()); err != nil {
		return nil, sdkerrors.ErrTxDecode.Wrapf("cannot unpack Tx interfaces: %v", err)
	}

	return &fallbackTx{tx: theTx, cdc: cdc}, nil
}
