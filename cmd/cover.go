package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		host := util.Host()

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
			addr := host + "/" + num
			p.Add(addr)
		}

		p.Wait()
	},
}
