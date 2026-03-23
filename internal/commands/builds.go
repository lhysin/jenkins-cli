package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Manage Jenkins builds",
}

var buildsListCmd = &cobra.Command{
	Use:   "list [job-name]",
	Short: "List builds for a job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		builds, err := client.ListBuilds(args[0])
		if err != nil {
			return fmt.Errorf("failed to list builds: %w", err)
		}

		fmt.Printf("%-10s %-15s %-15s %s\n", "NUMBER", "STATUS", "RESULT", "DURATION")
		fmt.Println(strings.Repeat("-", 50))
		for _, build := range builds {
			fmt.Printf("%-10d %-15s %-15s %dms\n", build.Number, build.Status, build.Result, build.Duration)
		}
		return nil
	},
}

var buildsInfoCmd = &cobra.Command{
	Use:   "info [job-name] [build-number]",
	Short: "Get build information",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		var buildNumber int
		fmt.Sscanf(args[1], "%d", &buildNumber)

		build, err := client.GetBuildInfo(args[0], buildNumber)
		if err != nil {
			return fmt.Errorf("failed to get build info: %w", err)
		}

		fmt.Printf("Number: %d\nURL: %s\nResult: %s\nDuration: %dms\n", build.Number, build.URL, build.Result, build.Duration)
		return nil
	},
}

var buildsTriggerCmd = &cobra.Command{
	Use:   "trigger [job-name]",
	Short: "Trigger a build",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		jobName := args[0]

		if err := client.TriggerBuild(jobName); err != nil {
			return fmt.Errorf("failed to trigger build: %w", err)
		}

		fmt.Printf("Build triggered for job '%s'\n", jobName)
		return nil
	},
}

var buildsLogsCmd = &cobra.Command{
	Use:   "logs [job-name] [build-number]",
	Short: "Get build logs",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		var buildNumber int
		fmt.Sscanf(args[1], "%d", &buildNumber)

		logs, err := client.GetBuildLogs(args[0], buildNumber)
		if err != nil {
			return fmt.Errorf("failed to get build logs: %w", err)
		}

		fmt.Print(logs)
		return nil
	},
}

func init() {
	buildsCmd.AddCommand(buildsListCmd)
	buildsCmd.AddCommand(buildsInfoCmd)
	buildsCmd.AddCommand(buildsTriggerCmd)
	buildsCmd.AddCommand(buildsLogsCmd)
	rootCmd.AddCommand(buildsCmd)
}
