# Test Case 1: Admin user creation - Capture the ID
POST http://localhost:8080/v1/user
Content-Type: application/json

{
  "user_name": "admin",
  "email": "admin@gmail.com",
  "pass": "test1@123",
  "role": "admin"
}

HTTP 200
[Captures]
user_id: jsonpath "$.user_id"

# Test Case 2: Manager user creation
POST http://localhost:8080/v1/user
Content-Type: application/json

{
  "user_name": "soorya",
  "email": "inst@gmail.com",
  "pass": "test1@123",
  "role": "manager"
}

HTTP 200

# Test Case 3: Placement Officer user creation
POST http://localhost:8080/v1/user
Content-Type: application/json

{
  "user_name": "soorya",
  "email": "inst1@gmail.com",
  "pass": "test1@123",
  "role": "user"
}

# Test get user by ID endpoint
GET http://localhost:8080/v1/user/id/{{user_id}}

HTTP 200
[Asserts]
jsonpath "$.user_name" exists
jsonpath "$.email" exists
jsonpath "$.role" exists

# Test get user by email endpoint
GET http://localhost:8080/v1/user/email/admin@gmail.com

HTTP 200
[Asserts]
jsonpath "$.user_name" exists
jsonpath "$.email" exists
jsonpath "$.role" exists
