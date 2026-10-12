package archive

import (
	"archive/zip"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// updateGolden overwrites the golden files with the current output when
// passed as a flag: go test ./internal/archive/... -run TestArchiveGolden -update
var updateGolden = flag.Bool("update", false, "overwrite golden files with current output")

const (
	goldenDir   = "testdata/golden"
	fixturesDir = "testdata/golden/fixtures"
	goldenAPW   = "GoldenWorkspace.apw"
)

// setupGoldenWorkspace copies the fixture directory into a temp dir, parses
// the APW, and changes the working directory to the temp dir so that Build()
// writes its zip there.
func setupGoldenWorkspace(t *testing.T) *apw.APW {
	t.Helper()

	dst := t.TempDir()
	require.NoError(t, copyFixtures(fixturesDir, dst))

	apwPath := filepath.Join(dst, goldenAPW)
	data, err := os.ReadFile(apwPath)
	require.NoError(t, err)

	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	oldWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(dst))

	return a
}

// copyFixtures recursively copies the src fixture directory into dst.
func copyFixtures(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0o644)
	})
}

// zipEntryNames opens the named zip and returns its entry names, sorted.
func zipEntryNames(t *testing.T, zipPath string) []string {
	t.Helper()

	r, err := zip.OpenReader(zipPath)
	require.NoError(t, err)
	defer func() { _ = r.Close() }()

	names := make([]string, 0, len(r.File))
	for _, f := range r.File {
		names = append(names, f.Name)
	}

	sort.Strings(names)

	return names
}

// checkGolden compares got against the golden file at goldenPath.
// When -update is set the file is overwritten instead.
func checkGolden(t *testing.T, goldenPath string, got []string) {
	t.Helper()

	content := strings.Join(got, "\n") + "\n"

	if *updateGolden {
		require.NoError(t, os.WriteFile(goldenPath, []byte(content), 0o644))
		t.Logf("updated %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	require.NoError(t, err,
		"golden file not found — regenerate with: go test -run %s -update", t.Name())
	// Normalise line endings so the comparison is CRLF-agnostic (the golden
	// files are text files committed on Windows but the builder produces LF).
	wantNorm := strings.ReplaceAll(string(want), "\r\n", "\n")
	require.Equal(t, wantNorm, content,
		"zip entries differ from golden — regenerate with: go test -run %s -update", t.Name())
}

// ---------------------------------------------------------------------------
// Golden tests
// ---------------------------------------------------------------------------

func TestArchiveGolden(t *testing.T) {
	cases := []struct {
		name        string
		goldenFile  string
		archiveFile string
		projectID   string
		systemID    string
	}{
		{
			name:        "full-workspace",
			goldenFile:  "full.txt",
			archiveFile: "GoldenWorkspace.archive.zip",
		},
		{
			name:        "project-scoped",
			goldenFile:  "project-a.txt",
			archiveFile: "GoldenWorkspace-ProjectA.archive.zip",
			projectID:   "ProjectA",
		},
		{
			name:        "system-scoped",
			goldenFile:  "system-a1.txt",
			archiveFile: "GoldenWorkspace-ProjectA-SystemA1.archive.zip",
			projectID:   "ProjectA",
			systemID:    "SystemA1",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Resolve the golden path before setupGoldenWorkspace changes the wd.
			goldenPath, err := filepath.Abs(filepath.Join(goldenDir, tc.goldenFile))
			require.NoError(t, err)

			a := setupGoldenWorkspace(t)
			opts := &Options{
				OutputFileSuffix: "archive.zip",
				ProjectID:        tc.projectID,
				SystemID:         tc.systemID,
			}
			require.NoError(t, NewBuilder(a, opts).Build())

			entries := zipEntryNames(t, tc.archiveFile)
			checkGolden(t, goldenPath, entries)
		})
	}
}
