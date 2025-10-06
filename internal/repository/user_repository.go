package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// User represents a minimal user record for verification state.
// Stored in Redis as JSON at key user:{userID}
//
type User struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	EmailVerified  bool      `json:"email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
}

type UserRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) key(userID string) string { return fmt.Sprintf("user:%s", userID) }

// Upsert sets/updates user fields
func (r *UserRepository) Upsert(ctx context.Context, user User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(user.UserID), b, 0).Err()
}

// Get retrieves a user; returns (nil, nil) if not found
func (r *UserRepository) Get(ctx context.Context, userID string) (*User, error) {
	val, err := r.client.Get(ctx, r.key(userID)).Bytes()
	if err != nil {
		if err == redis.Nil { return nil, nil }
		return nil, err
	}
	var u User
	if err := json.Unmarshal(val, &u); err != nil { return nil, err }
	return &u, nil
}

// SetEmailVerified marks the user as verified and persists
func (r *UserRepository) SetEmailVerified(ctx context.Context, userID, email string, verifiedAt time.Time) error {
	u, err := r.Get(ctx, userID)
	if err != nil { return err }
	if u == nil {
		u = &User{UserID: userID, Email: email}
	}
	u.Email = email
	u.EmailVerified = true
	u.EmailVerifiedAt = verifiedAt
	return r.Upsert(ctx, *u)
	"database/sql"
	"fmt"

	"github.com/codeZe-us/vestroll-backend/internal/models"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.FullName,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, full_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, password, full_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	
	return exists, nil
}
