package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"testing"
)

func TestPassthrough(t *testing.T) {
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	header := true
	in := `first_name,last_name,email,comment
name 1,surname 1,name@name1.local,我说中文
name 2,surname 2,name@name2.local,I speak English
name 3,surname 3,name@name3.local,I don't speak
`
	desired := `first_name,last_name,email,comment
name 1,surname 1,name@name1.local,我说中文
name 2,surname 2,name@name2.local,I speak English
name 3,surname 3,name@name3.local,I don't speak
`
	var mcols, dcols columns

	err := cleanCSV(strings.NewReader(in), out, mcols, dcols, header)
	if err != nil {
		t.Error(err)
	}

	if desired != b.String() {
		t.Errorf("unexpected result\nexpected:\n%v\ngot:\n%v", desired, b.String())
	}
}

func TestDeleteColumn(t *testing.T) {
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	header := true
	in := `first_name,last_name,email,comment
name 1,surname 1,name@name1.local,我说中文
name 2,surname 2,name@name2.local,I speak English
name 3,surname 3,name@name3.local,I don't speak
`
	desired := `first_name,last_name,email,comment
name 1,surname 1,,我说中文
name 2,surname 2,,I speak English
name 3,surname 3,,I don't speak
`

	dcols := columns{2}
	mcols := columns{}

	err := cleanCSV(strings.NewReader(in), out, dcols, mcols, header)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if desired != b.String() {
		t.Errorf("unexpected result\nexpected:\n%v\ngot:\n%v", desired, b.String())
	}
}

func TestMaskColumn(t *testing.T) {
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	header := true
	in := `first_name,last_name,email,comment
name 1,surname 1,name@name1.local,我说中文
name 2,surname 2,name@name2.local,I speak English
name 3,surname 3,name@name3.local,I don't speak
`

	dcols := columns{}
	mcols := columns{3}

	err := cleanCSV(strings.NewReader(in), out, dcols, mcols, header)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ir := csv.NewReader(strings.NewReader(in))
	or := csv.NewReader(bytes.NewReader(b.Bytes()))

	for {
		irow, err := ir.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("unexpected error when reading expected data: %v", err)
		}

		orow, err := or.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("unexpected error when reading results of processing: %v", err)
		}

		if header {
			header = false
			continue
		}

		for _, mcol := range mcols {
			if irow[mcol] == orow[mcol] {
				t.Errorf("column not masked: got %v", orow[mcol])
			}
		}
	}
}

func TestMaskAndDeleteColumn(t *testing.T) {
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	in := `first_name,last_name,email,comment
name 1,surname 1,name@name1.local,我说中文
name 2,surname 2,name@name2.local,I speak English
name 3,surname 3,name@name3.local,I don't speak
`
	desired := `first_name,last_name,email,comment
name 1,surname 1,,我说中文
name 2,surname 2,,I speak English
name 3,surname 3,,I don't speak
`

	dcols := columns{2}
	mcols := columns{2}

	err := cleanCSV(strings.NewReader(in), out, dcols, mcols, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if desired != b.String() {
		t.Errorf("unexpected result\nexpected:\n%v\ngot:\n%v", desired, b.String())
	}
}
