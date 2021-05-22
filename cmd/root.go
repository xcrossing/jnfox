package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	threads int
	cfgFile string
	rootCmd = &cobra.Command{Use: "jnfox"}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.jnfox.json)")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 2, "threads to get cover")
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(os.Getenv("HOME"))
		viper.SetConfigName(".jnfox")
	}
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
