package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfox/mdir"
	"github.com/xcrossing/jnfox/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	rootCmd.AddCommand(cmdCache)
}

var cmdCache = &cobra.Command{
	Use:   "cache [nums]",
	Short: "Get Cover from cache first, then from web",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, nums []string) {
		picDir := config.Pics
		if picDir == "" {
			fmt.Fprintln(os.Stderr, "no pics config")
			return
		}

		ext := ".jpg"

		mg, err := util.NewMgInstance(config.Mongo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		defer mg.Close()

		p := util.MakePool(threads, func(num string) {
			_, err := mg.Fetch(num)
			if err == mongo.ErrNoDocuments {
				fmt.Fprintf(os.Stderr, "%s : %s\n", num, err.Error())
				return
			}

			path, _ := mdir.PathOfName(num, 3)
			picPath := filepath.Join(picDir, path, num+ext)
			fmt.Println(picPath, num)
		})

		for _, num := range nums {
			p.Add(num)
		}

		p.Wait()
	},
}
