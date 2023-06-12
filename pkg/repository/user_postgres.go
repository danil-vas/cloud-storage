package repository

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUser(userId int) (cloud_storage.UserInfo, error) {
	var userInfo cloud_storage.UserInfo
	err := r.db.QueryRow("SELECT login, name, username, available_memory FROM users WHERE id=$1", userId).Scan(&userInfo.Login, &userInfo.Name, &userInfo.Username, &userInfo.AvailableMemory)
	if err != nil {
		return cloud_storage.UserInfo{}, err
	}
	userInfo.Id = userId
	return userInfo, nil
}
