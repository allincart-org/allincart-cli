---
trigger: glob
globs: *_test.go
---

- Use testify assert
- Prefer assert.ElementsMatch on lists to ignore ordering issues
- Use t.Setenv for environment variables
- Use t.Context() for Context creation in tests