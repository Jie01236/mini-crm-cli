package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"mini-crm-cli/internal/contacts"
	"mini-crm-cli/internal/storage"
)

var (
	updateID    uint
	updateName  string
	updateEmail string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing contact",
	RunE: func(cmd *cobra.Command, _ []string) error {
		id, name, email, err := resolveUpdateInput(cmd)
		if err != nil {
			return err
		}

		svc, err := getService()
		if err != nil {
			return err
		}
		ctx := commandContext(cmd)

		contact, err := svc.UpdateContact(ctx, id, name, email)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				return fmt.Errorf("contact %d not found: %w", id, storage.ErrNotFound)
			case errors.Is(err, contacts.ErrInvalidEmail), errors.Is(err, contacts.ErrInvalidName):
				return err
			default:
				return fmt.Errorf("update contact: %w", err)
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Contact %d updated: %s <%s>\n", contact.ID, contact.Name, contact.Email)
		return nil
	},
}

func init() {
	updateCmd.Flags().UintVarP(&updateID, "id", "i", 0, "identifier of the contact to update")
	updateCmd.Flags().StringVar(&updateName, "name", "", "new name value")
	updateCmd.Flags().StringVar(&updateEmail, "email", "", "new email value")
}

func resolveUpdateInput(cmd *cobra.Command) (uint, string, string, error) {
	reader := bufio.NewReader(cmd.InOrStdin())

	id := updateID
	if id == 0 {
		for {
			fmt.Fprint(cmd.OutOrStdout(), "ID: ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return 0, "", "", err
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
			id = uint(value)
			break
		}
	}

	name := strings.TrimSpace(updateName)
	email := strings.TrimSpace(updateEmail)

	if name == "" && email == "" {
		fmt.Fprint(cmd.OutOrStdout(), "New name (leave empty to keep): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, "", "", err
		}
		name = strings.TrimSpace(line)

		fmt.Fprint(cmd.OutOrStdout(), "New email (leave empty to keep): ")
		line, err = reader.ReadString('\n')
		if err != nil {
			return 0, "", "", err
		}
		email = strings.TrimSpace(line)
	}

	if name == "" && email == "" {
		return 0, "", "", errors.New("provide at least one field to update")
	}

	return id, name, email, nil
}
