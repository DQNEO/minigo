#!/usr/bin/env bash
set -eu
code=$(cat internalcode/runtime.go)

cat > internalcode.go <<EOF
package main

var internalRuntimeCode = \`
$code
\`
EOF
