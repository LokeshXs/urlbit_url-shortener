package webhook

import (
	"net/http"

	"github.com/Lokeshxs/url-shortener/db"
	"github.com/gin-gonic/gin"
)

type ClerEvent struct {
	Data struct {
		ID             string `json:"id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		EmailAddresses []struct {
			EmailAddress string `json:"email_address"`
		} `json:"email_addresses"`
	} `json:"data"`
}

func ClerkSignUp(c *gin.Context) {

	var payload ClerEvent

	err := c.ShouldBindJSON(&payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Sign up user",
		})

		return
	}

	firstName := payload.Data.FirstName
	lastName := payload.Data.LastName
	fullName := firstName + " " + lastName
	userID := payload.Data.ID
	email := payload.Data.EmailAddresses[0].EmailAddress

	query := `
	INSERT INTO users (id, name, email)
	VALUES ($1, $2, $3);
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})

		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(userID, fullName, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sign up success",
	})

}
