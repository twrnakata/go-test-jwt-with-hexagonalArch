package models

type UserRepository interface {
	Create(User) (*User, error)
	Find(int) (*User, error)
	Login(User) (*User, error)
}

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`
}
