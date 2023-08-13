package database

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDbSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	return db
}
func testUserCreate(t *testing.T, db *gorm.DB) {
	tx := db.Begin()
	defer tx.Rollback()

	// Create a single user
	user_in := User{
		Provider: "github",
		Name:     "U1",
	}
	tx.Create(&user_in)

	var user_out User

	tx.First(&user_out, "Provider = ? AND Name = ?", user_in.Provider, user_in.Name)
	assert.Equal(t, user_in.Provider, user_out.Provider)
	assert.Equal(t, user_in.Name, user_out.Name)
}

func testUserUnique(t *testing.T, db *gorm.DB) {
	tx := db.Begin()
	defer tx.Rollback()

	// First user
	user_1 := User{
		Provider: "github",
		Name:     "U1",
	}
	tx.Create(&user_1)

	// Same provider, other name is OK
	user_2 := User{
		Provider: "github",
		Name:     "U2",
	}
	tx.Create(&user_2)

	// Name collision!
	user_3 := User{
		Provider: "github",
		Name:     "U2",
	}
	err := tx.Create(&user_3).Error
	assert.NotNil(err)
}
func testUser(t *testing.T, db *gorm.DB) {
	t.Run(
		"create", func(t *testing.T) { testUserCreate(t, db) })
	t.Run(
		"unique", func(t *testing.T) { testUserCreate(t, db) })
}

func TestDatabase(t *testing.T) {
	databases := map[string]*gorm.DB{
		"sqlite": getDbSqlite(),
	}

	for name, db := range databases {
		t.Run(
			name,
			func(t *testing.T) {
				testUser(t, db)
			},
		)

	}
}
