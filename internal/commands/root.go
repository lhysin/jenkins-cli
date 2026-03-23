package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "jenkins-cli",
	Short: "Jenkins CLI is a command-line interface for Jenkins REST API",
	Long: `Jenkins CLI is a tool for interacting with Jenkins through its REST API.
It allows you to manage jobs, builds, nodes, and more from the command line.

Examples:
  jenkins-cli jobs list
  jenkins-cli builds list my-job
  jenkins-cli nodes list`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jenkins-cli.yaml)")
	rootCmd.PersistentFlags().String("url", "", "Jenkins server URL")
	rootCmd.PersistentFlags().String("user", "", "Jenkins username")
	rootCmd.PersistentFlags().String("token", "", "Jenkins API token")
}
