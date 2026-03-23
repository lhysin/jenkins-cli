package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: "Manage Jenkins views",
}

var viewsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all views",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		url := client.URL + "/api/json?tree=views[name,url,owner[class]]"
		
		hostsRaw := viper.Get("hosts")
		current := viper.GetString("current")
		var user, token string
		
		if hostsRaw != nil && current != "" {
			hosts := hostsRaw.(map[string]interface{})
			if host, ok := hosts[current]; ok {
				hostMap := host.(map[string]interface{})
				user = hostMap["user"].(string)
				token = hostMap["token"].(string)
			}
		}

		req, _ := http.NewRequest("GET", url, nil)
		if user != "" && token != "" {
			req.SetBasicAuth(user, token)
		}
		req.Header.Set("Accept", "application/json")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			return fmt.Errorf("failed to list views: %w", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		viewsRaw := result["views"]
		if viewsRaw == nil {
			fmt.Println("No views found")
			return nil
		}
		
		views := viewsRaw.([]interface{})
		
		fmt.Printf("%-30s %-50s\n", "NAME", "URL")
		fmt.Println(strings.Repeat("-", 80))
		for _, v := range views {
			view := v.(map[string]interface{})
			fmt.Printf("%-30s %-50s\n", view["name"], view["url"])
		}
		return nil
	},
}

var viewsInfoCmd = &cobra.Command{
	Use:   "info [view-name]",
	Short: "Get view information and jobs in it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		viewName := args[0]
		apiUrl := client.URL + "/view/" + url.QueryEscape(viewName) + "/api/json?tree=jobs[name,url,color]"

		hostsRaw := viper.Get("hosts")
		current := viper.GetString("current")
		var user, token string

		if hostsRaw != nil && current != "" {
			hosts := hostsRaw.(map[string]interface{})
			if host, ok := hosts[current]; ok {
				hostMap := host.(map[string]interface{})
				user = hostMap["user"].(string)
				token = hostMap["token"].(string)
			}
		}

		req, _ := http.NewRequest("GET", apiUrl, nil)
		if user != "" && token != "" {
			req.SetBasicAuth(user, token)
		}
		req.Header.Set("Accept", "application/json")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get view info: %w", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		fmt.Printf("View: %s\n", viewName)
		fmt.Printf("URL: %s\n\n", result["url"])

		jobsRaw := result["jobs"]
		if jobsRaw == nil {
			fmt.Println("No jobs in this view")
			return nil
		}

		jobs := jobsRaw.([]interface{})
		fmt.Printf("%-30s %-50s %s\n", "NAME", "URL", "COLOR")
		fmt.Println(strings.Repeat("-", 90))
		for _, j := range jobs {
			job := j.(map[string]interface{})
			fmt.Printf("%-30s %-50s %s\n", job["name"], job["url"], job["color"])
		}
		return nil
	},
}

func init() {
	viewsCmd.AddCommand(viewsListCmd)
	viewsCmd.AddCommand(viewsInfoCmd)
	rootCmd.AddCommand(viewsCmd)
}
