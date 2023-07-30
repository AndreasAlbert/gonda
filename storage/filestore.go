package storage

import (
	"io"
)

type FileStore interface {
	Get(string) (io.Reader, error)
	GetLink(string) (string, error)
	Put(string, io.Reader, bool) error
}
