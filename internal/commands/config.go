package commands

import (
	"fmt"
	"os"

	"github.com/jenkins-cli/jenkins-cli/pkg/jenkins"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getClient() *jenkins.Client {
	url, _ := rootCmd.Flags().GetString("url")
	user, _ := rootCmd.Flags().GetString("user")
	token, _ := rootCmd.Flags().GetString("token")

	if url == "" || user == "" || token == "" {
		currentHost := viper.GetString("current")
		
		hostsRaw := viper.Get("hosts")
		if hostsRaw != nil {
			hosts := hostsRaw.(map[string]interface{})
			if currentHost != "" {
				if host, ok := hosts[currentHost]; ok {
					hostMap := host.(map[string]interface{})
					if url == "" {
						url = hostMap["url"].(string)
					}
					if user == "" {
						user = hostMap["user"].(string)
					}
					if token == "" {
						token = hostMap["token"].(string)
					}
				}
			}
		}
		
		if url == "" {
			url = viper.GetString("jenkins.url")
		}
		if user == "" {
			user = viper.GetString("jenkins.user")
		}
		if token == "" {
			token = viper.GetString("jenkins.token")
		}
	}

	return jenkins.NewClient(url, user, token)
}

func getCurrentHost() (string, map[string]interface{}, error) {
	currentHost := viper.GetString("current")
	if currentHost == "" {
		return "", nil, fmt.Errorf("no host selected. Use 'jenkins-cli login' to add a host")
	}
	
	hostsRaw := viper.Get("hosts")
	if hostsRaw == nil {
		return "", nil, fmt.Errorf("no hosts configured")
	}
	
	hosts := hostsRaw.(map[string]interface{})
	host, ok := hosts[currentHost]
	if !ok {
		return "", nil, fmt.Errorf("host '%s' not found", currentHost)
	}
	
	return currentHost, host.(map[string]interface{}), nil
}

func saveConfig() error {
	configDir := os.Getenv("HOME") + "/.jenkins-cli"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	configPath := configDir + "/config.yaml"
	return viper.WriteConfigAs(configPath)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(homeDir + "/.jenkins-cli")
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func init() {
	cobra.OnInitialize(initConfig)
}
