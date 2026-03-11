package archive

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// Package-level color instances used by logVerbose and displayZippedFiles.
var (
	logBlue  = color.New(color.FgBlue)
	logCyan  = color.New(color.FgCyan)
	logGreen = color.New(color.FgGreen)
)

// sanitizeSegment replaces spaces with hyphens so that archive filenames are
// shell-friendly, then collapses any run of consecutive hyphens down to one.
// It is applied only to each segment of the output filename; nothing inside
// the archive (file paths, APW content) is affected.
func sanitizeSegment(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

// Options configures the archive build process, merging config file values with
// CLI flags.
type Options struct {
	OutputFileSuffix           string
	IncludeCompiledSourceFiles bool
	IncludeCompiledModuleFiles bool
	IncludeFilesNotInWorkspace bool
	ExtraFileSearchLocations   []string
	All                        bool
	IgnoredFiles               []string
	Verbose                    bool
	// ProjectID, when non-empty, restricts the archive to files in this project.
	ProjectID string
	// SystemID, when non-empty (and ProjectID is set), further restricts the
	// archive to files in this system within the named project.
	SystemID string
}

// Builder creates a zip archive from a parsed APW workspace.
type Builder struct {
	apw                 *apw.APW
	opts                *Options
	zipWriter           *zip.Writer
	entries             []string        // tracks zip entry names for verbose display
	seenEntries         map[string]bool // guards against duplicate zip entries
	extraFilesOnDisk    []string
	extraFileReferences []string
	locatedExtraRefs    []string
	// scopeFiles is the set of workspace files belonging to the targeted
	// project/system. When non-nil, extra-file scanning uses a scope-aware
	// check so that files listed in the workspace under a different scope are
	// still discovered and included rather than silently dropped.
	scopeFiles []apw.File
	// outputFile holds the file path of the last successfully built archive.
	outputFile string
}

// OutputFile returns the file path of the last successfully built archive.
// It is only valid after a successful call to Build.
func (b *Builder) OutputFile() string {
	return b.outputFile
}

// NewBuilder returns a new Builder for the given workspace and options.
func NewBuilder(workspace *apw.APW, opts *Options) *Builder {
	return &Builder{
		apw:         workspace,
		opts:        opts,
		seenEntries: make(map[string]bool),
	}
}

// Build produces the zip archive and writes it to disk.
func (b *Builder) Build() error {
	b.logVerbose(logBlue, "Creating archive...")

	outputName := sanitizeSegment(b.apw.ID())
	if b.opts.ProjectID != "" {
		outputName += "-" + sanitizeSegment(b.opts.ProjectID)
	}
	if b.opts.SystemID != "" {
		outputName += "-" + sanitizeSegment(b.opts.SystemID)
	}

	outputFile := fmt.Sprintf("%s.%s", outputName, b.opts.OutputFileSuffix)

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	b.zipWriter = zip.NewWriter(f)

	// addExtraFiles MUST run before addWorkspaceFiles. addWorkspaceFiles mutates
	// each FileRef's FilePathName to a bare filename (flat-layout path rewrite)
	// so that the marshalled APW stored in the archive is correct. The scoped
	// variants of GetExtraFileReferences* re-read those FilePathName values to
	// locate source files on disk — if they run after the mutation the paths no
	// longer resolve and no extra refs are found.
	if err := b.addExtraFiles(); err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return err
	}

	if err := b.addWorkspaceFiles(); err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return err
	}

	// Marshal a scoped copy of the workspace — only the targeted project/system
	// is included so the APW in the archive is self-consistent with its contents.
	scopedWS, err := b.apw.ScopedWorkspace(b.opts.ProjectID, b.opts.SystemID)
	if err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return fmt.Errorf("failed to scope workspace: %w", err)
	}

	apwData, err := apw.Marshal(scopedWS)
	if err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return fmt.Errorf("failed to marshal workspace: %w", err)
	}

	apwEntry := zipEntryPath(filepath.Base(b.apw.FilePath()))
	if err := b.writeEntry(apwEntry, apwData); err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return fmt.Errorf("failed to write workspace to archive: %w", err)
	}

	if err := b.zipWriter.Close(); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to finalise zip: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close output file: %w", err)
	}

	b.outputFile = outputFile

	if b.opts.Verbose {
		b.displayZippedFiles()
	}

	return nil
}

