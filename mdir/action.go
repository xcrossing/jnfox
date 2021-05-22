package mdir

import (
	"fmt"
	"io"
	"os"
)

func _mvFiles(src string, dest string) error {
	return os.Rename(src, dest)
}

func _cpFiles(src string, dest string) error {
	oldFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	neoFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer neoFile.Close()

	_, err = io.Copy(neoFile, oldFile)
	if err != nil {
		return err
	}

	return nil
}

func _dryRunFiles(src string, dest string) error {
	fmt.Printf("%s -> %s\n", src, dest)
	return nil
}
