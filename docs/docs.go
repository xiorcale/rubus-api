// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-05-29 12:01:30.013569535 +0200 CEST m=+0.044485400

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Quentin Vaucher",
            "email": "quentin.vaucher3@master.hes-so.ch"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/device": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Add a ` + "`" + `Device` + "`" + ` into the database and prepare the necessary directory structure for deploying it.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "operationId": "createDevice",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The hostname of the device",
                        "name": "hostname",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The device's switch port",
                        "name": "port",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Device"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Delete a ` + "`" + `Device` + "`" + ` from the database and remove its directory structure used for deployment.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Delete a device",
                "operationId": "deleteDevice",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The hostname of the device",
                        "name": "hostname",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "The device's switch port",
                        "name": "deviceId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {}
                }
            }
        },
        "/admin/user": {
            "get": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Return a list containing all the ` + "`" + `User` + "`" + `",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "List all the users",
                "operationId": "listUser",
                "responses": {
                    "200": {
                        "description": "A JSON array listing all the users",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Create a new Rubus ` + "`" + `User` + "`" + ` and save it into the database.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Create a new user",
                "operationId": "createUser",
                "parameters": [
                    {
                        "description": "All the fields are required, except for the ` + "`" + `role` + "`" + ` which will default to ` + "`" + `user` + "`" + ` if not specified, and the expiration date which can be null.",
                        "name": "RequestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NewUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            }
        },
        "/admin/user/{id}": {
            "delete": {
                "description": "Delete the ` + "`" + `User` + "`" + ` with the given id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Delete a user",
                "operationId": "deleteUser",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The id from the user to delete",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {}
                }
            }
        },
        "/admin/user/{id}/expiration": {
            "post": {
                "description": "Update the expiration date of a the` + "`" + `User` + "`" + ` with the given id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Set a new expiration date for a ` + "`" + `User` + "`" + `",
                "operationId": "updateUser",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The id from the user to update",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The new expiration date",
                        "name": "expiration",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {}
                }
            }
        },
        "/device": {
            "get": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "List all the ` + "`" + `Device` + "`" + `",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "list all the devices",
                "operationId": "listDevice",
                "responses": {
                    "200": {
                        "description": "A JSON array listing all the devices",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Device"
                            }
                        }
                    }
                }
            }
        },
        "/device/{id}": {
            "get": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Return the ` + "`" + `Device` + "`" + ` with the given ` + "`" + `id` + "`" + `.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "get a device by id",
                "operationId": "getDevice",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The id of the ` + "`" + `Device` + "`" + ` to get",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Device"
                        }
                    }
                }
            }
        },
        "/device/{id}/acquire": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Set the ` + "`" + `User` + "`" + ` who made the request as the owner of the ` + "`" + `Device` + "`" + `.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "acquire a device",
                "operationId": "acquire",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The id of the ` + "`" + `Device` + "`" + ` to acquire",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Device"
                        }
                    }
                }
            }
        },
        "/device/{id}/deploy": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Configure the PXE boot for the ` + "`" + `Device` + "`" + ` and reboot it.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "deploy a device",
                "operationId": "deploy",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The device id to deploy",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {}
                }
            }
        },
        "/device/{id}/off": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Shuts down the ` + "`" + `Device` + "`" + ` on the given ` + "`" + `port` + "`" + `",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "Shut down a device",
                "operationId": "powerOff",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The device id to turn off",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {}
                }
            }
        },
        "/device/{id}/on": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Boot the ` + "`" + `Device` + "`" + ` with the given ` + "`" + `id` + "`" + `.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "Boot a device",
                "operationId": "powerOn",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The device id to turn on",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {}
                }
            }
        },
        "/device/{id}/release": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Remove the ` + "`" + `Device` + "`" + `'s ownership from the ` + "`" + `User` + "`" + ` who made the request.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device"
                ],
                "summary": "release a device",
                "operationId": "release",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The device port to release",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Device"
                        }
                    }
                }
            }
        },
        "/login": {
            "get": {
                "description": "Log a ` + "`" + `User` + "`" + ` into the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Log a user in",
                "operationId": "login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The username used to login",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The password used to login",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {}
                }
            }
        },
        "/user/me": {
            "get": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "Return the ` + "`" + `User` + "`" + ` who made the request",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get the authenticated user",
                "operationId": "getMe",
                "responses": {
                    "200": {
                        "description": "A JSON object describing a user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            },
            "put": {
                "description": "Update the ` + "`" + `User` + "`" + ` who made the request.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "update the authenticated user",
                "operationId": "updateMe",
                "parameters": [
                    {
                        "description": "the ` + "`" + `User` + "`" + ` fields which can be updated. Giving all the fields is not mendatory, but at least one of them is required.",
                        "name": "RequestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PutUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A JSON object describing a user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete the ` + "`" + `User` + "`" + ` who made the request.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "delethe the autenticated user",
                "operationId": "deleteMe",
                "responses": {
                    "204": {}
                }
            }
        }
    },
    "definitions": {
        "models.Device": {
            "type": "object",
            "properties": {
                "hostname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "isTurnedOn": {
                    "type": "boolean"
                },
                "owner": {
                    "type": "integer"
                }
            }
        },
        "models.NewUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "rubus@mail.com"
                },
                "expiration": {
                    "type": "string",
                    "example": "2020-05-18"
                },
                "password": {
                    "type": "string",
                    "example": "rubus_secret"
                },
                "role": {
                    "type": "string",
                    "example": "administrator"
                },
                "username": {
                    "type": "string",
                    "example": "rubus"
                }
            }
        },
        "models.PutUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "rubus@mail.com"
                },
                "password": {
                    "type": "string",
                    "example": "rubus_secret"
                },
                "username": {
                    "type": "string",
                    "example": "rubus"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "rubus@mail.com"
                },
                "expiration": {
                    "type": "string",
                    "example": "2020-05-18"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "role": {
                    "type": "string",
                    "example": "administrator"
                },
                "username": {
                    "type": "string",
                    "example": "rubus"
                }
            }
        }
    },
    "securityDefinitions": {
        "jwt": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "Operations about user authentication",
            "name": "authentication"
        },
        {
            "description": "Operations which require administrative rights",
            "name": "admin"
        },
        {
            "description": "Operations about devices, such as provisioning or deployment",
            "name": "device"
        },
        {
            "description": "Operations about Users",
            "name": "user"
        }
    ]
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:1323",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Rubus API",
	Description: "Rubus API exposes provisioning services to manage an edge cluster system (i.e. Raspberry pi). This API takes advantage of various HTTP features like authentication, verbs or status code. All requests and response bodies are JSON encoded, including error responses.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
