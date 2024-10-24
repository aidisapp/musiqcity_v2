package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/aidisapp/musiqcity_v2/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var app *config.AppConfig

// Setup app config for new helpers
func NewHelpers(helper *config.AppConfig) {
	app = helper
}

func ClientError(responseWriter http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(responseWriter, http.StatusText(status), status)
}

func ServerError(responseWriter http.ResponseWriter, err error) {
	// err.Error() prints the nature of the error. debug.Stack() prints the detailed info about the nature of the error
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	// Send feedback to the user
	http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Check if a user is exists
func IsAuthenticated(request *http.Request) bool {
	exists := app.Session.Exists(request.Context(), "user_id")
	return exists
}

// Check if a user is admin
func IsAdmin(request *http.Request) bool {
	accessLevel := app.Session.GetInt(request.Context(), "access_level")
	return accessLevel == 3
}

func GenerateJWTToken(userID int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	jwtToken := os.Getenv("JWTSECRET")

	// Define the claims for the JWT token
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"iat": time.Now().Unix(),
	}

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(jwtToken))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
