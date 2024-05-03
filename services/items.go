package itemsService

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/boy52hz/go-demo-shop-api/db"
)

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Create(newItem *Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := db.Instance.ExecContext(ctx, "INSERT INTO items (name, description) VALUES (?, ?)", newItem.Name, newItem.Description)

	if err != nil {
		log.Println(err)
		return err
	}

	insertId, err := result.LastInsertId()

	if err != nil {
		log.Println(err)
		return err
	}

	newItem.Id = int(insertId)

	return nil
}

func FindOne(id int) (*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.Instance.QueryRowContext(ctx, "SELECT id, name, description FROM items WHERE id = ?", id)

	item := &Item{}
	err := row.Scan(
		&item.Id,
		&item.Name,
		&item.Description,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return item, nil
}

func FindAll() ([]Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := db.Instance.QueryContext(ctx, `SELECT id, name, description FROM items`)

	if err != nil {
		return nil, err
	}

	defer results.Close()

	items := make([]Item, 0)

	for results.Next() {
		var item Item
		results.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
		)
		items = append(items, item)
	}

	return items, nil
}

func Update(id int, updateItem *Item) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row, err := db.Instance.ExecContext(ctx, "UPDATE items SET name = ?, description = ? WHERE id = ?", updateItem.Name, updateItem.Description, id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	affectedRows, err := row.RowsAffected()

	if affectedRows == 0 {
		return 0, nil
	} else if err != nil {
		log.Println(err)
		return 0, err
	}

	updateItem.Id = id

	return int(affectedRows), nil
}

func Delete(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := db.Instance.ExecContext(ctx, "DELETE FROM items WHERE id = ?", id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(affectedRows), nil
}
