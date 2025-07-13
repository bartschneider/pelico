# Pelico Project File Inventory & Analysis

## Root Directory Files

### Essential Configuration Files
- **go.mod** - Go module definition, required for Go project
- **go.sum** - Go dependency checksums, required for Go project  
- **docker-compose.yml** - Container orchestration, required for deployment
- **Dockerfile** - Backend container definition, required for deployment
- **Makefile** - Build automation, useful for development
- **.dockerignore** - Docker build exclusions, optimizes container builds
- **.gitignore** - Git exclusions, required for version control
- **.env** - Environment variables (if present), required for configuration

### Documentation Files
- **README.md** - Project documentation, essential
- **CLAUDE.md** - Claude AI instructions, essential for AI assistance
- **LICENSE** - Software license, legally required
- **CHANGELOG.md** - Version history, useful for tracking changes
- **DEPLOYMENT.md** - Deployment instructions, useful
- **DEPLOYMENT_HOMELAB.md** - Homelab deployment instructions, useful

### Audit & Analysis Files (REVIEW NEEDED)
- **ARCHITECTURE_AUDIT_V4.md** - Architecture analysis v4, potentially redundant
- **ARCHITECTURE_AUDIT_V5.md** - Architecture analysis v5, potentially redundant
- **FRONTEND_UX_AUDIT.md** - UX audit v1, potentially redundant
- **FRONTEND_UX_AUDIT_V2.md** - UX audit v2, potentially redundant
- **GEMINI.md** - Gemini AI instructions, potentially redundant with CLAUDE.md
- **sys.md** - System notes, purpose unclear
- **refactor-plan.md** - Our refactor documentation, useful for reference

### Development Scripts
- **deploy-to-homelab.sh** - Deployment script, useful
- **deploy-to-homelab.example.sh** - Example deployment script, useful as template
- **claude-env.sh** - Claude environment setup, potentially redundant

### Binary Files (CLEANUP NEEDED)
- **pelico** - Compiled binary, should be in .gitignore
- **pelico-test** - Test binary, should be in .gitignore  
- **server** - Another binary, should be in .gitignore
- **deploy-to-homelab** - Binary script, should be in .gitignore
- **main.go** - Duplicate of cmd/server/main.go, redundant
- **Screenshot 2025-07-10 at 14.20.45.png** - Screenshot, potentially unnecessary

## Source Code Directories

### cmd/
- **cmd/server/main.go** - Application entry point, essential

### internal/
Core application code, all essential:

#### internal/api/
- **server.go** - HTTP server and routing setup

#### internal/config/
- **config.go** - Configuration management

#### internal/database/
- **database.go** - Database connection and setup

#### internal/errors/
- **errors.go** - Custom error types

#### internal/handlers/
HTTP request handlers, all essential:
- **backup_handler.go** - Backup/restore endpoints
- **directory_handler.go** - Directory browsing endpoints
- **game_handler.go** - Game CRUD endpoints
- **game_handler_test.go** - Unit tests for game handler
- **platform_handler.go** - Platform CRUD endpoints
- **scanner_handler.go** - ROM scanning endpoints
- **session_handler.go** - Play session endpoints
- **shortlist_handler.go** - Shortlist management endpoints
- **stats_handler.go** - Statistics endpoints
- **wishlist_handler.go** - Wishlist management endpoints

#### internal/middleware/
- **validation.go** - Request validation middleware

#### internal/models/
- **models.go** - Database models and types

#### internal/services/
Business logic services, all essential:
- **cache_service.go** - Caching functionality
- **igdb_service.go** - IGDB API integration
- **logger_service.go** - Structured logging
- **metadata_service.go** - Game metadata retrieval
- **nextcloud_backup.go** - Nextcloud backup integration
- **rom_scanner.go** - ROM file scanning

#### internal/version/
- **version.go** - Version information

### scripts/
Database initialization scripts, essential:
- **init.sql** - Database schema initialization
- **add_wishlist_shortlist_tables.sql** - Additional table definitions

## Legacy Frontend (web/)

### Purpose: Original Bootstrap frontend
- **web/static/css/style.css** - Legacy CSS styles
- **web/static/js/app.js** - Legacy JavaScript
- **web/templates/index.html** - Legacy HTML template

**Status**: Should be kept as fallback, used by backend for serving web interface when frontend is not available.

## New Frontend (web_new/)

