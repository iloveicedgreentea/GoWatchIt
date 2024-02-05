#!/bin/bash

unset LOG_LEVEL

set -euo pipefail
readarray -t dirs < <(find . -type d -not -path "./.git*")

run_tests() {
    local dir=$1
    echo "Testing $dir"
    (cd "$dir" && go test -cover -coverprofile=coverage.out)
}


for dir in "${dirs[@]}"; do
    run_tests "$dir"
done

"$HOME/go/bin/gocovmerge" "$(find . -name 'coverage*.out')" > merged_coverage.out

go tool cover -func=merged_coverage.out  > coverage.txt
