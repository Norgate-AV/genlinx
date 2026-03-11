package cfg

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// updateGolden overwrites golden files with current output when passed as a
// flag: go test ./internal/cfg/... -run TestCfgGolden -update
var updateGolden = flag.Bool("update", false, "overwrite golden files with current output")

const (
	cfgGoldenDir   = "testdata/golden"
	cfgFixturesDir = "testdata/golden/fixtures"
	cfgFixtureAPW  = "CfgWorkspace.apw"
)

// setupCfgGoldenWorkspace copies the fixture directory into a temp dir and
// parses the APW. It returns the parsed APW and the absolute path of the temp
// copy so the caller can normalise away the volatile temp directory prefix.
func setupCfgGoldenWorkspace(t *testing.T) (*apw.APW, string) {
	t.Helper()

	dst := t.TempDir()
	require.NoError(t, copyCfgFixtures(cfgFixturesDir, dst))

	apwPath := filepath.Join(dst, cfgFixtureAPW)
	data, err := os.ReadFile(apwPath)
	require.NoError(t, err)

	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	return a, dst
}

func copyCfgFixtures(src, dst string) error {
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

// normaliseCfgOutput makes the Build() output stable for golden comparison by:
//  1. Replacing the absolute fixture directory prefix in AXSFILE lines with
//     ${FIXTURES} so the golden files are portable across machines.
//  2. Normalising all backslashes to forward slashes.
func normaliseCfgOutput(output, fixturesDir string) string {
	// Normalise the fixtures path itself to forward slashes so both the
	// replacement key and the output use the same separator.
	fixturesNorm := filepath.ToSlash(fixturesDir)
	output = strings.ReplaceAll(output, fixturesNorm, "${FIXTURES}")

	// Also handle the backslash form in case the OS produced it.
	output = strings.ReplaceAll(output, fixturesDir, "${FIXTURES}")

	// Normalise all remaining backslashes to forward slashes.
	output = strings.ReplaceAll(output, "\\", "/")

	return output
}

// checkCfgGolden compares got against the golden file at goldenPath.
// When -update is set the file is overwritten instead.
func checkCfgGolden(t *testing.T, goldenPath, got string) {
	t.Helper()

	if *updateGolden {
		require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0o644))
		t.Logf("updated %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	require.NoError(t, err,
		"golden file not found — regenerate with: go test -run %s -update", t.Name())
	// Normalise line endings so the comparison is CRLF-agnostic.
	wantNorm := strings.ReplaceAll(string(want), "\r\n", "\n")
	require.Equal(t, wantNorm, got,
		"cfg output differs from golden — regenerate with: go test -run %s -update", t.Name())
}

// ---------------------------------------------------------------------------
// Golden tests
// ---------------------------------------------------------------------------

func TestCfgGolden(t *testing.T) {
	cases := []struct {
		name       string
		goldenFile string
		opts       *Options
	}{
		{
			name:       "default",
			goldenFile: "default.txt",
			opts: &Options{
				OutputFileSuffix:          "build.cfg",
				OutputLogFileSuffix:       "build.log",
				OutputLogFileOption:       "N",
				OutputLogConsoleOption:    true,
				BuildWithDebugInformation: false,
				BuildWithSource:           false,
			},
		},
		{
			name:       "debug-and-source",
			goldenFile: "debug-and-source.txt",
			opts: &Options{
				OutputFileSuffix:          "build.cfg",
				OutputLogFileSuffix:       "build.log",
				OutputLogFileOption:       "N",
				OutputLogConsoleOption:    true,
				BuildWithDebugInformation: true,
				BuildWithSource:           true,
			},
		},
		{
			name:       "with-paths",
			goldenFile: "with-paths.txt",
			opts: &Options{
				OutputFileSuffix:          "build.cfg",
				OutputLogFileSuffix:       "build.log",
				OutputLogFileOption:       "A",
				OutputLogConsoleOption:    false,
				BuildWithDebugInformation: false,
				BuildWithSource:           false,
				IncludePath:               []string{`C:\AMX\AXIs`},
				ModulePath:                []string{`C:\AMX\Modules`},
				LibraryPath:               []string{`C:\AMX\Duet`},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			goldenPath, err := filepath.Abs(filepath.Join(cfgGoldenDir, tc.goldenFile))
			require.NoError(t, err)

			a, fixturesDir := setupCfgGoldenWorkspace(t)
			output := NewBuilder(a, tc.opts).Build()
			normalised := normaliseCfgOutput(output, fixturesDir)

			checkCfgGolden(t, goldenPath, normalised)
		})
	}
}
