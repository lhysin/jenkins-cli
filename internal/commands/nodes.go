package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage Jenkins nodes",
}

var nodesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		nodes, err := client.ListNodes()
		if err != nil {
			return fmt.Errorf("failed to list nodes: %w", err)
		}

		fmt.Printf("%-30s %s\n", "NAME", "OFFLINE")
		fmt.Println(strings.Repeat("-", 40))
		for _, node := range nodes {
			offlineStatus := "false"
			if node.Offline {
				offlineStatus = "true"
			}
			fmt.Printf("%-30s %s\n", node.DisplayName, offlineStatus)
		}
		return nil
	},
}

func init() {
	nodesCmd.AddCommand(nodesListCmd)
	rootCmd.AddCommand(nodesCmd)
}
