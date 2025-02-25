package user

import (
	"database/sql"
	"errors"
	"log"
	"web-app/app/services/core"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CreatedAt string `json:"created_at"`
}

func NewUserModel() *User {
	return &User{}
}

// Create implements the Model interface Create method
func (u *User) Create() error {
	db, _ := core.NewPostgresService()

	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id`

	result, err := db.Create(query, u.Username, u.Password)
	if err != nil {
		log.Printf("error creating user: %v", err)
		return err
	}

	defer db.Close()

	err = result.Scan(&u.ID)
	if err != nil {
		log.Printf("error scanning user: %v", err)
	}

	return nil
}

// Find implements the Model interface Find method
func (u *User) Find() error {
	db, _ := core.NewPostgresService()

	if u.Username == "" {
		return errors.New("username is required")
	}

	query := `
        SELECT id, username, password, created_at 
        FROM users 
        WHERE username = $1`

	rows, err := db.Read(query, u.Username)
	if err != nil {
		log.Printf("error finding user: %v", err)
		return err
	}
	defer db.Close()

	if rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			log.Printf("error scanning user: %v", err)
			return err
		}
		return nil
	}

	return sql.ErrNoRows
}

// FindByUsername implements the Model interface FindByUsername method
func (u *User) FindByUsername() error {
	db, _ := core.NewPostgresService()

	if u.Username == "" {
		return errors.New("username is required")
	}

	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = $1`

	rows, err := db.Read(query, u.Username)
	if err != nil {
		log.Printf("error finding user: %v", err)
		return err
	}
	defer db.Close()

	if rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			log.Printf("error scanning user: %v", err)
			return err
		}
		return nil
	}

	return sql.ErrNoRows
}

// Update implements the Model interface Update method
func (u *User) Update() error {
	db, _ := core.NewPostgresService()

	if u.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        UPDATE users 
        SET username = $1, password = $2
        WHERE id = $3`

	defer db.Close()

	_, err := db.Update(query, u.Username, u.Password, u.ID)
	if err != nil {
		log.Printf("error updating user: %v", err)
		return err
	}

	return nil
}

// Delete implements the Model interface Delete method
func (u *User) Delete() error {
	db, _ := core.NewPostgresService()

	if u.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        DELETE FROM users 
        WHERE id = $1`

	defer db.Close()

	_, err := db.Delete(query, u.ID)
	if err != nil {
		log.Printf("error deleting user: %v", err)
		return err
	}

	return nil
}

func (u *User) Paginate(limit, page int) ([]User, error) {
	db, _ := core.NewPostgresService()

	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	query := `
        SELECT id, username, password, created_at
        FROM users
        ORDER BY id DESC
        LIMIT $1 OFFSET $2`

	defer db.Close()

	rows, err := db.Read(query, limit, offset)
	if err != nil {
		log.Printf("error paginating users: %v", err)
		return nil, err
	}

	users := make([]User, 0, limit) // Pre-allocate slice with capacity
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			log.Printf("error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("error iterating over rows: %v", err)
		return nil, err
	}

	return users, nil
}
