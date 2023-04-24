package schema

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

//nolint:gochecknoglobals
var validate *validator.Validate = NewSchemaValidator()

var (
	ErrInvalidPermission = errors.New("invalid permission")
	ErrInvalidSchema     = errors.New("invalid schema")
	ErrUnusedPermission  = errors.New("unused permission")
)

func (s *Schema) Validate() error {
	var result *multierror.Error

	err := validate.Struct(s)
	if err != nil {
		errors.As(err, &validator.ValidationErrors{})

		//nolint:forcetypeassert,errorlint
		for _, err := range err.(validator.ValidationErrors) {
			result = multierror.Append(result, fmt.Errorf("%w: %s: %s", ErrInvalidSchema, err.Field(), err.Tag()))
		}
	}

	for _, resource := range s.Resources {
		err := resource.Validate()
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func (resource *ResourceSchema) Validate() error {
	var result *multierror.Error

	isValidPerm := func(perm string) bool {
		return lo.Contains(resource.Permissions, perm)
	}

	for _, attr := range resource.Attributes {
		// Catch invalid attribute permissions
		for _, permission := range attr.Permissions {
			if !isValidPerm(permission) {
				result = multierror.Append(
					result, fmt.Errorf("%w: %q in attribute %q", ErrInvalidPermission, permission, attr.Name),
				)
			}
		}
	}

	for _, role := range resource.Roles {
		// Catch invalid role permissions
		for _, permission := range role.Permissions {
			if !isValidPerm(permission) {
				result = multierror.Append(result, fmt.Errorf("%w: %q in role %q", ErrInvalidPermission, permission, role.Name))
			}
		}
	}

	for _, permission := range resource.Permissions {
		// Catch unused permissions
		hasAttribute := lo.SomeBy(resource.Attributes, func(attr AttributeSchema) bool {
			return lo.Contains(attr.Permissions, permission)
		})

		hasRole := lo.SomeBy(resource.Roles, func(role RoleSchema) bool {
			return lo.Contains(role.Permissions, permission)
		})

		if !hasAttribute && !hasRole {
			result = multierror.Append(result, fmt.Errorf("%w: %q", ErrUnusedPermission, permission))
		}
	}

	return nil
}

func NewSchemaValidator() *validator.Validate {
	validate := validator.New()

	lo.Must0(validate.RegisterValidation("validName", nameValidator))
	lo.Must0(validate.RegisterValidation("notReserved", notReservedValidator))
	lo.Must0(validate.RegisterValidation("validQualName", qualifierNameValidator))
	lo.Must0(validate.RegisterValidation("validQualPath", qualifierPathValidator))

	return validate
}

func validName(name string) bool {
	// - must begin with a letter, and can have any number of additional letters and numbers.
	// - cannot start with a number.
	// - cannot contain spaces.
	// - cannot contain (very) special characters.
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

	return re.MatchString(name)
}

func qualifierPathValidator(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	re := regexp.MustCompile(`^([\w|.|/]+)$`)

	return re.MatchString(fieldLevel.Field().String())
}

func qualifierNameValidator(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	re := regexp.MustCompile(`^([A-Z][a-zA-Z0-9_-]*)$`)

	return re.MatchString(fieldLevel.Field().String())
}

func notReservedValidator(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return !lo.Contains([]string{"resolver", "authorizer"}, strcase.ToLowerCamel(field.String()))
}

func nameValidator(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return validName(field.String())
}
