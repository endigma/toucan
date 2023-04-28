// Code generated by toucan. DO NOT EDIT.
package toucan

import (
	"errors"
	"fmt"
	strcase "github.com/iancoleman/strcase"
	"strings"
)

type RepositoryPermission string

const (
	RepositoryPermissionRead      RepositoryPermission = "read"
	RepositoryPermissionPush      RepositoryPermission = "push"
	RepositoryPermissionDelete    RepositoryPermission = "delete"
	RepositoryPermissionSnakeCase RepositoryPermission = "snake_case"
)

func (s RepositoryPermission) String() string {
	return string(s)
}

func (s RepositoryPermission) Valid() bool {
	_, err := ParseRepositoryPermission(string(s))
	return err == nil
}

var (
	ErrInvalidRepositoryPermission = fmt.Errorf("not a valid RepositoryPermission, try [%s]", strings.Join(repositoryPermissionNames, ", "))
	ErrNilRepositoryPermission     = errors.New("value is nil")
)

var (
	repositoryPermissionMap = map[string]RepositoryPermission{
		"delete":     RepositoryPermissionDelete,
		"push":       RepositoryPermissionPush,
		"read":       RepositoryPermissionRead,
		"snake_case": RepositoryPermissionSnakeCase,
	}
	repositoryPermissionNames = []string{string(RepositoryPermissionRead), string(RepositoryPermissionPush), string(RepositoryPermissionDelete), string(RepositoryPermissionSnakeCase)}
)

func ParseRepositoryPermission(s string) (RepositoryPermission, error) {
	if x, ok := repositoryPermissionMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := repositoryPermissionMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return RepositoryPermission(""), fmt.Errorf("%s is %w", s, ErrInvalidRepositoryPermission)
}

type UserPermission string

const (
	UserPermissionRead   UserPermission = "read"
	UserPermissionWrite  UserPermission = "write"
	UserPermissionDelete UserPermission = "delete"
)

func (s UserPermission) String() string {
	return string(s)
}

func (s UserPermission) Valid() bool {
	_, err := ParseUserPermission(string(s))
	return err == nil
}

var (
	ErrInvalidUserPermission = fmt.Errorf("not a valid UserPermission, try [%s]", strings.Join(userPermissionNames, ", "))
	ErrNilUserPermission     = errors.New("value is nil")
)

var (
	userPermissionMap = map[string]UserPermission{
		"delete": UserPermissionDelete,
		"read":   UserPermissionRead,
		"write":  UserPermissionWrite,
	}
	userPermissionNames = []string{string(UserPermissionRead), string(UserPermissionWrite), string(UserPermissionDelete)}
)

func ParseUserPermission(s string) (UserPermission, error) {
	if x, ok := userPermissionMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := userPermissionMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return UserPermission(""), fmt.Errorf("%s is %w", s, ErrInvalidUserPermission)
}