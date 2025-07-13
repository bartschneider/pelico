# Pelico Refactor Recovery & Improvement Plan

## Executive Summary

The previous engineer attempted to modernize Pelico by introducing SvelteKit frontend and Docker containerization but encountered critical build issues that caused the project to be abandoned. This plan outlines a systematic approach to recover and complete the refactor successfully.

## Project Analysis

### Current State Mind Map

```
Pelico System Architecture
├── Backend (Go 1.24)
│   ├── API Server (Gin) 
│   ├── Database Layer (GORM + PostgreSQL)
│   ├── Services
│   │   ├── ROM Scanner
│   │   ├── Metadata (IGDB/TheGamesDB)
│   │   ├── Cache Service
│   │   └── Backup (Nextcloud)
│   └── Models (Game, Platform, PlaySession, etc.)
├── Frontend (Dual State)
│   ├── Legacy (web/) - Bootstrap + Vanilla JS ✅ Working
│   └── New (web_new/) - SvelteKit + TypeScript ❌ Broken
├── Deployment
│   ├── Docker Compose (3 services)
│   ├── Backend Dockerfile ✅ Working
│   └── Frontend Dockerfile ❌ Build Issues
└── Infrastructure
    ├── PostgreSQL Container
    ├── ROM File Storage
    └── Health Checks
```

### Critical Issues Identified

#### 1. Go Module Dependency Issues
- **Problem**: Missing dependencies in go.sum for core packages
- **Error**: `gorm.io/driver/postgres` not in requirements but used in code
- **Root Cause**: Incomplete dependency management during refactor

#### 2. Frontend Build Pipeline Failure
- **Problem**: SvelteKit dependencies not installed (node_modules missing)
- **Error**: All npm packages show "UNMET DEPENDENCY"
- **Root Cause**: package-lock.json exists but node_modules was not committed/built

#### 3. Docker Configuration Mismatch
- **Problem**: Frontend Dockerfile expects built SvelteKit output
- **Issue**: Build fails before Docker stage due to missing dependencies

#### 4. Adapter Misconfiguration
- **Problem**: Using `@sveltejs/adapter-auto` without proper production adapter
- **Issue**: Auto-adapter may not work correctly in containerized environment

## Recovery Strategy

### Phase 1: Immediate Fixes (1-2 hours)

#### 1.1 Fix Go Dependencies
```bash
# Add missing PostgreSQL driver
go get gorm.io/driver/postgres

# Update go.mod to include all required packages
go mod tidy

# Verify build
go build ./cmd/server
```

#### 1.2 Install Frontend Dependencies
```bash
cd web_new
npm install
npm audit fix
```

#### 1.3 Test Local Development
```bash
# Backend
go run cmd/server/main.go

# Frontend (separate terminal)
cd web_new && npm run dev
```

### Phase 2: Frontend Architecture Improvements (2-3 hours)

#### 2.1 Fix SvelteKit Configuration
- Replace `@sveltejs/adapter-auto` with `@sveltejs/adapter-node`
- Configure proper API proxy for development
- Add environment variable handling

#### 2.2 Improve Component Architecture
- Standardize component props and events
- Add proper TypeScript interfaces
- Implement consistent error handling
- Add loading states and user feedback

#### 2.3 Backend API Integration
- Create API client service with proper error handling
- Add response type definitions
- Implement caching for repeated requests

### Phase 3: Docker & Deployment Fixes (1-2 hours)

#### 3.1 Fix Frontend Dockerfile
```dockerfile
# Multi-stage build with proper node adapter
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY package.json ./
CMD ["node", "build"]
```

#### 3.2 Update Docker Compose
- Add build arguments for environment
- Configure proper networking between services
- Add volume mounts for development

### Phase 4: Feature Completion & Enhancement (3-4 hours)

#### 4.1 Complete Unfinished Features
- Finalize wishlist and shortlist functionality
- Implement proper session management
- Add comprehensive error handling

#### 4.2 UX Improvements
- Add responsive design improvements
- Implement proper loading states
- Add form validation feedback
- Improve navigation and user flow