// ---------------------------------------------------------------------------
// Logging helpers
// ---------------------------------------------------------------------------

func (b *Builder) logVerbose(c *color.Color, format string, args ...any) {
	if b.opts.Verbose {
		fmt.Println(c.Sprintf(format, args...))
	}
}

// warnf prints a yellow warning to stderr so it is always visible regardless
// of verbose mode and does not pollute stdout.
func warnf(format string, args ...any) {
	fmt.Fprintln(os.Stderr, color.YellowString("WARNING: "+format, args...))
}

// ---------------------------------------------------------------------------
// Zip entry helpers
// ---------------------------------------------------------------------------

// zipEntryPath converts an OS path to a forward-slash zip entry path.
func zipEntryPath(p string) string {
	return path.Clean(strings.ReplaceAll(p, "\\", "/"))
}

func (b *Builder) writeEntry(entryName string, data []byte) error {
	if b.seenEntries[entryName] {
		return nil
	}

	w, err := b.zipWriter.Create(entryName)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	b.seenEntries[entryName] = true
	b.entries = append(b.entries, entryName)

	return nil
}

func (b *Builder) addDiskFile(diskPath, entryName string) error {
	data, err := os.ReadFile(diskPath)
	if err != nil {
		return err
	}

	return b.writeEntry(entryName, data)
}

// ---------------------------------------------------------------------------
// Per-type add methods
// ---------------------------------------------------------------------------

// addGeneralItem adds a file to the archive root using just its filename.
// Mirrors GeneralItem.addToArchive() — flat layout, no subdirectory.
func (b *Builder) addGeneralItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Base(file.Path))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	return nil
}

// addSourceItem adds the source file and, when requested, its compiled .tkn.
// Both are placed at the archive root.
func (b *Builder) addSourceItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Base(file.Path))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	if !b.opts.IncludeCompiledSourceFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeSource]) +
		apw.AmxCompiledExtensions[apw.FileTypeSource]
	compiledEntry := zipEntryPath(filepath.Base(compiledPath))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		warnf("compiled file not found, skipping: %s", compiledPath)
	}

	return nil
}

// addModuleItem adds the module source and, when requested, its compiled .tko.
// Both are placed at the archive root.
func (b *Builder) addModuleItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Base(file.Path))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	if !b.opts.IncludeCompiledModuleFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeModule]) +
		apw.AmxCompiledExtensions[apw.FileTypeModule]
	compiledEntry := zipEntryPath(filepath.Base(compiledPath))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		warnf("compiled file not found, skipping: %s", compiledPath)
	}

	return nil
}

// addFileToArchive dispatches to the correct per-type method.
func (b *Builder) addFileToArchive(file apw.File) error {
	switch file.Type {
	case apw.FileTypeModule:
		return b.addModuleItem(file)
	case apw.FileTypeSource, apw.FileTypeMasterSrc:
		return b.addSourceItem(file)
	default:
		return b.addGeneralItem(file)
	}
}

// ---------------------------------------------------------------------------
// Workspace files pass
// ---------------------------------------------------------------------------

