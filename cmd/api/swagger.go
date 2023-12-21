//go:build swagger

// Rúnar Himmel is an enchanting Golang project inspired by Norse mythology
//
// API documents for Rúnar Himmel.
//
// ## Authentication
// All API endpoints with the lock icon require authentication token.
// Firstly, grab the **access_token** from the response of `/auth/login`. Then include this header in all API calls:
// ```
// Authorization: Bearer ${access_token}
// ```
//
// Terms Of Service: N/A
//
// Version: %{VERSION}
// Contact: m15t <nguyen.ndk@outlook.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearer: []
//
// SecurityDefinitions:
// bearer:
//	 type: apiKey
//	 name: Authorization
//	 in: header
//
// swagger:meta
package main

import (
	"embed"
)

//go:embed swagger-ui
var embedSwaggerui embed.FS

func init() {
	enableSwagger = true
	swaggerui = embedSwaggerui
}

// OK - No Content
// swagger:response ok
type swaggOKResp struct{}
