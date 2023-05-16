package repository

import (
	"fmt"
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user cloud_storage.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, name, username, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Login, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateMainDirectory(id int, login string) error {
	query := fmt.Sprintf("INSERT INTO %s (name, type_object_id, user_id) values ($1, $2, $3)", objectsTable)
	_, err := r.db.Exec(query, login, 3, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) GetUser(login, password string) (cloud_storage.User, error) {
	var user cloud_storage.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, login, password)

	return user, err
}
