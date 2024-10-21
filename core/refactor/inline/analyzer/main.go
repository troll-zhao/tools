// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// The inline command applies the inliner to the specified packages of
// Go source code. Run with:
//
//	$ go run ./core/refactor/inline/analyzer/main.go -fix packages...
package main

import (
	inlineanalyzer "golang.custom.org/x/tools/core/refactor/inline/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(inlineanalyzer.Analyzer) }
