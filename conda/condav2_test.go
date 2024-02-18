package conda

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestV2PackageFile_GetIndexEntry(t *testing.T) {
	file, err := os.Open("testdata/typing-extensions-4.8.0-hd8ed1ab_0.conda")
	assert.NoError(t, err)
	defer file.Close()

	f, err := NewLocalV2PackageFile(file)
	assert.Nil(t, err)
	entry, err := f.GetIndexEntry()
	assert.Nil(t, err)

	assert.Equal(t, entry.Name, "typing-extensions")
	assert.Equal(t, entry.Version, "4.8.0")
	assert.Equal(t, entry.BuildNumber, 0)
	assert.Equal(t, entry.Build, "hd8ed1ab_0")
}
