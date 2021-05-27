package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfox/crawler"
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

		links := []string{}

		w.Start(func(e *crawler.Element) {
			links = append(links, e.Attr("href"))
		})

		fmt.Println(links)
	},
}
