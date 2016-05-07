// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
