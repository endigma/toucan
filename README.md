# toucan

toucan is a heavily-alpha authorization library for Go. It makes use of code generation to help you build authorization for your Go app.

- **toucan is fully type safe** -- type assertions are only used for the conveinence methods, i.e. `Authorize()`, there are fully typed methods like `AuthorizeResource` and all resolvers are fully typed.
- **toucan meets you where you're at** -- with full RBAC/ABAC support and planned ReBAC helpers, `toucan` can fully model your app's authorization.
- **toucan uses code generation to cut down on boilerplate** -- you only implement the resolvers for facts toucan can't infer, all the solver logic is handled automatically.
- **toucan only resolves what it needs to** -- toucan only evaluates the relevant resolvers, instead of building a full authorization context.

## Quick start

**TODO**

Look in `_examples` for a basic demo of how it can be used right now.

## Roadmap

- [ ] CLI interface for basic usage
- [ ] Automatic resolver generation
- [ ] Escape-hatch for custom resolvers
- [ ] Generator output test coverage
- [ ] Multi-file configs / Globbing
- [ ] Memoization?
- [ ] Magic / Automations
  - [ ] Actor proxies
  - [ ] Permissions implication
  - [ ] Role implication
  - [ ] Easy testing contexts

## Policy Reference

`toucan` is configured through an HCL file describing the authorization setup for your app. A basic config looks like this:

```hcl
actor = "github.com/endigma/toucan/_examples/basic/models.User"

output {
  path    = "./gen/policy/"
  package = "policy"
}

resource "repository" {
  model = "github.com/endigma/toucan/_examples/basic/models.Repository"

  // Explicit permissions
  permissions = ["read", "push", "delete"]

  // Roles are a way to group permissions together
  role "owner" {
    permissions = ["read", "push", "delete"]
  }

  role "editor" {
    permissions = ["read", "push"]
  }

  role "viewer" {
    permissions = ["read"]
  }

  // Attributes are a way to force a resolver
  // for a particular group of permissions
  attribute " public " {
    permissions = [" read "]
  }
}
```

### Important

The config has some transformations applied to make it somewhat mistake-tolerant.

- All permissions, roles and resources are transformed to `snake_case`, then re-transformed to other cases as necessary.

  - This means your permissions may not be how you expect them to be: `updateUsername`/`update username`/etc will become `update_username`. The code names of your permissions for string usage are stored as consts in the generated package.
  - It is suggested to author in `snake_case` to prevent unexpected behavior.

- All names are trimmed for whitespace, so `foo   ` will become `foo`.

## Why `toucan`?

I've drawn up a comparison table for the different authorization options I'm aware of in Go. If you see something that's wrong, please open an issue or PR!

|                        | [toucan](https://github.com/endigma/toucan) | [casbin](https://casbin.org/) | [oso](https://www.osohq.com/)      |
| ---------------------- | ------------------------------------------- | ----------------------------- | ---------------------------------- |
| Kind                   | resolver-based                              | fact-based                    | fact-based                         |
| Language               | Go                                          | Go (Polyglot)                 | Rust (FFI)                         |
| Implementation         | codegen                                     | runtime                       | runtime (FFI/CGO)                  |
| Type Safety            | fully                                       | none (just strings)           | none (FFI magic)                   |
| RBAC                   | 👍                                          | 👍                            | 👍                                 |
| ABAC                   | 👍                                          | 👍                            | 👍                                 |
| ReBAC                  | 👍 (manual) / 🚧 (magic)                    | 👍                            | 👍                                 |
| Embedded / Local first | 👍                                          | 👍 (requires fact store)      | 👍 (oso-library) / ⛔️ (oso-cloud) |
| Custom Roles           | 🚧                                          | ⛔️                           | 👍                                 |
| Field Level Auth       | 🚧                                          | ❓                            | 👍                                 |
| Testable               | 👍                                          | 👍                            | 👍                                 |

### "resolver-based"

`toucan` doesn't store anything about your authorization context anywhere, it just asks you for the facts it needs to make a decision. This means you can use any data store or structure for your role associations or attributes.

### "fact-based"

These libraries use a "fact store" that holds associations like:

(casbin)

```
p, alice, data1, read
p, bob, data2, write
```

(oso)

```
has_role(user: User, "admin")
has_role(user: User, "editor")

has_permission(user: User, "read", repo: Repository)
has_permission(user: User, "push", repo: Repository)

is_public(repo: Repository)
```

The library then uses a rules engine to evaluate the facts and make a decision. This can make lookup times faster, as it is able to optimize the query for this specific data store, but it also means you have to update the facts when your data changes, creating the potential for divergence from your app's state.
