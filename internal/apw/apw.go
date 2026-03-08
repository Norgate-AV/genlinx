package apw

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// APW represents a parsed AMX NetLinx workspace (.apw) file.
type APW struct {
	id    string
	name  string
	path  string
	dir   string
	files map[string]File
	ws    *Workspace
}

// Parse validates path and data, then constructs a fully populated APW.
// The path must have a .apw extension. The data must be the contents of that
// file — the caller is responsible for reading it.
func Parse(path string, data []byte) (*APW, error) {
	if strings.ToLower(filepath.Ext(path)) != FileExtensionAPW {
		return nil, fmt.Errorf("not a NetLinx Workspace file: %s", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	if !bytes.Contains(data, []byte("<!DOCTYPE Workspace [")) {
		return nil, fmt.Errorf("not a NetLinx Workspace file: %s", absPath)
	}

	clean, err := stripDTD(data)
	if err != nil {
		return nil, err
	}

	ws := &Workspace{}
	if err := xml.Unmarshal(clean, ws); err != nil {
		return nil, fmt.Errorf("failed to parse APW file: %w", err)
	}

	name := filepath.Base(absPath)

	a := &APW{
		path:  absPath,
		name:  name,
		id:    strings.TrimSuffix(name, filepath.Ext(name)),
		dir:   filepath.Dir(absPath),
		files: make(map[string]File),
		ws:    ws,
	}

	a.buildFiles()

	return a, nil
}

// stripDTD removes the inline <!DOCTYPE ... ]> block so that encoding/xml
// can parse the remaining document without error.
func stripDTD(data []byte) ([]byte, error) {
	start := bytes.Index(data, []byte("<!DOCTYPE"))
	if start == -1 {
		return data, nil
	}

	rest := data[start:]
	_, after, ok := bytes.Cut(rest, []byte("]>"))

	if !ok {
		return nil, fmt.Errorf("malformed APW file: no closing ]> for DOCTYPE")
	}

	return append(data[:start], after...), nil
}

// buildFiles walks the parsed Workspace tree and populates a.files with
// absolute path → File entries. Duplicate paths are deduplicated by nature of the map.
func (a *APW) buildFiles() {
	for _, project := range a.ws.Projects {
		for _, system := range project.Systems {
			for _, fr := range system.Files {
				relPath := filepath.FromSlash(strings.ReplaceAll(fr.FilePathName, `\`, `/`))
				absPath := filepath.Join(a.dir, relPath)

				_, statErr := os.Stat(absPath)

				a.files[absPath] = File{
					ID:      fr.Identifier,
					Type:    fr.Type,
					Path:    absPath,
					Exists:  statErr == nil,
					IsExtra: false,
				}
			}
		}
	}
}

// FileIsReadable reports whether the file can be scanned for #include / define_module references.
func FileIsReadable(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == FileExtensionAXS || ext == FileExtensionAXI
}

// FileIsOfInterest reports whether the file should be considered as a candidate
// during an extra-file disk search (.axs, .axi, .jar, .xdd).
func FileIsOfInterest(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == FileExtensionAXS || ext == FileExtensionAXI || ext == FileExtensionJAR || ext == FileExtensionXDD
}

// GetFileType returns the FileType constant for a file path based on its extension.
func GetFileType(file string) FileType {
	ext := strings.ToLower(filepath.Ext(file))

	for fileType, fileExt := range AmxExtensions {
		if fileExt != "" && fileExt == ext {
			return fileType
		}
	}

	return FileTypeOther
}

// isInWorkspace reports whether the given id/name already appears in the workspace.
func (a *APW) isInWorkspace(id string) bool {
	for _, f := range a.files {
		if f.ID == id || strings.Contains(f.Path, id) {
			return true
		}
	}

	return false
}

func (a *APW) searchForExtraFileReferences(file string, pattern *regexp.Regexp) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var refs []string

	for _, match := range pattern.FindAllStringSubmatch(string(data), -1) {
		if len(match) < 2 {
			continue
		}

		id := match[1]

		if a.isInWorkspace(id) {
			continue
		}

		if !seen[id] {
			refs = append(refs, id)
			seen[id] = true
		}
	}

	return refs, nil
}

// GetExtraFileReferencesFromFile returns the unique extra file IDs referenced
// via #include and define_module directives in the given source file.
func (a *APW) GetExtraFileReferencesFromFile(file string) ([]string, error) {
	if !FileIsReadable(file) {
		return nil, nil
	}

	includePattern := regexp.MustCompile(`(?i)#(?:include)\s+'(.+)'`)
	modulePattern := regexp.MustCompile(`(?im)^(?:define_module)\s+'(.+)'`)

	includeRefs, err := a.searchForExtraFileReferences(file, includePattern)
	if err != nil {
		return nil, err
	}

	moduleRefs, err := a.searchForExtraFileReferences(file, modulePattern)
	if err != nil {
		return nil, err
	}

	combined := append(includeRefs, moduleRefs...)

	seen := make(map[string]bool)
	var unique []string

	for _, r := range combined {
		if !seen[r] {
			unique = append(unique, r)
			seen[r] = true
		}
	}

	return unique, nil
}

// GetExtraFileReferences returns all unique extra file IDs referenced across
// every existing file in the workspace.
func (a *APW) GetExtraFileReferences() ([]string, error) {
	seen := make(map[string]bool)
	var refs []string

	for _, f := range a.files {
		if !f.Exists {
			continue
		}

		fileRefs, err := a.GetExtraFileReferencesFromFile(f.Path)
		if err != nil {
			continue
		}

		for _, r := range fileRefs {
			if !seen[r] {
				refs = append(refs, r)
				seen[r] = true
			}
		}
	}

	return refs, nil
}

// AllFiles returns all unique workspace file references plus the workspace file itself,
// sorted by path.
func (a *APW) AllFiles() []File {
	files := make([]File, 0, len(a.files)+1)
	for _, f := range a.files {
		files = append(files, f)
	}

	files = append(files, File{
		ID:      a.id,
		Type:    FileTypeWorkspace,
		Path:    a.path,
		Exists:  true,
		IsExtra: false,
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files
}

// MasterSrcPath returns the directories that contain MasterSrc files.
func (a *APW) MasterSrcPath() []string {
	return a.getFileDirectories(a.filesByType(FileTypeMasterSrc))
}

// filesByType returns all files matching t, sorted by path.
func (a *APW) filesByType(t FileType) []File {
	var files []File

	for _, f := range a.files {
		if f.Type == t {
			files = append(files, f)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files
}

// ModuleFiles returns all unique file references with type Module.
func (a *APW) ModuleFiles() []File {
	return a.filesByType(FileTypeModule)
}

// MasterSrcFiles returns all unique file references with type MasterSrc.
func (a *APW) MasterSrcFiles() []File {
	return a.filesByType(FileTypeMasterSrc)
}

// getFileDirectories returns the unique directory paths for the given files.
// Paths in File are already absolute.
func (a *APW) getFileDirectories(files []File) []string {
	seen := make(map[string]bool)
	var dirs []string

	for _, f := range files {
		dir := filepath.Dir(f.Path)
		if !seen[dir] {
			dirs = append(dirs, dir)
			seen[dir] = true
		}
	}

	return dirs
}

// IncludePath returns the unique directories containing Include files.
func (a *APW) IncludePath() []string {
	return a.getFileDirectories(a.filesByType(FileTypeInclude))
}

// ModulePath returns the unique directories containing Module, Duet, and XDD files.
func (a *APW) ModulePath() []string {
	var files []File
	for _, f := range a.files {
		if f.Type == FileTypeModule || f.Type == FileTypeDuet || f.Type == FileTypeXDD {
			files = append(files, f)
		}
	}

	return a.getFileDirectories(files)
}

// ID returns the workspace identifier derived from the filename stem.
func (a *APW) ID() string {
	return a.id
}

// FilePath returns the absolute path of the workspace file.
func (a *APW) FilePath() string {
	return a.path
}
