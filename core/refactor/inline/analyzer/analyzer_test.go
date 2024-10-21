// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analyzer_test

import (
	"testing"

	inlineanalyzer "github.com/troll-zhao/tools/core/refactor/inline/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), inlineanalyzer.Analyzer, "a", "b")
}
