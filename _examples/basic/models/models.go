package models

import (
	"github.com/rs/xid"
)

type User struct {
	ID xid.ID `json:"id"`

	Name string `json:"name"`

	GlobalAdmin bool `json:"admin"`

	Roles []RepositoryRole `json:"roles"`
}

func NewUser(name string, admin bool, roles ...RepositoryRole) *User {
	return &User{
		ID:    xid.New(),
		Name:  name,
		Roles: roles,
	}
}

type Repository struct {
	ID     xid.ID `json:"id"`
	Label  string `json:"label"`
	Public bool   `json:"public"`
}

type RepositoryRole struct {
	Role string
	Repo xid.ID
}

func NewRepository(label string, public bool) *Repository {
	return &Repository{
		ID:     xid.New(),
		Label:  label,
		Public: public,
	}
}

func (u *User) ToucanKey() string       { return u.ID.String() }
func (r *Repository) ToucanKey() string { return r.ID.String() }
