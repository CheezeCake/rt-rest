package web

import (
	"crypto/sha512"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CheezeCake/rt-rest/config"
	"github.com/CheezeCake/rt-rest/data"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	authSecret string
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) error {
	var credentials credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return err
	}

	shaPass := sha512.Sum512([]byte(credentials.Password))
	passHash, err := bcrypt.GenerateFromPassword(shaPass[:], bcrypt.DefaultCost)
	err = data.AddUser(credentials.Username, string(passHash[:]))
	return err
}

func Login(w http.ResponseWriter, r *http.Request) error {
	var credentials credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return err
	}

	dbPassword, err := data.GetPassword(credentials.Username)
	if err != nil {
		if err == data.ErrUserNotFound {
			return ErrLoginFailed
		}
		return err
	}
	shaPass := sha512.Sum512([]byte(credentials.Password))
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), shaPass[:])
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrLoginFailed
		}
		return err
	}

	token, err := createToken(credentials)
	if err == nil {
		json.NewEncoder(w).Encode(&struct{ Token string }{token})
	}
	return err
}

func createToken(credentials credentials) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": credentials.Username,
		"password": credentials.Password,
	})
	return token.SignedString([]byte(authSecret))
}

func Auth() Middleware {
	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			authorizationHeader := r.Header.Get("authorization")
			if authorizationHeader == "" {
				return ErrNeedLogin
			}
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) != 2 {
				return ErrMalformedAuthorization
			}
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrSigningMethod
				}
				return []byte(authSecret), nil
			})
			if err != nil {
				return err
			}
			if token.Valid {
				return h(w, r)
			} else {
				return ErrNeedLogin
			}
		}
	}
}

func Init(cfg config.Cfg) {
	authSecret = cfg.Secret
}
