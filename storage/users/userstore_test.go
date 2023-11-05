package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testUserCreate(t *testing.T, us UserStore) {
	user := User{
		Name:     "abc",
		Provider: "github",
		Email:    "x@y.z",
	}

	// User should not exist initially
	_, err := us.RetrieveUser(user.Provider, user.Name)
	assert.NotNil(t, err)

	// Create and assert that readback succeeds
	us.CreateUser(user)
	readback_user, err := us.RetrieveUser(user.Provider, user.Name)
	assert.Nil(t, err)
	assert.Equal(t, user, readback_user)

	// Update
	user.Email = "xyz"
	err = us.UpdateUser(user)
	assert.Nil(t, err)
	readback_user, err = us.RetrieveUser(user.Provider, user.Name)
	assert.Nil(t, err)
	assert.Equal(t, user, readback_user)

	// Delete
	err = us.DeleteUser(user)
	assert.Nil(t, err)

}

func getDbSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&DBUser{})
	return db
}
func TestDBUserStore(t *testing.T) {
	databases := map[string]func() *gorm.DB{
		"sqlite": getDbSqlite,
	}

	for name, dbfunc := range databases {
		t.Run(
			name,
			func(t *testing.T) {
				testUserCreate(t, NewDBUserStore(dbfunc())})
			},
		)

	}
}
