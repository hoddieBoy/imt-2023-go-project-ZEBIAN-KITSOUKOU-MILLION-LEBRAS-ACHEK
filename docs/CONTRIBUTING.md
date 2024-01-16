# Project Contribution Guide

## Commit Convention

This project follows the commit convention to standardize commit messages. We adhere to the [Conventional Commit specification](https://www.conventionalcommits.org) to ensure clear and consistent semantics in our commit messages.

### Commit Types

- **feat**: for a new feature
- **fix**: for a bug fix
- **docs**: for documentation changes
- **style**: for changes that do not affect the code's meaning (formatting, white spaces, etc.)
- **refactor**: for code changes that neither add a feature nor fix a bug
- **test**: for adding or modifying tests
- **chore**: for maintenance tasks or other tasks not related to the source code

### Commit Message Structure

Each commit message should adhere to the following structure:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples of Valid Commits

- `feat(user-auth): add user authentication feature`
- `fix(database): fix database connection issue`
- `docs(readme): update project documentation`

### Full Commit Messages

```
feat(user-auth): add user authentication feature

  This feature allows users to log in using their email and password.

```

## Branch Naming

To maintain a clear and consistent branch structure, please follow the conventions below when creating branches:

- `feature/<feature-name>` for ongoing feature development
- `fix/<bug-name>` for bug fix branches
- `docs/<documentation-name>` for documentation changes
- `chore/<task-name>` for maintenance tasks or other tasks not related to the source code

### Example Branch Names

- `feature/user-auth`
- `fix/database-connection`
- `docs/readme-updates`
- `chore/cleanup-codebase`

## Submitting a Pull Request

Before submitting a Pull Request, ensure that your commits follow the Conventional Commit convention, and your branch name adheres to the branch naming convention.

This facilitates understanding of the changes made and expedites the review process.