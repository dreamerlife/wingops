package auth

type Permission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Role struct {
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Permissions []Permission `json:"permissions"`
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
