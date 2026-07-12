package errors

import (
	"fmt"
	"strings"
)

// TeamErrors holds all errors encountered during the dispatch, keyed by team name
type TeamErrors map[string][]error

// Implement the error interface so you can return it as a standard error
func (te TeamErrors) Error() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("failed to dispatch status updates for %d teams:\n", len(te)))

	for team, errs := range te {
		msg.WriteString(fmt.Sprintf("  • Team %s:\n", team))
		for _, err := range errs {
			msg.WriteString(fmt.Sprintf("    - %v\n", err))
		}
	}
	return msg.String()
}
