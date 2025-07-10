# Pelico Homelab Deployment Guide

## 🔐 Security Notice

**NEVER commit credentials to Git!** This repository is configured to protect sensitive information.

## Prerequisites

1. **Docker and Docker Compose** installed on homelab server
2. **sshpass** for automated deployment: `brew install hudochenkov/sshpass/sshpass` (macOS)
3. **Environment variables** set for deployment

## Environment Setup

### Required Environment Variables

Set these in your shell or a `.env.local` file (which is git-ignored):

```bash
export SERVER_PASSWORD="your_server_password"
export TWITCH_CLIENT_ID="your_twitch_client_id"
export TWITCH_CLIENT_SECRET="your_twitch_client_secret"
export NEXTCLOUD_URL="https://your-nextcloud-instance.com"
export NEXTCLOUD_USER="your_nextcloud_username"
export NEXTCLOUD_PASS="your_nextcloud_app_password"
```

### Homelab Server Configuration

Create a `.env` file on your homelab server at `/home/bartosz/pelico/.env`:

```bash
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=pelico
DB_PASSWORD=secure_db_password
DB_NAME=pelico

# IGDB/Twitch API (for game metadata)
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret

# Nextcloud Backup Integration
NEXTCLOUD_URL=https://your-nextcloud-instance.com
NEXTCLOUD_USER=your_nextcloud_username
NEXTCLOUD_PASS=your_nextcloud_app_password

# Application Configuration
PORT=8081
LOG_LEVEL=info
```

## Deployment Commands

### Full Deployment
```bash
# Build and deploy to homelab
make deploy
```

### Check Status
```bash
# View container status
make homelab-status

# View application logs
make homelab-logs
```

### Manual Steps (if needed)

1. **Build Docker image:**
   ```bash
   make docker-build
   ```

2. **Transfer to homelab:**
   ```bash
   docker save pelico:latest | ssh bartosz@192.168.1.52 'docker load'
   ```

3. **Deploy on server:**
   ```bash
   ssh bartosz@192.168.1.52 'cd pelico && docker compose up -d'
   ```

## Application URLs

After deployment:
- **Pelico Web Interface**: http://192.168.1.52:8081
- **Health Check**: http://192.168.1.52:8081/api/v1/health

## Troubleshooting

### Common Issues

1. **Permission denied (publickey)**
   - Ensure SSH key is set up: `ssh-copy-id bartosz@192.168.1.52`

2. **Container won't start**
   - Check logs: `make homelab-logs`
   - Verify .env file exists on server

3. **Database connection error**
   - Ensure PostgreSQL container is running
   - Check database credentials in .env

4. **Backup/Nextcloud integration not working**
   - Verify Nextcloud credentials and URL
   - Test WebDAV connection manually

### Health Checks

```bash
# Test application health
curl http://192.168.1.52:8081/api/v1/health

# Test cache status  
curl http://192.168.1.52:8081/api/v1/cache/stats
```

## Post-Deployment Tasks

1. **Verify all services are running:**
   ```bash
   make homelab-status
   ```

2. **Test core functionality:**
   - Access web interface
   - Add a game manually
   - Test ROM scanning
   - Verify backup functionality

3. **Set up automated backups** (optional):
   - Configure cron job for daily backups
   - Test Nextcloud integration

## Security Best Practices

- ✅ All credentials stored in environment variables
- ✅ .env files excluded from Git
- ✅ No hardcoded passwords in codebase
- ✅ Database isolated in Docker network
- ✅ Application runs on internal homelab network only

## Updates

To update the application:

1. Pull latest changes: `git pull origin main`
2. Redeploy: `make deploy`

The deployment process automatically handles:
- Building new Docker image
- Stopping old containers
- Starting updated containers
- Preserving database and configuration