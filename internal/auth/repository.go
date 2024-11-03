package auth

type UserRepository interface {
	GetUserByValue(field, value string) (*User, error)
	CreateUser(
		username string,
		email string,
		password string,
	) (*User, error)
}
