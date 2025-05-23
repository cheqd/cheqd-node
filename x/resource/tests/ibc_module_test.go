package tests

import (
	"encoding/json"
	"fmt"

	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
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
		PortID:              types.ResourcePortID,
		ChannelID:           "some-channel",
		ChanCap:             capabilitytypes.Capability{Index: 1},
		CounterpartyType:    channeltypes.Counterparty{PortId: "counterparty-port-id", ChannelId: "counterparty-channel-id"},
		CounterpartyVersion: types.IBCVersion,
	}
	return params
}

func DefaultPacket(collectionID string, resourceID string) channeltypes.Packet {
	packet := types.ResourceReqPacket{
		ResourceId:   resourceID,
		CollectionId: collectionID,
	}

	jsonPacket, err := json.Marshal(packet)
	if err != nil {
		panic(err)
	}

	return channeltypes.Packet{
		// number corresponds to the order of sends and receives, where a Packet
		// with an earlier sequence number must be sent and received before a Packet
		// with a later sequence number.
		Sequence: 1,
		// identifies the port on the sending chain.
		SourcePort: "source-port",
		// identifies the channel end on the sending chain.
		SourceChannel: "source-channel",
		// identifies the port on the receiving chain.
		DestinationPort: types.ResourcePortID,
		// identifies the channel end on the receiving chain.
		DestinationChannel: "dest-channel",
		// actual opaque bytes transferred directly to the application module
		Data: jsonPacket,
		// block height after which the packet times out
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: 1,
			RevisionHeight: 10,
		},
		// block timestamp (in nanoseconds) after which the packet times out
		TimeoutTimestamp: 123,
	}
}

var _ = Describe("Resource-IBC", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo
	var resource *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
		resource = setup.CreateSimpleResource(alice.CollectionID, SchemaData, "Resource 1", CLSchemaType, []didsetup.SignInput{alice.SignInput})
	})

	It("OnRecvPacket with existing resource returns resource", func() {
		setup.StorePortWithGenesis()
		packet := DefaultPacket(alice.CollectionID, resource.Resource.Id)
		ack := setup.IBCModule.OnRecvPacket(setup.SdkCtx, packet, []byte{0})
		Expect(ack.Success()).To(Equal(true))

		var ackResult channeltypes.Acknowledgement_Result
		var ackData types.ResourceWithMetadata

		err := json.Unmarshal(ack.Acknowledgement(), &ackResult)
		Expect(err).To(BeNil())

		err = json.Unmarshal(ackResult.Result, &ackData)
		Expect(err).To(BeNil())

		Expect(ackData.Metadata.CollectionId).To(Equal(resource.Resource.CollectionId))
		Expect(ackData.Metadata.Name).To(Equal(resource.Resource.Name))
		Expect(ackData.Resource.GetData()).To(Equal([]byte(SchemaData)))
	})

	It("OnRecvPacket without existing resource returns ack error", func() {
		setup.StorePortWithGenesis()
		packet := DefaultPacket(alice.CollectionID, "not-a-resource-id")
		ack := setup.IBCModule.OnRecvPacket(setup.SdkCtx, packet, []byte{0})
		Expect(ack.Success()).To(Equal(false))

		var ackResult channeltypes.Acknowledgement_Error
		err := json.Unmarshal(ack.Acknowledgement(), &ackResult)

		Expect(err).To(BeNil())
		Expect(ackResult.Error).To(ContainSubstring(fmt.Sprint(types.ErrResourceNotAvail.ABCICode())))
	})

	It("OnRecvPacket incorrect packet returns ack error", func() {
		setup.StorePortWithGenesis()
		packet := DefaultPacket(alice.CollectionID, "not-a-resource-id")
		packet.Data = []byte{1, 2, 3}
		ack := setup.IBCModule.OnRecvPacket(setup.SdkCtx, packet, []byte{0})
		Expect(ack.Success()).To(Equal(false))

		var ackResult channeltypes.Acknowledgement_Error
		err := json.Unmarshal(ack.Acknowledgement(), &ackResult)

		Expect(err).To(BeNil())
		Expect(ackResult.Error).To(ContainSubstring(fmt.Sprint(types.ErrUnexpectedPacket.ABCICode())))
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

	It("OnChanOpenAck (Genesis setup) returns no error ", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		err := setup.IBCModule.OnChanOpenAck(setup.SdkCtx, p.PortID, p.ChannelID, p.CounterpartyType.ChannelId, p.CounterpartyVersion)
		Expect(err).To(BeNil())
	})

	It("OnChanOpenAck (Genesis setup) with wrong version returns error ", func() {
		setup.StorePortWithGenesis()
		p := DefaultParams()
		p.CounterpartyVersion = "invalid-version"
		err := setup.IBCModule.OnChanOpenAck(setup.SdkCtx, p.PortID, p.ChannelID, p.CounterpartyType.ChannelId, p.CounterpartyVersion)
		Expect(err.Error()).To(ContainSubstring("invalid counterparty version"))
	})
})
