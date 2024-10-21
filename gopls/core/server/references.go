// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"

	"github.com/troll-zhao/tools/core/event"
	"github.com/troll-zhao/tools/gopls/core/file"
	"github.com/troll-zhao/tools/gopls/core/golang"
	"github.com/troll-zhao/tools/gopls/core/label"
	"github.com/troll-zhao/tools/gopls/core/protocol"
	"github.com/troll-zhao/tools/gopls/core/telemetry"
	"github.com/troll-zhao/tools/gopls/core/template"
)

func (s *server) References(ctx context.Context, params *protocol.ReferenceParams) (_ []protocol.Location, rerr error) {
	recordLatency := telemetry.StartLatencyTimer("references")
	defer func() {
		recordLatency(ctx, rerr)
	}()

	ctx, done := event.Start(ctx, "lsp.Server.references", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()
	switch snapshot.FileKind(fh) {
	case file.Tmpl:
		return template.References(ctx, snapshot, fh, params)
	case file.Go:
		return golang.References(ctx, snapshot, fh, params.Position, params.Context.IncludeDeclaration)
	}
	return nil, nil // empty result
}
