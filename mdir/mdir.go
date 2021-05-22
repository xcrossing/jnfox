package mdir

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

type void struct{}

var pathMember void

type Cmd struct {
	CopyFile bool
	Force    bool
	Progress bool
	Src      string
	Dest     string
	Segments []int

	_destRoot []string
	_paths    map[string]void
}

func (cmd *Cmd) MvFiles() error {
	if cmd.Src == "" || cmd.Dest == "" {
		return errors.New("no src or dest")
	}

	list, err := listFiles(cmd.Src)
	if err != nil {
		return err
	}

	// batch create dirs once
	cmd._paths = make(map[string]void)
	for _, file := range list {
		if err := cmd.generateNewPath(file); err != nil {
			return err
		}
	}

	action := cmd.action()

	// progress bar
	var bar progress
	if cmd.Progress {
		bar = &realProgress{pb.StartNew(len(list))}
	} else {
		bar = &fakeProgress{}
	}

	// mv files
	for _, file := range list {
		if err := action(file.oldPath, file.newPath); err != nil {
			return err
		}
		bar.increment()
	}
	bar.finish()

	return nil
}

func (cmd *Cmd) action() func(src string, dest string) error {
	if !cmd.Force {
		return _dryRunFiles
	} else if cmd.CopyFile {
		return _cpFiles
	} else {
		return _mvFiles
	}
}

func (cmd *Cmd) destRoot() []string {
	if cmd._destRoot != nil {
		return cmd._destRoot
	}
	// add two more cap for root and file
	root := make([]string, 0, len(cmd.Segments)+2)
	root = append(root, cmd.Dest)
	cmd._destRoot = root
	return root
}

func (cmd *Cmd) generateNewPath(file *fileInfo) error {
	md5path := md5Str(file.baseNameNoExt)
	dirs, err := splitStr(md5path, cmd.Segments...)
	if err != nil {
		return err
	}

	newDir, err := cmd.mkdirs(append(cmd.destRoot(), dirs...))
	if err != nil {
		return err
	}

	file.newPath = filepath.Join(newDir, file.baseName)
	return nil
}

func (cmd *Cmd) mkdirs(path []string) (string, error) {
	pathStr := filepath.Join(path...)
	if cmd.Force {
		if _, exists := cmd._paths[pathStr]; !exists {
			if err := os.MkdirAll(pathStr, os.ModePerm); err != nil {
				return "", err
			}
		}
	}
	cmd._paths[pathStr] = pathMember
	return pathStr, nil
}
