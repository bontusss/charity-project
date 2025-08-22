package auth

import (
	db "charity/db/sqlc"
	"charity/internal/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db     *db.Queries
	config *config.Config
}

func NewAuthService(db *db.Queries, c *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: c,
	}
}

func (a *AuthService) CreateAdminUser(ctx context.Context, email, password string) (*db.CreateAdminRow, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := a.db.CreateAdmin(ctx, db.CreateAdminParams{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type UserObjectWithToken struct {
	User  *db.Admin
	Token string
}

func (a *AuthService) LoginAdmin(ctx context.Context, email, password string) (string, *db.Admin, error) {
	// Get admin by email
	admin, err := a.db.GetAdminByEmail(ctx, email)
	if err != nil {
		log.Printf("Failed to get admin by email %s: %v", email, err)
		return "", nil, errors.New("invalid credentials")
	}

	log.Printf("Found admin with email: %s", admin.Email)

	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		log.Printf("Password comparison failed for admin %s: %v", admin.Email, err)
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token for admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		log.Printf("Failed to sign JWT token: %v", err)
		return "", nil, err
	}

	log.Printf("Successfully generated token for admin %s", admin.Email)
	return signedToken, &admin, nil
}

func (a *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.config.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
