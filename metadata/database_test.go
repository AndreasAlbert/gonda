package metadata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDbSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
func testCreateFileUpload(t *testing.T, mstore DBMetaDataStore) {
	filename := "testfile"
	uploadtime := time.Now().UTC()

	_, err := mstore.GetFileUploadByName(filename)
	assert.Error(t, err)
	mstore.CreateFileUpload(filename, uploadtime)
	upload, err := mstore.GetFileUploadByName(filename)
	assert.NoError(t, err)
	assert.Equal(t, filename, upload.FileName)
	assert.Equal(t, uploadtime, upload.UploadTime)
	assert.Equal(t, Pending, upload.Status)

}
func testUpdateFileUploadStatus(t *testing.T, mstore DBMetaDataStore) {
	filename := "testfile"
	uploadtime := time.Now().UTC()
	assert.NoError(t, mstore.CreateFileUpload(filename, uploadtime))
	assert.NoError(t, mstore.UpdateFileUploadStatus(filename, Done))

	upload, err := mstore.GetFileUploadByName(filename)
	assert.NoError(t, err)
	assert.Equal(t, filename, upload.FileName)
	assert.Equal(t, uploadtime, upload.UploadTime)
	assert.Equal(t, Done, upload.Status)
}

func TestDatabase(t *testing.T) {
	databases := map[string]func() *gorm.DB{
		"sqlite": getDbSqlite,
	}
	tests := []func(*testing.T, DBMetaDataStore){
		testCreateFileUpload,
		testUpdateFileUploadStatus,
	}
	for name, getDb := range databases {
		t.Run(
			name,
			func(t *testing.T) {
				for _, test := range tests {
					db := getDb()
					db.AutoMigrate(&DBFileUpload{})
					mstore := DBMetaDataStore{db: db}
					test(t, mstore)
					db.Migrator().DropTable(&DBFileUpload{})
				}
			},
		)
	}
}
