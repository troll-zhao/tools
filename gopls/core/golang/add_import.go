// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"context"

	"github.com/troll-zhao/tools/core/imports"
	"github.com/troll-zhao/tools/gopls/core/cache"
	"github.com/troll-zhao/tools/gopls/core/cache/parsego"
	"github.com/troll-zhao/tools/gopls/core/file"
	"github.com/troll-zhao/tools/gopls/core/protocol"
)

// AddImport adds a single import statement to the given file
func AddImport(ctx context.Context, snapshot *cache.Snapshot, fh file.Handle, importPath string) ([]protocol.TextEdit, error) {
	pgf, err := snapshot.ParseGo(ctx, fh, parsego.Full)
	if err != nil {
		return nil, err
	}
	return ComputeOneImportFixEdits(snapshot, pgf, &imports.ImportFix{
		StmtInfo: imports.ImportInfo{
			ImportPath: importPath,
		},
		FixType: imports.AddImport,
	})
}
