package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/src-d/code-annotation/server/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// Strips 'Bearer ' prefix from bearer token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

var extractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromTokenString,
}

// JWTConfig defines enviroment variables for JWT
type JWTConfig struct {
	SigningKey string `envconfig:"SIGNING_KEY" required:"true"`
}

// JWT service abstracts JWT implementation
type JWT struct {
	signingKey []byte
}

// NewJWT return new JWT service
func NewJWT(signingKey string) *JWT {
	return &JWT{signingKey: []byte(signingKey)}
}

type userIDContext int

const userIDKey userIDContext = 1

type jwtClaim struct {
	ID int
	jwt.StandardClaims
}

// MakeToken generates token string for a user
func (j *JWT) MakeToken(user *model.User) (string, error) {
	claims := &jwtClaim{ID: user.ID}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := t.SignedString(j.signingKey)
	if err != nil {
		return "", fmt.Errorf("can't sign jwt token: %s", err)
	}
	return ss, nil
}

// Middleware return http.Handler which validates token and set user id in context
func (j *JWT) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims jwtClaim
		_, err := request.ParseFromRequestWithClaims(r, extractor, &claims, func(token *jwt.Token) (interface{}, error) {
			return j.signingKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIDKey, claims.ID))
		next.ServeHTTP(w, r)
	})
}

// getUserInt gets the value stored in the Context for the key userIDKey, bool
// is true on success
func getUserInt(ctx context.Context) (int, bool) {
	i, ok := ctx.Value(userIDKey).(int)
	return i, ok
}

// GetUserID gets the user ID set by the JWT middleware in the Context
func GetUserID(ctx context.Context) (int, error) {
	id, ok := getUserInt(ctx)
	if !ok {
		return 0, fmt.Errorf("User ID is not set in the context")
	}

	return id, nil
}
