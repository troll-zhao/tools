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
	"github.com/troll-zhao/tools/gopls/core/mod"
	"github.com/troll-zhao/tools/gopls/core/protocol"
	"github.com/troll-zhao/tools/gopls/core/settings"
	"github.com/troll-zhao/tools/gopls/core/telemetry"
	"github.com/troll-zhao/tools/gopls/core/template"
	"github.com/troll-zhao/tools/gopls/core/work"
)

func (s *server) Hover(ctx context.Context, params *protocol.HoverParams) (_ *protocol.Hover, rerr error) {
	recordLatency := telemetry.StartLatencyTimer("hover")
	defer func() {
		recordLatency(ctx, rerr)
	}()

	ctx, done := event.Start(ctx, "lsp.Server.hover", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()

	switch snapshot.FileKind(fh) {
	case file.Mod:
		return mod.Hover(ctx, snapshot, fh, params.Position)
	case file.Go:
		var pkgURL func(path golang.PackagePath, fragment string) protocol.URI
		if snapshot.Options().LinksInHover == settings.LinksInHover_Gopls {
			web, err := s.getWeb()
			if err != nil {
				event.Error(ctx, "failed to start web server", err)
			} else {
				pkgURL = func(path golang.PackagePath, fragment string) protocol.URI {
					return web.PkgURL(snapshot.View().ID(), path, fragment)
				}
			}
		}
		return golang.Hover(ctx, snapshot, fh, params.Position, pkgURL)
	case file.Tmpl:
		return template.Hover(ctx, snapshot, fh, params.Position)
	case file.Work:
		return work.Hover(ctx, snapshot, fh, params.Position)
	}
	return nil, nil // empty result
}