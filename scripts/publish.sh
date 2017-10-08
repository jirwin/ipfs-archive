#!/usr/bin/env bash

VERSION=$(git describe --tags)
echo "Publishing $VERSION..."

mkdir dist
mkdir releases
gox -osarch="linux/amd64" -osarch="linux/386" -osarch="darwin/amd64" -osarch="freebsd/amd64" -osarch="freebsd/386" -ldflags "-X main.Version=$VERSION" -output "dist/{{.OS}}_{{.Arch}}/ipfs-archive"

for i in dist/* ; do
  if [ -d "$i" ]; then
   ARCH=$(basename "$i")
   mkdir "ipfs-archive_$VERSION"
   cp "dist/$ARCH/ipfs-archive ipfs-archive_$VERSION"
   zip -r "releases/ipfs-archive_$VERSION-$ARCH.zip" "ipfs-archive_$VERSION"
   rm -rf "ipfs-archive_$VERSION"
  fi
done

ghr -t "$GITHUB_TOKEN" -u jirwin -r ipfs-archive --replace "$VERSION" releases/

rm -rf dist
rm -rf releases
