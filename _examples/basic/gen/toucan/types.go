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

func (s RepositoryPermission) Valid() bool {
	_, err := ParseRepositoryPermission(string(s))
	return err == nil
}

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

func (s UserPermission) Valid() bool {
	_, err := ParseUserPermission(string(s))
	return err == nil
}

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

type GlobalPermission string

const (
	GlobalPermissionReadAllUsers    GlobalPermission = "read_all_users"
	GlobalPermissionWriteAllUsers   GlobalPermission = "write_all_users"
	GlobalPermissionReadAllProfiles GlobalPermission = "read_all_profiles"
)

var (
	ErrInvalidGlobalPermission = fmt.Errorf("not a valid GlobalPermission, try [%s]", strings.Join(globalPermissionNames, ", "))
	ErrNilGlobalPermission     = errors.New("value is nil")
)

var (
	globalPermissionMap = map[string]GlobalPermission{
		"read_all_profiles": GlobalPermissionReadAllProfiles,
		"read_all_users":    GlobalPermissionReadAllUsers,
		"write_all_users":   GlobalPermissionWriteAllUsers,
	}
	globalPermissionNames = []string{string(GlobalPermissionReadAllUsers), string(GlobalPermissionWriteAllUsers), string(GlobalPermissionReadAllProfiles)}
)

func (s GlobalPermission) Valid() bool {
	_, err := ParseGlobalPermission(string(s))
	return err == nil
}

func ParseGlobalPermission(s string) (GlobalPermission, error) {
	if x, ok := globalPermissionMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := globalPermissionMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return GlobalPermission(""), fmt.Errorf("%s is %w", s, ErrInvalidGlobalPermission)
}
