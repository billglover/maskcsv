package main

import (
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type columns []int

func (c *columns) String() string {
	return fmt.Sprintf("%v$", *c)
}

func (c *columns) Set(s string) error {
	a := strings.Split(s, ",")
	for _, vs := range a {
		n, err := strconv.Atoi(vs)
		if err != nil {
			return fmt.Errorf("unable to parse column number: %v", err)
		}
		*c = append(*c, n)
	}

	return nil
}

func main() {

	var dcols, mcols columns

	flag.Var(&dcols, "d", "comma separated list of columns to remove (remove takes precidence over mask)")
	flag.Var(&mcols, "m", "comma separated list of columns to mask")
	in := flag.String("i", "", "input file (csv formatted)")
	out := flag.String("o", "", "output file (csv formatted)")
	header := flag.Bool("header", false, "indicate whether the csv containers a header row")

	flag.Parse()

	inf, err := openInput(*in)
	if err != nil {
		log.Fatalf("unable to open input file: %v", err)
	}
	defer inf.Close()

	outf, err := openOutput(*out)
	if err != nil {
		log.Fatalf("unable to open output file: %v", err)
	}
	defer outf.Close()

	err = cleanCSV(inf, outf, dcols, mcols, *header)
	if err != nil {
		log.Fatalf("unable to process CSV: %v", err)
	}
}

func openInput(path string) (f *os.File, err error) {
	if path == "" {
		return os.Stdin, nil
	}

	f, err = os.Open(path)
	return
}

func openOutput(path string) (f *os.File, err error) {
	if path == "" {
		return os.Stdout, nil
	}

	f, err = os.Create(path)
	return
}

func cleanCSV(inf io.Reader, outf io.Writer, dcols, mcols columns, header bool) error {

	salt := time.Time.String(time.Now())

	r := csv.NewReader(inf)
	w := csv.NewWriter(outf)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to read record: %v", err)
		}

		if header {
			err := w.Write(record)
			if err != nil {
				return fmt.Errorf("unable to write header record: %v", err)
			}
			header = false
			continue
		}

		err = maskCols(mcols, &record, salt)
		if err != nil {
			return fmt.Errorf("unable to mask columns: %v", err)
		}

		err = deleteCols(dcols, &record)
		if err != nil {
			return fmt.Errorf("unable to delete columns: %v", err)
		}

		err = w.Write(record)
		if err != nil {
			return fmt.Errorf("unable to write record: %v", err)
		}
	}
	w.Flush()
	err := w.Error()
	if err != nil {
		return fmt.Errorf("unable to flush output: %v", err)
	}

	return nil
}

func deleteCols(dcols columns, r *[]string) error {
	for _, c := range dcols {
		if c < 0 || c >= len(*r) {
			return fmt.Errorf("column for deletion is out of bounds: %d is not in the range [%d,%d]", c, 0, len(*r))
		}

		(*r)[c] = ""
	}
	return nil
}

func maskCols(mcols columns, r *[]string, salt string) error {
	for _, c := range mcols {
		if c < 0 || c >= len(*r) {
			return fmt.Errorf("column for deletion is out of bounds: %d is not in the range [%d,%d]", c, 0, len(*r))
		}

		val := []string{salt, (*r)[c]}
		sum := sha256.Sum256([]byte(strings.Join(val, "|")))
		(*r)[c] = fmt.Sprintf("%x", sum)
	}

	return nil
}
