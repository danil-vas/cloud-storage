package repository

import (
	"github.com/jmoiron/sqlx"
	"strconv"
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}

func (r *FilePostgres) PathUploadFile(userId int, objectId int) (string, error) {
	path := ""
	query := "WITH RECURSIVE parent_objects AS (SELECT id, name, parent_id FROM objects WHERE id =" + strconv.Itoa(objectId) + " UNION SELECT o.id," +
		" o.name, o.parent_id FROM objects o INNER JOIN parent_objects p ON p.parent_id = o.id) SELECT name FROM parent_objects;"
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
