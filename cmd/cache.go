package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfox/mdir"
	"github.com/xcrossing/jnfox/util"
)

const ext = ".jpg"

type cache struct {
	bango       string
	picName     string
	picPath     string
	inDb        bool
	hasPicCache bool
}

func init() {
	rootCmd.AddCommand(cmdCache)
}

var cmdCache = &cobra.Command{
	Use:   "cache [nums]",
	Short: "Get Cover from cache first, then from web",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, nums []string) {
		if config.Pics.Root == "" || config.Pics.Sep == 0 {
			fmt.Fprintln(os.Stderr, "no pics config")
			return
		}

		mg, err := util.NewMgInstance(config.Mongo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		defer mg.Close()

		caches, err := checkCache(mg, nums)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}

		p := util.MakeFuncPool(threads)
		for _, c := range caches {
			p.Add(func(c cache) func() {
				return func() {
					fmt.Println(c)
				}
			}(c))
		}

		p.Wait()
	},
}

func checkCache(mongo *util.MgInstance, nums []string) ([]cache, error) {
	docs, err := mongo.BatchFetch(nums)
	if err != nil {
		return nil, err
	}
	mgDocMap := make(map[string]*util.MgDoc)
	for _, doc := range *docs {
		mgDocMap[doc.Bango] = &doc
	}

	caches := make([]cache, 0, len(nums))
	for _, num := range nums {
		path, _ := mdir.PathOfName(num, config.Pics.Sep)
		picPath := filepath.Join(config.Pics.Root, path, num+ext)

		macth, _ := filepath.Glob(picPath)
		hasPicCache := (len(macth) > 0)

		_, inDb := mgDocMap[num]

		c := cache{
			bango:       num,
			picPath:     picPath,
			inDb:        inDb,
			hasPicCache: hasPicCache,
		}

		caches = append(caches, c)
	}

	return caches, nil
}
