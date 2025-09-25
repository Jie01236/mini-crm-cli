package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"mini-crm-cli/internal/contacts"
)

var (
	addName  string
	addEmail string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new contact",
	RunE: func(cmd *cobra.Command, _ []string) error {
		svc, err := getService()
		if err != nil {
			return err
		}
		ctx := commandContext(cmd)

		name, email, err := resolveAddInput(cmd)
		if err != nil {
			return err
		}

		contact, err := svc.AddContact(ctx, name, email)
		if err != nil {
			if errors.Is(err, contacts.ErrInvalidName) || errors.Is(err, contacts.ErrInvalidEmail) {
				return err
			}
			return fmt.Errorf("add contact: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Contact %q created with ID %d\n", contact.Name, contact.ID)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "name of the contact")
	addCmd.Flags().StringVarP(&addEmail, "email", "e", "", "email of the contact")
}

func resolveAddInput(cmd *cobra.Command) (string, string, error) {
	name := strings.TrimSpace(addName)
	email := strings.TrimSpace(addEmail)

	reader := bufio.NewReader(cmd.InOrStdin())

	if name == "" {
		fmt.Fprint(cmd.OutOrStdout(), "Name: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		name = strings.TrimSpace(line)
	}

	if email == "" {
		fmt.Fprint(cmd.OutOrStdout(), "Email: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		email = strings.TrimSpace(line)
	}

	return name, email, nil
}
