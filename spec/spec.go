package spec

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/endigma/toucan/config"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type Spec struct {
	Actor     QualifierSpec  `validate:"required"`
	Resources []ResourceSpec `validate:"required,unique=Model,unique=Name,dive"`

	Output OutputSpec `validate:"required,dive"`
}

type OutputSpec struct {
	Path    string `validate:"required"`
	Package string `validate:"required,validName"`
}

type ResourceSpec struct {
	Name        string          `validate:"required,validName,notReserved"`
	Model       QualifierSpec   `validate:"required,dive"`
	Permissions []string        `validate:"unique,dive,required"`
	Roles       []RoleSpec      `validate:"unique=Name,dive,required"`
	Attributes  []AttributeSpec `validate:"unique=Name,dive,required"`
}

type RoleSpec struct {
	Name        string   `validate:"required,validName"`
	Permissions []string `validate:"required,dive,required"`
}

type AttributeSpec struct {
	Name        string   `validate:"required,validName"`
	Permissions []string `validate:"required,unique,dive,required"`
}

var validate *validator.Validate = validator.New()

func init() {
	lo.Must0(validate.RegisterValidation("validName", nameValidator))
	lo.Must0(validate.RegisterValidation("notReserved", notReservedValidator))
	lo.Must0(validate.RegisterValidation("validQualName", qualifierNameValidator))
	lo.Must0(validate.RegisterValidation("validQualPath", qualifierPathValidator))
}

type QualifierSpec struct {
	Path string `validate:"required,validQualPath"`
	Name string `validate:"required,validName,validQualName"`
}

func (s *Spec) Validate() error {
	var result *multierror.Error

	for _, resource := range s.Resources {
		isValidPerm := func(perm string) bool {
			return lo.Contains(resource.Permissions, perm)
		}

		for _, attr := range resource.Attributes {
			// Catch invalid attribute permissions
			for _, permission := range attr.Permissions {
				if !isValidPerm(permission) {
					result = multierror.Append(result, fmt.Errorf("invalid permission %q in attribute %q", permission, attr.Name))
				}
			}
		}

		for _, role := range resource.Roles {
			// Catch invalid role permissions
			for _, permission := range role.Permissions {
				if !isValidPerm(permission) {
					result = multierror.Append(result, fmt.Errorf("invalid permission %q in role %q", permission, role.Name))
				}
			}
		}

		for _, permission := range resource.Permissions {
			// Catch unused permissions
			hasAttribute := lo.SomeBy(resource.Attributes, func(attr AttributeSpec) bool {
				return lo.Contains(attr.Permissions, permission)
			})

			hasRole := lo.SomeBy(resource.Roles, func(role RoleSpec) bool {
				return lo.Contains(role.Permissions, permission)
			})

			if !hasAttribute && !hasRole {
				result = multierror.Append(result, fmt.Errorf("unused permission %q", permission))
			}
		}
	}

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			result = multierror.Append(result, fmt.Errorf("invalid %s: %s", err.Field(), err.Tag()))
		}
	}

	return result.ErrorOrNil()
}

func FromConfig(config *config.Config) (*Spec, error) {
	actor, err := parseQualifier(config.Actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor %q: %w", config.Actor, err)
	}

	spec := &Spec{Actor: actor}

	for _, cfgResource := range config.Resources {
		roles := []RoleSpec{}

		for _, cfgRole := range cfgResource.Roles {
			roles = append(roles, RoleSpec{
				Name:        cfgRole.Name,
				Permissions: cfgRole.Permissions,
			})
		}

		attributes := []AttributeSpec{}
		for _, cfgAttr := range cfgResource.Attributes {
			attributes = append(attributes, AttributeSpec{
				Name:        cfgAttr.Name,
				Permissions: cfgAttr.Permissions,
			})
		}

		model, err := parseQualifier(cfgResource.Model)
		if err != nil {
			return nil, fmt.Errorf("invalid model %q: %w", cfgResource.Model, err)
		}

		spec.Resources = append(spec.Resources, ResourceSpec{
			Name:        cfgResource.Name,
			Model:       model,
			Permissions: cfgResource.Permissions,
			Roles:       roles,
			Attributes:  attributes,
		})
	}

	spec.Output = OutputSpec{
		Path:    config.Output.Path,
		Package: config.Output.Package,
	}

	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}

	return spec, nil
}

func parseQualifier(qualifier string) (QualifierSpec, error) {
	re := regexp.MustCompile(`^([\w|.|/]+)(?:\.)([A-Z][a-zA-Z0-9_-]*)$`)
	matches := re.FindStringSubmatch(qualifier)

	if err := lo.Validate(len(matches) == 3, "Expected 2 matches, path and name"); err != nil {
		return QualifierSpec{}, fmt.Errorf("invalid qualifier %q: %w", qualifier, err)
	}

	return QualifierSpec{
		Path: matches[1],
		Name: matches[2],
	}, nil
}

func validName(name string) bool {
	// - must begin with a letter, and can have any number of additional letters and numbers.
	// - cannot start with a number.
	// - cannot contain spaces.
	// - cannot contain (very) special characters.
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	return re.MatchString(name)
}

func qualifierPathValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.String {
		return false
	}

	re := regexp.MustCompile(`^([\w|.|/]+)$`)
	return re.MatchString(fl.Field().String())
}

func qualifierNameValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.String {
		return false
	}

	re := regexp.MustCompile(`^([A-Z][a-zA-Z0-9_-]*)$`)
	return re.MatchString(fl.Field().String())
}

func notReservedValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return !lo.Contains([]string{"resolver", "authorizer"}, strcase.ToLowerCamel(field.String()))
}

func nameValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return validName(field.String())
}

type PermissionSource struct {
	Type string // role, attribute
	Name string
}

// GetPermissionSources returns all sources of a permission, and a boolean `true` if the
func (resource ResourceSpec) GetPermissionSources(permission string) []PermissionSource {
	sources := lo.Union(
		lo.FilterMap(resource.Attributes, func(attr AttributeSpec, _ int) (PermissionSource, bool) {
			if lo.Contains(attr.Permissions, permission) {
				return PermissionSource{
					Type: "attribute",
					Name: strcase.ToCamel(resource.Name + "Attribute" + strcase.ToCamel(attr.Name)),
				}, true
			}

			return PermissionSource{}, false
		}),
		lo.FilterMap(resource.Roles, func(role RoleSpec, _ int) (PermissionSource, bool) {
			if lo.Contains(role.Permissions, permission) {
				return PermissionSource{
					Type: "role",
					Name: strcase.ToCamel(resource.Name + "Role" + strcase.ToCamel(role.Name)),
				}, true
			}

			return PermissionSource{}, false
		}),
	)

	return sources
}
