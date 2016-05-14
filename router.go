// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"time"

	"github.com/husobee/vestigo"
)

func NewRouter() *vestigo.Router {
	router := vestigo.NewRouter()
	vestigo.AllowTrace = true

	// Global CORS policy
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*", "nerds.berlin"}, //TODO
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Header", "X-Y-Header"},
		MaxAge:           3600 * time.Second,
		AllowHeaders:     []string{"X-Header", "X-Y-Header"},
	})

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		if route.NeedsAuthentication {
			handler = authMiddleware(handler)
		}

		router.Add(
			route.Method,
			route.Pattern,
			handler.ServeHTTP)
	}

	return router
}
