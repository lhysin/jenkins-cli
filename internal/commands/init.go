package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize and login to Jenkins server",
	Long: `Initialize Jenkins CLI by saving your Jenkins server credentials.
This command will prompt for your Jenkins URL, username, and API token.

The credentials will be saved to ~/.jenkins-cli/config.yaml

Example:
  jenkins-cli init`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var url, user, token string

		fmt.Print("Jenkins URL (e.g., http://localhost:8080): ")
		fmt.Scanln(&url)

		fmt.Print("Username: ")
		fmt.Scanln(&user)

		fmt.Print("API Token: ")
		fmt.Scanln(&token)

		if url == "" || user == "" || token == "" {
			return fmt.Errorf("URL, username, and token are required")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := fmt.Sprintf("%s/.jenkins-cli", homeDir)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		configPath := fmt.Sprintf("%s/config.yaml", configDir)
		
		viper.SetConfigFile(configPath)
		viper.Set("jenkins.url", url)
		viper.Set("jenkins.user", user)
		viper.Set("jenkins.token", token)

		if err := viper.WriteConfig(); err != nil {
			if _, writeErr := os.Stat(configPath); os.IsNotExist(writeErr) {
				if err := viper.SafeWriteConfig(); err != nil {
					return fmt.Errorf("failed to write config: %w", err)
				}
			} else {
				return fmt.Errorf("failed to write config: %w", err)
			}
		}

		fmt.Printf("\n✓ Configuration saved to %s\n", configPath)
		fmt.Println("✓ You can now use jenkins-cli without --url, --user, --token flags")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
