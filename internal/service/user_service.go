package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"example.com/EVENT-MANAGEMENT-SYSTEM/pkg/utils"
	"github.com/gofrs/uuid"
)

type UserService interface {
	RegisterUser(username, email, password, first_name, last_name string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	// DeactivateUser(userID uint) error
	// ActivateUser(userID uint) error
	DeleteUser(userID uuid.UUID) error
	GetUserByID(userID uuid.UUID) (*entity.User, error)
	// GetUserByEmail(email string) (*entity.User, error)
	ListUsers() ([]*entity.User, error)
	AuthenticateUser(email, password string) (*entity.User, error)
}

// userServiceImpl is the implementation of UserService.
type UserServiceImpl struct {
	repo      repository.UserRepository
	tokenRepo repository.TokenRepository
}

// ListUsers implements UserService.
func (s *UserServiceImpl) ListUsers() ([]*entity.User, error) {
	users, err := s.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}

	return users, nil

}

// ListUsers implements UserService.

// GetUserByID implements UserService.
func (s *UserServiceImpl) GetUserByID(userID uuid.UUID) (*entity.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with ID %s: %v", userID, err)
	}

	return user, nil
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) UserService {
	return &UserServiceImpl{
		repo:      userRepo,
		tokenRepo: tokenRepo,
	}

}

func (s *UserServiceImpl) RegisterUser(username, email, password, first_name, last_name string) (*entity.User, error) {

	// Check if the user already exists
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash the password using bcrypt
	HashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Username:  username,
		Email:     email,
		Password:  HashPassword,
		FirstName: first_name,
		LastName:  last_name,
		IsActive:  true,
	}

	// Save the user to the repository
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user by email and password.
func (s *UserServiceImpl) AuthenticateUser(email, password string) (*entity.User, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token (UUID in this case, but you can replace it with JWT or any other mechanism)
	newToken, err := uuid.NewV4()
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	token := &entity.Token{
		ID:        newToken,
		UserID:    user.ID,
		Token:     newToken.String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Store the token in the database
	if err := s.tokenRepo.Create(token); err != nil {
		return nil, errors.New("failed to save token")
	}

	// Attach the token to the user entity

	// Return the user and the token
	return user, nil
}

// update user
func (s *UserServiceImpl) UpdateUser(user *entity.User) error {
	// Check if the user exists by their ID
	_, err := s.repo.FindByID(user.ID)
	if err != nil {
		// If the user does not exist, return the error
		return fmt.Errorf("could not find user with ID %s", user.ID)
	}

	// Call the repository to update the user
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update user with ID %s: %v", user.ID, err)
	}

	return nil
}

// delete user
func (s *UserServiceImpl) DeleteUser(userID uuid.UUID) error {
	// Check if the user exists by their ID
	_, err := s.repo.FindByID(userID)
	if err != nil {
		// If the user does not exist, return the error
		log.Printf("Could not find user with ID %s: %v", userID, err)
		return fmt.Errorf("could not find user with ID %s", userID)
	}

	// Call the repository to delete the user
	if err := s.repo.Delete(userID); err != nil {
		log.Printf("Failed to delete user with ID %s: %v", userID, err)
		return fmt.Errorf("failed to delete user with ID %s: %v", userID, err)
	}

	log.Printf("Successfully deleted user with ID %s", userID)
	return nil
}
