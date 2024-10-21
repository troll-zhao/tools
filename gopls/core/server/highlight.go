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
	"github.com/troll-zhao/tools/gopls/core/template"
)

func (s *server) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	ctx, done := event.Start(ctx, "lsp.Server.documentHighlight", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()

	switch snapshot.FileKind(fh) {
	case file.Tmpl:
		return template.Highlight(ctx, snapshot, fh, params.Position)
	case file.Go:
		rngs, err := golang.Highlight(ctx, snapshot, fh, params.Position)
		if err != nil {
			event.Error(ctx, "no highlight", err)
		}
		return rngs, nil
	}
	return nil, nil // empty result
}
