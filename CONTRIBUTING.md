# Contributing Guidelines

Thank you for contributing to this project! To keep the codebase organized, please follow these guidelines when making changes.

## General Rules

1. **Do not** push code directly to the `main` branch. All changes must go through a pull request (PR).
2. Create a **new branch** from `main` for each change. The branch name should be descriptive, e.g., `feature/add-auth` or `fix/api-bug`.
3. The **merge commit** message must start with the issue number from Yandex Tracker, for example:
   ```
   PLAN-1: Start development
   ```
4. PRs must be reviewed and approved by at least one reviewer before merging into `main`.

## How to Contribute

1. **Create a new branch**
   ```sh
   git checkout main
   git pull origin main
   git checkout -b feature/your-feature-name
   ```

2. **Make changes and commit**
   ```sh
   git add .
   git commit -m "PLAN-42: Implement new feature"
   ```

3. **Push changes to the remote repository**
   ```sh
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request (PR)** on GitHub and request a review.

5. **After approval**, merge **only via PR**, using the commit message format `PLAN-XX: description`.

## Code Style
- Follow the project's coding standards.
- Ensure your changes are clear and well-structured.

Thank you for your contribution!
