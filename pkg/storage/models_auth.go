package storage

import (
	"github.com/vpmv/bibl.io/pkg/dto"
	"gorm.io/gorm"
)

type APIToken struct {
	gorm.Model
	Token       string `json:"token"`
	ExternalID  string `json:"application"`
	Description string `json:"description"`

	Permissions []*Permission `gorm:"many2many:api_permissions;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
}

func (p Permission) DTO() *dto.Permission {
	return &dto.Permission{
		Name:        p.Name,
		Description: p.Description,
	}
}

func (t *APIToken) HasPermission(name string) bool {
	for _, perm := range t.Permissions {
		if perm.Name == name {
			return true
		}
	}
	return false
}

func (t *APIToken) DTO() *dto.Authorization {
	permissions := make([]dto.Permission, len(t.Permissions))
	for i, perm := range t.Permissions {
		permissions[i] = *perm.DTO()
	}
	return &dto.Authorization{
		Token:       t.Token,
		Description: t.Description,
		Permissions: permissions,
	}
}
