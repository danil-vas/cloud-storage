package repository

import (
	"encoding/json"
	"fmt"
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
	"time"
)

type DirectoryPostgres struct {
	db *sqlx.DB
}

var s *DirectoryPostgres

func NewDirectoryPostgres(db *sqlx.DB) *DirectoryPostgres {
	return &DirectoryPostgres{db: db}
}

func (r *DirectoryPostgres) AddDirectory(userId int, objectId int, nameDirectory string) (int, error) {
	var objId int
	query := fmt.Sprintf("INSERT INTO %s (name, server_name, create_date, user_id, parent_id, type_object_id, size) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id", objectsTable)
	row := r.db.QueryRow(query, nameDirectory, nameDirectory, time.Now(), userId, objectId, 2, 0)
	if err := row.Scan(&objId); err != nil {
		return 0, err
	}
	return objId, nil
}

func (r *DirectoryPostgres) GetDirectoriesAndFiles(userId int, objectId int) ([]cloud_storage.Node, error) {
	resp := make([]*cloud_storage.Node, 0)
	mainResp := make([]cloud_storage.Node, 0)
	firstQuery := fmt.Sprintf("SELECT id, name, server_name, size, create_date, type_object_id FROM objects WHERE id=$1")
	firstRow, err := r.db.Query(firstQuery, objectId)
	if err != nil {
		return nil, err
	}
	defer firstRow.Close()
	for firstRow.Next() {
		var root cloud_storage.Node
		if err := firstRow.Scan(&root.ID, &root.Name, &root.ServerName, &root.Size, &root.CreateDate, &root.Type); err != nil {
			return nil, err
		}
		mainResp = append(mainResp, root)
	}

	query := fmt.Sprintf("WITH RECURSIVE nested_objects AS (\n            SELECT id, name, server_name, size, create_date, user_id, parent_id, type_object_id\n            FROM %s\n            WHERE parent_id = $1\n            UNION ALL\n            SELECT o.id, o.name, o.server_name, o.size, o.create_date, o.user_id, o.parent_id, o.type_object_id\n            FROM %s o\n            JOIN nested_objects n ON o.parent_id = n.id\n        )\n        SELECT json_build_object(\n            'id', no.id,\n            'name', no.name,\n            'server_name', no.server_name,\n            'size', no.size,\n            'create_date', no.create_date,\n            'type', (CASE WHEN no.type_object_id = 1 THEN 'file' ELSE 'directory' END),\n            'children', (\n                SELECT json_agg(json_build_object(\n                    'id', child.id,\n                    'name', child.name,\n                    'server_name', child.server_name,\n                    'size', child.size,\n                    'create_date', child.create_date,\n                    'type', (CASE WHEN child.type_object_id = 1 THEN 'file' ELSE 'directory' END),\n                    'children', (\n                        SELECT json_agg(json_build_object(\n                            'id', sub_child.id,\n                            'name', sub_child.name,\n                            'server_name', sub_child.server_name,\n                            'size', sub_child.size,\n                            'create_date', sub_child.create_date,\n                            'type', (CASE WHEN sub_child.type_object_id = 1 THEN 'file' ELSE 'directory' END),\n                            'children', '[]'\n                        ))\n                        FROM nested_objects sub_child\n                        WHERE sub_child.parent_id = child.id\n                    ))\n                ) FROM nested_objects child\n                WHERE child.parent_id = no.id\n            ) \n        )\n        FROM nested_objects no\n        WHERE no.parent_id = $1\n        GROUP BY no.id, no.name, no.server_name, no.size, no.create_date, no.type_object_id;", objectsTable, objectsTable)
	rows, err := r.db.Query(query, objectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var nestedObjectJSON string
		if err := rows.Scan(&nestedObjectJSON); err != nil {
			return mainResp, err
		}
		var root *cloud_storage.Node
		err = json.Unmarshal([]byte(nestedObjectJSON), &root)
		resp = append(resp, root)
	}
	if err := rows.Err(); err != nil {
		return mainResp, err
	}
	mainResp[0].Children = resp
	return mainResp, nil
}

func (r *DirectoryPostgres) DeleteDirectory(objectId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 ", objectsTable)
	_, err := r.db.Exec(query, objectId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DirectoryPostgres) GetIdMainDirectory(userId int) (int, error) {
	var objId int
	query := fmt.Sprintf("SELECT id FROM objects WHERE user_id = $1 AND type_object_id = 3")
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&objId); err != nil {
		return 0, err
	}
	return objId, nil
}
