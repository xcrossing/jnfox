package util

import (
	"fmt"
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

func Cp(src string, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
