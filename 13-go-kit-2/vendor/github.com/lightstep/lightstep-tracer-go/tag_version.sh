#!/bin/sh

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "version must be set."
    exit 1
fi

echo $VERSION > ./VERSION

cat > version.go <<EOF
package lightstep

// TracerVersionValue provides the current version of the lightstep-tracer-go release
const TracerVersionValue = "$VERSION"
EOF
