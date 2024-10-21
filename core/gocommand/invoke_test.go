// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocommand_test

import (
	"context"
	"testing"

	"golang.custom.org/x/tools/core/gocommand"
	"golang.custom.org/x/tools/core/testenv"
)

func TestGoVersion(t *testing.T) {
	testenv.NeedsTool(t, "go")

	inv := gocommand.Invocation{
		Verb: "version",
	}
	gocmdRunner := &gocommand.Runner{}
	if _, err := gocmdRunner.Run(context.Background(), inv); err != nil {
		t.Error(err)
	}
}
