package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start a new project with interactive setup",
	Run: func(cmd *cobra.Command, args []string) {
		// Project name
		projectPrompt := promptui.Prompt{
			Label: "🚀 Project name",
		}
		projectName, _ := projectPrompt.Run()

		// Language select
		langPrompt := promptui.Select{
			Label: "📦 Choose Language",
			Items: []string{"JavaScript", "TypeScript"},
		}
		_, lang, _ := langPrompt.Run()

		// Tailwind yes/no
		tailwindPrompt := promptui.Select{
			Label: "🎨 Use Tailwind?",
			Items: []string{"Yes", "No"},
		}
		_, tailwind, _ := tailwindPrompt.Run()

		// Testing setup yes/no
		testPrompt := promptui.Select{
			Label: "🧪 Add testing setup (Jest)?",
			Items: []string{"Yes", "No"},
		}
		_, testing, _ := testPrompt.Run()

		// Output summary
		fmt.Println("\n📄 Configuration Summary:")
		fmt.Println("Project Name:", strings.TrimSpace(projectName))
		fmt.Println("Language:", lang)
		fmt.Println("Tailwind:", tailwind)
		fmt.Println("Testing:", testing)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
