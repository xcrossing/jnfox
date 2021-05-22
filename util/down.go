package util

import (
	"io"
	"net/http"
	"os"
)

func Download(src string, dest string) error {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Fetch file
	resp, err := client.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create blank file
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write file
	_, err = io.Copy(file, resp.Body)
	return err
}
