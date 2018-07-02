#!/bin/bash
set -e
PROJECT=$(basename $(dirname $(readlink -f $0)))

go generate

NAMES=$(ls cmd/* -d | xargs -n1 basename | grep -v agilemarkdown)
for NAME in $NAMES; do
	OSES=${OSS:-"linux darwin windows"}
	ARCHS=${ARCHS:-"amd64"}
	for ARCH in $ARCHS; do
		for OS in $OSES; do
			echo $OS $ARCH $NAME
			GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=0 GOARM=7 go build -o build/${NAME}-${OS}-${ARCH} cmd/${NAME}/*.go
			if [ $? -eq 0 ]; then
				echo OK
			fi
			if [ "$OS" == "windows" ]; then
				mv build/${NAME}-${OS}-${ARCH} build/${NAME}-${OS}-${ARCH}.exe
			fi
		done
	done
done

echo "Resulting files:"
find build -type f