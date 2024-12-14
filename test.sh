#!/usr/bin/env bash

# Strict mode
set -euo pipefail

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Store the root directory
ROOT_DIR="$(pwd)"

# Array to store test failures
declare -a FAILURES=()

# Find all directories excluding .git and vendor
find_test_dirs() {
    find . -type d \( ! -path "*/\.*" ! -path "*/vendor/*" ! -path "*/web/*" \)
}

# Run tests in a directory
run_tests() {
    local dir="$1"
    echo "Testing $dir"
    
    # Always return to root directory before cd-ing into target
    cd "$ROOT_DIR"
    
    if ! cd "$dir"; then
        echo -e "${RED}Failed to enter directory: $dir${NC}"
        return 0 # Continue with next directory
    fi

    # Run the tests and capture output
    output=$(go test -v -cover -coverprofile=coverage.out 2>&1) || true
    exit_code=$?

    # Check for real test failures
    # Skip directories with no Go files or only main packages
    if echo "$output" | grep -q "no Go files in" || \
       echo "$output" | grep -q "package .* is a program, not an importable package"; then
        return 0
    fi

    # If tests actually failed, store the failure
    if [ $exit_code -ne 0 ]; then
        FAILURES+=("${dir}:${output}")
        echo -e "${RED}Error testing $dir${NC}"
        echo "$output"
    fi

    return 0
}

# Print coverage report with colorization
print_coverage_report() {
    echo -e "\nCoverage Report:"
    echo "----------------------------------------"
    
    # Process coverage output line by line
    while IFS= read -r line; do
        if [[ $line =~ ^total: ]]; then
            # Extract coverage percentage from total line
            percent=$(echo "$line" | grep -o '[0-9.]\+%' | sed 's/%//')
            if (( $(echo "$percent < 80" | bc -l) )); then
                echo -e "${YELLOW}$line${NC}"
            elif (( $(echo "$percent > 80" | bc -l) )); then
                echo -e "${GREEN}$line${NC}"
            else
                echo -e "$line"
            fi
        elif [[ $line =~ [[:space:]]0.0% ]]; then
            # Lines with 0% coverage
            echo -e "${RED}$line${NC}"
        elif [[ $line =~ [[:space:]]100.0% ]]; then
            # Lines with 100% coverage
            echo -e "${GREEN}$line${NC}"
        else
            echo "$line"
        fi
    done < <(go tool cover -func=merged_coverage.out)
}

# Print test failure summary
print_failure_summary() {
    if [ ${#FAILURES[@]} -eq 0 ]; then
        echo -e "\n${GREEN}All tests passed!${NC}"
        return 0
    fi

    echo -e "\n${RED}Test Failure Summary:${NC}"
    echo "----------------------------------------"
    
    for failure in "${FAILURES[@]}"; do
        dir="${failure%%:*}"
        output="${failure#*:}"
        
        echo -e "${RED}Failed: $dir${NC}"
        # Extract test name and line number using grep
        echo "$output" | grep -E "^--- FAIL:|Error Trace:" | while read -r line; do
            echo "  $line"
        done
        echo "----------------------------------------"
    done

    return 1
}

main() {
    # Clear previous coverage files
    find . -name 'coverage.out' -delete

    # Run tests in each directory
    while IFS= read -r dir; do
        run_tests "$dir" || true # Continue even if a test fails
    done < <(find_test_dirs)

    # Return to root directory for coverage processing
    cd "$ROOT_DIR"

    # Merge coverage reports if gocovmerge is available
    if command -v gocovmerge >/dev/null 2>&1; then
        if coverage_files=$(find . -name 'coverage.out'); then
            gocovmerge $coverage_files > merged_coverage.out 2>/dev/null || true
            
            # Print colorized coverage report
            if [ -f merged_coverage.out ]; then
                print_coverage_report
            fi
        fi
    fi

    # Print failure summary
    print_failure_summary
}

# Run the script
main