package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity/readmodels"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, *readmodels.UserAccount, error)
	LoginByEmail(ctx context.Context, email string) (string, *readmodels.UserAccount, error)
	ValidateToken(tokenString string) (*readmodels.UserAccount, error)
}

type service struct {
	repo      Repository
	jwtSecret string
	ttl       time.Duration
	now       func() time.Time
}

type Options struct {
	JWTSecret string
	TTL       time.Duration
	Now       func() time.Time
}

func NewService(repo Repository, opts Options) Service {
	ttl := opts.TTL
	if ttl == 0 {
		ttl = 7 * 24 * time.Hour
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}

	return &service{
		repo:      repo,
		jwtSecret: opts.JWTSecret,
		ttl:       ttl,
		now:       now,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (string, *readmodels.UserAccount, error) {
	usr, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		// Keep the exact error mapping consistent across repos.
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if !usr.User.IsActive {
		return "", nil, fmt.Errorf("user account is inactive")
	}

	hash, err := s.repo.GetPasswordHash(ctx, usr.User.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}
	if hash == "" {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.generateJWT(usr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, usr, nil
}

func (s *service) LoginByEmail(ctx context.Context, email string) (string, *readmodels.UserAccount, error) {
	usr, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
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

func (s *service) ValidateToken(tokenString string) (*readmodels.UserAccount, error) {
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

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	usr, err := s.repo.GetAccountByID(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return usr, nil
}

func (s *service) generateJWT(usr *readmodels.UserAccount) (string, error) {
	now := s.now()

	claims := jwt.MapClaims{
		"user_id":   usr.User.ID.String(),
		"email":     usr.Person.Email,
		"orcid_id":  usr.Person.ORCiD,
		"is_active": usr.User.IsActive,
		"exp":       now.Add(s.ttl).Unix(),
		"iat":       now.Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.jwtSecret))
}
