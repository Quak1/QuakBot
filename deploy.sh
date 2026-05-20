#!/bin/bash
set -e

VPS_USER="deploy"
VPS_IP="vps"
VPS_PATH="/opt/bot"
BINARY_NAME="quakmate"

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o $BINARY_NAME ./cmd/run 

echo "Copying to VPS..."
scp $BINARY_NAME $VPS_USER@$VPS_IP:/tmp/myapp

echo "Installing..."
ssh -t $VPS_USER@$VPS_IP "
    sudo mv /tmp/myapp $VPS_PATH/$BINARY_NAME &&
    sudo chown botuser:botgroup $VPS_PATH/$BINARY_NAME &&
    sudo chmod 750 $VPS_PATH/$BINARY_NAME &&
    sudo systemctl restart $BINARY_NAME &&
    sudo systemctl status $BINARY_NAME --no-pager
"

rm $BINARY_NAME

echo "Done"
