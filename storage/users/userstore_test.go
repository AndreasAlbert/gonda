package storage

import (
	"fmt"
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
	readback_user, err := us.CreateUser(user)

	// ID was not known before, set now for easy comparison
	user.Id = readback_user.Id
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

func testUserUnique(t *testing.T, us UserStore) {
	user := User{
		Name:     "abc",
		Provider: "github",
		Email:    "x@y.z",
	}
	_, err := us.CreateUser(user)
	assert.Nil(t, err)

	// Same name, different provider = OK
	user2 := User{
		Name:     user.Name,
		Provider: "another-provider",
		Email:    "another-email",
	}
	_, err = us.CreateUser(user2)
	assert.Nil(t, err)

	// Different name, same provider = OK
	user3 := User{
		Name:     "another-name",
		Provider: user.Provider,
		Email:    "yet-another-email",
	}
	_, err = us.CreateUser(user3)
	assert.Nil(t, err)

	// Same name, same provider = Not OK
	user4 := User{
		Name:     user.Name,
		Provider: user.Provider,
		Email:    "yet-yet-another-email",
	}
	_, err = us.CreateUser(user4)
	assert.NotNil(t, err)
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
		us, err := NewDBUserStore(dbfunc())
		if err != nil {
			panic(err)
		}

		t.Run(
			fmt.Sprintf("%s/%s", name, "create"),
			func(t *testing.T) {
				testUserCreate(t, us)
			},
		)
		t.Run(
			fmt.Sprintf("%s/%s", name, "unique"),
			func(t *testing.T) {
				testUserUnique(t, us)
			},
		)

	}
}
