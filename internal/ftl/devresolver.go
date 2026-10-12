package ftl

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// DPS is a NetLinx device:port:system address triple.
type DPS struct {
	Device int
	Port   int
	System int
}

// deviceRegistry maps symbolic variable names to their DPS addresses as
// extracted from NetLinx source files.
type deviceRegistry map[string]DPS

// customAddrRx matches a DevAddr in the form  "Custom [D:P:S]".
var customAddrRx = regexp.MustCompile(`^Custom\s*\[(\d+):(\d+):(\d+)\]$`)

// deviceDefRx matches one device definition line inside a DEFINE_DEVICE block:
//
//	dvName = D:P:S
var deviceDefRx = regexp.MustCompile(`(?i)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(\d+):(\d+):(\d+)\s*(?:$|;)`)

// includeRx matches a NetLinx  #INCLUDE  directive.
var includeRx = regexp.MustCompile(`(?i)^\s*#\s*INCLUDE\s+'([^']+)'`)

// sectionStartRx matches any NetLinx keyword that ends a DEFINE_DEVICE block.
var sectionStartRx = regexp.MustCompile(
	`(?i)^\s*(DEFINE_CONSTANT|DEFINE_TYPE|DEFINE_VARIABLE|DEFINE_LATCHING|` +
		`DEFINE_MUTUALLY_EXCLUSIVE|DEFINE_CALL|DEFINE_START|DEFINE_EVENT|` +
		`DEFINE_FUNCTION|DEFINE_PROGRAM|PROGRAM_NAME|DEFINE_MODULE|MODULE_NAME)\b`,
)

const maxIncludeDepth = 10

// resolveDevAddr resolves a single DevAddr string to a DPS triple.
//
//   - "Custom [D:P:S]" is parsed directly without any source-file lookup.
//   - Any other string is treated as a symbolic variable name and looked up in
//     the registry; if not found, the second return value is false.
func resolveDevAddr(devAddr string, reg deviceRegistry) (DPS, bool) {
	if m := customAddrRx.FindStringSubmatch(strings.TrimSpace(devAddr)); m != nil {
		d, _ := strconv.Atoi(m[1])
		p, _ := strconv.Atoi(m[2])
		s, _ := strconv.Atoi(m[3])
		return DPS{d, p, s}, true
	}

	if reg != nil {
		if dps, ok := reg[devAddr]; ok {
			return dps, true
		}
	}

	return DPS{}, false
}

// buildDeviceRegistry scans the system's MasterSrc file and all transitively
// #INCLUDEd .axi/.axs files for DEFINE_DEVICE entries and returns the
// resulting registry.  Returns an empty (non-nil) registry if no MasterSrc is
// present or the file cannot be read.
func buildDeviceRegistry(system *apw.System, apwDir string) deviceRegistry {
	reg := make(deviceRegistry)
	visited := make(map[string]bool)

	for _, fr := range system.Files {
		if fr.Type == apw.FileTypeMasterSrc {
			path := filepath.Join(apwDir, normalisePath(fr.FilePathName))
			scanNetLinxFile(path, apwDir, reg, visited, 0)
			break // only one MasterSrc per system
		}
	}

	return reg
}

// scanNetLinxFile parses a single NetLinx source or include file for
// DEFINE_DEVICE entries and recurses into any #INCLUDE files it discovers.
func scanNetLinxFile(path, apwDir string, reg deviceRegistry, visited map[string]bool, depth int) {
	if depth > maxIncludeDepth {
		return
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}

	if visited[abs] {
		return
	}

	visited[abs] = true

	data, err := os.ReadFile(abs)
	if err != nil {
		return
	}

	src := stripNetLinxBlockComments(string(data))
	inDefineDevice := false

	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		line := scanner.Text()

		// Strip inline // line comments.
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// #INCLUDE is valid anywhere in the file.
		if m := includeRx.FindStringSubmatch(trimmed); m != nil {
			incPath := resolveIncludePath(m[1], apwDir)
			scanNetLinxFile(incPath, apwDir, reg, visited, depth+1)
			continue
		}

		// Detect DEFINE_DEVICE section start.
		if strings.EqualFold(trimmed, "DEFINE_DEVICE") {
			inDefineDevice = true
			continue
		}

		// Any other top-level section keyword ends DEFINE_DEVICE.
		if inDefineDevice && sectionStartRx.MatchString(trimmed) {
			inDefineDevice = false
			continue
		}

		if inDefineDevice {
			if m := deviceDefRx.FindStringSubmatch(trimmed); m != nil {
				d, _ := strconv.Atoi(m[2])
				p, _ := strconv.Atoi(m[3])
				s, _ := strconv.Atoi(m[4])
				// First definition wins; don't overwrite if already present.
				if _, exists := reg[m[1]]; !exists {
					reg[m[1]] = DPS{d, p, s}
				}
			}
		}
	}
}

// stripNetLinxBlockComments removes (* ... *) and /* ... */ block comments
// (which may span multiple lines) from NetLinx source text.  Newlines inside
// removed blocks are preserved so that line-based scanning stays accurate.
func stripNetLinxBlockComments(src string) string {
	src = stripDelimited(src, "(*", "*)")
	src = stripDelimited(src, "/*", "*/")
	return src
}

func stripDelimited(src, open, close string) string {
	var b strings.Builder

	for {
		start := strings.Index(src, open)
		if start < 0 {
			b.WriteString(src)
			break
		}

		b.WriteString(src[:start])
		rest := src[start+len(open):]
		end := strings.Index(rest, close)

		if end < 0 {
			// Unterminated comment — discard the rest of the file.
			break
		}

		// Preserve newlines inside the comment block so line numbers stay
		// consistent for any subsequent diagnostic output.
		for _, c := range rest[:end] {
			if c == '\n' {
				b.WriteByte('\n')
			}
		}

		src = rest[end+len(close):]
	}

	return b.String()
}

// resolveIncludePath returns the absolute path of a #INCLUDE target, searching:
//  1. As a path relative to the APW directory.
//  2. As a bare filename inside the Include/ subdirectory.
//  3. As a bare filename inside the Source/ subdirectory.
//
// The first candidate that exists on disk is returned; otherwise candidate [0]
// is returned and scanNetLinxFile will fail gracefully on os.ReadFile.
func resolveIncludePath(includePath, apwDir string) string {
	cleaned := normalisePath(includePath)

	candidates := []string{
		filepath.Join(apwDir, cleaned),
		filepath.Join(apwDir, "Include", filepath.Base(cleaned)),
		filepath.Join(apwDir, "Source", filepath.Base(cleaned)),
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	return candidates[0]
}
