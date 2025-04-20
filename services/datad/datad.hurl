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
admin_user_id: jsonpath "$.user_id"

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
[Captures]
manager_user_id: jsonpath "$.user_id"

# Test Case 3: Placement Officer user creation
POST http://localhost:8080/v1/user
Content-Type: application/json

{
  "user_name": "soorya",
  "email": "inst1@gmail.com",
  "pass": "test1@123",
  "role": "user"
}

HTTP 200
[Captures]
officer_user_id: jsonpath "$.user_id"

# Get JWT by login endpoint
POST http://localhost:8080/v1/login
Content-Type: application/json

{
  "user_name": "admin",
  "email": "admin@gmail.com",
  "pass": "test1@123"
}

HTTP 200
[Captures]
admin_jwt: jsonpath "$.jwt_token"

# Get JWT by login endpoint
POST http://localhost:8080/v1/login
Content-Type: application/json

{
  "user_name": "soorya",
  "email": "inst@gmail.com",
  "pass": "test1@123"
}

HTTP 200
[Captures]
manager_jwt: jsonpath "$.jwt_token"

# Get JWT by login endpoint
POST http://localhost:8080/v1/login
Content-Type: application/json

{
  "user_name": "soorya",
  "email": "inst1@gmail.com",
  "pass": "test1@123"
}

HTTP 200
[Captures]
officer_jwt: jsonpath "$.jwt_token"

# Test Company data creation with admin user
POST http://localhost:8080/v1/data
Content-Type: application/json
Authorization: Bearer {{admin_jwt}}

{
    "jwt": "{{admin_jwt}}",
    "company_data": {
        "name": "Test Company",
        "employees": 100,
        "active": true
    }
}

HTTP 200
[Captures]
data_id: jsonpath "$.data_id"

# Test Company data creation with manager user
POST http://localhost:8080/v1/data
Content-Type: application/json
Authorization: Bearer {{manager_jwt}}

{
    "jwt": "{{manager_jwt}}",
    "company_data": {
        "name": "Test Company",
        "employees": 100,
        "active": true
    }
}

HTTP 200

# Test Company data creation with officer user (unauthorized role)
POST http://localhost:8080/v1/data
Content-Type: application/json

{
    "jwt": "{{officer_jwt}}",
    "company_data": {
        "name": "Test Company",
        "employees": 100,
        "active": true
    }
}

HTTP 403