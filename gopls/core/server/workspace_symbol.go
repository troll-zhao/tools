// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"

	"github.com/troll-zhao/tools/core/event"
	"github.com/troll-zhao/tools/gopls/core/cache"
	"github.com/troll-zhao/tools/gopls/core/golang"
	"github.com/troll-zhao/tools/gopls/core/protocol"
	"github.com/troll-zhao/tools/gopls/core/telemetry"
)

func (s *server) Symbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) (_ []protocol.SymbolInformation, rerr error) {
	recordLatency := telemetry.StartLatencyTimer("symbol")
	defer func() {
		recordLatency(ctx, rerr)
	}()

	ctx, done := event.Start(ctx, "lsp.Server.symbol")
	defer done()

	views := s.session.Views()
	matcher := s.Options().SymbolMatcher
	style := s.Options().SymbolStyle

	var snapshots []*cache.Snapshot
	for _, v := range views {
		snapshot, release, err := v.Snapshot()
		if err != nil {
			continue // snapshot is shutting down
		}
		// If err is non-nil, the snapshot is shutting down. Skip it.
		defer release()
		snapshots = append(snapshots, snapshot)
	}
	return golang.WorkspaceSymbols(ctx, matcher, style, snapshots, params.Query)
}
