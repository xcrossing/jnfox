package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfo"
	"github.com/xcrossing/jnfox/mdir"
	"github.com/xcrossing/jnfox/util"
)

const ext = ".jpg"

type cache struct {
	bango string

	picName      string
	picCachePath string

	hasDbCache  bool
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
					c.process(mg)
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
	mgDocMap := make(map[string]util.MgDoc)
	for _, doc := range *docs {
		mgDocMap[doc.Bango] = doc
	}

	caches := make([]cache, 0, len(nums))
	for _, num := range nums {
		path, _ := mdir.PathOfName(num, config.Pics.Sep)
		picCachePath := filepath.Join(config.Pics.Root, path, num+ext)

		c := cache{bango: num, picCachePath: picCachePath}

		// hasDbCache
		doc, inDb := mgDocMap[num]
		c.hasDbCache = inDb
		if c.hasDbCache {
			c.picName = doc.PicName()
		}

		// hasPicCache
		macth, _ := filepath.Glob(picCachePath)
		c.hasPicCache = (len(macth) > 0)

		caches = append(caches, c)
	}

	return caches, nil
}

func (c *cache) process(mg *util.MgInstance) {
	var (
		picLink string
		picName string
	)

	if !c.hasDbCache || !c.hasPicCache {
		nfo, err := jnfo.New(config.Host + "/" + c.bango)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s -> %s\n", c.bango, err)
			return
		}
		picLink = nfo.PicLink
		picName = nfo.NumCastPicName()

		if !c.hasDbCache {
			if err := mg.InsertOne(&util.MgDocInsertion{*nfo}); err != nil {
				fmt.Fprintf(os.Stderr, "%s -> %s\n", c.bango, err)
				return
			}
		}
	} else {
		picName = c.picName
	}

	if !c.hasPicCache {
		if err := os.MkdirAll(filepath.Dir(c.picCachePath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "fail to mkdir %s -> %s\n", c.bango, err)
			return
		}
		if err := util.Download(picLink, c.picCachePath); err != nil {
			fmt.Fprintf(os.Stderr, "fail to download %s from %s -> %s\n", c.bango, picLink, err)
			return
		}
	}

	if err := util.Cp(c.picCachePath, picName); err != nil {
		fmt.Fprintf(os.Stderr, "fail to cp %s -> %s\n", c.bango, err)
	}
}
