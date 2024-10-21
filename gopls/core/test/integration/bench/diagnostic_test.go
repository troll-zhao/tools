// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bench

import (
	"testing"

	"golang.custom.org/x/tools/gopls/core/protocol"
	. "golang.custom.org/x/tools/gopls/core/test/integration"
	"golang.custom.org/x/tools/gopls/core/test/integration/fake"
)

// BenchmarkDiagnosePackageFiles measures how long it takes to request
// diagnostics for 10 files in a single package, following a change to that
// package.
//
// This can be used to measure the efficiency of pull diagnostics
// (golang/go#53275).
func BenchmarkDiagnosePackageFiles(b *testing.B) {
	if testing.Short() {
		b.Skip("pull diagnostics are not supported by the benchmark dashboard baseline")
	}

	env := getRepo(b, "kubernetes").newEnv(b, fake.EditorConfig{
		Settings: map[string]any{
			"pullDiagnostics": true, // currently required for pull diagnostic support
		},
	}, "diagnosePackageFiles", false)

	// 10 arbitrary files in a single package.
	files := []string{
		"pkg/kubelet/active_deadline.go",      // 98 lines
		"pkg/kubelet/active_deadline_test.go", // 95 lines
		"pkg/kubelet/kubelet.go",              // 2439 lines
		"pkg/kubelet/kubelet_pods.go",         // 2061 lines
		"pkg/kubelet/kubelet_network.go",      // 70 lines
		"pkg/kubelet/kubelet_network_test.go", // 46 lines
		"pkg/kubelet/pod_workers.go",          // 1323 lines
		"pkg/kubelet/pod_workers_test.go",     // 1758 lines
		"pkg/kubelet/runonce.go",              // 175 lines
		"pkg/kubelet/volume_host.go",          // 297 lines
	}

	env.Await(InitialWorkspaceLoad)

	for _, file := range files {
		env.OpenFile(file)
	}

	env.AfterChange()

	edit := makeEditFunc(env, files[0])

	if stopAndRecord := startProfileIfSupported(b, env, qualifiedName("kubernetes", "diagnosePackageFiles")); stopAndRecord != nil {
		defer stopAndRecord()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		edit()
		var diags []protocol.Diagnostic
		for _, file := range files {
			fileDiags := env.Diagnostics(file)
			for _, d := range fileDiags {
				if d.Severity == protocol.SeverityError {
					diags = append(diags, d)
				}
			}
		}
		if len(diags) != 0 {
			b.Fatalf("got %d error diagnostics, want 0\ndiagnostics:\n%v", len(diags), diags)
		}
	}
}
