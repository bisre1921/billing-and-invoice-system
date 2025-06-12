# test_all.ps1

Write-Host "`n=============================================" -ForegroundColor Cyan
Write-Host "  BILLING AND INVOICE SYSTEM TEST RUNNER" -ForegroundColor Cyan
Write-Host "=============================================`n" -ForegroundColor Cyan

Write-Host "Starting tests..." -ForegroundColor Yellow
Write-Host "---------------------------------------------" -ForegroundColor DarkGray

# Run the tests with verbose output and process line by line for real-time output
$testIndex = 0
$failedTests = @()

go test -v ./tests/... 2>&1 | ForEach-Object {
    $line = $_

    # Display "RUN" messages with test name
    if ($line -match "=== RUN\s+([^\s/]+)(?:/.*)?$") {
        $testName = $matches[1]
        # Only show main test functions (not subtests)
        if ($line -match "=== RUN\s+([^\s/]+)$") {
            $testIndex++
            Write-Host "`n[$testIndex] Running Test: " -ForegroundColor Blue -NoNewline
            Write-Host "$testName" -ForegroundColor White
        } elseif ($line -match "=== RUN\s+[^\s/]+/(.+)$") {
            $subtestName = $matches[1]
            # Show subtests with proper indentation
            Write-Host "  - " -ForegroundColor Gray -NoNewline
            Write-Host "$subtestName" -ForegroundColor DarkCyan
        }
    }
    # Display passed tests
    if ($line -match "--- PASS: ([^\s/]+)(?:/.*)?.*\((.+?)\)") {
        $testName = $matches[1]
        $duration = $matches[2]

        # Only show main test results (not subtests)
        if ($line -match "--- PASS: ([^\s/]+)\s+\((.+?)\)") {
            Write-Host "  PASS " -ForegroundColor Green -NoNewline
            Write-Host "$testName" -ForegroundColor White -NoNewline
            Write-Host " ($duration)" -ForegroundColor DarkGray
        }
    }

    # Display failed tests
    if ($line -match "--- FAIL: ([^\s/]+)(?:/.*)?.*\((.+?)\)") {
        $testName = $matches[1]
        $duration = $matches[2]

        # Only show main test failures (not subtests)
        if ($line -match "--- FAIL: ([^\s/]+)\s+\((.+?)\)") {
            $failedTests += $testName
            Write-Host "  FAIL " -ForegroundColor Red -NoNewline
            Write-Host "$testName" -ForegroundColor White -NoNewline
            Write-Host " ($duration)" -ForegroundColor DarkGray
        }
    }
    # Display GIN routes being tested
    if ($line -match "\[GIN\].+\| (\d+) \|.+\| ([A-Z]+)\s+""(.+)""") {
        $status = $matches[1]
        $method = $matches[2]
        $path = $matches[3]

        $statusColor = "Green"
        if ([int]$status -ge 300 -and [int]$status -lt 400) {
            $statusColor = "Yellow"
        } elseif ([int]$status -ge 400) {
            $statusColor = "Red"
        }

        Write-Host "    -> " -ForegroundColor Gray -NoNewline
        Write-Host "$method" -ForegroundColor Magenta -NoNewline
        Write-Host " $path" -ForegroundColor DarkCyan -NoNewline
        Write-Host " [$status]" -ForegroundColor $statusColor
    }

    # Display errors in red
    if ($line -match "FAIL|panic:|fatal:") {
        Write-Host $line -ForegroundColor Red
    }
}

Write-Host "`n---------------------------------------------" -ForegroundColor DarkGray

# Display summary
if ($failedTests.Count -gt 0) {
    Write-Host "`nTEST SUMMARY: " -ForegroundColor White -NoNewline
    Write-Host "FAILED" -ForegroundColor Red
    Write-Host "Failed tests: $($failedTests -join ", ")" -ForegroundColor Red
} else {
    Write-Host "`nTEST SUMMARY: " -ForegroundColor White -NoNewline
    Write-Host "ALL TESTS PASSED SUCCESSFULLY" -ForegroundColor Green
}

Write-Host "`n=============================================" -ForegroundColor Cyan
Write-Host "        TESTING COMPLETED" -ForegroundColor Cyan
Write-Host "=============================================`n" -ForegroundColor Cyan