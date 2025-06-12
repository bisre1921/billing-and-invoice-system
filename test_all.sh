#!/bin/bash
# filepath: c:\Users\Nafyad\Desktop\final Year project\billing-and-invoice-system\test_all.sh

# ANSI color codes
CYAN='\033[1;36m'
YELLOW='\033[1;33m'
GREEN='\033[1;32m'
RED='\033[1;31m'
BLUE='\033[1;34m'
WHITE='\033[1;37m'
GRAY='\033[1;30m'
DARK_CYAN='\033[0;36m'
MAGENTA='\033[1;35m'
NC='\033[0m' # No Color

echo -e "\n${CYAN}============================================="
echo -e "  BILLING AND INVOICE SYSTEM TEST RUNNER"
echo -e "=============================================\n${NC}"

echo -e "${YELLOW}Starting tests...${NC}"
echo -e "${GRAY}---------------------------------------------${NC}"

# Run tests and capture output with test numbering
test_index=0
go test -v ./tests/... | while IFS= read -r line; do
    # Display "RUN" messages with test name
    if [[ $line =~ ^===\ RUN\ +([^[:space:]/]+)$ ]]; then
        test_name="${BASH_REMATCH[1]}"
        ((test_index++))
        echo -e "\n${BLUE}[$test_index] Running Test: ${WHITE}$test_name${NC}"
    elif [[ $line =~ ^===\ RUN\ +.+/([^/]+)$ ]]; then
        subtest_name="${BASH_REMATCH[1]}"
        echo -e "  ${GRAY}- ${DARK_CYAN}$subtest_name${NC}"
    fi

    # Display passed tests
    if [[ $line =~ ^---\ PASS:\ ([^[:space:]]+)\ \((.+)\)$ ]]; then
        test_name="${BASH_REMATCH[1]}"
        duration="${BASH_REMATCH[2]}"
        # Only show main test results (not subtests)
        if [[ ! $test_name =~ / ]]; then
            echo -e "  ${GREEN}PASS${NC} ${WHITE}$test_name${NC} ${GRAY}($duration)${NC}"
        fi
    fi

    # Display failed tests
    if [[ $line =~ ^---\ FAIL:\ ([^[:space:]]+)\ \((.+)\)$ ]]; then
        test_name="${BASH_REMATCH[1]}"
        duration="${BASH_REMATCH[2]}"
        # Only show main test failures (not subtests)
        if [[ ! $test_name =~ / ]]; then
            echo -e "  ${RED}FAIL${NC} ${WHITE}$test_name${NC} ${GRAY}($duration)${NC}"
        fi
    fi

    # Display GIN routes being tested
    if [[ $line =~ "\[GIN\]" ]]; then
        if [[ $line =~ "[0-9]+" ]]; then
            status="${BASH_REMATCH[0]}"
            if [[ $line =~ "POST|GET|PUT|DELETE|PATCH" ]]; then
                method="${BASH_REMATCH[0]}"
                if [[ $line =~ \"([^\"]+)\" ]]; then
                    path="${BASH_REMATCH[1]}"

                    status_color=$GREEN
                    if (( status >= 300 && status < 400 )); then
                        status_color=$YELLOW
                    elif (( status >= 400 )); then
                        status_color=$RED
                    fi

                    echo -e "    ${GRAY}->${NC} ${MAGENTA}$method${NC} ${DARK_CYAN}$path${NC} ${status_color}[$status]${NC}"
                fi
            fi
        fi
    fi

    # Display errors in red
    if [[ $line =~ (FAIL|panic:|fatal:) ]]; then
        echo -e "${RED}$line${NC}"
    fi
done

echo -e "\n${GRAY}---------------------------------------------${NC}"

# Check the exit code of the last test run
if [ $? -eq 0 ]; then
    echo -e "\n${WHITE}TEST SUMMARY: ${GREEN}ALL TESTS PASSED SUCCESSFULLY${NC}"
else
    echo -e "\n${WHITE}TEST SUMMARY: ${RED}FAILED${NC}"
fi

echo -e "\n${CYAN}============================================="
echo -e "        TESTING COMPLETED"
echo -e "=============================================\n${NC}"