package auth

import (
	"context"
	"errors"
	"sort"
	"strconv"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, user User, password string, roleNames []string) (User, error)
	UpdateUser(ctx context.Context, user User, password string) (User, error)
	DeleteUser(ctx context.Context, id string) error
	ListRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, role Role) (Role, error)
	UpdateRole(ctx context.Context, role Role) (Role, error)
	DeleteRole(ctx context.Context, name string) error
	ListPermissions(ctx context.Context) ([]Permission, error)
}

type MemoryRepository struct {
	nextUserID  int
	users       map[string]User
	roles       map[string]Role
	permissions map[string]Permission
}

func NewMemoryRepository(users ...User) *MemoryRepository {
	repo := &MemoryRepository{
		nextUserID:  1,
		users:       make(map[string]User, len(users)),
		roles:       defaultRoles(),
		permissions: defaultPermissions(),
	}
	for _, user := range users {
		if user.ID == "" {
			user.ID = strconv.Itoa(repo.nextUserID)
			repo.nextUserID++
		}
		repo.users[user.Username] = user
		for _, role := range user.Roles {
			repo.roles[role.Name] = role
		}
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

func (r *MemoryRepository) ListUsers(_ context.Context) ([]User, error) {
	users := make([]User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].Username < users[j].Username })
	return users, nil
}

func (r *MemoryRepository) ListRoles(_ context.Context) ([]Role, error) {
	roles := make([]Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Name < roles[j].Name })
	return roles, nil
}

func (r *MemoryRepository) CreateUser(_ context.Context, user User, password string, roleNames []string) (User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}
	user.ID = strconv.Itoa(r.nextUserID)
	r.nextUserID++
	user.PasswordHash = hash
	user.Roles = r.rolesByName(roleNames)
	if user.Status == "" {
		user.Status = "active"
	}
	r.users[user.Username] = user
	return user, nil
}

func (r *MemoryRepository) UpdateUser(_ context.Context, user User, password string) (User, error) {
	var existingKey string
	var existing User
	for username, candidate := range r.users {
		if candidate.ID == user.ID {
			existingKey = username
			existing = candidate
			break
		}
	}
	if existing.ID == "" {
		return User{}, ErrUserNotFound
	}
	if password != "" {
		hash, err := HashPassword(password)
		if err != nil {
			return User{}, err
		}
		existing.PasswordHash = hash
	}
	existing.Username = user.Username
	existing.DisplayName = user.DisplayName
	existing.Status = user.Status
	existing.Roles = r.rolesByName(roleNames(user.Roles))
	delete(r.users, existingKey)
	r.users[existing.Username] = existing
	return existing, nil
}

func (r *MemoryRepository) DeleteUser(_ context.Context, id string) error {
	for username, user := range r.users {
		if user.ID == id {
			delete(r.users, username)
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *MemoryRepository) CreateRole(_ context.Context, role Role) (Role, error) {
	role.Permissions = r.permissionsByCode(role.Permissions)
	r.roles[role.Name] = role
	return role, nil
}

func (r *MemoryRepository) UpdateRole(_ context.Context, role Role) (Role, error) {
	if _, ok := r.roles[role.Name]; !ok {
		return Role{}, ErrRoleNotFound
	}
	role.Permissions = r.permissionsByCode(role.Permissions)
	r.roles[role.Name] = role
	return role, nil
}

func (r *MemoryRepository) DeleteRole(_ context.Context, name string) error {
	if _, ok := r.roles[name]; !ok {
		return ErrRoleNotFound
	}
	delete(r.roles, name)
	for username, user := range r.users {
		filtered := user.Roles[:0]
		for _, role := range user.Roles {
			if role.Name != name {
				filtered = append(filtered, role)
			}
		}
		user.Roles = filtered
		r.users[username] = user
	}
	return nil
}

func (r *MemoryRepository) ListPermissions(_ context.Context) ([]Permission, error) {
	permissions := make([]Permission, 0, len(r.permissions))
	for _, permission := range r.permissions {
		permissions = append(permissions, permission)
	}
	sort.Slice(permissions, func(i, j int) bool { return permissions[i].Code < permissions[j].Code })
	return permissions, nil
}

func (r *MemoryRepository) rolesByName(names []string) []Role {
	roles := make([]Role, 0, len(names))
	for _, name := range names {
		if role, ok := r.roles[name]; ok {
			roles = append(roles, role)
		}
	}
	return roles
}

func (r *MemoryRepository) permissionsByCode(items []Permission) []Permission {
	permissions := make([]Permission, 0, len(items))
	for _, item := range items {
		if permission, ok := r.permissions[item.Code]; ok {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}

func roleNames(roles []Role) []string {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return names
}

func defaultPermissions() map[string]Permission {
	descriptions := map[string]string{
		"cmdb.asset.read":     "查看 CMDB 资产",
		"cmdb.asset.write":    "管理 CMDB 资产",
		"cmdb.model.read":     "查看 CMDB 模型",
		"cmdb.model.write":    "管理 CMDB 模型",
		"cmdb.apikey.read":    "查看 CMDB API Key",
		"cmdb.apikey.write":   "管理 CMDB API Key",
		"auth.user.read":      "查看用户",
		"auth.user.write":     "管理用户",
		"auth.role.read":      "查看角色",
		"auth.role.write":     "管理角色和授权",
		"audit.log.read":      "查看审计日志",
		"system.config.read":  "查看系统配置",
		"system.config.write": "管理系统配置",
	}
	permissions := make(map[string]Permission, len(descriptions))
	for code, description := range descriptions {
		permissions[code] = Permission{Code: code, Description: description}
	}
	return permissions
}

func defaultRoles() map[string]Role {
	permissions := defaultPermissions()
	rolePermissions := map[string][]string{
		"system_admin": {"cmdb.asset.read", "cmdb.asset.write", "cmdb.model.read", "cmdb.model.write", "cmdb.apikey.read", "cmdb.apikey.write", "auth.user.read", "auth.user.write", "auth.role.read", "auth.role.write", "audit.log.read", "system.config.read", "system.config.write"},
		"ops_admin":    {"cmdb.asset.read", "cmdb.asset.write", "cmdb.model.read", "cmdb.model.write", "cmdb.apikey.read", "cmdb.apikey.write", "audit.log.read"},
		"ops_operator": {"cmdb.asset.read", "cmdb.asset.write", "cmdb.model.read"},
		"readonly":     {"cmdb.asset.read", "cmdb.model.read"},
	}
	displayNames := map[string]string{
		"system_admin": "系统管理员",
		"ops_admin":    "运维管理员",
		"ops_operator": "运维操作员",
		"readonly":     "只读用户",
	}
	roles := make(map[string]Role, len(rolePermissions))
	for name, codes := range rolePermissions {
		role := Role{Name: name, DisplayName: displayNames[name]}
		for _, code := range codes {
			role.Permissions = append(role.Permissions, permissions[code])
		}
		roles[name] = role
	}
	return roles
}
