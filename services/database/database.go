package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	db *pgx.Conn
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func InitDatabase(connString string) (Database, error) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return Database{}, fmt.Errorf("unable to connect to the database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(context.Background()); err != nil {
		return Database{}, fmt.Errorf("failed to ping database: %w", err)
	}

	// create a user table
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	)`)

	if err != nil {
		return Database{}, fmt.Errorf("failed to create table: %w", err)
	}

	// Assign the connection to the global variable
	db := Database{db: conn}
	log.Println("Database connection initialized successfully.")
	return db, nil
}

func (d *Database) CreateUser(user User) (int64, error) {
	var id int64
	err := d.db.QueryRow(context.Background(), `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`, user.Email, user.Name, user.Password).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert user: %w", err)
	}
	return id, nil
}

func (d *Database) GetUserById(id string) (User, error) {
	var user User
	row := d.db.QueryRow(context.Background(), `SELECT * FROM users WHERE id = $1`, id)
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (d *Database) GetUserByEmail(email string) (User, error) {
	var user User
	row := d.db.QueryRow(context.Background(), `SELECT * FROM users WHERE email = $1`, email)
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (d *Database) Close() {
	d.db.Close(context.Background())
	log.Println("Database connection closed.")
}
