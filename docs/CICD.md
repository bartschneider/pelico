# Pelico CI/CD Pipeline Documentation

## Overview

Pelico uses an automated continuous deployment pipeline that triggers on git commits and pushes. The pipeline automatically manages versioning, runs tests, builds Docker images, and deploys to your homelab.

## Prerequisites

1. **Environment Configuration**: Ensure `.env` file contains all required variables:
   ```bash
   HOMELAB_USER="your-username"
   HOMELAB_IP="192.168.1.52"
   HOMELAB_SSH_PORT="22"
   HOMELAB_SSH_PASSWORD="your-password"
   DEPLOY_TO_PORT="8081"
   ```

2. **SSH Access**: Password-based SSH access to your homelab server
3. **Docker**: Docker must be installed on both development and homelab machines
4. **Git Hooks**: Git hooks are automatically installed when you clone the repository

## Pipeline Components

### 1. Automatic Version Management

- **Trigger**: Every commit to the main branch
- **Hook**: `.git/hooks/post-commit`
- **Script**: `scripts/version-bump.sh`
- **Behavior**: 
  - Commits with `[major]` increment major version (1.0.0 → 2.0.0)
  - Commits with `[minor]` increment minor version (1.0.0 → 1.1.0)
  - All other commits increment patch version (1.0.0 → 1.0.1)
  - Creates git tags automatically

### 2. Pre-Push Validation

- **Trigger**: Before pushing to remote repository
- **Hook**: `.git/hooks/pre-push`
- **Behavior**:
  - Runs Go tests to ensure code quality
  - Prompts for deployment confirmation
  - Aborts push if tests fail

### 3. Continuous Deployment

- **Trigger**: After successful push to main branch
- **Script**: `scripts/deploy-continuous.sh`
- **Process**:
  1. Runs comprehensive test suite
  2. Builds Docker image with version injection
  3. Saves current deployment as rollback point
  4. Transfers image to homelab
  5. Deploys using docker-compose
  6. Verifies deployment health
  7. Performs automatic rollback on failure

### 4. Deployment Verification

- **Script**: `scripts/verify-deployment.sh`
- **Checks**:
  - Container running status
  - Container health check
  - Version correctness
  - API endpoint availability
  - Database connectivity
  - Frontend accessibility
- **Rollback**: Automatic rollback on critical failures

## Usage Guide

### Basic Workflow

1. **Make Changes**: Edit your code as needed
   ```bash
   # Edit files
   vim internal/handlers/game_handler.go
   ```

2. **Commit Changes**: Version bump happens automatically
   ```bash
   git add .
   git commit -m "feat: Add game search functionality"
   # Version automatically bumps to 1.2.3
   ```

3. **Push to Deploy**: Tests run and deployment triggers
   ```bash
   git push origin main
   # Tests run automatically
   # Prompts for deployment confirmation
   # Deployment starts in background
   ```

4. **Monitor Deployment**: Check deployment progress
   ```bash
   tail -f deployment.log
   # Or check status on homelab
   make homelab-status
   ```

### Version Control

Control version bumping with commit message tags:

```bash
# Major version bump (1.2.3 → 2.0.0)
git commit -m "feat: [major] Complete API redesign"

# Minor version bump (1.2.3 → 1.3.0)
git commit -m "feat: [minor] Add multiplayer support"

# Patch version bump (1.2.3 → 1.2.4) - default
git commit -m "fix: Correct game loading issue"
```

### Manual Operations

#### Deploy Without Push
```bash
bash scripts/deploy-continuous.sh
```

#### Verify Deployment
```bash
bash scripts/verify-deployment.sh
# With automatic rollback
bash scripts/verify-deployment.sh --auto-rollback
```

#### Check Version
```bash
curl http://192.168.1.52:8081/api/v1/version
```

#### View Logs
```bash
make homelab-logs
```

## Troubleshooting

### Deployment Fails

1. Check deployment log:
   ```bash
   cat deployment.log
   ```

2. Verify environment variables:
   ```bash
   grep HOMELAB .env
   ```

3. Test SSH connection:
   ```bash
   sshpass -p $HOMELAB_SSH_PASSWORD ssh -p $HOMELAB_SSH_PORT $HOMELAB_USER@$HOMELAB_IP
   ```

### Version Not Updating

1. Ensure you're on main branch:
   ```bash
   git branch --show-current
   ```

2. Check if hooks are executable:
   ```bash
   ls -la .git/hooks/
   ```

3. Run version bump manually:
   ```bash
   bash scripts/version-bump.sh
   ```

### Tests Failing

1. Run tests locally:
   ```bash
   go test ./... -v
   ```

2. Skip deployment on push:
   ```bash
   # Answer 'N' when prompted for deployment
   git push origin main
   ```

### Rollback Required

1. Automatic rollback (if verification fails):
   ```bash
   bash scripts/verify-deployment.sh --auto-rollback
   ```

2. Manual rollback:
   ```bash
   ssh $HOMELAB_USER@$HOMELAB_IP
   cd pelico
   docker tag pelico:rollback pelico:latest
   docker-compose down
   docker-compose up -d
   ```

## Security Considerations

1. **Credentials**: Store sensitive data in `.env` file (git-ignored)
2. **SSH Keys**: Consider using SSH keys instead of passwords for production
3. **Network**: Ensure homelab is accessible only from trusted networks
4. **Backups**: Regular backups are created before each deployment

## Extending the Pipeline

### Add Custom Checks

Edit `scripts/verify-deployment.sh` to add custom verification:

```bash
# Add new check function
check_custom_metric() {
    info "Checking custom metric..."
    # Your check logic here
}

# Call in main()
check_custom_metric
```

### Modify Version Strategy

Edit `scripts/version-bump.sh` to change versioning logic:

```bash
# Add custom version tags
if echo "$COMMIT_MSG" | grep -q "\[hotfix\]"; then
    # Custom hotfix versioning
fi
```

### Add Deployment Environments

Extend `scripts/deploy-continuous.sh` for multiple environments:

```bash
# Add environment parameter
ENVIRONMENT=${1:-"homelab"}

case $ENVIRONMENT in
    "staging")
        HOMELAB_IP=$STAGING_IP
        ;;
    "production")
        HOMELAB_IP=$PROD_IP
        ;;
esac
```

## Best Practices

1. **Commit Messages**: Use conventional commits for clear version bumping
2. **Testing**: Always ensure tests pass before pushing
3. **Monitoring**: Check deployment logs after each push
4. **Rollback Plan**: Keep previous versions available for quick rollback
5. **Documentation**: Update this guide when modifying the pipeline

## Pipeline Scripts Reference

- `scripts/version-bump.sh` - Manages semantic versioning
- `scripts/deploy-continuous.sh` - Main deployment orchestrator
- `scripts/verify-deployment.sh` - Health checks and rollback
- `scripts/test-cicd-pipeline.sh` - Pipeline testing utility
- `.git/hooks/post-commit` - Triggers version bumping
- `.git/hooks/pre-push` - Runs tests before push
- `.git/hooks/post-merge` - Triggers deployment after push