package helpers

import (
	"github.com/tendermint/tendermint/abci/types"
)

type HumanReadableEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index,omitempty"`
}

type HumanReadableEvent struct {
	Type       string                        `json:"type"`
	Attributes []HumanReadableEventAttribute `json:"attributes"`
}

func ReadableEvents(events []types.Event) []HumanReadableEvent {
	readableEvents := make([]HumanReadableEvent, len(events))
	for _, event := range events {
		readableAttributes := make([]HumanReadableEventAttribute, len(event.Attributes))
		for i, attribute := range event.Attributes {
			readableAttributes[i] = HumanReadableEventAttribute{
				Key:   string(attribute.Key),
				Value: string(attribute.Value),
				Index: attribute.Index,
			}
		}
		readableEvents = append(readableEvents, HumanReadableEvent{
			Type:       event.Type,
			Attributes: readableAttributes,
		})
	}
	return readableEvents
}
