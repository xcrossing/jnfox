package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfo"
)

func main() {
	var threads int
	var cmdCover = &cobra.Command{
		Use:   "cover [nums]",
		Short: "Get Cover directly",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, nums []string) {
			host := os.Getenv("JNFOX_HOST")
			u, err := url.Parse(host)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s\n", host, err.Error())
				return
			}
			path := u.Path

			p := makePool(threads, func(addr string) {
				nfo, err := jnfo.New(addr)
				if err == nil {
					err = download(nfo.PicLink, nfo.NumCastPicName())
				}

				if err != nil {
					fmt.Fprintf(os.Stderr, "%s %s\n", addr, err.Error())
				}
			})

			for _, num := range nums {
				u.Path = filepath.Join(path, num)
				addr := u.String()
				p.add(addr)
			}

			p.wait()
		},
	}

	var cmdCacheCover = &cobra.Command{
		Use:   "cache-cover [nums]",
		Short: "Get Cover from cache first, then from web",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, nums []string) {
			mg, err := newMgInstance()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			defer mg.close()

			p := makePool(threads, func(num string) {
				doc, err := mg.fetch(num)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					return
				}
				fmt.Println(doc.picName())
			})

			for _, num := range nums {
				p.add(num)
			}

			p.wait()

			fmt.Println(nums)
		},
	}

	var rootCmd = &cobra.Command{Use: "jnfox"}
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 2, "threads to get cover")
	rootCmd.AddCommand(cmdCover, cmdCacheCover)
	rootCmd.Execute()
}
