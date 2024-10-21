// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/troll-zhao/tools/core/event"
	"github.com/troll-zhao/tools/core/gocommand"
	"github.com/troll-zhao/tools/gopls/core/cache"
	"github.com/troll-zhao/tools/gopls/core/cache/metadata"
	"github.com/troll-zhao/tools/gopls/core/protocol"
	"github.com/troll-zhao/tools/gopls/core/settings"
)

// GCOptimizationDetails invokes the Go compiler on the specified
// package and reports its log of optimizations decisions as a set of
// diagnostics.
//
// TODO(adonovan): this feature needs more consistent and informative naming.
// Now that the compiler is cmd/compile, "GC" now means only "garbage collection".
// I propose "(Toggle|Display) Go compiler optimization details" in the UI,
// and CompilerOptimizationDetails for this function and compileropts.go for the file.
func GCOptimizationDetails(ctx context.Context, snapshot *cache.Snapshot, mp *metadata.Package) (map[protocol.DocumentURI][]*cache.Diagnostic, error) {
	if len(mp.CompiledGoFiles) == 0 {
		return nil, nil
	}
	pkgDir := filepath.Dir(mp.CompiledGoFiles[0].Path())
	outDir, err := os.MkdirTemp("", fmt.Sprintf("gopls-%d.details", os.Getpid()))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(outDir); err != nil {
			event.Error(ctx, "cleaning gcdetails dir", err)
		}
	}()

	tmpFile, err := os.CreateTemp(os.TempDir(), "gopls-x")
	if err != nil {
		return nil, err
	}
	tmpFile.Close() // ignore error
	defer os.Remove(tmpFile.Name())

	outDirURI := protocol.URIFromPath(outDir)
	// GC details doesn't handle Windows URIs in the form of "file:///C:/...",
	// so rewrite them to "file://C:/...". See golang/go#41614.
	if !strings.HasPrefix(outDir, "/") {
		outDirURI = protocol.DocumentURI(strings.Replace(string(outDirURI), "file:///", "file://", 1))
	}
	inv, cleanupInvocation, err := snapshot.GoCommandInvocation(false, &gocommand.Invocation{
		Verb: "build",
		Args: []string{
			fmt.Sprintf("-gcflags=-json=0,%s", outDirURI),
			fmt.Sprintf("-o=%s", tmpFile.Name()),
			".",
		},
		WorkingDir: pkgDir,
	})
	if err != nil {
		return nil, err
	}
	defer cleanupInvocation()
	_, err = snapshot.View().GoCommandRunner().Run(ctx, *inv)
	if err != nil {
		return nil, err
	}
	files, err := findJSONFiles(outDir)
	if err != nil {
		return nil, err
	}
	reports := make(map[protocol.DocumentURI][]*cache.Diagnostic)
	opts := snapshot.Options()
	var parseError error
	for _, fn := range files {
		uri, diagnostics, err := parseDetailsFile(fn, opts)
		if err != nil {
			// expect errors for all the files, save 1
			parseError = err
		}
		fh := snapshot.FindFile(uri)
		if fh == nil {
			continue
		}
		if pkgDir != filepath.Dir(fh.URI().Path()) {
			// https://github.com/golang/go/issues/42198
			// sometimes the detail diagnostics generated for files
			// outside the package can never be taken back.
			continue
		}
		reports[fh.URI()] = diagnostics
	}
	return reports, parseError
}

func parseDetailsFile(filename string, options *settings.Options) (protocol.DocumentURI, []*cache.Diagnostic, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return "", nil, err
	}
	var (
		uri         protocol.DocumentURI
		i           int
		diagnostics []*cache.Diagnostic
	)
	type metadata struct {
		File string `json:"file,omitempty"`
	}
	for dec := json.NewDecoder(bytes.NewReader(buf)); dec.More(); {
		// The first element always contains metadata.
		if i == 0 {
			i++
			m := new(metadata)
			if err := dec.Decode(m); err != nil {
				return "", nil, err
			}
			if !strings.HasSuffix(m.File, ".go") {
				continue // <autogenerated>
			}
			uri = protocol.URIFromPath(m.File)
			continue
		}
		d := new(protocol.Diagnostic)
		if err := dec.Decode(d); err != nil {
			return "", nil, err
		}
		d.Tags = []protocol.DiagnosticTag{} // must be an actual slice
		msg := d.Code.(string)
		if msg != "" {
			msg = fmt.Sprintf("%s(%s)", msg, d.Message)
		}
		if !showDiagnostic(msg, d.Source, options) {
			continue
		}
		var related []protocol.DiagnosticRelatedInformation
		for _, ri := range d.RelatedInformation {
			// TODO(rfindley): The compiler uses LSP-like JSON to encode gc details,
			// however the positions it uses are 1-based UTF-8:
			// https://github.com/golang/go/blob/master/src/cmd/compile/internal/logopt/log_opts.go
			//
			// Here, we adjust for 0-based positions, but do not translate UTF-8 to UTF-16.
			related = append(related, protocol.DiagnosticRelatedInformation{
				Location: protocol.Location{
					URI:   ri.Location.URI,
					Range: zeroIndexedRange(ri.Location.Range),
				},
				Message: ri.Message,
			})
		}
		diagnostic := &cache.Diagnostic{
			URI:      uri,
			Range:    zeroIndexedRange(d.Range),
			Message:  msg,
			Severity: d.Severity,
			Source:   cache.OptimizationDetailsError, // d.Source is always "go compiler" as of 1.16, use our own
			Tags:     d.Tags,
			Related:  related,
		}
		diagnostics = append(diagnostics, diagnostic)
		i++
	}
	return uri, diagnostics, nil
}

// showDiagnostic reports whether a given diagnostic should be shown to the end
// user, given the current options.
func showDiagnostic(msg, source string, o *settings.Options) bool {
	if source != "go compiler" {
		return false
	}
	if o.Annotations == nil {
		return true
	}
	switch {
	case strings.HasPrefix(msg, "canInline") ||
		strings.HasPrefix(msg, "cannotInline") ||
		strings.HasPrefix(msg, "inlineCall"):
		return o.Annotations[settings.Inline]
	case strings.HasPrefix(msg, "escape") || msg == "leak":
		return o.Annotations[settings.Escape]
	case strings.HasPrefix(msg, "nilcheck"):
		return o.Annotations[settings.Nil]
	case strings.HasPrefix(msg, "isInBounds") ||
		strings.HasPrefix(msg, "isSliceInBounds"):
		return o.Annotations[settings.Bounds]
	}
	return false
}

// The range produced by the compiler is 1-indexed, so subtract range by 1.
func zeroIndexedRange(rng protocol.Range) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      rng.Start.Line - 1,
			Character: rng.Start.Character - 1,
		},
		End: protocol.Position{
			Line:      rng.End.Line - 1,
			Character: rng.End.Character - 1,
		},
	}
}

func findJSONFiles(dir string) ([]string, error) {
	ans := []string{}
	f := func(path string, fi os.FileInfo, _ error) error {
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".json") {
			ans = append(ans, path)
		}
		return nil
	}
	err := filepath.Walk(dir, f)
	return ans, err
}
