package cmd

import "github.com/spf13/cobra"

var (
	threads int
	rootCmd = &cobra.Command{Use: "jnfox"}
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 2, "threads to get cover")
}

func Execute() error {
	return rootCmd.Execute()
}
