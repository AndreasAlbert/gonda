package conda

import (
	"archive/tar"
	"fmt"
	"io"
)

func ReadFileFromTar(reader *tar.Reader, name string, maxSize int64) ([]byte, error) {
	for {
		hdr, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if hdr.Name != name {
			continue
		}
		if maxSize > 0 && hdr.FileInfo().Size() > maxSize {
			panic("")
		}

		fmt.Println(hdr.Name)
		b, err := io.ReadAll(reader)
		return b, err
	}
	panic("")
}
