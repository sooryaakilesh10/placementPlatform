#!/bin/bash

# Define the base URL
BASE_URL="http://localhost:8080/v1/user"

# Initialize counters
PASSED=0
FAILED=0
# Variable to store the ID of the first successfully created user
CREATED_USER_ID="" 

# Function to make POST request and check response
test_post() {
    local test_name="$1"
    local json_data="$2"
    local expected_status="$3"
    local capture_id_flag="$4" # Flag to indicate if we should capture the user_id

    echo "Running Test: $test_name"
    echo "Request Data: $json_data"
    
    # Make the POST request and capture the HTTP status code and response body
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$json_data" \
        "$BASE_URL")
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract response body (all but the last line)
    body=$(echo "$response" | sed '$d')
    
    echo "Response Body: $body"
    echo "Expected Status: $expected_status"
    echo "Actual Status: $status_code"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo "✓ Test Passed"
        ((PASSED++))
        # Capture the user_id if flag is set, status is 200, and CREATED_USER_ID is not already set
        if [ "$capture_id_flag" == "capture" ] && [ "$status_code" -eq 200 ] && [ -z "$CREATED_USER_ID" ]; then
             # Check if jq is installed
            if command -v jq &> /dev/null; then
                # Attempt to parse user_id using jq
                parsed_id=$(echo "$body" | jq -r '.user_id // empty')
                if [ -n "$parsed_id" ] && [ "$parsed_id" != "null" ]; then
                    CREATED_USER_ID="$parsed_id"
                    echo "Captured User ID: $CREATED_USER_ID"
                else
                    echo "Warning: Could not parse user_id from response body: $body"
                fi
            else
                 echo "Warning: jq is not installed. Cannot capture user_id."
            fi
        fi
    else
        echo "✗ Test Failed"
        ((FAILED++))
    fi
    echo "------------------------"
}

# Function to make GET request and check response
test_get() {
    local test_name="$1"
    local endpoint="$2"
    local expected_status="$3"

    echo "Running Test: $test_name"
    echo "Endpoint: $endpoint"
    
    # Make the GET request and capture the HTTP status code and response body
    response=$(curl -s -w "\n%{http_code}" -X GET "$endpoint")
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract response body (all but the last line)
    body=$(echo "$response" | sed '$d')
    
    echo "Response Body: $body"
    echo "Expected Status: $expected_status"
    echo "Actual Status: $status_code"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo "✓ Test Passed"
        ((PASSED++))
    else
        echo "✗ Test Failed"
        ((FAILED++))
    fi
    echo "------------------------"
}

# Test health endpoint
echo "Testing health endpoint"
test_get "Health endpoint" "$BASE_URL/health" 200

# Test Case 1: Admin user creation - Capture the ID
test_post "Admin user creation" \
    '{"user_name": "admin", "email": "admin@gmail.com", "pass": "test1@123", "role": "admin"}' \
    200 \
    "capture" # Add flag to capture ID

# Test Case 2: manager user creation
test_post "manager user creation" \
    '{"user_name": "soorya", "email": "inst@gmail.com", "pass": "test1@123", "role": "manager"}' \
    200

# Test Case 3: Valid user creation
test_post "Valid user creation" \
    '{"user_name": "testuser", "email": "testuser@gmail.com", "pass": "test1@123", "role": "admin"}' \
    200

# Test Case 4: Duplicate user
test_post "Duplicate user" \
    '{"user_name": "testuser", "email": "testuser@gmail.com", "pass": "test1@123", "role": "admin"}' \
    500

# Test Case 5: New email, same username
test_post "New email, same username" \
    '{"user_name": "testuser", "email": "testuser2@gmail.com", "pass": "test1@123", "role": "admin"}' \
    200

# Test Case 6: Duplicate email, different role
test_post "Duplicate email, different role" \
    '{"user_name": "testuser3", "email": "testuser2@gmail.com", "pass": "test1@123", "role": "tester"}' \
    500

# Test Case 7: Invalid email format (missing domain)
test_post "Invalid email format 1" \
    '{"user_name": "testuser4", "email": "testuser4@.com", "pass": "test1@123", "role": "admin"}' \
    500

# Test Case 8: Invalid email format (incomplete)
test_post "Invalid email format 2" \
    '{"user_name": "testuser5", "email": "@s.com", "pass": "test1@123", "role": "admin"}' \
    500

# Test Case 9: Empty password
test_post "Empty password" \
    '{"user_name": "testuser6", "email": "testuser6@gmail.com", "pass": "", "role": "admin"}' \
    500

# Test Case 10: Empty username
test_post "Empty username" \
    '{"user_name": "", "email": "testuser7@gmail.com", "pass": "test1@123", "role": "admin"}' \
    500

# Test Case 11: Get user by ID (using captured ID)
echo "Testing Get User By ID using captured ID: $CREATED_USER_ID"
if [ -n "$CREATED_USER_ID" ]; then
    test_get "Get user by ID" "$BASE_URL/id/$CREATED_USER_ID" 200
else
    echo "Skipping Test: Get user by ID (No user ID captured)"
    ((FAILED++)) # Count as failed if we couldn't capture an ID
fi

# Test Case 12: Get user by non-existent ID (use zero UUID)
test_get "Get user by non-existent ID" "$BASE_URL/id/00000000-0000-0000-0000-000000000000" 500

# Test Case 13: Get user by invalid ID format
test_get "Get user by invalid ID format" "$BASE_URL/id/abc" 400

# Test Case 14: Get user by email (assuming admin@gmail.com exists from test case 1)
test_get "Get user by email" "$BASE_URL/email/admin@gmail.com" 200

# Test Case 15: Get user by non-existent email
test_get "Get user by non-existent email" "$BASE_URL/email/nonexistent@gmail.com" 500

# Print summary
echo "Test Summary:"
echo "Total Tests: $((PASSED + FAILED))"
echo "Tests Passed: $PASSED"
echo "Tests Failed: $FAILED"
echo "------------------------"
echo "All tests completed!"