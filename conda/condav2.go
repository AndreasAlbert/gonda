package conda

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/AndreasAlbert/gonda/metadata"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
	"io"
	"os"
	"strings"
)

type V2PackageFile struct {
	PackageFile
	reader io.ReaderAt
	size   int64
}

func (f V2PackageFile) getInfoReader() (*tar.Reader, error) {

	reader, err := zip.NewReader(f.reader, f.size)
	if err != nil {
		return nil, err
	}
	// Step 1: Find the info*.tar.zst file
	var infoPaths []string
	for _, file := range reader.File {
		fmt.Println(file.Name)
		if strings.HasPrefix(file.Name, "info-") && strings.HasSuffix(file.Name, ".tar.zst") {
			infoPaths = append(infoPaths, file.Name)
		}
	}
	if len(infoPaths) != 1 {
		panic("x")
	}
	infoPath := infoPaths[0]

	// Step 2: Open the info*.tar.zst file attach a zstd reader
	infoFile, err := reader.Open(infoPath)
	if err != nil {
		return nil, err
	}
	zReader, err := zstd.NewReader(infoFile)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(zReader)
	return tarReader, nil
}

func (f V2PackageFile) readInfoFile(name string) ([]byte, error) {
	tarReader, err := f.getInfoReader()
	if err != nil {
		return nil, err
	}
	return ReadFileFromTar(tarReader, name, 1e6)
}

func (f V2PackageFile) GetIndexEntry() (metadata.IndexEntry, error) {
	entry := metadata.IndexEntry{}
	b, err := f.readInfoFile("info/index.json")
	if err != nil {
		return metadata.IndexEntry{}, err
	}
	err = json.Unmarshal(b, &entry)
	return entry, err
}

func NewLocalV2PackageFile(file *os.File) (V2PackageFile, error) {
	finfo, err := file.Stat()
	if err != nil {
		return V2PackageFile{}, err
	}
	f := V2PackageFile{reader: file, size: finfo.Size()}
	return f, nil
}
