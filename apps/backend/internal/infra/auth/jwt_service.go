package auth

import (
	"context"
	"fmt"
	"time"

	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/apikey"
	userent "github.com/SURF-Innovatie/MORIS/ent/user"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	client    *ent.Client
	userSvc   user.Service
	personSvc person.Service
	jwtSecret string
}

func NewJWTService(client *ent.Client, userSvc user.Service, personSvc person.Service, jwtSecret string) coreauth.Service {
	return &service{
		client:    client,
		userSvc:   userSvc,
		personSvc: personSvc,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error) {
	// Find user by email
	usr, err := s.userSvc.GetAccountByEmail(ctx, email)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	u, err := s.client.User.
		Query().
		Where(userent.IDEQ(usr.User.ID)).
		Only(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Check if user has a password set (OAuth-only users don't)
	if u.Password == "" {
		return "", nil, fmt.Errorf("invalid credentials") // OAuth-only user, can't login with password
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(usr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, usr, nil
}

// LoginByEmail issues a JWT for an existing user identified by email.
// This is intended for external IdP logins (e.g. SURFconext) where MORIS is not
// handling a local password challenge.
func (s *service) LoginByEmail(ctx context.Context, email string) (string, *entities.UserAccount, error) {
	usr, err := s.userSvc.GetAccountByEmail(ctx, email)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !usr.User.IsActive {
		return "", nil, fmt.Errorf("user account is inactive")
	}

	token, err := s.generateJWT(usr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, usr, nil
}

// generateJWT creates a JWT token for the user
func (s *service) generateJWT(usr *entities.UserAccount) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   usr.User.ID,
		"email":     usr.Person.Email,
		"orcid_id":  usr.Person.ORCiD,
		"is_active": usr.User.IsActive,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiry
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user info
func (s *service) ValidateToken(tokenString string) (*entities.UserAccount, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
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
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	usr, err := s.userSvc.GetAccount(context.Background(), uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: question why do we take fields from token instead of DB; answer = fast lookup without DB call
	//email, ok := claims["email"].(string)
	//if !ok {
	//	return nil, fmt.Errorf("invalid email in token")
	//}
	//
	//orcidID, _ := claims["orcid_id"].(string) // Optional field

	//rolesInterface, ok := claims["roles"].([]interface{})
	//if !ok {
	//	return nil, fmt.Errorf("invalid roles in token")
	//}

	//roles := make([]string, len(rolesInterface))
	//for i, r := range rolesInterface {
	//	roles[i], ok = r.(string)
	//	if !ok {
	//		return nil, fmt.Errorf("invalid role format in token")
	//	}
	//}

	return usr, nil
}

// ValidateAPIKey validates an API key and returns the user info
func (s *service) ValidateAPIKey(ctx context.Context, plainKey string) (*entities.UserAccount, error) {
	const APIKeyPrefix = "moris_"

	// Remove prefix if present
	plainKey = strings.TrimPrefix(plainKey, APIKeyPrefix)

	keyHash := s.hashAPIKey(plainKey)

	key, err := s.client.APIKey.
		Query().
		Where(
			apikey.KeyHash(keyHash),
			apikey.IsActive(true),
		).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("invalid or inactive api key")
	}

	// Check expiration
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return nil, fmt.Errorf("api key has expired")
	}

	// Update last used timestamp
	_, _ = s.client.APIKey.
		UpdateOneID(key.ID).
		SetLastUsedAt(time.Now()).
		Save(ctx)

	return s.userSvc.GetAccount(ctx, key.UserID)
}

func (s *service) hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
