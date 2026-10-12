package apw

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ErrProjectNotFound is returned when a project identifier cannot be found in the workspace.
var ErrProjectNotFound = errors.New("project not found")

// ErrSystemNotFound is returned when a system identifier cannot be found in a project.
var ErrSystemNotFound = errors.New("system not found")

// ErrAmbiguousSystem is returned when a system identifier matches more than one
// project and the caller has not supplied a project to disambiguate.
var ErrAmbiguousSystem = errors.New("system name is ambiguous across multiple projects")

// SystemMatch records a project that contains a system with the queried identifier.
type SystemMatch struct {
	ProjectID string
}

// FindSystemAcrossProjects returns all projects that contain a system whose
// Identifier matches systemID. Use the slice length to detect ambiguity:
// 0 = not found, 1 = unique match, 2+ = ambiguous.
func (a *APW) FindSystemAcrossProjects(systemID string) []SystemMatch {
	var matches []SystemMatch

	for _, p := range a.ws.Projects {
		if _, ok := p.FindSystem(systemID); ok {
			matches = append(matches, SystemMatch{ProjectID: p.Identifier})
		}
	}

	return matches
}

// APW represents a parsed AMX NetLinx workspace (.apw) file.
type APW struct {
	id    string
	name  string
	path  string
	dir   string
	files map[string]File
	ws    *Workspace
}

func NewAPW(id, path string) *APW {
	ws := NewWorkspace(id)
	ws.CreateVersion = "4.0"
	ws.CurrentVersion = "4.0"

	absPath := path
	if path != "" {
		if abs, err := filepath.Abs(path); err == nil {
			absPath = abs
		}
	}

	return &APW{
		id:    id,
		name:  id + FileExtensionAPW,
		path:  absPath,
		dir:   filepath.Dir(absPath),
		files: make(map[string]File),
		ws:    ws,
	}
}

// Workspace returns the APW's inner Workspace so callers can build the
// project/system/file hierarchy directly.
func (a *APW) Workspace() *Workspace {
	return a.ws
}

// SetWorkspace replaces the APW's inner Workspace and rebuilds the internal
// files map from the new workspace tree.
func (a *APW) SetWorkspace(ws *Workspace) {
	a.ws = ws
	a.files = make(map[string]File)

	if ws != nil {
		a.buildFiles()
	}
}

