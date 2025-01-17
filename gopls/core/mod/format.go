// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"context"

	"github.com/troll-zhao/tools/core/diff"
	"github.com/troll-zhao/tools/core/event"
	"github.com/troll-zhao/tools/gopls/core/cache"
	"github.com/troll-zhao/tools/gopls/core/file"
	"github.com/troll-zhao/tools/gopls/core/protocol"
)

func Format(ctx context.Context, snapshot *cache.Snapshot, fh file.Handle) ([]protocol.TextEdit, error) {
	ctx, done := event.Start(ctx, "mod.Format")
	defer done()

	pm, err := snapshot.ParseMod(ctx, fh)
	if err != nil {
		return nil, err
	}
	formatted, err := pm.File.Format()
	if err != nil {
		return nil, err
	}
	// Calculate the edits to be made due to the change.
	diffs := diff.Bytes(pm.Mapper.Content, formatted)
	return protocol.EditsFromDiffEdits(pm.Mapper, diffs)
}
