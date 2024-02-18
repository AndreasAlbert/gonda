package conda

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestV2PackageFile_GetIndexEntry2(t *testing.T) {
	file, err := os.Open("testdata/pluggy-1.0.0-pyhd8ed1ab_5.tar.bz2")
	assert.NoError(t, err)
	f, err := NewLocalV1PackageFile(file)
	assert.NoError(t, err)
	entry, err := f.GetIndexEntry()
	assert.NoError(t, err)
	assert.Equal(t, "pluggy", entry.Name)
	assert.Equal(t, "1.0.0", entry.Version)
	assert.Equal(t, "pyhd8ed1ab_5", entry.Build)
	assert.Equal(t, 5, entry.BuildNumber)

}
