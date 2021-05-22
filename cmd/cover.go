package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xcrossing/jnfo"
	"github.com/xcrossing/jnfox/util"
)

func init() {
	rootCmd.AddCommand(cmdCover)
}

var cmdCover = &cobra.Command{
	Use:   "cover [nums]",
	Short: "Get Cover directly",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, nums []string) {
		host := viper.GetString("host")
		if host == "" {
			fmt.Fprintln(os.Stderr, "no host config")
			return
		}

		u, err := url.Parse(host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", host, err.Error())
			return
		}
		path := u.Path

		p := util.MakePool(threads, func(addr string) {
			nfo, err := jnfo.New(addr)
			if err == nil {
				err = util.Download(nfo.PicLink, nfo.NumCastPicName())
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s\n", addr, err.Error())
			}
		})

		for _, num := range nums {
			u.Path = filepath.Join(path, num)
			addr := u.String()
			p.Add(addr)
		}

		p.Wait()
	},
}
