package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	databasePath string
)

var rootCmd = &cobra.Command{
	Use:   "filebrowser",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("File Browser Version (UNTRACKED)")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&databasePath, "database", "d", "./filebrowser.db", "path to the database")

	rootCmd.AddCommand(versionCmd)
}

// Execute executes the commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
