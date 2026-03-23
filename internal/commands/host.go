package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login [name] [url] [user]",
	Short: "Login to a Jenkins host (user and token are optional for anonymous access)",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]
		var user, token string

		if len(args) >= 3 {
			user = args[2]
			fmt.Print("Token (press Enter for anonymous access): ")
			fmt.Scanln(&token)
		} else {
			fmt.Println("Anonymous access mode (read-only)")
		}

		hostsRaw := viper.Get("hosts")
		hosts := make(map[string]interface{})
		if hostsRaw != nil {
			hosts = hostsRaw.(map[string]interface{})
		}

		hosts[name] = map[string]interface{}{
			"url":   url,
			"user":  user,
			"token": token,
		}

		viper.Set("hosts", hosts)
		viper.Set("current", name)

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if user == "" {
			fmt.Printf("Logged in to '%s' (anonymous) and set as current host\n", name)
		} else {
			fmt.Printf("Logged in to '%s' and set as current host\n", name)
		}
		return nil
	},
}

var useCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch to a Jenkins host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		hostsRaw := viper.Get("hosts")
		if hostsRaw == nil {
			return fmt.Errorf("no hosts configured. Use 'jenkins-cli login' to add a host")
		}

		hosts := hostsRaw.(map[string]interface{})
		if _, ok := hosts[name]; !ok {
			return fmt.Errorf("host '%s' not found", name)
		}

		viper.Set("current", name)

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Switched to host '%s'\n", name)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout [name]",
	Short: "Remove a Jenkins host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		hostsRaw := viper.Get("hosts")
		if hostsRaw == nil {
			return fmt.Errorf("no hosts configured")
		}

		hosts := hostsRaw.(map[string]interface{})
		if _, ok := hosts[name]; !ok {
			return fmt.Errorf("host '%s' not found", name)
		}

		delete(hosts, name)
		viper.Set("hosts", hosts)

		current := viper.GetString("current")
		if current == name {
			viper.Set("current", "")
			if len(hosts) > 0 {
				fmt.Println("\nRemaining hosts:")
				hostList := make([]string, 0, len(hosts))
				for k := range hosts {
					hostList = append(hostList, k)
				}
				for i, h := range hostList {
					fmt.Printf("  %d) %s\n", i+1, h)
				}
				fmt.Printf("\nSelect a host to use (1-%d) or press Enter to skip: ", len(hosts))
				var selection string
				fmt.Scanln(&selection)
				
				if selection != "" {
					var idx int
					fmt.Sscanf(selection, "%d", &idx)
					if idx >= 1 && idx <= len(hostList) {
						viper.Set("current", hostList[idx-1])
						fmt.Printf("Switched to host '%s'\n", hostList[idx-1])
					}
				}
				
				if viper.GetString("current") == "" {
					fmt.Println("No host selected. Use 'jenkins-cli use <name>' to select one")
				}
			}
		}

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Host '%s' removed\n", name)
		return nil
	},
}

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List all configured Jenkins hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsRaw := viper.Get("hosts")
		current := viper.GetString("current")

		if hostsRaw == nil {
			fmt.Println("No hosts configured. Use 'jenkins-cli login' to add a host")
			return nil
		}

		hosts := hostsRaw.(map[string]interface{})
		if len(hosts) == 0 {
			fmt.Println("No hosts configured. Use 'jenkins-cli login' to add a host")
			return nil
		}

		fmt.Printf("%-20s %-40s %s\n", "HOST", "URL", "CURRENT")
		fmt.Println(strings.Repeat("-", 70))
		for name, hostData := range hosts {
			hostMap := hostData.(map[string]interface{})
			url := hostMap["url"].(string)
			marker := ""
			if name == current {
				marker = "*"
			}
			fmt.Printf("%-20s %-40s %s\n", name+marker, url, marker)
		}

		if current == "" {
			fmt.Println("\nNo current host selected. Use 'jenkins-cli use <name>' to select one")
		} else {
			fmt.Printf("\nCurrent host: %s\n", current)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(hostsCmd)
}
