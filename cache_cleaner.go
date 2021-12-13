package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var index = map[string]string{
	"node_modules":      "JS",
	"target":            "Rust",
	"_build":            "Erlang,Elxir,OCaml",
	"cmake-build-debug": "C++",
}

const defaultList = "node_modules,target,.cache,_build,.dart_tool"

var (
	dryRun      = flag.Bool("dry-run", true, "only analyze and prepare list")
	dirToRemove = flag.String("patter", "", "list of directories to remove; comma separated")

	rootPath = os.Args[len(os.Args)-1]
)

func _main() {
	flag.Parse()

	out := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(out, "Project	Full Path\n")

	err := filepath.Walk(rootPath,
		func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}

			if info.Name() != ".git" {
				return nil
			}

			fmt.Printf("\r\033[K ... %s", currentPath)
			fmt.Fprintf(out, "%s	%s\n",
				extractParent(currentPath), currentPath)

			return filepath.SkipDir
		})

	if err != nil {
		fmt.Println(err)
	}

}

func __main() {
	flag.Parse()

	if *dirToRemove != "" {
		index = map[string]string{}
		for _, dir := range strings.Split(*dirToRemove, ",") {
			index[dir] = ""
		}
	}

	out := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(out, "Project	Directory	Size	Full Path\n")

	var removeList = map[string]int64{}
	var totalSize int64
	err := filepath.Walk(rootPath,
		func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}

			if _, ok := index[info.Name()]; !ok {
				return nil
			}

			size := calculateDirSize(currentPath)
			totalSize += size
			removeList[currentPath] = size

			fmt.Printf("\r\033[K ... %s	%s 	%s", info.Name(), prettyPrintSize(size), currentPath)
			fmt.Fprintf(out, "%s	%s	%s 	%s\n",
				extractParent(currentPath), info.Name(), prettyPrintSize(size), currentPath)

			return filepath.SkipDir
		})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\r\033[K")
	fmt.Fprintf(out, "\nTotal Size:	%s\n", prettyPrintSize(totalSize))
	out.Flush()

	if *dryRun {
		return
	}

	for dirPath := range removeList {
		fmt.Printf("\r\033[KDeleting: %s", dirPath)
		os.RemoveAll(dirPath)
	}

	fmt.Printf("\r\033[K")
	fmt.Fprintf(out, "\nTotal Space Reclaimed:	%s\n", prettyPrintSize(totalSize))
}

func extractParent(fullPath string) string {
	_, name := path.Split(path.Dir(fullPath))
	return name
}

func calculateDirSize(basePath string) int64 {
	var size int64
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			size += info.Size()
			return nil
		})

	if err != nil {
		fmt.Println(err)
	}
	return size
}

func prettyPrintSize(size int64) string {
	switch {
	case size > 1024*1024*1024:
		return fmt.Sprintf("%.1fG", float64(size)/(1024*1024*1024))
	case size > 1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/(1024*1024))
	case size > 1024:
		return fmt.Sprintf("%.1fK", float64(size)/1024)
	default:
		return fmt.Sprintf("%d", size)
	}

}
