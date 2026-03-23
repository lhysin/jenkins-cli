package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Jenkins server connection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		
		if client.URL == "" {
			fmt.Println("✗ Not configured. Run 'jenkins-cli init' first.")
			return nil
		}

		fmt.Printf("✓ Connected to: %s\n", client.URL)
		fmt.Printf("✓ User: %s\n", client.User)

		jobs, err := client.ListJobs()
		if err != nil {
			fmt.Printf("✗ Failed to connect: %v\n", err)
			return err
		}

		fmt.Printf("✓ Successfully connected! Found %d jobs.\n", len(jobs))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
