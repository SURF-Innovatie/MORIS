package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (*ent.User, error)
	Login(ctx context.Context, email, password string) (string, *AuthenticatedUser, error)
	ValidateToken(tokenString string) (*AuthenticatedUser, error)
	GetUserByID(ctx context.Context, id int) (*ent.User, error)
}

type service struct {
	client *ent.Client
}

func NewService(client *ent.Client) Service {
	return &service{client: client}
}

// Register creates a new user with hashed password
func (s *service) Register(ctx context.Context, name, email, password string) (*ent.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with default "user" role
	usr, err := s.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(string(hashedPassword)).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return usr, nil
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(ctx context.Context, email, password string) (string, *AuthenticatedUser, error) {
	// Find user by email
	usr, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(usr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	authUser := &AuthenticatedUser{
		ID:    usr.ID,
		Email: usr.Email,
	}

	return token, authUser, nil
}

// generateJWT creates a JWT token for the user
func (s *service) generateJWT(usr *ent.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production" // Fallback for development
	}

	claims := jwt.MapClaims{
		"user_id": usr.ID,
		"email":   usr.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiry
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user info
func (s *service) ValidateToken(tokenString string) (*AuthenticatedUser, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract user info from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email in token")
	}

	rolesInterface, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid roles in token")
	}

	roles := make([]string, len(rolesInterface))
	for i, r := range rolesInterface {
		roles[i], ok = r.(string)
		if !ok {
			return nil, fmt.Errorf("invalid role format in token")
		}
	}

	return &AuthenticatedUser{
		ID:    int(userID),
		Email: email,
		Roles: roles,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (s *service) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	usr, err := s.client.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return usr, nil
}
