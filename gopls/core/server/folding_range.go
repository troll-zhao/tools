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
)

func (s *server) FoldingRange(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	ctx, done := event.Start(ctx, "lsp.Server.foldingRange", label.URI.Of(params.TextDocument.URI))
	defer done()

	fh, snapshot, release, err := s.fileOf(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()
	if snapshot.FileKind(fh) != file.Go {
		return nil, nil // empty result
	}
	ranges, err := golang.FoldingRange(ctx, snapshot, fh, snapshot.Options().LineFoldingOnly)
	if err != nil {
		return nil, err
	}
	return toProtocolFoldingRanges(ranges)
}

func toProtocolFoldingRanges(ranges []*golang.FoldingRangeInfo) ([]protocol.FoldingRange, error) {
	result := make([]protocol.FoldingRange, 0, len(ranges))
	for _, info := range ranges {
		rng := info.MappedRange.Range()
		result = append(result, protocol.FoldingRange{
			StartLine:      rng.Start.Line,
			StartCharacter: rng.Start.Character,
			EndLine:        rng.End.Line,
			EndCharacter:   rng.End.Character,
			Kind:           string(info.Kind),
		})
	}
	return result, nil
}
