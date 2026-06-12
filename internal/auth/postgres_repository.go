package auth

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type userRow struct {
	ID           string
	Username     string
	PasswordHash string
	DisplayName  string
	Status       string
}

type roleRow struct {
	Name        string
	DisplayName string
}

type permissionRow struct {
	Code        string
	Description string
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindByUsername(ctx context.Context, username string) (User, error) {
	var row userRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, username, password_hash, display_name, status
FROM users
WHERE username = ?`, username).Scan(&row).Error; err != nil {
		return User{}, err
	}
	if row.ID == "" {
		return User{}, ErrUserNotFound
	}
	roles, err := r.rolesForUser(ctx, row.ID)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		DisplayName:  row.DisplayName,
		Status:       row.Status,
		Roles:        roles,
	}, nil
}

func (r *PostgresRepository) ListUsers(ctx context.Context) ([]User, error) {
	var rows []userRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, username, password_hash, display_name, status
FROM users
ORDER BY username`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		roles, err := r.rolesForUser(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, User{
			ID:          row.ID,
			Username:    row.Username,
			DisplayName: row.DisplayName,
			Status:      row.Status,
			Roles:       roles,
		})
	}
	return users, nil
}

func (r *PostgresRepository) ListRoles(ctx context.Context) ([]Role, error) {
	var rows []roleRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT name, display_name
FROM roles
ORDER BY name`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	roles := make([]Role, 0, len(rows))
	for _, row := range rows {
		permissions, err := r.permissionsForRole(ctx, row.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, Role{Name: row.Name, DisplayName: row.DisplayName, Permissions: permissions})
	}
	return roles, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user User, password string, roleNames []string) (User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}
	if user.Status == "" {
		user.Status = "active"
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row userRow
		if err := tx.Raw(`
INSERT INTO users (username, password_hash, display_name, status)
VALUES (?, ?, ?, ?)
RETURNING id::text, username, password_hash, display_name, status`,
			user.Username, hash, user.DisplayName, user.Status).Scan(&row).Error; err != nil {
			return err
		}
		user.ID = row.ID
		return saveUserRoles(tx, user.ID, roleNames)
	})
	if err != nil {
		return User{}, err
	}
	return r.userByID(ctx, user.ID)
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user User, password string) (User, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if password != "" {
			hash, err := HashPassword(password)
			if err != nil {
				return err
			}
			if err := tx.Exec(`
UPDATE users SET username = ?, password_hash = ?, display_name = ?, status = ?, updated_at = now()
WHERE id = ?::uuid`, user.Username, hash, user.DisplayName, user.Status, user.ID).Error; err != nil {
				return err
			}
		} else if err := tx.Exec(`
UPDATE users SET username = ?, display_name = ?, status = ?, updated_at = now()
WHERE id = ?::uuid`, user.Username, user.DisplayName, user.Status, user.ID).Error; err != nil {
			return err
		}
		return saveUserRoles(tx, user.ID, roleNames(user.Roles))
	})
	if err != nil {
		return User{}, err
	}
	return r.userByID(ctx, user.ID)
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?::uuid", id).Error; err != nil {
			return err
		}
		result := tx.Exec("DELETE FROM users WHERE id = ?::uuid", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrUserNotFound
		}
		return nil
	})
}

func (r *PostgresRepository) CreateRole(ctx context.Context, role Role) (Role, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
INSERT INTO roles (name, display_name)
VALUES (?, ?)`, role.Name, role.DisplayName).Error; err != nil {
			return err
		}
		return saveRolePermissions(tx, role.Name, role.PermissionCodes())
	})
	if err != nil {
		return Role{}, err
	}
	return r.roleByName(ctx, role.Name)
}

func (r *PostgresRepository) UpdateRole(ctx context.Context, role Role) (Role, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec("UPDATE roles SET display_name = ? WHERE name = ?", role.DisplayName, role.Name)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrRoleNotFound
		}
		return saveRolePermissions(tx, role.Name, role.PermissionCodes())
	})
	if err != nil {
		return Role{}, err
	}
	return r.roleByName(ctx, role.Name)
}

