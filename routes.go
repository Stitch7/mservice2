// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

var routes = Routes{
	Route{
		"Index",
		"GET",
		false,
		"/",
		Index,
	},
	Route{
		"BoardIndex",
		"GET",
		true,
		"/test-login",
		LoginTest,
	},
	Route{
		"BoardIndex",
		"GET",
		false,
		"/boards",
		BoardIndex,
	},
	Route{
		"BoardIndex",
		"GET",
		false,
		"/board/:boardId/threads",
		ThreadIndex,
	},
}
