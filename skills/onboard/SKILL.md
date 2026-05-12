---
name: onboard
description: Summarize the stack, entry points, and architecture of a codebase for quick orientation
---

**Discovery Phase** (run these in parallel):
1. **Tech Stack Discovery**: Identify primary language (check extensions, go.mod, package.json, Cargo.toml, pom.xml), framework, and package manager.
2. **Project Structure Overview**: Map directory layout (src/, pkg/, internal/, apps/, services/, etc.) to understand code organization.
3. **Version Gathering**: Note language version (go.mod, package.json engines, .ruby-version, .python-version) and framework versions.
4. **Infrastructure Check**: Look for Dockerfiles, docker-compose.yml, Terraform files, Kubernetes manifests, or CI/CD YAMLs (.github/, .gitlab-ci.yml).
5. **Configuration Check**: Look for .env.example, config.*, settings.*, or config directories to understand runtime requirements.
6. **Entry Point Hunt**: Locate main entry point (main.py, index.ts, main.go, app.rb, Program.cs). Identify dependencies, env vars, and ports.
7. **Dependencies Analysis**: List main third-party libraries/frameworks from lockfiles or imports.
8. **Testing & Build Discovery**: Check for Makefile, test configs (jest.config.js, pytest.ini, vitest.config.ts), and build scripts.
9. **API & Documentation Check**: Look for OpenAPI/Swagger specs, route files, or API documentation.
10. **Security & Infrastructure**: Check .gitignore, security policies, dependency scanning configs (snyk, dependabot), and secrets handling.

**Documentation Phase**:
11. **README Analysis**: Read README and summarize "Getting Started" steps, including how to run tests.
12. **Report**: Write summary to project_map.md with these sections:
   - Stack & Versions
   - Project Structure
   - Infrastructure & Deployment
   - Configuration
   - Entry Points
   - Dependencies
   - Development (setup, tests, build)
   - API Documentation (if applicable)
   - Security Notes

**Error Handling**:
- If Dockerfile missing → note deployment may be serverless/PAAS
- If lockfile missing → note dependencies may not be reproducible
- If no README → mark "Getting Started" as "not documented"
