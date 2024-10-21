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
	"golang.custom.org/x/tools/gopls/core/template"
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
