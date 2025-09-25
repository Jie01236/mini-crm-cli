package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"mini-crm-cli/internal/storage"
)

var deleteID uint

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a contact by ID",
	RunE: func(cmd *cobra.Command, _ []string) error {
		id, err := resolveDeleteID(cmd)
		if err != nil {
			return err
		}
		svc, err := getService()
		if err != nil {
			return err
		}
		ctx := commandContext(cmd)

		if err := svc.DeleteContact(ctx, id); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return fmt.Errorf("contact %d not found: %w", id, storage.ErrNotFound)
			}
			return fmt.Errorf("delete contact: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Contact %d deleted\n", id)
		return nil
	},
}

func init() {
	deleteCmd.Flags().UintVarP(&deleteID, "id", "i", 0, "identifier of the contact to delete")
}

func resolveDeleteID(cmd *cobra.Command) (uint, error) {
	if deleteID != 0 {
		return deleteID, nil
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	for {
		fmt.Fprint(cmd.OutOrStdout(), "ID: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		value, err := strconv.ParseUint(line, 10, 0)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "invalid id: %v\n", err)
			continue
		}
		if value == 0 {
			fmt.Fprintln(cmd.ErrOrStderr(), "id must be greater than zero")
			continue
		}
		return uint(value), nil
	}
}
