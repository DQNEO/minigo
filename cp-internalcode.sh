#!/usr/bin/env bash
set -eu
code=$(cat internalcode/runtime.go)

cat > internalcode.go <<EOF
package main

var internalRuntimeCode string = \`
$code
\`
EOF

code=$(cat internalcode/universe.go)

cat > universe.go <<EOF
package main

var internalUniverseCode string = \`
$code
\`
EOF
