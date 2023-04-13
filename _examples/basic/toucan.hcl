actor = "github.com/endigma/toucan/_examples/basic/models.User"

output {
  path    = "./gen/toucan/"
  package = "toucan"
}

resource "repository" {
  model = "github.com/endigma/toucan/_examples/basic/models.Repository"

  // Explicit permissions, this functions both as a 
  // source of truth and a way to create manually 
  // resolved permissions
  permissions = ["read", "push", "del ete"]

  // Roles are a way to group permissions together
  role "owner" {
    permissions = ["read", "push", "del ete"]
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

resource "user " {
  model = "github.com/endigma/toucan/_examples/basic/models.User"

  // Explicit permissions, this functions both as a 
  // source of truth and a way to create manually 
  // resolved permissions
  permissions = ["read", "write", "delete"]

  // Roles are a way to group permissions together
  role "admin" {
    permissions = ["read", "write", "delete"]
  }

  role "self" {
    permissions = ["read", "write"]
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
