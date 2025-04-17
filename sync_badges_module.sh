#!/bin/bash

# Create backup directory if it doesn't exist
mkdir -p ../badges-module/backups

# Backup codec.go if it exists
if [ -f "../badges-module/x/badges/types/codec.go" ]; then
    cp "../badges-module/x/badges/types/codec.go" "../badges-module/backups/codec.go.bak"
fi

# Sync the badges module
rm -rf ../badges-module/x/badges
cp -r ./x/badges ../badges-module/x/badges

# Restore codec.go if it was backed up
if [ -f "../badges-module/backups/codec.go.bak" ]; then
    cp "../badges-module/backups/codec.go.bak" "../badges-module/x/badges/types/codec.go"
    rm "../badges-module/backups/codec.go.bak"
fi

# Sync the badges module types
rm -rf ../badges-module/api/badges
cp -r ./api/badges ../badges-module/api/badges

# Sync the badges module keeper
rm -rf ../badges-module/docs
cp -r ./docs ../badges-module/docs

rm -rf ../badges-module/proto/badges
cp -r ./proto/badges ../badges-module/proto/badges

rm -rf ../badges-module/x/badges/types/expected_ibc_keeper.go


#Recursively replace all instances of github.com/bitbadges/bitbadgeschain with github.com/bitbadges/badges-module
find ../badges-module -type f -exec sed -i 's/github.com\/bitbadges\/bitbadgeschain/github.com\/bitbadges\/badges-module/g' {} +

