package storage

type User struct {
	Provider string
	Name     string
	Email    string
}

type UserStore interface {
	CreateUser(User) error
	RetrieveUser(string, string) (User, error)
	UpdateUser(User) error
	DeleteUser(User) error
}

// type ApiKey struct {
// 	gorm.Model
// 	User User
// }
