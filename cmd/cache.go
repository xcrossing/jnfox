package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfox/util"
)

func init() {
	rootCmd.AddCommand(cmdCache)
}

var cmdCache = &cobra.Command{
	Use:   "cache [nums]",
	Short: "Get Cover from cache first, then from web",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, nums []string) {
		mg, err := util.NewMgInstance()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		defer mg.Close()

		p := util.MakePool(threads, func(num string) {
			doc, err := mg.Fetch(num)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				return
			}
			fmt.Println(doc.PicName())
		})

		for _, num := range nums {
			p.Add(num)
		}

		p.Wait()

		fmt.Println(nums)
	},
}
