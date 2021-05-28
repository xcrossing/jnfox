package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfo"
	"github.com/xcrossing/jnfox/crawler"
	"github.com/xcrossing/jnfox/mdir"
	"github.com/xcrossing/jnfox/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	rootCmd.AddCommand(cmdSync)
}

var cmdSync = &cobra.Command{
	Use:   "sync [code]",
	Short: "Sync jnfo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, path []string) {
		w := crawler.Walk{
			Uri:  config.Host + "/" + path[0],
			Next: "#next",
			Item: "a.movie-box",
			Cookies: []*http.Cookie{
				&http.Cookie{Name: "existmag", Value: "all"},
			},
		}

		// mongo client init
		mg, err := util.NewMgInstance(config.Mongo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		defer mg.Close()

		// fetch info in threads
		p := util.MakeFuncPool(threads)
		w.Start(func(e *crawler.Element) {
			p.Add(func(uri string) func() {
				return func() {
					syncNum(uri, mg)
				}
			}(e.Attr("href")))
		})
		p.Wait()
	},
}

func syncNum(uri string, mg *util.MgInstance) {
	<-time.After(2 * time.Second)

	u, _ := url.Parse(uri)
	num := u.Path[1:]

	_, err := mg.Fetch(num)
	if err != mongo.ErrNoDocuments {
		fmt.Println(time.Now(), num, "found")
		return
	}

	nfo, err := jnfo.New(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	// down cover
	path, _ := mdir.PathOfName(nfo.Num, config.Pics.Sep)
	picCachePath := filepath.Join(config.Pics.Root, path, nfo.Num+ext)
	if err := util.Download(nfo.PicLink, picCachePath); err != nil {
		fmt.Fprintf(os.Stderr, "%s -> %s\n", nfo.Num, err)
		return
	}

	// save info
	if err := mg.InsertOne(nfo); err != nil {
		fmt.Fprintf(os.Stderr, "%s -> %s\n", nfo.Num, err)
		return
	}

	fmt.Println(time.Now(), num, "synced => "+picCachePath)
}
