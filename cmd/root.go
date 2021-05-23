package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfox/util"
)

var (
	threads int
	cfgFile string
	config  *util.Config

	rootCmd = &cobra.Command{
		Use: "jnfox",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cfgFile == "" {
				cfgFile = os.Getenv("HOME") + "/.jnfox"
			}
			cfg, err := util.LoadConfig(cfgFile)
			if err != nil {
				return err
			}
			if cfg.Host == "" {
				return errors.New("host not defined")
			}
			config = cfg
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.jnfox.json)")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 2, "threads to get cover")
}

func Execute() error {
	return rootCmd.Execute()
}
