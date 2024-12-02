package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type UserModel struct {
	DB   *pgx.Conn
	User User
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func (u *UserModel) Create(user User) (string, error) {
	var id string
	var err error
	var query string
	// Create a query string. If the user.Role is not an empty string, add the role to the query; else, do not add it.
	if user.Role != "" {
		query = `INSERT INTO users (email, name, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
		err = u.DB.QueryRow(context.Background(), query, user.Email, user.Name, user.Password, user.Role).Scan(&id)
	} else {
		query = `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`
		err = u.DB.QueryRow(context.Background(), query, user.Email, user.Name, user.Password).Scan(&id)
	}


	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
			if pgErr.Code == "23505" { // unique_violation
				return "", fmt.Errorf("email already in use")
			}
		} else if strings.Contains(err.Error(), "23505") {
			return "", fmt.Errorf("email already in use")
		}
		fmt.Printf("error: %v\n", err)
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	return id, nil
}

func (u *UserModel) ById(id string) (User, error) {
	var user User
	row := u.DB.QueryRow(context.Background(), `SELECT id, email, name, password, role FROM users WHERE id = $1`, id)
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Role)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *UserModel) ByEmail(email string) (User, error) {
	var user User
	row := u.DB.QueryRow(context.Background(), `SELECT id, email, name, password, role FROM users WHERE email = $1`, email)
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Role)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
