---
trigger: always_on
description: 
globs: 
---
# General rules

- Use convential commit messages
- Use always Go 1.24 syntax, use slices package for slice containing
- Inline error into if condition like this:

```go
if err := foo(); err != nil {
    // something
}
```

