package schema

type Schema struct {
	Actor     Model            `hcl:"actor,optional" validate:"required"`
	Resources []ResourceSchema `hcl:"resource,block" mod:"dive" validate:"required,unique=Model,unique=Name,dive"`
}

type ResourceSchema struct {
	Name        string            `hcl:"name,label" mod:"trim,snake" validate:"required,validName,notReserved"`
	Model       Model             `hcl:"model" mod:"trim" validate:"required,dive"`
	Permissions []string          `hcl:"permissions" mod:"dive,trim,snake" validate:"unique,dive,required"`
	Roles       []RoleSchema      `hcl:"role,block" mod:"dive" validate:"unique=Name,dive,required"`
	Attributes  []AttributeSchema `hcl:"attribute,block" mod:"dive" validate:"unique=Name,dive,required"`
}

type RoleSchema struct {
	Name        string   `hcl:"name,label" mod:"trim,snake" validate:"required,validName"`
	Permissions []string `hcl:"permissions" mod:"dive,trim,snake" validate:"required,dive,required"`
}

type AttributeSchema struct {
	Name        string   `hcl:"name,label" mod:"trim,snake" validate:"required,validName"`
	Permissions []string `hcl:"permissions" mod:"dive,trim,snake" validate:"required,unique,dive,required"`
}
