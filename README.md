# URL Shortener Service API

HTTP API for URL shortening service with the following endpoints:

## URL Operations

### POST /api/shorten
Request body:
{
    "url": "string"    // Original URL to be shortened
}
Arguments:
- url: required field, must be a valid URL

Response: 
{
    "result": "string"  // Shortened URL
}

### POST /api/shorten/batch
Request body:
[
    {
        "correlation_id": "string",  // Client-defined ID
        "original_url": "string"     // URL to be shortened
    }
]

Response:
[
    {
        "correlation_id": "string",  // Matching client ID
        "short_url": "string"       // Shortened URL
    }
]

### GET /{id}
Arguments:
- id: shortened URL identifier

Response: Redirects to original URL

### GET /api/user/urls
Response:
[
    {
        "short_url": "string",     // Shortened URL
        "original_url": "string"   // Original URL
    }
]

### DELETE /api/user/urls
Request body:
[
    "string"   // Array of shortened URL IDs to delete
]

### GET /ping
Response: Database connection status

## Response Codes
- 200: Successful operation
- 201: URL successfully created
- 307: Temporary redirect
- 400: Invalid request format
- 401: Authentication required
- 404: URL not found
- 409: URL already exists
- 500: Internal server error

## Authentication
Protected endpoints require JWT token in header:
Authorization: Bearer <token>
