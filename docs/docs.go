// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/shorten": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Creates a shortened version of a URL provided in JSON format",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Shorten URL via JSON",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "URL to shorten",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ShortenURLRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Shortened URL",
                        "schema": {
                            "$ref": "#/definitions/models.ShortenURLResponse"
                        }
                    },
                    "409": {
                        "description": "URL already exists",
                        "schema": {
                            "$ref": "#/definitions/models.ShortenURLResponse"
                        }
                    }
                }
            }
        },
        "/api/shorten/batch": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Creates shortened versions for multiple URLs in a single request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Shorten multiple URLs in batch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Array of URLs to shorten",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.BatchUnitURLRequest"
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Array of shortened URLs",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.BatchUnitURLResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Error reading body!/Error unmarshalling body!/Empty or malformed body sent!/Error saving URLs!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/url": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Creates a shortened version of a provided URL",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Create shortened URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Original URL to shorten",
                        "name": "url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Shortened URL",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Can't read body!/Empty body!/Malformed URI!/Couldn't encode URL!",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "URL already exists",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/urls": {
            "delete": {
                "description": "Delete multiple URLs for a specific user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Delete URLs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Array of URLs to delete",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Error reading body!/Error unmarshalling body!/Empty or malformed body sent!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/urls": {
            "get": {
                "description": "Retrieves all URLs associated with the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Get user's URLs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of user's URLs",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.UserURLResponse"
                            }
                        }
                    },
                    "204": {
                        "description": "No URLs found!",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Error finding URLs!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Check if database connection is alive",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Ping database",
                "responses": {
                    "200": {
                        "description": "Live",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Can't connect to the Database!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/{id}": {
            "get": {
                "description": "Retrieves and redirects to the original URL from a shortened URL ID",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Get original URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Shortened URL ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "307": {
                        "description": "Temporary Redirect",
                        "schema": {
                            "type": "string"
                        },
                        "headers": {
                            "Location": {
                                "type": "string",
                                "description": "Original URL for redirect"
                            }
                        }
                    },
                    "400": {
                        "description": "URL not found!",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "410": {
                        "description": "URL was deleted!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.BatchUnitURLRequest": {
            "type": "object",
            "properties": {
                "correlation_id": {
                    "type": "string"
                },
                "original_url": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "models.BatchUnitURLResponse": {
            "type": "object",
            "properties": {
                "correlation_id": {
                    "type": "string"
                },
                "short_url": {
                    "type": "string"
                }
            }
        },
        "models.ShortenURLRequest": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "models.ShortenURLResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "string"
                }
            }
        },
        "models.UserURLResponse": {
            "type": "object",
            "properties": {
                "original_url": {
                    "type": "string"
                },
                "short_url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL Shortener API",
	Description:      "A URL shortening service API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
