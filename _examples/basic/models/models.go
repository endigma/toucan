package models

import (
	"github.com/rs/xid"
)

type User struct {
	ID xid.ID `json:"id"`

	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

func NewUser(name string) *User {
	return &User{
		ID:   xid.New(),
		Name: name,
	}
}

type Repository struct {
	ID      xid.ID `json:"id"`
	Label   string `json:"label"`
	Public  bool   `json:"public"`
	Secret0 string `json:"secret0"` // barely restricted
	Secret1 string `json:"secret1"` // restricted
	Secret2 string `json:"secret2"` // highly restricted
}

// func (r *Repository) HasRole(user *User, role string) bool {
// 	return lo.ContainsBy(r.Roles, func(r RepositoryRole) bool {
// 		return r.User.ID == user.ID && r.Role == role
// 	})
// }

func NewRepository(label string, public bool) *Repository {
	return &Repository{
		ID:     xid.New(),
		Label:  label,
		Public: public,
	}
}
