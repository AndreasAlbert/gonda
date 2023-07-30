package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type LocalFileStore struct {
	Path string
	FileStore
}

func (fs LocalFileStore) Get(path string) (io.Reader, error) {

	fullPath := filepath.Join(fs.Path, path)
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("File does not exists")
	}

	reader, fileOpenErr := os.Open(fullPath)
	if fileOpenErr != nil {
		return nil, fileOpenErr
	}

	return reader, nil
}
func (fs LocalFileStore) Put(path string, reader io.Reader, overwrite bool) error {

	if !overwrite {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			return errors.New("File exists")
		}
	}
	of, fileCreateErr := os.Create(filepath.Join(fs.Path, path))
	if fileCreateErr != nil {
		return fileCreateErr
	}
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := of.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	return nil
}
func NewLocalFileStore(path string) (LocalFileStore, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(path, os.ModePerm)
	}
	fs := LocalFileStore{
		Path: path,
	}

	return fs, nil
}
