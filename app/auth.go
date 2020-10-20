package app

import (
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var TempSecret = "1234567890"

// https://levelup.gitconnected.com/crud-restful-api-with-go-gorm-jwt-postgres-mysql-and-testing-460a85ab7121

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Error compare hash and pssword. %v", err)
		err = errors.New("Wrong username or password")
	}
	return err == nil, err
}

/*
	Создаем и подписываем пару токенов
*/
func CreateTokenPair(user_id int64) (*TokenPair, error) {
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["sub"] = user_id
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() //Token expires after 15 minutes

	refreshClaims := jwt.MapClaims{}
	refreshClaims["sub"] = user_id
	refreshClaims["exp"] = time.Now().Add(time.Hour * 12).Unix() //Token expires after 12 hour

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		log.Printf("Create signed access tocken string %v", err)
	}
	refreshTokenSting, err := refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		log.Printf("Create signed refresh tocken string %v", err)
	}
	return &TokenPair{accessTokenString, refreshTokenSting}, err
}

/*
	Проверяем токен
*/
func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(TempSecret), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(r *http.Request) (int64, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(TempSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["sub"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return int64(uid), nil
	}
	return 0, nil
}

//Pretty display the claims licely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}