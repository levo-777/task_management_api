#!/bin/bash

echo "Testing Taskify API with all new features..."

# Test rate limiting
echo "Testing rate limiting..."
for i in {1..5}; do
    echo "Request $i:"
    curl -s -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username": "admin", "password": "admin123"}' \
        -w "Status: %{http_code}\n"
    echo ""
done

echo "Waiting 3 seconds for rate limit reset..."
sleep 3

# Test admin login
echo "Testing admin login..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "admin123"}')

echo "Login response: $RESPONSE"

# Extract token if login successful
if [[ $RESPONSE == *"access_token"* ]]; then
    TOKEN=$(echo $RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    echo "Token extracted: ${TOKEN:0:50}..."
    
    # Test pagination and filtering
    echo "Testing pagination and filtering..."
    curl -s -X GET "http://localhost:8080/api/v1/tasks?page=1&page_size=5&search=test&sort_by=created_at&sort_order=desc" \
        -H "Authorization: Bearer $TOKEN" \
        -w "Status: %{http_code}\n"
    
    # Test caching (same request should be faster)
    echo "Testing caching (second request)..."
    curl -s -X GET "http://localhost:8080/api/v1/tasks?page=1&page_size=5" \
        -H "Authorization: Bearer $TOKEN" \
        -w "Status: %{http_code}\n"
        
    # Test task creation
    echo "Testing task creation..."
    curl -s -X POST http://localhost:8080/api/v1/tasks \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"title": "Test Task", "description": "Testing pagination and caching", "status": "pending", "priority": "medium"}' \
        -w "Status: %{http_code}\n"
else
    echo "Login failed, cannot test other endpoints"
fi

echo "API testing completed!"
