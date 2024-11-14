#!/bin/bash

#TODO: Add your own paths here
homeDir="/home/trevormil/.bitbadgeschain"
binaryDir="/home/trevormil/go/bin"
binaryName="bitbadgeschaind"
dataPath="$homeDir/data"
backupPath="/tmp/711315-data.bak"

currDir=$(pwd)
git clone https://github.com/BitBadges/bitbadgeschain.git

# We assume you currently have 1 - 711315 saved in the data folder


# Backup the current data folder (only if not already backed up)
# Helps avoid the case where we overwrite the backup with new non-snapshot data folder
if [ ! -d "$backupPath" ]; then
    mv $dataPath $backupPath
else
    echo "Backup already exists, skipping to avoid backup corruption..."
fi

cd $homeDir/config
rm genesis.json
curl -o genesis.json https://raw.githubusercontent.com/BitBadges/bitbadgeschain/master/genesis-711316.json

cd $homeDir
# Reset and start chain as new from 711316+
$binaryDir/$binaryName comet unsafe-reset-all --home $homeDir
# Sync a few blocks (it will give a halt error once done but we just ignore that)
$binaryDir/$binaryName start --pruning=nothing --rpc.laddr=tcp://0.0.0.0:26657 --rpc.unsafe --halt-height=711320 --home $homeDir

cd $currDir
cd ./bitbadgeschain/scripts
go run migrate.go -source $backupPath -target $dataPath

# Cleanup
cd $currDir
rm -rf bitbadgeschain