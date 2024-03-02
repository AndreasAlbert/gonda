package conda

import (
	"github.com/AndreasAlbert/gonda/metadata"
)

type PackageFile interface {
	GetIndexEntry() metadata.IndexEntry
}
