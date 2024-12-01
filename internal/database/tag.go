package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TagModel struct {
	DB  *pgx.Conn
	Tag Tag
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TagCreation struct {
	Name string `json:"name"`
}

func (t *TagModel) Create(name string) (int64, error) {
	query := "INSERT INTO tags (name) VALUES (?) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id) RETURNING id"
	var newId int64
	err := t.DB.QueryRow(context.Background(), query, name).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (t *TagModel) ByName(name string) (Tag, error) {
	var tag Tag
	query := "SELECT id, name FROM tags WHERE name = ?"
	err := t.DB.QueryRow(context.Background(), query, name).Scan(&tag.Id, &tag.Name)
	if err != nil {
		return tag, err
	}
	return tag, nil
}
