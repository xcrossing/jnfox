package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xcrossing/jnfo"
	"github.com/xcrossing/jnfox/pool"
)

func main() {
	var coverThreads int
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

			p := pool.MakePool(coverThreads, func(addr string) {
				nfo, err := jnfo.New(addr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s %s\n", addr, err.Error())
				} else {
					fmt.Println(nfo.NumCastPicName())
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

	cmdCover.Flags().IntVarP(&coverThreads, "threads", "t", 2, "threads to get cover")

	var rootCmd = &cobra.Command{Use: "jnfox"}
	rootCmd.AddCommand(cmdCover)
	rootCmd.Execute()
}
