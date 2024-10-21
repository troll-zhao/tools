// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The unusedparams command runs the unusedparams analyzer.
package main

import (
	"golang.custom.org/x/tools/gopls/core/analysis/unusedparams"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(unusedparams.Analyzer) }
