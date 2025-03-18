package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func configCmd(cdk *CDK) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage stencil CLI configuration",
	}
	cmd.AddCommand(configInitCommand(cdk))
	cmd.AddCommand(configListCommand(cdk))
	return cmd
}

func configInitCommand(cdk *CDK) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize CLI configuration",
		Example: heredoc.Doc(`
			$ stencil config init
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cdk.Config.Init(&ClientConfig{}); err != nil {
				return err
			}

			fmt.Printf("Config created\n")
			return nil
		},
	}
}

func configListCommand(cdk *CDK) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List client configuration settings",
		Example: heredoc.Doc(`
			$ stencil config list
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cdk.Config.View()
			if err != nil {
				return ErrClientConfigNotFound
			}

			fmt.Println(data)
			return nil
		},
	}
	return cmd
}
