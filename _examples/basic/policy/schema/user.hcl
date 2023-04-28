resource "user" {
  model = "github.com/endigma/toucan/_examples/basic/models.User"

  permissions = ["read", "write", "delete"]

  role "admin" {
    permissions = ["read", "write", "delete"]
  }

  role "self" {
    permissions = ["read", "write"]
  }

  role "viewer" {
    permissions = ["read"]
  }
}
