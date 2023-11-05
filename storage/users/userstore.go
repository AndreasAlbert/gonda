package storage

type User struct {
	Id       uint
	Provider string
	Name     string
	Email    string
}

type UserStore interface {
	CreateUser(User) (User, error)
	RetrieveUser(string, string) (User, error)
	UpdateUser(User) error
	DeleteUser(User) error
}

// type ApiKey struct {
// 	gorm.Model
// 	User User
// }
