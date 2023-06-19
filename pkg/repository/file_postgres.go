package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}

func (r *FilePostgres) PathUploadFile(userId int, objectId int) (string, error) {
	path := ""
	query := "WITH RECURSIVE parent_objects AS (SELECT id, server_name, parent_id FROM objects WHERE id =" + strconv.Itoa(objectId) + " UNION SELECT o.id," +
		" o.server_name, o.parent_id FROM objects o INNER JOIN parent_objects p ON p.parent_id = o.id) SELECT server_name FROM parent_objects;"
	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	i := 0
	for rows.Next() {
		str := ""
		if err := rows.Scan(&str); err != nil {
			return "", err
		}
		if i == 0 {
			path = str + path
		} else {
			path = str + "/" + path
		}
		i++
	}
	return path, nil
}

func (r *FilePostgres) AddUploadFileToUser(userId int, objectId int, originalName string, serverName string, size int64, create_time time.Time) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, server_name, size, create_date, user_id, parent_id, type_object_id) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id", objectsTable)
	//_, err := r.db.Exec(query, originalName, serverName, size, create_time, userId, objectId, 1)
	row := r.db.QueryRow(query, originalName, serverName, size, create_time, userId, objectId, 1)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	query = fmt.Sprintf("UPDATE users SET available_memory = available_memory - $1 WHERE id = $2")
	_, err := r.db.Exec(query, size, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *FilePostgres) OriginalFileName(objectId int) (string, error) {
	var name string
	query := fmt.Sprintf("SELECT name FROM %s WHERE id = $1", objectsTable)
	row := r.db.QueryRow(query, objectId)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func (r *FilePostgres) DeleteFile(userId int, objectId int) error {
	var sizeFile int
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING size", objectsTable)
	row := r.db.QueryRow(query, objectId)
	if err := row.Scan(&sizeFile); err != nil {
		return err
	}
	query = fmt.Sprintf("UPDATE users SET available_memory = available_memory + $1 WHERE id = $2")
	_, err := r.db.Exec(query, sizeFile, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *FilePostgres) GetAvailableMemory(userId int) (int, error) {
	var availableMemory int
	query := fmt.Sprintf("SELECT available_memory FROM users WHERE id = $1")
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&availableMemory); err != nil {
		return 0, err
	}
	return availableMemory, nil
}

func (r *FilePostgres) CheckAccessToObject(userId int, objectId int) (bool, error) {
	var objectUserID int
	err := r.db.QueryRow("SELECT user_id FROM objects WHERE id=$1", objectId).Scan(&objectUserID)
	if err != nil {
		return false, nil
	}
	if objectUserID == userId {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *FilePostgres) GetTypeObject(objectId int) (string, error) {
	var typeFile int
	err := r.db.QueryRow("SELECT type_object_id FROM objects WHERE id=$1", objectId).Scan(&typeFile)
	if err != nil {
		return "", err
	}
	if typeFile == 1 {
		return "file", nil
	} else {
		return "directory", nil
	}
}

func (r *FilePostgres) OriginalFileNameThroughServerName(serverName string) (string, error) {
	var name string
	query := fmt.Sprintf("SELECT name FROM %s WHERE server_name = $1", objectsTable)
	row := r.db.QueryRow(query, serverName)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}
