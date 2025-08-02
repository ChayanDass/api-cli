// File: cmd/apitest.go
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var testName string

var apitestCmd = &cobra.Command{
	Use:   "apitest",
	Short: "Test APIs interactively or from saved configs",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg APIConfig

		// If --test is passed, load directly
		if testName != "" {
			found, err := loadAPIConfigByName(testName)
			if err != nil {
				fmt.Println("❌ Failed to load API config:", err)
				return
			}
			cfg = found
			fmt.Println("🧪 Testing saved API:", cfg.Name)
		} else {
			// Show all saved APIs to select
			all, err := loadAllAPIConfigs()
			if err != nil {
				fmt.Println("❌ Failed to load saved APIs:", err)
				return
			}
			if len(all) == 0 {
				fmt.Println("⚠️ No saved APIs found. Use `saveapi` command first.")
				return
			}

			items := make([]string, len(all))
			for i, a := range all {
				items[i] = fmt.Sprintf("%s (%s %s)", a.Name, a.Method, a.URL)
			}

			selector := promptui.Select{
				Label: "📚 Select API to test",
				Items: items,
			}
			index, _, err := selector.Run()
			if err != nil {
				fmt.Println("❌ Selection error:", err)
				return
			}
			cfg = all[index]
		}

		// Make the request
		client := &http.Client{}
		req, err := http.NewRequest(cfg.Method, cfg.URL, bytes.NewBuffer([]byte(cfg.Body)))
		if err != nil {
			fmt.Println("❌ Request creation failed:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if cfg.Token != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.Token)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Request failed:", err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		fmt.Println("\n✅ Response:")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("📡 Status: %s\n", resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("🔸 %s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Println(strings.Repeat("-", 60))

		// Try pretty-printing JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, respBody, "", "  "); err == nil {
			fmt.Println(prettyJSON.String())
		} else {
			fmt.Println(string(respBody))
		}
		fmt.Println(strings.Repeat("=", 60))
	},
}

func init() {
	apitestCmd.Flags().StringVarP(&testName, "test", "t", "", "Test a saved API by name")
	rootCmd.AddCommand(apitestCmd)
}

// loadAPIConfigByName loads a single API config by name
func loadAPIConfigByName(name string) (APIConfig, error) {
	all, err := loadAllAPIConfigs()
	if err != nil {
		return APIConfig{}, err
	}
	for _, c := range all {
		if c.Name == name {
			return c, nil
		}
	}
	return APIConfig{}, fmt.Errorf("API with name %q not found", name)
}

// loadAllAPIConfigs loads all saved configs
func loadAllAPIConfigs() ([]APIConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	file := filepath.Join(home, ".chayan", "apis.json")

	var configs []APIConfig
	data, err := os.ReadFile(file)
	if err != nil {
		return configs, err
	}
	if err := json.Unmarshal(data, &configs); err != nil {
		return configs, err
	}
	return configs, nil
}
