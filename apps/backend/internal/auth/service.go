package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (*ent.User, error)
	Login(ctx context.Context, email, password string) (string, *AuthenticatedUser, error)
	ValidateToken(tokenString string) (*AuthenticatedUser, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GenerateORCIDAuthURL(ctx context.Context) (string, error)
	LinkORCID(ctx context.Context, userID uuid.UUID, code string) error
	UnlinkORCID(ctx context.Context, userID uuid.UUID) error
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
		ID:      usr.ID,
		Email:   usr.Email,
		OrcidID: usr.OrcidID,
		//Roles:   usr.Roles,
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
		"user_id":  usr.ID,
		"email":    usr.Email,
		"orcid_id": usr.OrcidID,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiry
		"iat":      time.Now().Unix(),
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
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email in token")
	}

	orcidID, _ := claims["orcid_id"].(string) // Optional field

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

	return &AuthenticatedUser{
		ID:      uid,
		Email:   email,
		OrcidID: orcidID,
		//Roles:   roles,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	usr, err := s.client.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return usr, nil
}

// GenerateORCIDAuthURL generates the ORCID authorization URL
func (s *service) GenerateORCIDAuthURL(ctx context.Context) (string, error) {
	config, err := GetORCIDConfig()
	if err != nil {
		return "", err
	}
	return config.GenerateAuthURL(), nil
}

// LinkORCID links an ORCID ID to a user account
func (s *service) LinkORCID(ctx context.Context, userID uuid.UUID, code string) error {
	config, err := GetORCIDConfig()
	if err != nil {
		return err
	}

	orcidID, err := config.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	// Check if ORCID is already linked to another user
	exists, err := s.client.User.Query().Where(user.OrcidIDEQ(orcidID)).Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if ORCID is already linked: %w", err)
	}
	if exists {
		return fmt.Errorf("ORCID ID is already linked to another account")
	}

	// Update user with ORCID ID
	_, err = s.client.User.UpdateOneID(userID).SetOrcidID(orcidID).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to link ORCID ID: %w", err)
	}

	return nil
}

// UnlinkORCID removes the ORCID ID from a user account
func (s *service) UnlinkORCID(ctx context.Context, userID uuid.UUID) error {
	// Update user to remove ORCID ID
	_, err := s.client.User.UpdateOneID(userID).ClearOrcidID().Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unlink ORCID ID: %w", err)
	}

	return nil
}
