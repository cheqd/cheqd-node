package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	// porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	// host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	// ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"
)

type Params struct {
	Order               channeltypes.Order
	ConnectionHops      []string
	PortID              string
	ChannelID           string
	ChanCap             capabilitytypes.Capability
	CounterpartyType    channeltypes.Counterparty
	CounterpartyVersion string
}

func DefaultParams() Params {
	params := Params{
		Order:               channeltypes.UNORDERED,
		ConnectionHops:      []string{},
		PortID:              types.ResourcePortId,
		ChannelID:           "some-channel",
		ChanCap:             capabilitytypes.Capability{Index: 1},
		CounterpartyType:    channeltypes.Counterparty{PortId: "counterparty-port-id", ChannelId: "counterparty-channel-id"},
		CounterpartyVersion: types.IBCVersion,
	}
	return params
}

var _ = Describe("Resource-IBC", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo
	//var resource *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
		_ = setup.CreateSimpleResource(alice.CollectionID, SchemaData, "Resource 1", CLSchemaType, []didsetup.SignInput{alice.SignInput})
	})

	It("OnChanOpenInit (Genesis setup) with params returns correct version", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		version, err := setup.IBCModule.OnChanOpenInit(setup.SdkCtx, p.Order, p.ConnectionHops, p.PortID, p.ChannelID, &p.ChanCap, p.CounterpartyType, p.CounterpartyVersion)
		Expect(err).To(BeNil())
		Expect(version).To(Equal(types.IBCVersion))
	})

	It("OnChanOpenInit (Genesis setup) with wrong version fails", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		p.CounterpartyVersion = "invalid-version"
		version, err := setup.IBCModule.OnChanOpenInit(setup.SdkCtx, p.Order, p.ConnectionHops, p.PortID, p.ChannelID, &p.ChanCap, p.CounterpartyType, p.CounterpartyVersion)
		Expect(err.Error()).To(ContainSubstring("invalid ibc version"))
		Expect(version).To(Equal(""))
	})

	It("OnChanOpenInit (Genesis setup) with Ordered channel fails", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		p.Order = channeltypes.ORDERED
		version, err := setup.IBCModule.OnChanOpenInit(setup.SdkCtx, p.Order, p.ConnectionHops, p.PortID, p.ChannelID, &p.ChanCap, p.CounterpartyType, p.CounterpartyVersion)
		Expect(err.Error()).To(ContainSubstring("invalid channel ordering"))
		Expect(version).To(Equal(""))
	})

	It("OnChanOpenTry (Genesis setup) with right port returns correct version", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		version, err := setup.IBCModule.OnChanOpenTry(setup.SdkCtx, p.Order, p.ConnectionHops, p.PortID, p.ChannelID, &p.ChanCap, p.CounterpartyType, p.CounterpartyVersion)
		Expect(err).To(BeNil())
		Expect(version).To(Equal(types.IBCVersion))
	})

	It("OnChanOpenTry (Genesis setup) with wrong port fails", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		p.PortID = "invalid-port"
		version, err := setup.IBCModule.OnChanOpenTry(setup.SdkCtx, p.Order, p.ConnectionHops, p.PortID, p.ChannelID, &p.ChanCap, p.CounterpartyType, p.CounterpartyVersion)
		Expect(err.Error()).To(ContainSubstring("invalid port"))
		Expect(version).To(Equal(""))
	})

})
