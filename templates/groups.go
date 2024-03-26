package templates

import (
	"strings"

	"github.com/spf13/cobra"
)

// AddGroup adds a group to the passed command and adds the passed commands to the group.
func AddGroup(cmd *cobra.Command, id string, title string, groupCmds ...*cobra.Command) {
	if !strings.HasSuffix(title, ":") {
		title += ":"
	}
	cmd.AddGroup(&cobra.Group{ID: id, Title: title})
	for _, groupCmd := range groupCmds {
		groupCmd.GroupID = id
	}
	cmd.AddCommand(groupCmds...)
}
