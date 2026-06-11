package auth

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
}

type MemoryRepository struct {
	users map[string]User
}

func NewMemoryRepository(users ...User) *MemoryRepository {
	repo := &MemoryRepository{users: make(map[string]User, len(users))}
	for _, user := range users {
		repo.users[user.Username] = user
	}
	return repo
}

func NewDevelopmentRepository() (*MemoryRepository, error) {
	hash, err := HashPassword("admin123")
	if err != nil {
		return nil, err
	}
	return NewMemoryRepository(User{
		ID:           "00000000-0000-0000-0000-000000000001",
		Username:     "admin",
		PasswordHash: hash,
		DisplayName:  "管理员",
		Status:       "active",
		Roles: []Role{{
			Name: "system_admin",
			Permissions: []Permission{
				{Code: "cmdb.asset.read"},
				{Code: "cmdb.asset.write"},
				{Code: "cmdb.model.read"},
				{Code: "cmdb.model.write"},
				{Code: "auth.user.read"},
				{Code: "auth.role.read"},
				{Code: "audit.log.read"},
				{Code: "system.config.read"},
				{Code: "system.config.write"},
			},
		}},
	}), nil
}

func (r *MemoryRepository) FindByUsername(_ context.Context, username string) (User, error) {
	user, ok := r.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}
