package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage Jenkins jobs",
}

var jobsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		jobs, err := client.ListJobs()
		if err != nil {
			return fmt.Errorf("failed to list jobs: %w", err)
		}

		fmt.Printf("%-30s %-50s %s\n", "NAME", "URL", "COLOR")
		fmt.Println(strings.Repeat("-", 90))
		for _, job := range jobs {
			fmt.Printf("%-30s %-50s %s\n", job.Name, job.URL, job.Color)
		}
		return nil
	},
}

var jobsInfoCmd = &cobra.Command{
	Use:   "info [job-name]",
	Short: "Get job information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		job, err := client.GetJobInfo(args[0])
		if err != nil {
			return fmt.Errorf("failed to get job info: %w", err)
		}

		fmt.Printf("Name: %s\nURL: %s\nColor: %s\nDescription: %s\n", job.Name, job.URL, job.Color, job.Description)
		return nil
	},
}

var jobsCreateCmd = &cobra.Command{
	Use:   "create [job-name]",
	Short: "Create a new job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		name := args[0]

		configFlag, _ := cmd.Flags().GetString("config")
		config, err := readConfigFromInput(configFlag)
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if err := client.CreateJob(name, config); err != nil {
			return fmt.Errorf("failed to create job: %w", err)
		}

		fmt.Printf("Job '%s' created successfully\n", name)
		return nil
	},
}

var jobsDeleteCmd = &cobra.Command{
	Use:   "delete [job-name]",
	Short: "Delete a job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		name := args[0]

		fmt.Printf("Are you sure you want to delete job '%s'? (y/N): ", name)
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled")
			return nil
		}

		if err := client.DeleteJob(name); err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		fmt.Printf("Job '%s' deleted successfully\n", name)
		return nil
	},
}

var jobsUpdateCmd = &cobra.Command{
	Use:   "update [job-name]",
	Short: "Update a job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		name := args[0]

		configFlag, _ := cmd.Flags().GetString("config")
		config, err := readConfigFromInput(configFlag)
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if err := client.UpdateJob(name, config); err != nil {
			return fmt.Errorf("failed to update job: %w", err)
		}

		fmt.Printf("Job '%s' updated successfully\n", name)
		return nil
	},
}

func readConfigFromInput(configFlag string) (string, error) {
	if configFlag == "" {
		return "", fmt.Errorf("config is required (use -c flag, or pipe from stdin with -c -)")
	}

	if configFlag == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
		return string(data), nil
	}

	if _, err := os.Stat(configFlag); err == nil {
		data, err := os.ReadFile(configFlag)
		if err != nil {
			return "", fmt.Errorf("failed to read config file: %w", err)
		}
		return string(data), nil
	}

	return configFlag, nil
}

func init() {
	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsInfoCmd)
	jobsCmd.AddCommand(jobsCreateCmd)
	jobsCmd.AddCommand(jobsDeleteCmd)
	jobsCmd.AddCommand(jobsUpdateCmd)
	rootCmd.AddCommand(jobsCmd)

	jobsCreateCmd.Flags().StringP("config", "c", "", "Job configuration XML (\"-\" for stdin, or file path)")
	jobsUpdateCmd.Flags().StringP("config", "c", "", "Job configuration XML (\"-\" for stdin, or file path)")
}
