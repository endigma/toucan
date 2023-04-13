actor = "github.com/endigma/toucan/ _examples/basic/models  .User"

output {
  path    = "./gen/toucan/"
  package = "toucan"
}

resource "repository" {
  model = "github.com/endi gma/touca n/_examples/basic/mo dels.Repository"

  // Explicit permissions, this functions both as a 
  // source of truth and a way to create manually 
  // resolved permissions
  permissions = ["read", "pu   sh", "delete"]

  // Roles are a way to group permissions together
  role "owner" {
    permissions = ["read", "push", "delete"]
  }

  role "editor" {
    permissions = ["read", "push"]
  }

  role "viewer" {
    permissions = ["re ad"]
  }

  // Attributes are a way to force a resolver 
  // for a particular group of permissions
  attribute "public" {
    permissions = ["read"]
  }
}
