package cmd

import (
	"os"

	cmd "github.com/can3p/kleiner/shared/cmd/cobra"
	"github.com/can3p/kleiner/shared/published"
	"github.com/can3p/sackmesser/generated/buildinfo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sackmesser",
	Short: "json/yaml mutation tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	info := buildinfo.Info()

	cmd.Setup(info, rootCmd)
	published.MaybeNotifyAboutNewVersion(info)
}
