resource "repository" {
  model = "github.com/endigma/toucan/_examples/basic/models.Repository"

  // Explicit permissions, this functions both as a 
  // source of truth and a way to create manually 
  // resolved permissions
  permissions = ["read", "push", "delete", "snake_case"]

  // Roles are a way to group permissions together
  role "owner" {
    permissions = ["read", "push", "delete", "snake_case"]
  }

  role "editor" {
    permissions = ["read", "push"]
  }

  role "viewer" {
    permissions = ["read"]
  }

  // Attributes are a way to force a resolver 
  // for a particular group of permissions
  attribute "public" {
    permissions = ["read"]
  }
}
