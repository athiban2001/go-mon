package watch

import (
	"fmt"
	"io/fs"
	"strconv"
	"testing"
)

func TestIsValidDec(t *testing.T) {
	type isValidData struct {
		filename string
		isDir    bool
		output   bool
	}

	tests := []struct {
		name           string
		ignoreDotFiles bool
		extensions     []string
		data           []isValidData
	}{
		{"default values", true, []string{".go"},
			[]isValidData{{"1.go", false, true}, {".git", true, false}}},
		{"multiple extensions and dot files", false, []string{".go", ".html"},
			[]isValidData{{"1.go", false, true}, {"1.html", false, true}, {".git", true, true}, {"1.css", false, false}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := isValidDecorator(test.ignoreDotFiles, test.extensions)
			for _, subtest := range test.data {
				if subtest.output != isValid(subtest.filename, subtest.isDir) {
					t.Errorf("Expected %v, For {IgnoreDotFiles : %v,Extensions : %v,subtest filename : %v}",
						subtest.output, test.ignoreDotFiles, test.extensions, subtest.filename)
				}
			}
		})
	}
}

type entry struct {
	name  string
	fType fs.FileMode
	isDir bool
}

func (e *entry) Name() string {
	return e.name
}

func (e *entry) IsDir() bool {
	return e.isDir
}

func (e *entry) Type() fs.FileMode {
	return e.fType
}

func (e *entry) Info() (fs.FileInfo, error) {
	return nil, fmt.Errorf("You should not call this method")
}

func TestRemoveEntry(t *testing.T) {
	entries := [10]fs.DirEntry{}
	entriesLen, i := 10, 0
	for i = 0; i < entriesLen; i++ {
		entries[i] = &entry{name: strconv.FormatInt(int64(i), 10)}
	}
	middleOutput := make([]fs.DirEntry, 0)
	middleOutput = append(middleOutput, entries[:entriesLen/2]...)
	middleOutput = append(middleOutput, entries[(entriesLen/2)+1:]...)

	tests := []struct {
		name    string
		index   int
		entries []fs.DirEntry
		output  []fs.DirEntry
	}{
		{"beginning", 0, entries[:], entries[1:]},
		{"middle", entriesLen / 2, entries[:], middleOutput},
		{"end", entriesLen - 1, entries[:], entries[:entriesLen-1]},
	}

	for k, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := removeEntry(test.entries, test.index)
			if len(output) != len(test.output) {
				t.Errorf("Expected Length : %d, Received Length : %d", len(test.output), len(output))
			}
			for i = 0; i < len(output); i++ {
				if test.output[i] != output[i] {
					t.Errorf("#%d Case, Expected : %v, Received : %v", k+1, test.output[i], output[i])
				}
			}
		})
	}
}