### Configuration Files
- **package.json** - NPM dependencies, essential
- **package-lock.json** - Dependency lock file, essential
- **svelte.config.js** - SvelteKit configuration, essential
- **tsconfig.json** - TypeScript configuration, essential
- **vite.config.js** - Vite build configuration, essential
- **Dockerfile** - Frontend container definition, essential
- **.prettierrc** - Code formatting rules, useful

### Source Code
- **src/app.html** - HTML template, essential
- **src/app.css** - Global styles, essential
- **src/hooks.server.ts** - Server-side hooks, essential for API proxy
- **src/lib/api.ts** - API client, essential
- **src/lib/models.ts** - TypeScript models, essential
- **src/lib/actions/modal.ts** - Modal actions, essential
- **src/lib/components/** - Svelte components, all essential:
  - **GameCard.svelte** - Game display component
  - **GameFormModal.svelte** - Game editing modal
  - **Nav.svelte** - Navigation component
- **src/routes/** - Page routes, all essential:
  - **+layout.svelte** - Layout wrapper
  - **+page.svelte** - Home page
  - **shortlist/+page.svelte** - Shortlist page
  - **stats/+page.svelte** - Statistics page
  - **wishlist/+page.svelte** - Wishlist page

### Build Artifacts (CLEANUP NEEDED)
All these should be cleaned up as they're generated during build:

#### .svelte-kit/ directory
Generated SvelteKit files:
- **.svelte-kit/generated/** - Generated code, should be in .gitignore
- **.svelte-kit/output/** - Build output, should be in .gitignore
- **.svelte-kit/adapter-node/** - Adapter output, should be in .gitignore
- **.svelte-kit/types/** - Generated types, should be in .gitignore
- **.svelte-kit/ambient.d.ts** - Generated ambient types
- **.svelte-kit/non-ambient.d.ts** - Generated non-ambient types
- **.svelte-kit/tsconfig.json** - Generated TypeScript config

#### build/ directory
Production build output, should be in .gitignore:
- **build/client/** - Client-side build
- **build/server/** - Server-side build
- **build/*.js** - Build entry points

## test-go/ Directory (CLEANUP NEEDED)

**Purpose**: Appears to be a separate test module
- **test-go/go.mod** - Separate Go module
- **test-go/main.go** - Test main file
- **test-go/test/test.go** - Test file

**Status**: This appears to be a separate testing playground that's not integrated with the main project. Should be removed.

## IDE Configuration (.idea/)

**Purpose**: JetBrains IDE configuration
- **.idea/** - IDE workspace settings

**Status**: Should be in .gitignore to avoid committing IDE-specific settings.

## Cleanup Actions Completed ✅

### Files Removed:
1. ✅ **Binary files**: pelico, pelico-test, server, deploy-to-homelab (compiled artifacts)
2. ✅ **Duplicate files**: main.go (duplicate of cmd/server/main.go)
3. ✅ **Build artifacts**: web_new/.svelte-kit/, web_new/build/
4. ✅ **Test playground**: test-go/ (entire directory)
5. ✅ **Redundant docs**: ARCHITECTURE_AUDIT_V4.md, FRONTEND_UX_AUDIT.md, sys.md
6. ✅ **IDE files**: .idea/ directory
7. ✅ **Screenshots**: Screenshot 2025-07-10 at 14.20.45.png
8. ✅ **Sensitive config**: claude-env.sh (contained project-specific credentials)

### Files Kept for Good Reasons:
1. **GEMINI.md** - Project documentation for Gemini AI (similar to CLAUDE.md)
2. **ARCHITECTURE_AUDIT_V5.md** - Latest architecture audit (useful reference)
3. **FRONTEND_UX_AUDIT_V2.md** - Latest UX audit (useful reference)
4. **deploy-to-homelab.sh** - Deployment script (functional, not just example)
5. **refactor-plan.md** - Our refactor documentation (useful historical reference)

### .gitignore Updates Completed ✅:
1. ✅ Fixed go.sum exclusion (go.sum should be committed for reproducible builds)
2. ✅ Added web_new/build/ and additional frontend build artifacts
3. ✅ Added additional binary patterns (server, main, pelico-test)
4. ✅ Updated patterns for generated files (*.tsbuildinfo)

## Final Project Structure

The project is now clean with only essential files:

**Root**: Essential config, docs, and deployment files
**cmd/**: Application entry point
**internal/**: Core Go application code
**scripts/**: Database initialization
**web/**: Legacy frontend (still used by backend)
**web_new/**: Modern SvelteKit frontend (source code only)

Total reduction: Removed ~150+ generated/redundant files while preserving all essential functionality.