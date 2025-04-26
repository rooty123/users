package actions

import (
	"net/http"
	"strconv"
	"users/db"

	"github.com/labstack/echo/v4"
)

// Placeholder for User struct
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock in-memory "database"
var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
}

// GetUsers returns a list of users
func GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	// Extract the chat ID from the URL path
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	// Create a new database handler
	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	// Get the user from the database
	user, err := dbh.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Return the user as JSON
	return c.JSON(http.StatusOK, user)
}

// CreateUser adds a new user
func CreateUser(c echo.Context) error {
	var newUser struct {
		ChatID       int64  `json:"chatID"`
		TelegramID   int64  `json:"telegramID"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		LanguageCode string `json:"languageCode"`
		Username     string `json:"username"`
	}

	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	if err := dbh.CreateUser(newUser.ChatID, newUser.TelegramID, newUser.FirstName, newUser.LastName, newUser.LanguageCode, newUser.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, newUser)
}

// UpdateUser updates an existing user
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var updatedUser User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}
	for i, user := range users {
		if user.ID == id {
			users[i] = updatedUser
			return c.JSON(http.StatusOK, updatedUser)
		}
	}
	return c.JSON(http.StatusNotFound, "User not found")
}

// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusNotFound, "User not found")
}