func (r *PostgresRepository) DeleteRole(ctx context.Context, name string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var roleID string
		if err := tx.Raw("SELECT id::text FROM roles WHERE name = ?", name).Scan(&roleID).Error; err != nil {
			return err
		}
		if roleID == "" {
			return ErrRoleNotFound
		}
		if err := tx.Exec("DELETE FROM user_roles WHERE role_id = ?::uuid", roleID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?::uuid", roleID).Error; err != nil {
			return err
		}
		return tx.Exec("DELETE FROM roles WHERE id = ?::uuid", roleID).Error
	})
}

func (r *PostgresRepository) ListPermissions(ctx context.Context) ([]Permission, error) {
	var rows []permissionRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT code, description
FROM permissions
ORDER BY code`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	permissions := make([]Permission, 0, len(rows))
	for _, row := range rows {
		permissions = append(permissions, Permission{Code: row.Code, Description: row.Description})
	}
	return permissions, nil
}

func (r *PostgresRepository) rolesForUser(ctx context.Context, userID string) ([]Role, error) {
	var rows []roleRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT r.name, r.display_name
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = ?::uuid
ORDER BY r.name`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	roles := make([]Role, 0, len(rows))
	for _, row := range rows {
		permissions, err := r.permissionsForRole(ctx, row.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, Role{Name: row.Name, DisplayName: row.DisplayName, Permissions: permissions})
	}
	return roles, nil
}

func (r *PostgresRepository) permissionsForRole(ctx context.Context, roleName string) ([]Permission, error) {
	var rows []permissionRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT p.code, p.description
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN roles r ON r.id = rp.role_id
WHERE r.name = ?
ORDER BY p.code`, roleName).Scan(&rows).Error; err != nil {
		return nil, err
	}
	permissions := make([]Permission, 0, len(rows))
	for _, row := range rows {
		permissions = append(permissions, Permission{Code: row.Code, Description: row.Description})
	}
	return permissions, nil
}

func (r *PostgresRepository) userByID(ctx context.Context, id string) (User, error) {
	var row userRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, username, password_hash, display_name, status
FROM users
WHERE id = ?::uuid`, id).Scan(&row).Error; err != nil {
		return User{}, err
	}
	if row.ID == "" {
		return User{}, ErrUserNotFound
	}
	roles, err := r.rolesForUser(ctx, row.ID)
	if err != nil {
		return User{}, err
	}
	return User{ID: row.ID, Username: row.Username, DisplayName: row.DisplayName, Status: row.Status, Roles: roles}, nil
}

func (r *PostgresRepository) roleByName(ctx context.Context, name string) (Role, error) {
	var row roleRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT name, display_name
FROM roles
WHERE name = ?`, name).Scan(&row).Error; err != nil {
		return Role{}, err
	}
	if row.Name == "" {
		return Role{}, ErrRoleNotFound
	}
	permissions, err := r.permissionsForRole(ctx, row.Name)
	if err != nil {
		return Role{}, err
	}
	return Role{Name: row.Name, DisplayName: row.DisplayName, Permissions: permissions}, nil
}

func saveUserRoles(tx *gorm.DB, userID string, roleNames []string) error {
	if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?::uuid", userID).Error; err != nil {
		return err
	}
	for _, name := range roleNames {
		if err := tx.Exec(`
INSERT INTO user_roles (user_id, role_id)
SELECT ?::uuid, id FROM roles WHERE name = ?
ON CONFLICT DO NOTHING`, userID, name).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveRolePermissions(tx *gorm.DB, roleName string, permissionCodes []string) error {
	var roleID string
	if err := tx.Raw("SELECT id::text FROM roles WHERE name = ?", roleName).Scan(&roleID).Error; err != nil {
		return err
	}
	if roleID == "" {
		return ErrRoleNotFound
	}
	if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?::uuid", roleID).Error; err != nil {
		return err
	}
	for _, code := range permissionCodes {
		if err := tx.Exec(`
INSERT INTO role_permissions (role_id, permission_id)
SELECT ?::uuid, id FROM permissions WHERE code = ?
ON CONFLICT DO NOTHING`, roleID, code).Error; err != nil {
			return err
		}
	}
	return nil
}
