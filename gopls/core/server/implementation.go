// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"

	"golang.custom.org/x/tools/core/event"
	"golang.custom.org/x/tools/gopls/core/file"
	"golang.custom.org/x/tools/gopls/core/golang"
	"golang.custom.org/x/tools/gopls/core/label"
	"golang.custom.org/x/tools/gopls/core/protocol"
	"golang.custom.org/x/tools/gopls/core/telemetry"
)

func (s *server) Implementation(ctx context.Context, params *protocol.ImplementationParams) (_ []protocol.Location, rerr error) {
	recordLatency := telemetry.StartLatencyTimer("implementation")
	defer func() {
		recordLatency(ctx, rerr)
	}()

	ctx, done := event.Start(ctx, "lsp.Server.implementation", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()
	if snapshot.FileKind(fh) != file.Go {
		return nil, nil // empty result
	}
	return golang.Implementation(ctx, snapshot, fh, params.Position)
}
