package apw

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// File type constants matching the TypeScript AmxFileType enum.
const (
	FileTypeWorkspace = "Workspace"
	FileTypeModule    = "Module"
	FileTypeMasterSrc = "MasterSrc"
	FileTypeSource    = "Source"
	FileTypeInclude   = "Include"
	FileTypeIR        = "IR"
	FileTypeTP4       = "TP4"
	FileTypeTP5       = "TP5"
	FileTypeTPD       = "TPD"
	FileTypeKPD       = "KPD"
	FileTypeAXB       = "AXB"
	FileTypeTKO       = "TKO"
	FileTypeIRDB      = "IRDB"
	FileTypeIRNDB     = "IRNDB"
	FileTypeDuet      = "Duet"
	FileTypeTOK       = "TOK"
	FileTypeTKN       = "TKN"
	FileTypeKPB       = "KPB"
	FileTypeXDD       = "XDD"
	FileTypeOther     = "Other"
)

// AmxExtensions maps file type constants to their file extensions.
var AmxExtensions = map[string]string{
	FileTypeWorkspace: ".apw",
	FileTypeModule:    ".axs",
	FileTypeMasterSrc: ".axs",
	FileTypeSource:    ".axs",
	FileTypeInclude:   ".axi",
	FileTypeIR:        ".irl",
	FileTypeTP4:       ".tp4",
	FileTypeTP5:       ".tp5",
	FileTypeTPD:       ".tpd",
	FileTypeDuet:      ".jar",
	FileTypeXDD:       ".xdd",
	FileTypeKPD:       ".kpd",
	FileTypeAXB:       ".axb",
	FileTypeTKO:       ".tko",
	FileTypeIRDB:      ".irdb",
	FileTypeIRNDB:     ".irndb",
	FileTypeTOK:       ".tok",
	FileTypeTKN:       ".tkn",
	FileTypeKPB:       ".kpb",
	FileTypeOther:     "",
}

// AmxCompiledExtensions maps file type constants to their compiled output extensions.
var AmxCompiledExtensions = map[string]string{
	FileTypeWorkspace: ".apw",
	FileTypeModule:    ".tko",
	FileTypeMasterSrc: ".tkn",
	FileTypeSource:    ".tkn",
	FileTypeInclude:   ".tkn",
	FileTypeIR:        ".irl",
	FileTypeTP4:       ".tp4",
	FileTypeTP5:       ".tp5",
	FileTypeTPD:       ".tpd",
	FileTypeDuet:      ".jar",
	FileTypeXDD:       ".xdd",
	FileTypeKPD:       ".kpd",
	FileTypeAXB:       ".axb",
	FileTypeTKO:       ".tko",
	FileTypeIRDB:      ".irdb",
	FileTypeIRNDB:     ".irndb",
	FileTypeTOK:       ".tok",
	FileTypeTKN:       ".tkn",
	FileTypeKPB:       ".kpb",
	FileTypeOther:     "",
}

// File represents a file reference parsed from an APW workspace.
type File struct {
	ID      string
	Type    string
	Path    string
	Exists  bool
	IsExtra bool
	Content string // optional in-memory content (e.g. for .env files)
}

// APW represents a parsed AMX NetLinx workspace (.apw) file.
type APW struct {
	filePath             string
	id                   string
	fileReferences       []File
	uniqueFileReferences []File
}

// New returns a new APW instance with an absolute file path.
func New(filePath string) *APW {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	return &APW{filePath: absPath}
}

// Load reads and parses the APW file, populating all internal state.
func (a *APW) Load() error {
	data, err := a.read()
	if err != nil {
		return fmt.Errorf("failed to load APW file: %w", err)
	}

	id, err := parseID(data)
	if err != nil {
		return fmt.Errorf("failed to load APW file: %w", err)
	}

	a.id = id

	refs, err := a.parseFileReferences(data)
	if err != nil {
		return fmt.Errorf("failed to load APW file: %w", err)
	}

	a.fileReferences = refs
	a.uniqueFileReferences = deduplicateFileReferences(refs)

	return nil
}

func (a *APW) read() (string, error) {
	if _, err := os.Stat(a.filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist")
	}

	data, err := os.ReadFile(a.filePath)
	if err != nil {
		return "", err
	}

	content := string(data)

	matched, _ := regexp.MatchString(`<!DOCTYPE Workspace \[`, content)
	if !matched {
		return "", fmt.Errorf("not a Netlinx Workspace file")
	}

	return content, nil
}

// parseID extracts the workspace <Identifier> value.
// Mirrors: /<Workspace.+\r?\n?.*?<Identifier>(?<id>.+)<.+>/m
func parseID(data string) (string, error) {
	re := regexp.MustCompile(`<Workspace.+\r?\n?.*?<Identifier>(.+)<.+>`)

	match := re.FindStringSubmatch(data)
	if match == nil {
		return "", fmt.Errorf("no Workspace ID found")
	}

	return strings.TrimSpace(match[1]), nil
}

