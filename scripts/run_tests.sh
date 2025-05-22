#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running all tests...${NC}"

# Create directory for coverage reports
mkdir -p coverage

# Run all tests with code coverage
go test -coverprofile=coverage/coverage.out ./...

# Check test execution result
if [ $? -eq 0 ]; then
    echo -e "${GREEN}All tests passed successfully!${NC}"
    
    # Generate HTML code coverage report
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    echo -e "${GREEN}Coverage report saved to coverage/coverage.html${NC}"
    
    # Display brief coverage information
    COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
    echo -e "${GREEN}Total code coverage: ${COVERAGE}${NC}"
else
    echo -e "${RED}Some tests failed. Check the output above.${NC}"
    exit 1
fi 