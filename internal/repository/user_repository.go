package repository

import (
	"database/sql"
	"fmt"

	"github.com/codeZe-us/vestroll-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`
	
	row := r.db.QueryRow(query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	
	return &user, nil
}

// CreateUser creates a new user (for future registration functionality)
func (r *UserRepository) CreateUser(user models.User) (*models.User, error) {
	query := `
		INSERT INTO users (email, password) 
		VALUES ($1, $2) 
		RETURNING id, email, created_at, updated_at
	`
	
	var newUser models.User
	err := r.db.QueryRow(query, user.Email, user.Password).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return &newUser, nil
}