// Rebuild re-derives the internal files map from the current workspace tree.
// Call this after modifying the workspace hierarchy programmatically so that
// AllFiles, ModuleFiles, IncludePath etc. reflect the latest state.
func (a *APW) Rebuild() {
	a.files = make(map[string]File)

	if a.ws != nil {
		a.buildFiles()
	}
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

// isInFileSet reports whether id appears in the given file set.
// Used by scope-aware extra-file scanning to avoid treating files
// that belong to the current archive scope as "extra".
func isInFileSet(id string, files []File) bool {
	for _, f := range files {
		if f.ID == id || strings.Contains(f.Path, id) {
			return true
		}
	}

	return false
}

func (a *APW) searchForExtraFileReferences(file string, pattern *regexp.Regexp, inScope func(string) bool) ([]string, error) {
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

		if inScope(id) {
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
// It excludes references whose IDs already appear anywhere in the workspace.
func (a *APW) GetExtraFileReferencesFromFile(file string) ([]string, error) {
	if !FileIsReadable(file) {
		return nil, nil
	}

	includePattern := regexp.MustCompile(`(?i)#(?:include)\s+'(.+)'`)
	modulePattern := regexp.MustCompile(`(?im)^(?:define_module)\s+'(.+)'`)

	includeRefs, err := a.searchForExtraFileReferences(file, includePattern, a.isInWorkspace)
	if err != nil {
		return nil, err
	}

	moduleRefs, err := a.searchForExtraFileReferences(file, modulePattern, a.isInWorkspace)
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

// GetExtraFileReferencesFromFileInScope is like GetExtraFileReferencesFromFile
// but only considers scopeFiles (the files belonging to the current archive
// scope) as already accounted for. References to files that exist in the
// workspace under a different scope are returned rather than skipped, so
// they can be discovered and included in a scoped archive.
func (a *APW) GetExtraFileReferencesFromFileInScope(file string, scopeFiles []File) ([]string, error) {
	if !FileIsReadable(file) {
		return nil, nil
	}

	includePattern := regexp.MustCompile(`(?i)#(?:include)\s+'(.+)'`)
	modulePattern := regexp.MustCompile(`(?im)^(?:define_module)\s+'(.+)'`)

	inScope := func(id string) bool { return isInFileSet(id, scopeFiles) }

	includeRefs, err := a.searchForExtraFileReferences(file, includePattern, inScope)
	if err != nil {
		return nil, err
	}

	moduleRefs, err := a.searchForExtraFileReferences(file, modulePattern, inScope)
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

// collectExtraRefsMap scans the given files for #include and define_module
// references not already in the scope, returning a map of unique ref ID →
// the path of the first source file in which it was found.
func (a *APW) collectExtraRefsMap(files []File) (map[string]string, error) {
	result := make(map[string]string)

	for _, f := range files {
		if !f.Exists {
			continue
		}

		fileRefs, err := a.GetExtraFileReferencesFromFileInScope(f.Path, files)
		if err != nil {
			continue
		}

		for _, r := range fileRefs {
			if _, seen := result[r]; !seen {
				result[r] = f.Path
			}
		}
	}

	return result, nil
}

// collectExtraRefs scans the given files for #include and define_module
// references not already in the scope, returning unique IDs.
// Only files within the provided slice are treated as "already in scope";
// workspace files outside this slice are returned as extra references so
// that scoped archives can discover and include them.
func (a *APW) collectExtraRefs(files []File) ([]string, error) {
	m, err := a.collectExtraRefsMap(files)
	if err != nil {
		return nil, err
	}

	refs := make([]string, 0, len(m))
	for k := range m {
		refs = append(refs, k)
	}

	return refs, nil
}

// GetExtraFileReferences returns all unique extra file IDs referenced across
// every existing file in the workspace.
func (a *APW) GetExtraFileReferences() ([]string, error) {
	files := make([]File, 0, len(a.files))
	for _, f := range a.files {
		files = append(files, f)
	}

	return a.collectExtraRefs(files)
}

// GetExtraFileReferencesMap returns a map of unique extra file ID → source
// file path for every existing file in the workspace.
func (a *APW) GetExtraFileReferencesMap() (map[string]string, error) {
	files := make([]File, 0, len(a.files))
	for _, f := range a.files {
		files = append(files, f)
	}

	return a.collectExtraRefsMap(files)
}

// filesForProject returns the File slice for the named project, walking the
// workspace tree. Returns ErrProjectNotFound if the project does not exist.
func (a *APW) filesForProject(projectID string) ([]File, error) {
	proj, ok := a.ws.FindProject(projectID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProjectNotFound, projectID)
	}

	var files []File

	for _, sys := range proj.Systems {
		for _, fr := range sys.Files {
			relPath := filepath.FromSlash(strings.ReplaceAll(fr.FilePathName, `\`, `/`))
			absPath := filepath.Join(a.dir, relPath)
			_, statErr := os.Stat(absPath)
			files = append(files, File{
				ID:     fr.Identifier,
				Type:   fr.Type,
				Path:   absPath,
				Exists: statErr == nil,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files, nil
}

// filesForSystem returns the File slice for the named system within the named
// project. Returns ErrProjectNotFound or ErrSystemNotFound as appropriate.
func (a *APW) filesForSystem(projectID, systemID string) ([]File, error) {
	proj, ok := a.ws.FindProject(projectID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProjectNotFound, projectID)
	}

	sys, ok := proj.FindSystem(systemID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrSystemNotFound, systemID)
	}

	var files []File

	for _, fr := range sys.Files {
		relPath := filepath.FromSlash(strings.ReplaceAll(fr.FilePathName, `\`, `/`))
		absPath := filepath.Join(a.dir, relPath)
		_, statErr := os.Stat(absPath)
		files = append(files, File{
			ID:     fr.Identifier,
			Type:   fr.Type,
			Path:   absPath,
			Exists: statErr == nil,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files, nil
}

// FilesForProject returns all file references belonging to the named project,
// sorted by path. Returns ErrProjectNotFound if no project has that identifier.
func (a *APW) FilesForProject(projectID string) ([]File, error) {
	return a.filesForProject(projectID)
}

// FilesForSystem returns all file references belonging to the named system
// within the named project, sorted by path. Returns ErrProjectNotFound or
// ErrSystemNotFound if the identifiers cannot be resolved.
func (a *APW) FilesForSystem(projectID, systemID string) ([]File, error) {
	return a.filesForSystem(projectID, systemID)
}

// GetExtraFileReferencesForProject returns unique extra file IDs referenced
// from files belonging to the named project. Returns ErrProjectNotFound if
// the project does not exist.
func (a *APW) GetExtraFileReferencesForProject(projectID string) ([]string, error) {
	files, err := a.filesForProject(projectID)
	if err != nil {
		return nil, err
	}

	return a.collectExtraRefs(files)
}

// GetExtraFileReferencesForProjectMap returns a map of unique extra file ID →
// source file path for every file in the named project.
func (a *APW) GetExtraFileReferencesForProjectMap(projectID string) (map[string]string, error) {
	files, err := a.filesForProject(projectID)
	if err != nil {
		return nil, err
	}

	return a.collectExtraRefsMap(files)
}

// GetExtraFileReferencesForSystem returns unique extra file IDs referenced
// from files belonging to the named system within the named project. Returns
// ErrProjectNotFound or ErrSystemNotFound if the identifiers cannot be resolved.
func (a *APW) GetExtraFileReferencesForSystem(projectID, systemID string) ([]string, error) {
	files, err := a.filesForSystem(projectID, systemID)
	if err != nil {
		return nil, err
	}

	return a.collectExtraRefs(files)
}

// GetExtraFileReferencesForSystemMap returns a map of unique extra file ID →
// source file path for every file in the named system within the named project.
func (a *APW) GetExtraFileReferencesForSystemMap(projectID, systemID string) (map[string]string, error) {
	files, err := a.filesForSystem(projectID, systemID)
	if err != nil {
		return nil, err
	}

	return a.collectExtraRefsMap(files)
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

// ScopedWorkspace returns a Workspace value for marshalling that contains only
// the targeted scope. FileRef pointers are shared with the receiver so path
// rewrites previously applied via FileRef.SetPath are visible in the result.
//
//   - projectID="" and systemID=""	: full workspace returned unchanged.
//   - projectID set, systemID=""	: workspace trimmed to just that project (all its systems).
//   - projectID set, systemID set	: workspace trimmed to that project containing only that system.
func (a *APW) ScopedWorkspace(projectID, systemID string) (*Workspace, error) {
	if projectID == "" {
		return a.ws, nil
	}

	proj, ok := a.ws.FindProject(projectID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProjectNotFound, projectID)
	}

	// Shallow-copy the Workspace so we don't mutate a.ws.Projects.
	scoped := *a.ws

	if systemID == "" {
		scoped.Projects = []*Project{proj}
		return &scoped, nil
	}

	sys, ok := proj.FindSystem(systemID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrSystemNotFound, systemID)
	}

	// Shallow-copy the Project so we don't mutate proj.Systems.
	projCopy := *proj
	projCopy.Systems = []*System{sys}
	scoped.Projects = []*Project{&projCopy}

	return &scoped, nil
}

// ID returns the workspace identifier derived from the filename stem.
func (a *APW) ID() string {
	return a.id
}

// FilePath returns the absolute path of the workspace file.
func (a *APW) FilePath() string {
	return a.path
}