#### 4.3 Performance Optimization
- Implement component lazy loading
- Add image optimization
- Optimize API calls with proper caching

### Phase 5: Testing & Quality Assurance (2-3 hours)

#### 5.1 Backend Testing
```bash
go test ./...
go test -race ./...
```

#### 5.2 Frontend Testing
```bash
npm run test
npm run lint
npm run type-check
```

#### 5.3 Integration Testing
- Test API endpoints with frontend
- Verify Docker build and deployment
- Test database migrations and data persistence

## Implementation Examples

### 1. Fix Go Dependencies
```go
// go.mod additions needed
module pelico

require (
    // ... existing deps
    gorm.io/driver/postgres v1.5.4  // Add this
)
```

### 2. SvelteKit Adapter Fix
```javascript
// svelte.config.js
import adapter from '@sveltejs/adapter-node';

const config = {
  kit: {
    adapter: adapter({
      out: 'build',
      precompress: false,
      envPrefix: '',
    }),
  }
};
```

### 3. API Client Service
```typescript
// src/lib/api.ts
export class ApiClient {
  private baseUrl = '/api/v1';
  
  async get<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`);
    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }
    return response.json();
  }
  
  async post<T>(endpoint: string, data: unknown): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }
    return response.json();
  }
}
```

### 4. Component Error Handling
```svelte
<!-- GameCard.svelte improvement -->
<script lang="ts">
  import type { Game } from '$lib/models';
  
  export let game: Game;
  export let loading = false;
  export let onEdit: (game: Game) => void;
  
  let imageError = false;
  
  function handleImageError() {
    imageError = true;
  }
</script>

{#if loading}
  <div class="card game-card loading">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
{:else}
  <div class="card game-card">
    <!-- Content with proper error states -->
  </div>
{/if}
```

## Risk Assessment & Mitigation

### High Risk Issues
1. **Database Migration**: Ensure backward compatibility
   - **Mitigation**: Test with database backups
2. **API Breaking Changes**: Frontend expects specific API format
   - **Mitigation**: Maintain API contract, add versioning if needed

### Medium Risk Issues
1. **Docker Build Performance**: Multi-stage builds can be slow
   - **Mitigation**: Implement proper layer caching
2. **Frontend Bundle Size**: SvelteKit + Bootstrap might be large
   - **Mitigation**: Implement code splitting and tree shaking

## Success Criteria

### Phase 1 Success
- [x] Go build completes without errors
- [x] npm install completes successfully
- [x] Local development servers start

### Phase 2 Success
- [ ] Frontend builds and runs in production mode
- [ ] API integration works correctly
- [ ] All existing features functional in new UI

### Phase 3 Success
- [ ] Docker compose up completes successfully
- [ ] All containers healthy and communicating
- [ ] Application accessible via web browser

### Phase 4 Success
- [ ] All planned features implemented
- [ ] Responsive design works on mobile/desktop
- [ ] Performance meets acceptable standards

### Phase 5 Success
- [ ] All tests pass
- [ ] No critical security vulnerabilities
- [ ] Documentation updated

## Timeline Estimate

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Phase 1 | 1-2 hours | None |
| Phase 2 | 2-3 hours | Phase 1 complete |
| Phase 3 | 1-2 hours | Phase 1-2 complete |
| Phase 4 | 3-4 hours | Phase 1-3 complete |
| Phase 5 | 2-3 hours | All phases complete |
| **Total** | **9-14 hours** | Linear execution |

## Deployment Strategy

### Development
1. Local development with separate frontend/backend servers
2. Use docker-compose for full stack testing
3. Hot reload enabled for both frontend and backend

### Production
1. Multi-stage Docker builds for optimization
2. Reverse proxy (nginx) for serving static assets
3. Health checks and monitoring
4. Database backups and migration scripts

## Conclusion

The refactor is salvageable with systematic fixing of dependency issues and completion of the SvelteKit implementation. The previous engineer made good architectural decisions but left critical configuration and dependency issues unresolved. Following this plan should result in a modern, maintainable, and fully functional video game collection manager.