// addWorkspaceFiles walks the live workspace hierarchy directly, adds each
// referenced file to the archive root (filename only — flat layout), and
// rewrites the FileRef path to just the filename so the marshalled .apw
// written at the end of Build() reflects the flat layout NetLinx Studio
// expects after extraction.
//
// When b.opts.ProjectID is set, only that project is processed. When
// b.opts.SystemID is also set, only that system within the project is processed.
func (b *Builder) addWorkspaceFiles() error {
	workspaceDir := filepath.Dir(b.apw.FilePath())

	projects := b.apw.Workspace().Projects
	if b.opts.ProjectID != "" {
		proj, ok := b.apw.Workspace().FindProject(b.opts.ProjectID)
		if !ok {
			return fmt.Errorf("project %q not found in workspace", b.opts.ProjectID)
		}

		projects = []*apw.Project{proj}
	}

	for _, proj := range projects {
		systems := proj.Systems
		if b.opts.SystemID != "" {
			sys, ok := proj.FindSystem(b.opts.SystemID)
			if !ok {
				return fmt.Errorf("system %q not found in project %q", b.opts.SystemID, proj.Identifier)
			}

			systems = []*apw.System{sys}
		}

		for _, sys := range systems {
			for _, fr := range sys.Files {
				relPath := filepath.FromSlash(strings.ReplaceAll(fr.FilePathName, `\`, `/`))
				diskPath := filepath.Join(workspaceDir, relPath)

				if _, err := os.Stat(diskPath); err != nil {
					warnf("file referenced in workspace does not exist on disk, skipping: %s", diskPath)
					continue
				}

				file := apw.File{
					Type: fr.Type,
					Path: diskPath,
				}

				if err := b.addFileToArchive(file); err != nil {
					warnf("could not add %s: %v", diskPath, err)
					continue
				}

				fr.SetPath(filepath.Base(diskPath))
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Extra-files pass (files not in the workspace)
// ---------------------------------------------------------------------------

// getExtraFilesOnDisk walks the given locations and collects files that are of
// interest for extra-file matching (.axs, .axi, .jar, .xdd).
func (b *Builder) getExtraFilesOnDisk(locations []string) {
	b.logVerbose(logBlue, "Searching known locations for extra files...")

	for _, location := range locations {
		b.logVerbose(logCyan, "--> Searching %s", location)

		_ = filepath.WalkDir(location, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip inaccessible paths
			}

			if d.IsDir() {
				name := d.Name()
				if name == "node_modules" || name == ".git" || name == ".history" {
					return filepath.SkipDir
				}

				return nil
			}

			if apw.FileIsOfInterest(p) {
				b.extraFilesOnDisk = append(b.extraFilesOnDisk, p)
			}

			return nil
		})
	}
}

func (b *Builder) displayExtraFileReferences(refs []string) {
	b.logVerbose(logGreen, "Found %d extra files referenced...", len(refs))

	for _, r := range refs {
		b.logVerbose(logCyan, "--> %s", r)
	}
}

func (b *Builder) isIgnored(filePath string) bool {
	base := filepath.Base(filePath)

	return slices.Contains(b.opts.IgnoredFiles, base)
}

// searchForExtraFiles tries to locate each referenced ID on disk and records
// found files for later addition to the archive.
// refs maps each reference ID to the source file it was found in, so that
// warning messages can pinpoint exactly where an unresolved reference came from.
func (b *Builder) searchForExtraFiles(refs map[string]string) error {
	var newLocated []string

	for ref, sourceFile := range refs {
		var found string

		for _, diskFile := range b.extraFilesOnDisk {
			if strings.Contains(diskFile, ref) {
				found = diskFile
				break
			}
		}

		if found == "" {
			warnf("Could not find %s (referenced in %s)", ref, filepath.Base(sourceFile))
			continue
		}

		if b.isIgnored(found) {
			b.logVerbose(logBlue, "Ignoring %s", filepath.Base(found))
			continue
		}

		// Skip if this exact disk path was already located via a different
		// reference ID earlier in the pipeline — the same physical file can be
		// referenced with or without its extension, leading to duplicate entries
		// that cause redundant re-scanning and misleading log output.
		if slices.Contains(b.locatedExtraRefs, found) || slices.Contains(newLocated, found) {
			continue
		}

		b.logVerbose(logGreen, "Found %s: %s", ref, found)
		newLocated = append(newLocated, found)
	}

	if len(newLocated) == 0 {
		return nil
	}

	b.locatedExtraRefs = append(b.locatedExtraRefs, newLocated...)

	return b.getFileReferencesFromFiles(newLocated)
}

// getFileReferencesFromFiles scans source files for #include and define_module
// references not already tracked, then recurses to find those files on disk.
// Mirrors ArchiveBuilder.getFileReferences().
func (b *Builder) getFileReferencesFromFiles(files []string) error {
	for _, file := range files {
		if !apw.FileIsReadable(file) {
			continue
		}

		var refs []string
		var err error
		if b.scopeFiles != nil {
			refs, err = b.apw.GetExtraFileReferencesFromFileInScope(file, b.scopeFiles)
		} else {
			refs, err = b.apw.GetExtraFileReferencesFromFile(file)
		}
		if err != nil {
			continue
		}

		if len(refs) == 0 {
			continue
		}

		var newRefs []string

		for _, r := range refs {
			alreadyKnown := slices.Contains(b.extraFileReferences, r)

			if !alreadyKnown {
				newRefs = append(newRefs, r)
			}
		}

		if len(newRefs) == 0 {
			continue
		}

		b.logVerbose(logBlue, "Scanning %s...", filepath.Base(file))
		b.extraFileReferences = append(b.extraFileReferences, newRefs...)
		b.displayExtraFileReferences(newRefs)

		newRefsMap := make(map[string]string, len(newRefs))
		for _, r := range newRefs {
			newRefsMap[r] = file
		}

		if err := b.searchForExtraFiles(newRefsMap); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) addExtraFiles() error {
	if !b.opts.IncludeFilesNotInWorkspace {
		return nil
	}

	b.logVerbose(logBlue, "Searching for extra files that are not part of the workspace...")

	var (
		refsMap map[string]string
		err     error
	)

	switch {
	case b.opts.ProjectID != "" && b.opts.SystemID != "":
		refsMap, err = b.apw.GetExtraFileReferencesForSystemMap(b.opts.ProjectID, b.opts.SystemID)
		b.scopeFiles, _ = b.apw.FilesForSystem(b.opts.ProjectID, b.opts.SystemID)
	case b.opts.ProjectID != "":
		refsMap, err = b.apw.GetExtraFileReferencesForProjectMap(b.opts.ProjectID)
		b.scopeFiles, _ = b.apw.FilesForProject(b.opts.ProjectID)
	default:
		refsMap, err = b.apw.GetExtraFileReferencesMap()
		// scopeFiles stays nil — recursive scan uses full workspace check
	}

	if err != nil {
		return err
	}

	for k := range refsMap {
		b.extraFileReferences = append(b.extraFileReferences, k)
	}

	if len(b.extraFileReferences) == 0 {
		b.logVerbose(logBlue, "No extra file references found in the workspace file")
		return nil
	}

	searchLocations := append(
		append([]string{}, b.opts.ExtraFileSearchLocations...),
		filepath.Dir(b.apw.FilePath()),
	)

	b.logVerbose(logGreen, "Found %d extra files referenced...", len(b.extraFileReferences))
	b.getExtraFilesOnDisk(searchLocations)

	if err := b.searchForExtraFiles(refsMap); err != nil {
		return err
	}

	if len(b.locatedExtraRefs) == 0 {
		return nil
	}

	for _, ref := range b.locatedExtraRefs {
		fileType := apw.GetFileType(ref)

		// GetFileType is ambiguous for .axs: normalise to FileTypeModule since
		// extra .axs files are always discovered via define_module directives.
		if fileType == apw.FileTypeSource || fileType == apw.FileTypeMasterSrc {
			fileType = apw.FileTypeModule
		}

		file := apw.File{
			Type:    fileType,
			Path:    ref,
			Exists:  true,
			IsExtra: true,
		}

		if err := b.addFileToArchive(file); err != nil {
			warnf("could not add %s: %v", ref, err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Display helper
// ---------------------------------------------------------------------------

func (b *Builder) displayZippedFiles() {
	for _, entry := range b.entries {
		fmt.Println(logCyan.Sprintf("--> %s", entry))
	}
}
