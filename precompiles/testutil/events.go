package testutil

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// validateEvents checks if the provided event names are included as keys in the contract events.
func validateEvents(contractEvents map[string]abi.Event, events []string) ([]abi.Event, error) {
	expEvents := make([]abi.Event, 0, len(events))
	for _, eventStr := range events {
		event, found := contractEvents[eventStr]
		if !found {
			availableABIEvents := make([]string, 0, len(contractEvents))
			for event := range contractEvents {
				availableABIEvents = append(availableABIEvents, event)
			}
			availableABIEventsStr := strings.Join(availableABIEvents, ", ")
			return nil, fmt.Errorf("unknown event %q is not contained in given ABI events:\n%s", eventStr, availableABIEventsStr)
		}
		expEvents = append(expEvents, event)
	}
	return expEvents, nil
}
