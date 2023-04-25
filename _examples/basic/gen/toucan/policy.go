// Code generated by toucan. DO NOT EDIT.
package policy

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	models "github.com/endigma/toucan/_examples/basic/models"
	strcase "github.com/iancoleman/strcase"
	"strings"
)

// resource `repository`

// Enum RepositoryPermission
type RepositoryPermission string

func (s RepositoryPermission) String() string {
	return string(s)
}

func (s RepositoryPermission) Valid() bool {
	_, err := ParseRepositoryPermission(string(s))
	return err == nil
}

var ErrInvalidRepositoryPermission = fmt.Errorf("not a valid repositoryPermission, try [%s]", strings.Join(repositoryPermissionNames, ", "))
var ErrNilRepositoryPermission = errors.New("value is nil")

const (
	RepositoryPermissionRead      RepositoryPermission = "read"
	RepositoryPermissionPush      RepositoryPermission = "push"
	RepositoryPermissionDelete    RepositoryPermission = "delete"
	RepositoryPermissionSnakeCase RepositoryPermission = "snake_case"
)

var repositoryPermissionNames = []string{string(RepositoryPermissionRead), string(RepositoryPermissionPush), string(RepositoryPermissionDelete), string(RepositoryPermissionSnakeCase)}
var repositoryPermissionMap = map[string]RepositoryPermission{
	"delete":     RepositoryPermissionDelete,
	"push":       RepositoryPermissionPush,
	"read":       RepositoryPermissionRead,
	"snake_case": RepositoryPermissionSnakeCase,
}

func RepositoryPermissionNames() []string {
	tmp := make([]string, len(repositoryPermissionNames))
	copy(tmp, repositoryPermissionNames)
	return tmp
}

