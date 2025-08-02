package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listapisCmd = &cobra.Command{
	Use:   "listapis",
	Short: "List all saved APIs",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		file := filepath.Join(home, ".chayan", "apis.json")

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("❌ No saved APIs found.")
			return
		}

		var apis []APIConfig
		if err := json.Unmarshal(data, &apis); err != nil {
			fmt.Println("❌ Failed to parse API file.")
			return
		}

		if len(apis) == 0 {
			fmt.Println("🕳️ No APIs saved yet.")
			return
		}

		fmt.Println("📋 Saved APIs:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tMETHOD\tURL")
		for _, api := range apis {
			fmt.Fprintf(w, "%s\t%s\t%s\n", api.Name, api.Method, api.URL)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listapisCmd)
}
