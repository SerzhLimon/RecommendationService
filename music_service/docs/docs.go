// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/song": {
            "get": {
                "description": "Simulates listening to a song by providing the user ID and song ID as query parameters. Returns the song's text if found.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music"
                ],
                "summary": "Listen to a song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "song_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song details",
                        "schema": {
                            "$ref": "#/definitions/models.ListenSongResponse"
                        }
                    },
                    "400": {
                        "description": "error: invalid query parameters or incorrect IDs",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "error: song not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song/create": {
            "post": {
                "description": "This endpoint allows creating a new song with a songname and a songtext.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music"
                ],
                "summary": "Create a new song",
                "parameters": [
                    {
                        "description": "Song creation request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CreateSongResponse"
                        }
                    },
                    "400": {
                        "description": "error: invalid request or failed to create song",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song/delete": {
            "delete": {
                "description": "This endpoint allows deleting a song by providing their song ID. The song ID must be greater than 0 for a successful deletion.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music"
                ],
                "summary": "Delete an existing song",
                "parameters": [
                    {
                        "description": "Song deletion request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.DeleteSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success: true",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error: invalid request, incorrect song ID or failed to delete song",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song/like": {
            "post": {
                "description": "Allows a user to like a song by providing the user ID and song ID in the request body.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music"
                ],
                "summary": "Like a song",
                "parameters": [
                    {
                        "description": "Like Song Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LikeSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success: true",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error: invalid request or failed to like song",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song/update": {
            "patch": {
                "description": "This endpoint allows updating an existing song's details like name and/or text. You must provide the song's ID and the fields to be updated (name and/or text).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music"
                ],
                "summary": "Update an existing song",
                "parameters": [
                    {
                        "description": "Song update request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success: true",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error: invalid request, incorrect song ID or failed to update song",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CreateSongRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.CreateSongResponse": {
            "type": "object",
            "properties": {
                "song_id": {
                    "type": "integer"
                }
            }
        },
        "models.DeleteSongRequest": {
            "type": "object",
            "properties": {
                "song_id": {
                    "type": "integer"
                }
            }
        },
        "models.LikeSongRequest": {
            "type": "object",
            "properties": {
                "song_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.ListenSongResponse": {
            "type": "object",
            "properties": {
                "song": {
                    "type": "string"
                }
            }
        },
        "models.UpdateSongRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "song_id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Music Service",
	Description:      "This is a simple music service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
