package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var (
		dry      bool
		quiet    bool
		sizeHint int
	)
	flag.BoolVar(&dry, "d", false, "dry mode (no file modification)")
	flag.BoolVar(&quiet, "q", false, "quiet mode")
	flag.IntVar(&sizeHint, "n", 1024, "size hint")
	flag.Parse()
	w := io.Discard
	if !quiet {
		w = os.Stdout
	}
	var r io.Reader = &emptyReader{}
	if fn := flag.Arg(0); len(fn) != 0 {
		f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failure to open file %q\n", fn)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "failure to close file: %v\n", err)
				os.Exit(1)
			}
		}()
		if !dry {
			w = io.MultiWriter(f, w)
		}
		r = f
	}
	if err := Bnew(w, r, os.Stdin, sizeHint); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Bnew writes the unique lines present in r2 but absent from r1 to w, using n
// as a hint for the number of unique lines present in either r1 or r2. If an
// error occurs, Bnew returns it. Bnew buffers reads and writes internally.
func Bnew(w io.Writer, r1, r2 io.Reader, n int) error {
	if n < 0 {
		return fmt.Errorf("negative size hint: %d", n)
	}
	set := make(map[string]struct{}, n)
	var empty struct{}
	sc := bufio.NewScanner(r1)
	for sc.Scan() {
		set[sc.Text()] = empty
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("failure to scan: %v", err)
	}
	bw := bufio.NewWriter(w)
	sc = bufio.NewScanner(r2)
	for sc.Scan() {
		line := sc.Text()
		if _, ok := set[line]; ok {
			continue
		}
		set[line] = empty
		if _, err := fmt.Fprintln(bw, line); err != nil {
			return fmt.Errorf("failure to write line: %v", err)
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("failure to scan: %v", err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("failure to flush: %v", err)
	}
	return nil
}

type emptyReader struct{}

func (*emptyReader) Read(buf []byte) (int, error) {
	return 0, io.EOF
}
