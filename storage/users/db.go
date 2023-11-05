package storage

import "gorm.io/gorm"

type DBUser struct {
	Id   uint `gorm:"primaryKey"`
	User User `gorm:"embedded"`
}
type DBUserStore struct {
	UserStore
	db *gorm.DB
}

func (dbus DBUserStore) AutoMigrate() error {
	err := dbus.db.AutoMigrate(&DBUser{})
	return err
}
func NewDBUserStore(db *gorm.DB) DBUserStore {
	dbus := DBUserStore{db: db}
	dbus.AutoMigrate()
	return dbus
}
func (dbus DBUserStore) CreateUser(user User) error {
	dbuser := DBUser{User: user}
	dbus.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbuser).Error; err != nil {
			return err
		}
		return nil
	})
	return nil
}
func (dbus DBUserStore) RetrieveUser(provider string, name string) (User, error) {
	dbuser, err := dbus.retrieveDBUser(provider, name)
	return dbuser.User, err
}
func (dbus DBUserStore) retrieveDBUser(provider string, name string) (DBUser, error) {
	var dbuser DBUser
	err := dbus.db.First(&dbuser, "Provider = ? AND Name = ?", provider, name).Error
	return dbuser, err
}

func (dbus DBUserStore) UpdateUser(user User) error {
	dbuser, err := dbus.retrieveDBUser(user.Provider, user.Name)
	if err != nil {
		return err
	}
	dbuser.User = user
	err = dbus.db.Save(&dbuser).Error
	return err
}

func (dbus DBUserStore) DeleteUser(user User) error {
	dbuser, err := dbus.retrieveDBUser(user.Provider, user.Name)
	if err != nil {
		return err
	}
	err = dbus.db.Delete(&dbuser).Error
	return err
}
