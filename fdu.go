package main

import (
	"fmt"
	"os"
	"path/filepath"
	"flag"
)

func humanSize(size int64) string {
	units := []string{"", "KB", "MB", "GB", "TB"}
	var offset int
	for size > 1024 && offset < len(units) {
		size /= 1024
		offset++
	}
	return fmt.Sprintf("%d %s", size, units[offset])
}

func myGlob(name string) chan string {
	yield := make(chan string)
	go func() {
		ret, _ := filepath.Glob(name)
		for _, f := range ret {
			if string(name[0]) != "." && string(f[0]) == "." {
				continue
			}
			yield <- f
		}
		close(yield)
	}()
	return yield
}

func isRealDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return false
	}
	if s.IsDir() && (s.Mode()&os.ModeSymlink == 0) {
		return true
	}
	return false
}

func fileSize(fname string) int64 {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s\n", fname)
		return 0
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't stat %s\n", fname)
	}
	return s.Size()
}

func write(fname string, params params) int64 {
	params.args = []string{fname + "/*"}
	size := fdu(true, params)

	sname := fname + "/.SIZE"
	sizeFile, err := os.Create(sname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create %s\n", sname)
		return size
	}
	fmt.Fprintf(sizeFile, "%d\n", size)
	return size
}

func readInt(f *os.File) int64 {
	var size int64
	_, err := fmt.Fscanf(f, "%d\n", &size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read a file\n")
	}
	return size
}

type params struct {
	update *bool
	args   []string
}

func fdu(silent bool, params params) int64 {
	var totalSize int64
	for _, fnames := range params.args {
		for fname := range myGlob(fnames) {
			var size int64
			if isRealDir(fname) {
				if *params.update {
					size = write(fname, params)
				} else {
					f, err := os.Open(fname + "/.SIZE")
					if err != nil {
						size = write(fname, params)
					} else {
						size = readInt(f)
					}
				}
				indexLast := len(fname) - 1
				if fname[indexLast] != '/' {
					fname = fname + "/"
				}

			} else { // NOT dir
				size = fileSize(fname)
			}
			if !silent {
				fmt.Printf("%s %v\n", fname, humanSize(size))
			}
			totalSize += size
		}
	}
	if !silent {
		fmt.Printf("total: %v\n", humanSize(totalSize))
	}
	return totalSize
}

func main() {
	var params params
	params.update = flag.Bool("update", false, "recompute every .SIZE")
	flag.Parse()
	params.args = flag.Args()

	fdu(false, params)
}

// vim: set ts=2 sw=2 list ft=go:
