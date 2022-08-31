package cmd

import "github.com/spf13/cobra"

func IsClientCLI(cmd *cobra.Command) bool {
	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["client"] == "true" {
			return true
		}
	}
	return false
}
