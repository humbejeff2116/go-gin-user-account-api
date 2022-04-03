
package lib

import (
    
    "log"
	"go-gin-user-account-api/configs"
    "golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	
)

var serverConfigs = configs.SetServerConfigurations();

func GeneratePasswordHash(password string) string {

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

    if err != nil {

        log.Println(err)

    }

    return string(passwordHash)

}

func CheckPassword(userPasswordGuess, hashedPassword string) error {

    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPasswordGuess))
}

func GenerateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err :=  token.SignedString(serverConfigs.JwtSecret)

	if err != nil {

		log.Println("An error occured while generating JWT token")

		return "", err

	}

	return tokenString, nil

}