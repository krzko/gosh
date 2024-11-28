// internal/shell/builtins/ls.go
package builtins

import (
	"fmt"
	"os"
	"strings"

	"github.com/krzko/gosh/internal/utils/formatter"
	"golang.org/x/term"
)

type LsCommand struct {
	formatter *formatter.TableFormatter
}

type LsOptions struct {
	Long  bool
	All   bool
	Human bool
}

func NewLsCommand() *LsCommand {
	return &LsCommand{
		formatter: formatter.New(),
	}
}

func filterHiddenFiles(entries []formatter.FileInfo) []formatter.FileInfo {
	var filtered []formatter.FileInfo
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name, ".") {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (l *LsCommand) parseOptions(args []string) ([]string, LsOptions, error) {
	var opts LsOptions
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, c := range arg[1:] {
				switch c {
				case 'l':
					opts.Long = true
				case 'a':
					opts.All = true
				case 'h':
					opts.Human = true
				default:
					return nil, opts, fmt.Errorf("unknown option: %c", c)
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}

	return paths, opts, nil
}

func (l *LsCommand) readDirectory(path string) ([]formatter.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []formatter.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		stat := info.Sys()
		owner, group := getOwnerGroup(stat)

		fileInfo := formatter.FileInfo{
			Name:        entry.Name(),
			Size:        info.Size(),
			Mode:        info.Mode(),
			ModTime:     info.ModTime(),
			IsDir:       entry.IsDir(),
			Owner:       owner,
			Group:       group,
			Permissions: info.Mode().String(),
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (l *LsCommand) Execute(args []string) error {
	paths, opts, err := l.parseOptions(args)
	if err != nil {
		return err
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	for i, path := range paths {
		if i > 0 {
			fmt.Printf("\n%s:\n", path)
		}

		entries, err := l.readDirectory(path)
		if err != nil {
			return err
		}

		if !opts.All {
			entries = filterHiddenFiles(entries)
		}

		if opts.Long {
			err = l.formatter.FormatLongList(entries)
		} else {
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				width = 80 // fallback width
			}
			err = l.formatter.FormatCompact(entries, width)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LsCommand) Help() string {
	return `ls - list directory contents

Usage: ls [OPTIONS] [PATH]

Options:
    -l    use long listing format
    -a    show hidden files
    -h    human-readable sizes`
}
