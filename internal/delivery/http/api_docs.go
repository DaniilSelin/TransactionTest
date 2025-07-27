package api

/*
API Documentation

Base URL: /api

## Wallet Operations

### Create Wallet
POST /wallet/create
Content-Type: application/json

Request:
{
    "balance": 100.50
}

Response:
{
    "address": "550e8400-e29b-41d4-a716-446655440000"
}

### Get Wallet Balance
GET /wallet/{address}/balance

Response:
{
    "balance": 100.50
}

### Get Wallet Info
GET /wallet/{address}

Response:
{
    "address": "550e8400-e29b-41d4-a716-446655440000",
    "balance": 100.50,
    "created_at": "2024-01-01T12:00:00Z"
}

### Update Wallet Balance
PUT /wallet/{address}/balance
Content-Type: application/json

Request:
{
    "balance": 200.75
}

Response:
{
    "message": "Balance updated successfully"
}

### Remove Wallet
DELETE /wallet/{address}

Response:
{
    "message": "Wallet removed successfully"
}

## Transaction Operations

### Send Money
POST /send
Content-Type: application/json

Request:
{
    "from": "550e8400-e29b-41d4-a716-446655440000",
    "to": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 50.25
}

Response:
{
    "message": "Money sent successfully"
}

### Get Last Transactions
GET /transactions?count=10

Response:
[
    {
        "id": 1,
        "from": "550e8400-e29b-41d4-a716-446655440000",
        "to": "550e8400-e29b-41d4-a716-446655440001",
        "amount": 50.25,
        "created_at": "2024-01-01T12:00:00Z"
    }
]

### Get Transaction by ID
GET /transaction/{id}

Response:
{
    "id": 1,
    "from": "550e8400-e29b-41d4-a716-446655440000",
    "to": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 50.25,
    "created_at": "2024-01-01T12:00:00Z"
}

### Get Transaction by Info
GET /transaction/{from}/{to}/{createdAt}

Response:
{
    "id": 1,
    "from": "550e8400-e29b-41d4-a716-446655440000",
    "to": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 50.25,
    "created_at": "2024-01-01T12:00:00Z"
}

### Remove Transaction
DELETE /transaction/{id}

Response:
{
    "message": "Transaction removed successfully"
}

## Error Responses

All endpoints return consistent error responses:

{
    "error": "Bad Request",
    "code": "INVALID_REQUEST_BODY",
    "message": "Invalid request body"
}

### Validation Errors

For validation errors, the response includes detailed field errors:

{
    "is_valid": false,
    "errors": [
        {
            "field": "amount",
            "message": "Amount must be greater than zero"
        }
    ]
}

## Error Codes

- WALLET_NOT_FOUND: Wallet with specified address not found
- TRANSACTION_NOT_FOUND: Transaction with specified ID not found
- INSUFFICIENT_FUNDS: Insufficient funds for transaction
- DUPLICATE_WALLET: Wallet with this address already exists
- NEGATIVE_BALANCE: Negative balance not allowed
- NEGATIVE_AMOUNT: Amount must be positive
- INVALID_TRANSACTION: Invalid transaction (e.g., self-transfer)
- INVALID_LIMIT: Invalid limit parameter
- INVALID_REQUEST_BODY: Invalid request body format
- INTERNAL_ERROR: Internal server error

## Validation Rules

### Wallet Address
- Must be a valid UUID format
- Cannot be empty

### Amount
- Must be greater than zero
- Cannot exceed 1e15 (1 quadrillion)

### Balance
- Cannot be negative
- Cannot exceed 1e15 (1 quadrillion)

### Limit
- Must be greater than zero
- Cannot exceed 1000

### Transaction ID
- Must be greater than zero

## Middleware

The API includes the following middleware:

1. RequestIDMiddleware: Adds unique request ID to context
2. LoggingMiddleware: Logs all HTTP requests with timing
3. RecoveryMiddleware: Handles panics gracefully
4. CORSMiddleware: Adds CORS headers for cross-origin requests

## Logging

All requests are logged with:
- HTTP method
- Request path
- Remote address
- User agent
- Status code
- Request duration
- Request ID (for correlation)

## Security Considerations

- All wallet addresses are validated as UUIDs
- Amount and balance values are validated for reasonable ranges
- Self-transfers are prevented
- Input validation is performed on all endpoints
- Panic recovery prevents server crashes
- CORS headers are included for web client compatibility
*/ 