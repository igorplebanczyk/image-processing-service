package auth

type UserRepository interface {
	CreateUser(user *User) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
}
