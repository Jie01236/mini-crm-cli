package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contacts",
	RunE: func(cmd *cobra.Command, _ []string) error {
		svc, err := getService()
		if err != nil {
			return err
		}
		ctx := commandContext(cmd)
		contacts, err := svc.ListContacts(ctx)
		if err != nil {
			return err
		}
		if len(contacts) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No contacts found.")
			return nil
		}

		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tEMAIL")
		for _, contact := range contacts {
			fmt.Fprintf(tw, "%d\t%s\t%s\n", contact.ID, contact.Name, contact.Email)
		}
		return tw.Flush()
	},
}
