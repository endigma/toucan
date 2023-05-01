global {
  permissions = ["read_all_users", "write_all_users", "read_all_profiles"]

  role "admin" {
    permissions = ["write_all_users", "read_all_users"]
  }

  attribute "profiles_are_public" {
    permissions = ["read_all_profiles"]
  }
}
