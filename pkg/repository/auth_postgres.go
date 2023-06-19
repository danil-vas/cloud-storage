package repository

import (
	"fmt"
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

const UserMemory = 1073741824 // 1 Gb

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user cloud_storage.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, name, username, password_hash, available_memory) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Login, user.Name, user.Username, user.Password, UserMemory)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateMainDirectory(id int, login string) error {
	query := fmt.Sprintf("INSERT INTO %s (name, server_name, type_object_id, user_id, create_date, size) values ($1, $2, $3, $4, $5, $6)", objectsTable)
	_, err := r.db.Exec(query, login, login, 3, id, time.Now(), 0)
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
