#!/usr/bin/env bash

unset LOG_LEVEL

set -uo pipefail

while IFS=\= read -r dir; do
    dirs+=("$dir")
done < <(find . -type d -not -path "./.git*")

# readarray -t dirs < <(find . -type d -not -path "./.git*")

run_tests() {
    local dir=$1
    echo "Testing $dir"
    (
        cd "$dir" || exit 1
        out=$(go test -cover -coverprofile=coverage.out 2>&1) 
        code=$?
        if [ $code -ne 0 ]; then
            if echo "$out"| grep -q "files in"; then
                code=0
            else
                echo "Error testing $dir"
                echo "$out"
                exit $code
            fi
        fi

    )
}


for dir in "${dirs[@]}"; do
    run_tests "$dir"
done

"$HOME/go/bin/gocovmerge" "$(find . -name 'coverage*.out')" > merged_coverage.out

go tool cover -func=merged_coverage.out  > coverage.txt
