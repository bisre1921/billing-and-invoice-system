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
        "/auth/login": {
            "post": {
                "description": "Authenticates the user and returns a JWT token if credentials are valid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user and return JWT token",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginCredentials"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register/user": {
            "post": {
                "description": "Creates a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Registers a new user",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.GenericResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/company/create": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Registers a new company in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Company"
                ],
                "summary": "Creates a new company",
                "parameters": [
                    {
                        "description": "Company Data",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Company"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Company created successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/user/update/{id}": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates user details by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User Data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateUserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Company": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.GenericResponse": {
            "type": "object",
            "properties": {
                "id": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "models.LoginCredentials": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "models.UpdateUserInput": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Billing and Invoice System API",
	Description:      "This is the API documentation for the Billing and Invoice system.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
