// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nonewvars_test

import (
	"testing"

	"golang.custom.org/x/tools/gopls/core/analysis/nonewvars"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, nonewvars.Analyzer, "a", "typeparams")
}
