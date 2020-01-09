package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/LacledesLAN/ShellAndTest/pkg/test"
	"github.com/spf13/cobra"
)

var testFile string
var testDir string
var showOutput bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ShellAndTest",
	Short: "CLI utility for writing automated tests for 3rd-party, CLI binaries",
	Long:  `CLI utility for writing automated tests for 3rd-party, CLI binaries`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Support for piping via stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString(';')
			test.Handler(text, false, showOutput)
			return
		}

		// Test
		if testFile != "" {
			test.Handler(testFile, false, showOutput)
			return
		} else if testDir != "" {
			test.Handler(testDir, true, showOutput)

			return
		} else {
			fmt.Println("No input/arguments specified. Use -h for help.")
			return
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&testFile, "testfile", "f", "", "Provide a path to json: /path/to/file.json")
	rootCmd.PersistentFlags().StringVarP(&testDir, "testdir", "d", "", "Provide a path to directory: /path/to/dir/")
	rootCmd.PersistentFlags().BoolVarP(&showOutput, "output", "o", false, "Show command outputs")
}
