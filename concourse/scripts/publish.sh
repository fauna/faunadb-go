#!/bin/sh

set -eou

git clone fauna-go-repository fauna-go-repository-updated

cd fauna-go-repository-updated

CURRENT_VERSION=$(cat version)

echo "Current version: $CURRENT_VERSION"

echo "Publishing a new $CURRENT_VERSION version..."
git config --global user.email "nobody@concourse-ci.org"
git config --global user.name "Fauna, Inc"

git tag "$CURRENT_VERSION"
