// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"

	"golang.custom.org/x/tools/core/event"
	"golang.custom.org/x/tools/gopls/core/file"
	"golang.custom.org/x/tools/gopls/core/golang"
	"golang.custom.org/x/tools/gopls/core/label"
	"golang.custom.org/x/tools/gopls/core/mod"
	"golang.custom.org/x/tools/gopls/core/protocol"
)

func (s *server) InlayHint(ctx context.Context, params *protocol.InlayHintParams) ([]protocol.InlayHint, error) {
	ctx, done := event.Start(ctx, "lsp.Server.inlayHint", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()

	switch snapshot.FileKind(fh) {
	case file.Mod:
		return mod.InlayHint(ctx, snapshot, fh, params.Range)
	case file.Go:
		return golang.InlayHint(ctx, snapshot, fh, params.Range)
	}
	return nil, nil // empty result
}
