package auth

type Permission struct {
	Code string
}

type Role struct {
	Name        string
	Permissions []Permission
}

func (r Role) Allows(code string) bool {
	for _, permission := range r.Permissions {
		if permission.Code == code {
			return true
		}
	}
	return false
}

func (r Role) PermissionCodes() []string {
	codes := make([]string, 0, len(r.Permissions))
	for _, permission := range r.Permissions {
		codes = append(codes, permission.Code)
	}
	return codes
}
