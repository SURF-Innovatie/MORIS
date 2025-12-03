package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	userent "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type jwtService struct {
	client  *ent.Client
	userSvc user.Service
}

func NewJWTService(client *ent.Client, userSvc user.Service) auth.Service {
	return &jwtService{client: client, userSvc: userSvc}
}

// Register creates a new user with hashed password
func (s *jwtService) Register(ctx context.Context, req auth.RegisterRequest) (*entities.UserAccount, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with default "user" role
	usr, err := s.client.User.
		Create().
		SetPersonID(req.PersonID).
		SetPassword(string(hashedPassword)).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.userSvc.GetAccount(ctx, usr.ID)
}

// Login authenticates a user and returns a JWT token
func (s *jwtService) Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error) {
	// Find user by email
	usr, err := s.userSvc.GetAccountByEmail(ctx, email)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	usrPwd, err := s.client.User.
		Query().
		Where(userent.IDEQ(usr.User.ID)).
		Select(userent.FieldPassword).
		String(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to query user password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usrPwd), []byte(password))
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

// generateJWT creates a JWT token for the user
func (s *jwtService) generateJWT(usr *entities.UserAccount) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production" // Fallback for development
	}

	claims := jwt.MapClaims{
		"user_id":  usr.User.ID,
		"email":    usr.Person.Email,
		"orcid_id": usr.Person.ORCiD,
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
func (s *jwtService) ValidateToken(tokenString string) (*entities.UserAccount, error) {
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

	usr, err := s.userSvc.GetAccount(context.Background(), uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: question why do we take fields from token instead of DB
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
