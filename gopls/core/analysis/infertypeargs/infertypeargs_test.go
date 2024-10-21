// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package infertypeargs_test

import (
	"testing"

	"golang.custom.org/x/tools/gopls/core/analysis/infertypeargs"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, infertypeargs.Analyzer, "a")
}
