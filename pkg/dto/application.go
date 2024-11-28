package dto

type Authorization struct {
	Token       string       `json:"token"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

func (a *Authorization) HasPermission(permission string) bool {
	for _, perm := range a.Permissions {
		if perm.Name == permission {
			return true
		}
	}

	return false
}

type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
