#!/bin/bash

# Version Bump Script for Pelico
# Automatically increments version based on commit message patterns

set -e

VERSION_FILE="internal/version/version.go"

# Get current version
CURRENT_MAJOR=$(grep "Major = " $VERSION_FILE | sed 's/.*Major = //')
CURRENT_MINOR=$(grep "Minor = " $VERSION_FILE | sed 's/.*Minor = //')
CURRENT_PATCH=$(grep "Patch = " $VERSION_FILE | sed 's/.*Patch = //')

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%B)

# Determine version bump type
if echo "$COMMIT_MSG" | grep -q "\[major\]"; then
    NEW_MAJOR=$((CURRENT_MAJOR + 1))
    NEW_MINOR=0
    NEW_PATCH=0
    BUMP_TYPE="major"
elif echo "$COMMIT_MSG" | grep -q "\[minor\]"; then
    NEW_MAJOR=$CURRENT_MAJOR
    NEW_MINOR=$((CURRENT_MINOR + 1))
    NEW_PATCH=0
    BUMP_TYPE="minor"
else
    NEW_MAJOR=$CURRENT_MAJOR
    NEW_MINOR=$CURRENT_MINOR
    NEW_PATCH=$((CURRENT_PATCH + 1))
    BUMP_TYPE="patch"
fi

# Update version file
cat > $VERSION_FILE << EOF
package version

import (
	"fmt"
	"time"
)

const (
	Major = $NEW_MAJOR
	Minor = $NEW_MINOR
	Patch = $NEW_PATCH
)

var (
	// BuildTime will be set during build
	BuildTime = ""
	// GitCommit will be set during build  
	GitCommit = ""
)

type Info struct {
	Version   string    \`json:"version"\`
	BuildTime string    \`json:"build_time"\`
	GitCommit string    \`json:"git_commit"\`
	Timestamp time.Time \`json:"timestamp"\`
}

func GetVersion() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}

func GetInfo() Info {
	return Info{
		Version:   GetVersion(),
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		Timestamp: time.Now(),
	}
}
EOF

OLD_VERSION="$CURRENT_MAJOR.$CURRENT_MINOR.$CURRENT_PATCH"
NEW_VERSION="$NEW_MAJOR.$NEW_MINOR.$NEW_PATCH"

echo "✅ Version bumped from $OLD_VERSION to $NEW_VERSION ($BUMP_TYPE)"

# Export for use in other scripts
export PELICO_VERSION="$NEW_VERSION"
echo "$NEW_VERSION" > .version

# Create git tag (but don't push it yet)
git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"