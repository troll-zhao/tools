// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/troll-zhao/tools/gopls/core/vulncheck/scan"
)

// vulncheck implements the vulncheck command.
// TODO(hakim): hide from the public.
type vulncheck struct {
	app *Application
}

func (v *vulncheck) Name() string   { return "vulncheck" }
func (v *vulncheck) Parent() string { return v.app.Name() }
func (v *vulncheck) Usage() string  { return "" }
func (v *vulncheck) ShortHelp() string {
	return "run vulncheck analysis (core-use only)"
}
func (v *vulncheck) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), `
	WARNING: this command is for core-use only.

	By default, the command outputs a JSON-encoded
	github.com/troll-zhao/tools/gopls/core/protocol/command.VulncheckResult
	message.
	Example:
	$ gopls vulncheck <packages>

`)
}

func (v *vulncheck) Run(ctx context.Context, args ...string) error {
	if err := scan.Main(ctx, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}
