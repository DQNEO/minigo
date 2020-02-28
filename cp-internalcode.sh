#!/usr/bin/env bash
set -eu
code=$(cat internal/runtime/runtime.go)

cat > internal_runtime.go <<EOF
package main

var internalRuntimeCode string = string(\`
$code
\`)
EOF

code=$(cat internal/universe/universe.go)

cat > internal_universe.go <<EOF
package main

var internalUniverseCode string = string(\`
$code
\`)
EOF

code=$(cat stdlib/unsafe/unsafe.go)

cat > internal_unsafe.go <<EOF
package main

var internalUnsafeCode string = string(\`
$code
\`)
EOF
