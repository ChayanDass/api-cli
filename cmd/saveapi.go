package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Define the saveapi command
var saveapiCmd = &cobra.Command{
	Use:   "saveapi",
	Short: "Save API configuration for future use",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg APIConfig

		// Prompt for API name
		namePrompt := promptui.Prompt{Label: "🔖 API Name"}
		if name, err := namePrompt.Run(); err == nil {
			cfg.Name = name
		} else {
			fmt.Println("❌ Failed to read API name:", err)
			return
		}

		// Prompt for API URL
		urlPrompt := promptui.Prompt{Label: "🌐 API URL"}
		if url, err := urlPrompt.Run(); err == nil {
			cfg.URL = url
		} else {
			fmt.Println("❌ Failed to read URL:", err)
			return
		}

		// Prompt for HTTP method
		methodPrompt := promptui.Select{
			Label: "📦 HTTP Method",
			Items: []string{"GET", "POST", "PUT", "DELETE"},
		}
		if _, method, err := methodPrompt.Run(); err == nil {
			cfg.Method = method
		} else {
			fmt.Println("❌ Failed to select method:", err)
			return
		}

		// Optional token
		tokenPrompt := promptui.Prompt{Label: "🔐 Auth Token (press enter to skip)"}
		if token, err := tokenPrompt.Run(); err == nil {
			cfg.Token = token
		}

		// Optional body
		bodyPrompt := promptui.Prompt{Label: "🧾 Request Body (optional)"}
		if body, err := bodyPrompt.Run(); err == nil {
			cfg.Body = body
		}

		// Save config
		if err := saveAPIConfig(cfg); err != nil {
			fmt.Println("❌ Failed to save API:", err)
			return
		}
		fmt.Println("✅ API saved successfully!")
	},
}

// Function to save the API config to file
func saveAPIConfig(cfg APIConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".chayan")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file := filepath.Join(dir, "apis.json")

	var configs []APIConfig
	if data, err := os.ReadFile(file); err == nil {
		_ = json.Unmarshal(data, &configs)
	}

	configs = append(configs, cfg)

	finalData, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, finalData, 0644)
}

// Add the command to root
func init() {
	rootCmd.AddCommand(saveapiCmd)
}
