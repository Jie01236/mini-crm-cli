package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"mini-crm-cli/internal/contacts"
	"mini-crm-cli/internal/storage"
)

var (
	cfgFile string

	serviceOnce sync.Once
	serviceInst *contacts.Service
	serviceErr  error
)

type Config struct {
	Storage storage.Config `mapstructure:"storage"`
}

var rootCmd = &cobra.Command{
	Use:   "mini-crm",
	Short: "Mini-CRM CLI provides contact management commands",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		_, err := getService()
		return err
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to the configuration file (default is ./config.yaml)")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func getService() (*contacts.Service, error) {
	serviceOnce.Do(func() {
		v := viper.GetViper()
		if cfgFile != "" {
			v.SetConfigFile(cfgFile)
		} else if v.ConfigFileUsed() == "" {
			v.SetConfigName("config")
			v.SetConfigType("yaml")
			v.AddConfigPath(".")
		}
		if err := v.ReadInConfig(); err != nil {
			serviceErr = fmt.Errorf("read configuration: %w", err)
			return
		}

		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			serviceErr = fmt.Errorf("unmarshal configuration: %w", err)
			return
		}

		if cfg.Storage.Type == "" {
			cfg.Storage.Type = "memory"
		}
		store, err := storage.New(cfg.Storage)
		if err != nil {
			serviceErr = err
			return
		}
		serviceInst = contacts.NewService(store)
	})

	return serviceInst, serviceErr
}

func commandContext(cmd *cobra.Command) context.Context {
	ctx := cmd.Context()
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