func RepositoryPermissionValues() []RepositoryPermission {
	return []RepositoryPermission{RepositoryPermissionRead, RepositoryPermissionPush, RepositoryPermissionDelete, RepositoryPermissionSnakeCase}
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

func MustParseRepositoryPermission(s string) RepositoryPermission {
	x, err := ParseRepositoryPermission(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s RepositoryPermission) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *RepositoryPermission) UnmarshalText(data []byte) error {
	x, err := ParseRepositoryPermission(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *RepositoryPermission) Scan(value any) (err error) {
	if value == nil {
		*s = RepositoryPermission("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseRepositoryPermission(v)
	case []byte:
		*s, err = ParseRepositoryPermission(string(v))
	case RepositoryPermission:
		*s = v
	case *RepositoryPermission:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseRepositoryPermission(*v)
	default:
		return errors.New("invalid type for RepositoryPermission")
	}

	return
}

func (s RepositoryPermission) Value() (driver.Value, error) {
	return string(s), nil
}

// Enum RepositoryRole
type RepositoryRole string

func (s RepositoryRole) String() string {
	return string(s)
}

func (s RepositoryRole) Valid() bool {
	_, err := ParseRepositoryRole(string(s))
	return err == nil
}

var ErrInvalidRepositoryRole = fmt.Errorf("not a valid repositoryRole, try [%s]", strings.Join(repositoryRoleNames, ", "))
var ErrNilRepositoryRole = errors.New("value is nil")

const (
	RepositoryRoleOwner  RepositoryRole = "owner"
	RepositoryRoleEditor RepositoryRole = "editor"
	RepositoryRoleViewer RepositoryRole = "viewer"
)

var repositoryRoleNames = []string{string(RepositoryRoleOwner), string(RepositoryRoleEditor), string(RepositoryRoleViewer)}
var repositoryRoleMap = map[string]RepositoryRole{
	"editor": RepositoryRoleEditor,
	"owner":  RepositoryRoleOwner,
	"viewer": RepositoryRoleViewer,
}

func RepositoryRoleNames() []string {
	tmp := make([]string, len(repositoryRoleNames))
	copy(tmp, repositoryRoleNames)
	return tmp
}

func RepositoryRoleValues() []RepositoryRole {
	return []RepositoryRole{RepositoryRoleOwner, RepositoryRoleEditor, RepositoryRoleViewer}
}

func ParseRepositoryRole(s string) (RepositoryRole, error) {
	if x, ok := repositoryRoleMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := repositoryRoleMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return RepositoryRole(""), fmt.Errorf("%s is %w", s, ErrInvalidRepositoryRole)
}

func MustParseRepositoryRole(s string) RepositoryRole {
	x, err := ParseRepositoryRole(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s RepositoryRole) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *RepositoryRole) UnmarshalText(data []byte) error {
	x, err := ParseRepositoryRole(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *RepositoryRole) Scan(value any) (err error) {
	if value == nil {
		*s = RepositoryRole("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseRepositoryRole(v)
	case []byte:
		*s, err = ParseRepositoryRole(string(v))
	case RepositoryRole:
		*s = v
	case *RepositoryRole:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseRepositoryRole(*v)
	default:
		return errors.New("invalid type for RepositoryRole")
	}

	return
}

func (s RepositoryRole) Value() (driver.Value, error) {
	return string(s), nil
}

// Enum RepositoryAttribute
type RepositoryAttribute string

func (s RepositoryAttribute) String() string {
	return string(s)
}

func (s RepositoryAttribute) Valid() bool {
	_, err := ParseRepositoryAttribute(string(s))
	return err == nil
}

var ErrInvalidRepositoryAttribute = fmt.Errorf("not a valid repositoryAttribute, try [%s]", strings.Join(repositoryAttributeNames, ", "))
var ErrNilRepositoryAttribute = errors.New("value is nil")

const (
	RepositoryAttributePublic RepositoryAttribute = "public"
)

var repositoryAttributeNames = []string{string(RepositoryAttributePublic)}
var repositoryAttributeMap = map[string]RepositoryAttribute{"public": RepositoryAttributePublic}

func RepositoryAttributeNames() []string {
	tmp := make([]string, len(repositoryAttributeNames))
	copy(tmp, repositoryAttributeNames)
	return tmp
}

func RepositoryAttributeValues() []RepositoryAttribute {
	return []RepositoryAttribute{RepositoryAttributePublic}
}

func ParseRepositoryAttribute(s string) (RepositoryAttribute, error) {
	if x, ok := repositoryAttributeMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := repositoryAttributeMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return RepositoryAttribute(""), fmt.Errorf("%s is %w", s, ErrInvalidRepositoryAttribute)
}

func MustParseRepositoryAttribute(s string) RepositoryAttribute {
	x, err := ParseRepositoryAttribute(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s RepositoryAttribute) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *RepositoryAttribute) UnmarshalText(data []byte) error {
	x, err := ParseRepositoryAttribute(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *RepositoryAttribute) Scan(value any) (err error) {
	if value == nil {
		*s = RepositoryAttribute("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseRepositoryAttribute(v)
	case []byte:
		*s, err = ParseRepositoryAttribute(string(v))
	case RepositoryAttribute:
		*s = v
	case *RepositoryAttribute:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseRepositoryAttribute(*v)
	default:
		return errors.New("invalid type for RepositoryAttribute")
	}

	return
}

func (s RepositoryAttribute) Value() (driver.Value, error) {
	return string(s), nil
}

// Resolver for resource `repository`
type RepositoryResolver interface {
	HasRole(context context.Context, actor *models.User, role RepositoryRole, resource *models.Repository) bool
	HasAttribute(context context.Context, attribute RepositoryAttribute, resource *models.Repository) bool
}

// authorizer for resource `repository`

func (a Authorizer) AuthorizeRepository(ctx context.Context, actor *models.User, action RepositoryPermission, resource *models.Repository) bool {
	if !action.Valid() {
		return false
	}
	switch action {
	case RepositoryPermissionRead:
		return a.resolver.Repository().HasAttribute(ctx, RepositoryAttributePublic, resource) ||
			a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleOwner, resource) ||
			a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleEditor, resource) ||
			a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleViewer, resource)
	case RepositoryPermissionPush:
		return a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleOwner, resource) ||
			a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleEditor, resource)
	case RepositoryPermissionDelete:
		return a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleOwner, resource)
	case RepositoryPermissionSnakeCase:
		return a.resolver.Repository().HasRole(ctx, actor, RepositoryRoleOwner, resource)
	default:
		return false
	}
}

// resource `user`

// Enum UserPermission
type UserPermission string

func (s UserPermission) String() string {
	return string(s)
}

func (s UserPermission) Valid() bool {
	_, err := ParseUserPermission(string(s))
	return err == nil
}

var ErrInvalidUserPermission = fmt.Errorf("not a valid userPermission, try [%s]", strings.Join(userPermissionNames, ", "))
var ErrNilUserPermission = errors.New("value is nil")

const (
	UserPermissionRead   UserPermission = "read"
	UserPermissionWrite  UserPermission = "write"
	UserPermissionDelete UserPermission = "delete"
)

var userPermissionNames = []string{string(UserPermissionRead), string(UserPermissionWrite), string(UserPermissionDelete)}
var userPermissionMap = map[string]UserPermission{
	"delete": UserPermissionDelete,
	"read":   UserPermissionRead,
	"write":  UserPermissionWrite,
}

func UserPermissionNames() []string {
	tmp := make([]string, len(userPermissionNames))
	copy(tmp, userPermissionNames)
	return tmp
}

func UserPermissionValues() []UserPermission {
	return []UserPermission{UserPermissionRead, UserPermissionWrite, UserPermissionDelete}
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

func MustParseUserPermission(s string) UserPermission {
	x, err := ParseUserPermission(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s UserPermission) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *UserPermission) UnmarshalText(data []byte) error {
	x, err := ParseUserPermission(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *UserPermission) Scan(value any) (err error) {
	if value == nil {
		*s = UserPermission("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseUserPermission(v)
	case []byte:
		*s, err = ParseUserPermission(string(v))
	case UserPermission:
		*s = v
	case *UserPermission:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseUserPermission(*v)
	default:
		return errors.New("invalid type for UserPermission")
	}

	return
}

func (s UserPermission) Value() (driver.Value, error) {
	return string(s), nil
}

// Enum UserRole
type UserRole string

func (s UserRole) String() string {
	return string(s)
}

func (s UserRole) Valid() bool {
	_, err := ParseUserRole(string(s))
	return err == nil
}

var ErrInvalidUserRole = fmt.Errorf("not a valid userRole, try [%s]", strings.Join(userRoleNames, ", "))
var ErrNilUserRole = errors.New("value is nil")

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleSelf   UserRole = "self"
	UserRoleViewer UserRole = "viewer"
)

var userRoleNames = []string{string(UserRoleAdmin), string(UserRoleSelf), string(UserRoleViewer)}
var userRoleMap = map[string]UserRole{
	"admin":  UserRoleAdmin,
	"self":   UserRoleSelf,
	"viewer": UserRoleViewer,
}

func UserRoleNames() []string {
	tmp := make([]string, len(userRoleNames))
	copy(tmp, userRoleNames)
	return tmp
}

func UserRoleValues() []UserRole {
	return []UserRole{UserRoleAdmin, UserRoleSelf, UserRoleViewer}
}

func ParseUserRole(s string) (UserRole, error) {
	if x, ok := userRoleMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := userRoleMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return UserRole(""), fmt.Errorf("%s is %w", s, ErrInvalidUserRole)
}

func MustParseUserRole(s string) UserRole {
	x, err := ParseUserRole(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s UserRole) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *UserRole) UnmarshalText(data []byte) error {
	x, err := ParseUserRole(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *UserRole) Scan(value any) (err error) {
	if value == nil {
		*s = UserRole("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseUserRole(v)
	case []byte:
		*s, err = ParseUserRole(string(v))
	case UserRole:
		*s = v
	case *UserRole:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseUserRole(*v)
	default:
		return errors.New("invalid type for UserRole")
	}

	return
}

func (s UserRole) Value() (driver.Value, error) {
	return string(s), nil
}

// Enum UserAttribute
type UserAttribute string

func (s UserAttribute) String() string {
	return string(s)
}

func (s UserAttribute) Valid() bool {
	_, err := ParseUserAttribute(string(s))
	return err == nil
}

var ErrInvalidUserAttribute = fmt.Errorf("not a valid userAttribute, try [%s]", strings.Join(userAttributeNames, ", "))
var ErrNilUserAttribute = errors.New("value is nil")

const (
	UserAttributePublic UserAttribute = "public"
)

var userAttributeNames = []string{string(UserAttributePublic)}
var userAttributeMap = map[string]UserAttribute{"public": UserAttributePublic}

func UserAttributeNames() []string {
	tmp := make([]string, len(userAttributeNames))
	copy(tmp, userAttributeNames)
	return tmp
}

func UserAttributeValues() []UserAttribute {
	return []UserAttribute{UserAttributePublic}
}

func ParseUserAttribute(s string) (UserAttribute, error) {
	if x, ok := userAttributeMap[s]; ok {
		return x, nil
	}

	// Try to parse from snake case
	if x, ok := userAttributeMap[strcase.ToSnake(s)]; ok {
		return x, nil
	}

	return UserAttribute(""), fmt.Errorf("%s is %w", s, ErrInvalidUserAttribute)
}

func MustParseUserAttribute(s string) UserAttribute {
	x, err := ParseUserAttribute(s)
	if err != nil {
		panic(err)
	}

	return x
}

func (s UserAttribute) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

func (s *UserAttribute) UnmarshalText(data []byte) error {
	x, err := ParseUserAttribute(string(data))
	if err != nil {
		return err
	}

	*s = x
	return nil
}

func (s *UserAttribute) Scan(value any) (err error) {
	if value == nil {
		*s = UserAttribute("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*s, err = ParseUserAttribute(v)
	case []byte:
		*s, err = ParseUserAttribute(string(v))
	case UserAttribute:
		*s = v
	case *UserAttribute:
		if v == nil {
			return
		}
		*s = *v
	case *string:
		if v == nil {
			return
		}
		*s, err = ParseUserAttribute(*v)
	default:
		return errors.New("invalid type for UserAttribute")
	}

	return
}

func (s UserAttribute) Value() (driver.Value, error) {
	return string(s), nil
}

// Resolver for resource `user`
type UserResolver interface {
	HasRole(context context.Context, actor *models.User, role UserRole, resource *models.User) bool
	HasAttribute(context context.Context, attribute UserAttribute, resource *models.User) bool
}

// authorizer for resource `user`

func (a Authorizer) AuthorizeUser(ctx context.Context, actor *models.User, action UserPermission, resource *models.User) bool {
	if !action.Valid() {
		return false
	}
	switch action {
	case UserPermissionRead:
		return a.resolver.User().HasAttribute(ctx, UserAttributePublic, resource) ||
			a.resolver.User().HasRole(ctx, actor, UserRoleAdmin, resource) ||
			a.resolver.User().HasRole(ctx, actor, UserRoleSelf, resource) ||
			a.resolver.User().HasRole(ctx, actor, UserRoleViewer, resource)
	case UserPermissionWrite:
		return a.resolver.User().HasRole(ctx, actor, UserRoleAdmin, resource) ||
			a.resolver.User().HasRole(ctx, actor, UserRoleSelf, resource)
	case UserPermissionDelete:
		return a.resolver.User().HasRole(ctx, actor, UserRoleAdmin, resource)
	default:
		return false
	}
}

// Global resolver
type Resolver interface {
	Repository() RepositoryResolver
	User() UserResolver
}

// Global authorizer
type Authorizer struct {
	resolver Resolver
}

func (a Authorizer) Authorize(ctx context.Context, actor *models.User, permission string, resource any) bool {
	switch resource.(type) {
	case *models.Repository:
		perm, err := ParseRepositoryPermission(permission)
		if err == nil {
			return a.AuthorizeRepository(ctx, actor, perm, resource.(*models.Repository))
		}
	case *models.User:
		perm, err := ParseUserPermission(permission)
		if err == nil {
			return a.AuthorizeUser(ctx, actor, perm, resource.(*models.User))
		}
	}

	return false
}

func NewAuthorizer(resolver Resolver) *Authorizer {
	return &Authorizer{resolver: resolver}
}