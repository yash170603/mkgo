package services

import (
	"errors"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	Logout(token string) error
}

type authService struct {
	db  *database.DB
	cfg config.Config
}

func NewAuthService(db *database.DB, cfg config.Config) AuthService {
	return &authService{
		db:  db,
		cfg: cfg,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	var user models.User
	if err := s.db.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	// Generate JWT token
	claims := Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JwtConfig.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &user, nil
}

func (s *authService) Logout(token string) error {
	// handled at client side later, no need to implement
	return nil
}
