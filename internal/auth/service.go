package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type Service struct {
	repo      Repository
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewService(repo Repository, jwtSecret string, tokenTTL time.Duration) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  tokenTTL,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (TokenPair, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}
	if !user.Active() || !VerifyPassword(user.PasswordHash, password) {
		return TokenPair{}, ErrInvalidCredentials
	}

	ttl := s.tokenTTL
	if ttl <= 0 {
		ttl = time.Hour
	}
	expiresAt := time.Now().Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":          user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"permissions":  collectPermissionCodes(user.Roles),
		"exp":          expiresAt.Unix(),
	})
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken: signed,
		TokenType:   "Bearer",
	}, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) CreateUser(ctx context.Context, user User, password string, roleNames []string) (User, error) {
	return s.repo.CreateUser(ctx, user, password, roleNames)
}

func (s *Service) UpdateUser(ctx context.Context, user User, password string) (User, error) {
	return s.repo.UpdateUser(ctx, user, password)
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *Service) CreateRole(ctx context.Context, role Role) (Role, error) {
	return s.repo.CreateRole(ctx, role)
}

func (s *Service) UpdateRole(ctx context.Context, role Role) (Role, error) {
	return s.repo.UpdateRole(ctx, role)
}

func (s *Service) DeleteRole(ctx context.Context, name string) error {
	return s.repo.DeleteRole(ctx, name)
}

func (s *Service) ListPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func collectPermissionCodes(roles []Role) []string {
	seen := make(map[string]struct{})
	codes := make([]string, 0)
	for _, role := range roles {
		for _, code := range role.PermissionCodes() {
			if _, ok := seen[code]; ok {
				continue
			}
			seen[code] = struct{}{}
			codes = append(codes, code)
		}
	}
	return codes
}
