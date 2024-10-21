// Copyright 2018 The Go Authors. All rights reserved.
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
	"golang.custom.org/x/tools/gopls/core/work"
)

func (s *server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	ctx, done := event.Start(ctx, "lsp.Server.formatting", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()

	switch snapshot.FileKind(fh) {
	case file.Mod:
		return mod.Format(ctx, snapshot, fh)
	case file.Go:
		return golang.Format(ctx, snapshot, fh)
	case file.Work:
		return work.Format(ctx, snapshot, fh)
	}
	return nil, nil // empty result
}
