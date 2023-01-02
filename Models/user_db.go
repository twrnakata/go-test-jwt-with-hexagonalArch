package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type userRepositoryDB struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) UserRepository {
	return userRepositoryDB{db: db}
}

func (r userRepositoryDB) Create(u User) (*User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO jwtuser (username, password, role) VALUES (?, ?, ?)"
	result, err := tx.Exec(
		query,
		u.Username,
		u.Password,
		u.Role,
	)

	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected <= 0 {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.Id = int(id)
	return &u, nil

}

func (r userRepositoryDB) Find(id int) (*User, error) {
	user := User{}
	query := "SELECT id, username, role FROM jwtuser WHERE id=?"
	err := r.db.Get(&user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r userRepositoryDB) Login(u User) (*User, error) {

	user := User{}
	query := "SELECT id, username, password, role FROM jwtuser WHERE username=?"
	err := r.db.Get(&user, query, u.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))

	if err != nil {
		return nil, err
	}

	return &user, nil

}
