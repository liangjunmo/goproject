#!/usr/bin/env bash

version=$1

if [[ $version == "" ]]
then
    echo "please specify a version, like 1.2.3"
    exit
fi

if [[ $OSTYPE == darwin* ]]
then
    sed -i "" -E "s/^(const Version = \")(.*)(\")$/\1$version\3/" internal/version/version.go
else
    sed -i -E "s/^(const Version = \")(.*)(\")$/\1$version\3/" internal/version/version.go
fi

git add . && git commit -m "bump version to $version"
