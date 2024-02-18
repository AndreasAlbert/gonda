package conda

import (
	"archive/tar"
	"compress/bzip2"
	"encoding/json"
	"github.com/AndreasAlbert/gonda/metadata"
	"io"
	"os"
)

type V1PackageFile struct {
	PackageFile
	reader io.Reader
	size   int64
}

func (f V1PackageFile) readInfoFile(name string) ([]byte, error) {
	tarReader := tar.NewReader(bzip2.NewReader(f.reader))
	return ReadFileFromTar(tarReader, name, 1e6)
}
func (f V1PackageFile) GetIndexEntry() (metadata.IndexEntry, error) {
	entry := metadata.IndexEntry{}
	b, err := f.readInfoFile("info/index.json")
	if err != nil {
		return entry, nil
	}
	err = json.Unmarshal(b, &entry)
	return entry, nil
}
func NewLocalV1PackageFile(file *os.File) (V1PackageFile, error) {
	finfo, err := file.Stat()
	if err != nil {
		return V1PackageFile{}, err
	}
	f := V1PackageFile{reader: file, size: finfo.Size()}
	return f, nil
}
