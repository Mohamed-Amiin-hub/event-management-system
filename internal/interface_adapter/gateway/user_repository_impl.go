package gateway

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"github.com/gofrs/uuid"
	//"github.com/google/uuid"
)

// userRepositoryImpl is the implementation of UserRepository.
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl.
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create implements repository.UserRepository.
func (u *userRepositoryImpl) Create(user *entity.User) error {
	// Generate a new UUID for the user
	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return err
	}
	user.ID = newUUID

	// Log user creation
	log.Printf("Inserting User: ID=%s, Username=%s, Email=%s, Password=%s FirstName=%s, LastName=%s, IsActive=%v\n",
		user.ID, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.IsActive)

	// Corrected SQL query
	query := `INSERT INTO users (id, username, password, email, first_name, last_name, is_active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Execute the query
	result, err := u.db.Exec(query, user.ID, user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.IsActive, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
	} else {
		log.Printf("Rows affected: %d", rowsAffected)
	}

	// Fetch the inserted user (optional)
	var insertedUser entity.User
	err = u.db.QueryRow(`SELECT id, username, email, password, first_name, last_name, is_active, created_at, updated_at 
	                    FROM users WHERE email=$1`, user.Email).Scan(
		&insertedUser.ID, &insertedUser.Username, &insertedUser.Email, &insertedUser.Password,
		&insertedUser.FirstName, &insertedUser.LastName, &insertedUser.IsActive, &insertedUser.CreatedAt, &insertedUser.UpdatedAt)

	if err != nil {
		log.Printf("Error retrieving inserted user: %v", err)
		return err
	}

	log.Printf("Inserted User: %+v", insertedUser)

	return nil
}

// Update implements repository.UserRepository.
func (u *userRepositoryImpl) Update(user *entity.User) error {
	query := `UPDATE users
	          SET username = $1, email = $2, password = $3, first_name = $4, last_name =$5, is_active =$6, updated_at = $7
			  WHERE id $8`

	// Execute the update query with the user data
	result, err := u.db.Exec(query, user.Username, user.Email, user.FirstName, user.LastName, user.IsActive, time.Now(), user.ID)
	if err != nil {
		log.Printf("Error updating user with ID: %v, error: %v", user.ID, err)
		return err
	}

	// Check how many rows were affected by the update
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	// If no rows were affected, it means the user was not found
	if rowsaffected == 0 {
		log.Printf("No user found with ID: %v", user.ID)
		return fmt.Errorf("user not found")
	}

	log.Printf("Updated user with ID: %v", user.ID)
	return nil
}

// Delete implements repository.UserRepository.
func (u *userRepositoryImpl) Delete(userID uuid.UUID) error {
	// Define the SQL delete query
	query := `DELETE FROM users WHERE id = $1`

	// Execute the delete query with the user ID
	result, err := u.db.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user with ID: %v, error: %v", userID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	// If no rows were affected, it means the user was not found
	if rowsAffected == 0 {
		log.Printf("No user found with ID: %v", err)
		return fmt.Errorf("user not found")
	}

	log.Printf("Deleted user with ID: %v", userID)
	return nil
}

// FindByID implements repository.UserRepository.
func (u *userRepositoryImpl) FindByID(userID uuid.UUID) (*entity.User, error) {
	// Define the user entity to store the result
	var user entity.User

	// Fetch the user from the database using the provided ID
	err := u.db.QueryRow(`SELECT id, username, email, password, first_name, last_name, is_active, created_at, updated_at 
	              FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	// Check for errors in retrieving the user
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with ID: %v", userID)
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error retrieving user by ID: %v", err)
		return nil, err
	}

	// Return the user if found
	return &user, nil
}

// FindByEmail implements repository.UserRepository.
func (u *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	query := "SELECT id, username, email, password, first_name, last_name, is_active, created_at, updated_at FROM users WHERE email = $1"
	row := u.db.QueryRow(query, email)

	err := row.Scan(user.ID, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return user, nil
}

// ListAll implements repository.UserRepository.
func (u *userRepositoryImpl) ListAll() ([]*entity.User, error) {
	// Fetch all users from the database
	rows, err := u.db.Query("SELECT id, username, email, password, first_name, last_name, is_active, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Scan the rows into a slice of User structs
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	// Check for any errors during the scan
	return users, nil
}
