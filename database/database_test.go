package database

import (
	"testing"

	"github.com/go-playground/assert/v2"
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

func testUser(t *testing.T, db *gorm.DB) {
	tx := db.Begin()
	defer tx.Rollback()
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
