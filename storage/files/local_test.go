package storage

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestLocalFileStore(t *testing.T) {
	tmpDir1 := t.TempDir()

	data1 := []byte("test1")
	data2 := []byte("test2")
	fname1 := "file1.tar"
	fname2 := "file2.root"
	fpath1 := filepath.Join(tmpDir1, fname1)
	fpath2 := filepath.Join(tmpDir1, fname2)

	err := os.WriteFile(fpath1, data1, 0644)
	check(err)
	err = os.WriteFile(fpath2, data2, 0644)
	check(err)

	tmpDir2 := t.TempDir()
	fs := LocalFileStore{Path: tmpDir2}

	// Get without previous put will fail
	_, err = fs.Get(fname1)
	assert.NotNil(t, err)

	// Put a file
	reader, err := os.Open(fpath1)
	check(err)
	err = fs.Put(fname1, reader, false)
	assert.Nil(t, err)

	// Read it back
	reader1, err := fs.Get(fname1)
	buf, err := io.ReadAll(reader1)
	assert.Nil(t, err)
	assert.Equal(t, buf, data1)

	// Get on other name should still fial
	_, err = fs.Get(fname2)
	assert.NotNil(t, err)

	// Put a second file
	reader, err = os.Open(fpath2)
	check(err)
	err = fs.Put(fname2, reader, false)
	assert.Nil(t, err)

	// Read it back
	reader2, err := fs.Get(fname2)
	buf, err = io.ReadAll(reader2)
	assert.Nil(t, err)
	assert.Equal(t, buf, data2)

	// Read the first file again, just to be sure
	reader1, err = fs.Get(fname1)
	buf, err = io.ReadAll(reader1)
	assert.Nil(t, err)
	assert.Equal(t, buf, data1)

	// Try putting again, causing a collision
	// (overwrite is false)
	reader, err = os.Open(fpath1)
	check(err)
	err = fs.Put(fname1, reader, false)
	assert.NotNil(t, err)

	// Now overwrite true with content of fpath2
	reader, err = os.Open(fpath2)
	check(err)
	err = fs.Put(fname1, reader, true)
	assert.Nil(t, err)

	// Read back again, this time compare to data2
	reader1, err = fs.Get(fname1)
	buf, err = io.ReadAll(reader1)
	assert.Nil(t, err)
	assert.Equal(t, buf, data2)
}
