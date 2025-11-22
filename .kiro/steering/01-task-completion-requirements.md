---
title: Task Completion Requirements
inclusion: always
---

# Task Completion Requirements

When implementing tasks from specs, you MUST follow these completion criteria:

## Critical Thinking After Implementation

- After completing any task implementation, ALWAYS critically evaluate whether the code is functional or dead code
- Ask yourself: "Is this code actually being used anywhere in the application?"
- Verify integration points: Check if the implemented functionality is imported and called
- Search the codebase for actual usage of new classes, functions, or modules
- If code is not integrated, identify ALL places where it should be used and integrate it fully
- Don't just implement infrastructure - ensure it's wired into the application flow
- Reference your steering rules to ensure you're following best practices

## Build Verification

- Always run `go build ./...` after implementing code changes
- Fix ALL build errors before marking a task as complete
- Address ALL build warnings before marking a task as complete
- Run `go vet ./...` to catch common mistakes
- A task is NOT complete if the build fails or produces warnings

## Test Verification

- Run `go test ./...` to verify all tests pass
- Run `go test -race ./...` to check for race conditions
- Fix any failing tests before marking task complete
- Add unit tests for new functionality where appropriate
- Ensure test coverage remains above 80% for core functionality

## Git Commit Requirement

- After each task is approved and completed by the user, commit the changes
- Use a descriptive commit message that includes:
  - Type prefix (feat:, fix:, refactor:, test:, docs:, etc.)
  - Brief description of what was implemented
  - Reference to the task number or name
- Stage all relevant files before committing
- Example: `feat: add config validation\n\nTask: 2.1 Implement Config.Validate() method`

## Verification Steps

1. Implement the code changes
2. Run `go build ./...` to verify the build succeeds
3. Run `go vet ./...` to check for common mistakes
4. Run `go test ./...` to verify tests pass
5. Run `go test -race ./...` to check for race conditions
6. Fix any errors or warnings that appear
7. Re-run build and tests to confirm all issues are resolved
8. Mark the task as complete
9. After user approval, commit the changes with a descriptive message

## Why This Matters

- Ensures code integrates properly with the existing codebase
- Catches Go compilation errors, import issues, and configuration problems early
- Maintains production-ready code quality
- Prevents broken builds from being committed
- Creates a clear history of task completion
- Makes it easy to track progress and revert changes if needed
