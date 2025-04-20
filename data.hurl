# First get a JWT token by logging in as an admin
POST http://localhost:8080/v1/user/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "pass": "admin123"
}

HTTP 200
[Captures]
jwt_token: jsonpath "$.jwt_token"

# Create data using the JWT token
POST http://localhost:8080/v1/data
Content-Type: application/json
{
  "jwt": "{{jwt_token}}",
  "company_data": {
    "name": "Acme Corporation",
    "employees": 150,
    "departments": ["Engineering", "Sales", "Marketing"],
    "active": true,
    "headquarters": {
      "city": "San Francisco",
      "country": "USA"
    }
  }
}

HTTP 200
[Asserts]
jsonpath "$.data_id" exists 