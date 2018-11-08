// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2018-11-08 23:10:53.654594 +0700 +07 m=+0.034657107

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Smart Home System",
        "title": "Smart home API",
        "contact": {
            "name": "Alex Filippov",
            "url": "https://e154.github.io/smart-home/",
            "email": "support@e154.ru"
        },
        "license": {
            "name": "MIT License",
            "url": "https://raw.githubusercontent.com/e154/smart-home/master/LICENSE"
        },
        "version": "2"
    },
    "host": "{{.Host}}",
    "basePath": "/api/v2",
    "paths": {
        "/": {
            "get": {
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "index"
                ],
                "summary": "index page",
                "responses": {
                    "200": {}
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