// parseFileReferences extracts all <File> entries from the workspace data.
// Mirrors the TypeScript multi-line gm regex for file references.
func (a *APW) parseFileReferences(data string) ([]File, error) {
	// This pattern is intentionally close to the TS original, relying on `.`
	// not matching newlines and explicit \r?\n? for line crossings.
	re := regexp.MustCompile(
		`<File.+Type="(.+?)".+\r?\n?.*?` +
			`<Identifier>(.+?)<\/Identifier>\r?\n?.*?>(.+?)<.+` +
			`\r?\n?.*?\r?\n?.*?(?:<DeviceMap.+\r?\n?.*?\r?\n?.*?\r?\n?)?.*?<\/File>`,
	)

	var refs []File

	for _, match := range re.FindAllStringSubmatch(data, -1) {
		if len(match) < 4 {
			continue
		}

		fileType := strings.TrimSpace(match[1])
		id := strings.TrimSpace(match[2])
		filePath := strings.TrimSpace(match[3])

		_, statErr := os.Stat(filePath)
		exists := statErr == nil

		refs = append(refs, File{
			ID:      id,
			Type:    fileType,
			Path:    filePath,
			Exists:  exists,
			IsExtra: false,
		})
	}

	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Path < refs[j].Path
	})

	return refs, nil
}

func deduplicateFileReferences(refs []File) []File {
	seen := make(map[string]bool)
	var unique []File

	for _, f := range refs {
		if !seen[f.Path] {
			seen[f.Path] = true
			unique = append(unique, f)
		}
	}

	return unique
}

// FileIsReadable reports whether the file can be scanned for #include / define_module references.
func FileIsReadable(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".axs" || ext == ".axi"
}

// FileIsOfInterest reports whether the file should be considered as a candidate
// during an extra-file disk search (.axs, .axi, .jar, .xdd).
func FileIsOfInterest(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".axs" || ext == ".axi" || ext == ".jar" || ext == ".xdd"
}

// GetFileType returns the AmxFileType constant for a file path based on its extension.
func GetFileType(file string) string {
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
	for _, f := range a.uniqueFileReferences {
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

	for _, f := range a.uniqueFileReferences {
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

// AllFiles returns all unique workspace file references plus the workspace file itself.
func (a *APW) AllFiles() []File {
	files := make([]File, len(a.uniqueFileReferences))
	copy(files, a.uniqueFileReferences)

	files = append(files, File{
		ID:      a.id,
		Type:    FileTypeWorkspace,
		Path:    a.filePath,
		Exists:  true,
		IsExtra: false,
	})

	return files
}

// MasterSrcPath returns the directories that contain MasterSrc files.
func (a *APW) MasterSrcPath() []string {
	seen := make(map[string]bool)
	var dirs []string

	for _, f := range a.uniqueFileReferences {
		if f.Type != FileTypeMasterSrc {
			continue
		}

		rootDir := filepath.Dir(a.filePath)
		absPath := filepath.Join(rootDir, f.Path)
		dir := filepath.Dir(absPath)

		if !seen[dir] {
			dirs = append(dirs, dir)
			seen[dir] = true
		}
	}

	return dirs
}

// ModuleFiles returns all unique file references with type Module.
func (a *APW) ModuleFiles() []File {
	var files []File

	for _, f := range a.uniqueFileReferences {
		if f.Type == FileTypeModule {
			files = append(files, f)
		}
	}

	return files
}

// MasterSrcFiles returns all unique file references with type MasterSrc.
func (a *APW) MasterSrcFiles() []File {
	var files []File

	for _, f := range a.uniqueFileReferences {
		if f.Type == FileTypeMasterSrc {
			files = append(files, f)
		}
	}

	return files
}

// getFileDirectories returns the unique directory paths for the given files,
// resolving relative paths against the workspace file's directory.
func (a *APW) getFileDirectories(files []File) []string {
	seen := make(map[string]bool)
	var dirs []string

	rootDir := filepath.Dir(a.filePath)

	for _, f := range files {
		absPath := filepath.Join(rootDir, f.Path)
		dir := filepath.Dir(absPath)

		if !seen[dir] {
			dirs = append(dirs, dir)
			seen[dir] = true
		}
	}

	return dirs
}

// IncludePath returns the unique directories containing Include files.
func (a *APW) IncludePath() []string {
	var files []File

	for _, f := range a.uniqueFileReferences {
		if f.Type == FileTypeInclude {
			files = append(files, f)
		}
	}

	return a.getFileDirectories(files)
}

// ModulePath returns the unique directories containing Module, Duet, and XDD files.
func (a *APW) ModulePath() []string {
	var files []File

	for _, f := range a.uniqueFileReferences {
		if f.Type == FileTypeModule || f.Type == FileTypeDuet || f.Type == FileTypeXDD {
			files = append(files, f)
		}
	}

	return a.getFileDirectories(files)
}

// ID returns the workspace identifier with spaces replaced by hyphens.
func (a *APW) ID() string {
	return strings.ReplaceAll(a.id, " ", "-")
}

// FilePath returns the absolute path of the workspace file.
func (a *APW) FilePath() string {
	return a.filePath
}
