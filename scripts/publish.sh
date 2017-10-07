#!/bin/bash

VERSION=$(git describe --tags)
echo "Publishing $VERSION..."

mkdir dist
mkdir releases
gox -osarch="linux/amd64" -osarch="linux/386" -osarch="darwin/amd64" -osarch="freebsd/amd64" -osarch="freebsd/386" -ldflags "-X main.Version=$VERSION" -output "dist/{{.OS}}_{{.Arch}}/ipfs-archive"

for i in dist/* ; do
  if [ -d "$i" ]; then
   ARCH=$(basename "$i")
   zip releases/ipfs-archive_$VERSION_$ARCH.zip dist/$ARCH/ipfs-archive
  fi
done

ghr -t $GITHUB_TOKEN -u jirwin -r ipfs-archive --replace $VERSION releases/

rm -rf dist
rm -rf releases
