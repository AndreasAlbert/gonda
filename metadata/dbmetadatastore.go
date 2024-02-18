package metadata

import (
	"time"

	"gorm.io/gorm"
)

type DBFileUpload struct {
	gorm.Model
	FileName   string `gorm:"primaryKey; not null"`
	UploadTime time.Time
	Status     UploadStatus
}

type DBMetaDataStore struct {
	db *gorm.DB
	MetaDataStore
}

func NewDBMetaDataStore(db *gorm.DB) (DBMetaDataStore, error) {
	err := db.AutoMigrate(&DBFileUpload{})
	if err != nil {
		return DBMetaDataStore{}, err
	}
	return DBMetaDataStore{db: db}, nil
}

func (mstore DBMetaDataStore) CreateFileUpload(FileName string, UploadTime time.Time) error {
	upload := DBFileUpload{
		FileName:   FileName,
		UploadTime: UploadTime,
		Status:     Pending,
	}
	tx := mstore.db.Begin()
	defer tx.Rollback()
	err := tx.Create(&upload).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
func (mstore DBMetaDataStore) GetFileUploadByName(FileName string) (FileUpload, error) {
	tx := mstore.db.Begin()
	defer tx.Rollback()
	return mstore.getFileUploadByName(FileName, tx)
}
func (mstore DBMetaDataStore) getFileUploadByName(FileName string, tx *gorm.DB) (FileUpload, error) {
	upload := DBFileUpload{FileName: FileName}
	err := tx.First(&upload).Error
	if err != nil {
		return FileUpload{}, err
	}
	return FileUpload{
		FileName:   upload.FileName,
		UploadTime: upload.UploadTime,
		Status:     upload.Status,
	}, nil
}
func (mstore DBMetaDataStore) UpdateFileUploadStatus(FileName string, Status UploadStatus) error {
	return mstore.db.Connection(
		func(tx *gorm.DB) error {
			upload, err := mstore.getFileUploadByName(FileName, tx)
			if err != nil {
				return err
			}
			upload.Status = Status
			err = tx.Save(&upload).Error
			if err != nil {
				return err
			}
			tx.Commit()
			return nil
		},
	)
}
