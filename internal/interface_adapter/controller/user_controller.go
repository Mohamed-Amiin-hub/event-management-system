package controller

import (
	"log"
	"net/http"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// User Controller handles user requests
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser creates a new user
func (uc *UserController) RegisterUser(c *gin.Context) {
	var user entity.User

	// Bind incoming JSON to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Bound User Struct: %+v", user)

	// Call the service layer to handle user registration
	createdUser, err := uc.userService.RegisterUser(user.Username, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, createdUser)
}

// Authenticate User
func (uc *UserController) AuthenticateUser(c *gin.Context) {
	var user entity.User

	// Bind incoming JSON to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Bound User Struct: %+v", user)

	// Call the service layer to handle user authentication
	authenticateUser, err := uc.userService.AuthenticateUser(user.Email, user.Password)
	if err != nil {
		log.Printf("error to authenticate user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, authenticateUser)

}

// update user
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user entity.User

	// Get the user ID from the URL parameter
	userIDParam := c.Param("id")
	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Bind incoming JSON to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	log.Printf("Bound User Struct: %+v", user)

	// Set the user ID from the params into the user struct
	user.ID = userID

	// Call the service layer to handle user update
	if err := uc.userService.UpdateUser(&user); err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// delete user
func (uc *UserController) DeleteUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	userIDParam := c.Param("id")
	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	}

	// Call the service layer to handle user deletion
	if err := uc.userService.DeleteUser(userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	// Get the user ID from the URL parameter
	userIDParam := c.Param("id")
	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the service layer to handle user deletion
	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) ListUsers(c *gin.Context) {

	// Call the service layer to handle user deletion
	user, err := uc.userService.ListUsers()
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"user": user})
}
