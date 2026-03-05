package archive

import (
	"archive/zip"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Norgate-AV/genlinx-go/internal/apw"
)

//go:embed scripts/symlink.bat
var symlinkBat []byte

//go:embed scripts/symlink.ps1
var symlinkPS1 []byte

// Options configures the archive build process, merging config file values with
// CLI flags. It is the Go equivalent of the TypeScript ArchiveOptions type.
type Options struct {
	OutputFileSuffix           string
	IncludeCompiledSourceFiles bool
	IncludeCompiledModuleFiles bool
	IncludeFilesNotInWorkspace bool
	ExtraFileSearchLocations   []string
	ExtraFileArchiveLocation   string
	All                        bool
	IgnoredFiles               []string
	Verbose                    bool
}

// Builder creates a zip archive from a parsed APW workspace.
type Builder struct {
	apw                 *apw.APW
	opts                *Options
	zipWriter           *zip.Writer
	entries             []string // tracks zip entry names for verbose display
	extraFilesOnDisk    []string
	extraFileReferences []string
	locatedExtraRefs    []string
}

// NewBuilder returns a new Builder for the given workspace and options.
func NewBuilder(workspace *apw.APW, opts *Options) *Builder {
	return &Builder{
		apw:  workspace,
		opts: opts,
	}
}

// Build produces the zip archive and writes it to disk.
func (b *Builder) Build() error {
	b.logVerbose("Creating archive...")

	outputFile := fmt.Sprintf("%s.%s", b.apw.ID(), b.opts.OutputFileSuffix)

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	b.zipWriter = zip.NewWriter(f)

	if err := b.addWorkspaceFiles(); err != nil {
		b.zipWriter.Close()
		f.Close()
		return err
	}

	if err := b.addExtraFiles(); err != nil {
		b.zipWriter.Close()
		f.Close()
		return err
	}

	if err := b.zipWriter.Close(); err != nil {
		f.Close()
		return fmt.Errorf("failed to finalise zip: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close output file: %w", err)
	}

	fmt.Printf("Created archive: %s\n", outputFile)

	if b.opts.Verbose {
		b.displayZippedFiles()
	}

	return nil
}

// ---------------------------------------------------------------------------
// Logging helpers
// ---------------------------------------------------------------------------

func (b *Builder) logVerbose(format string, args ...any) {
	if b.opts.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// ---------------------------------------------------------------------------
// Zip entry helpers
// ---------------------------------------------------------------------------

// zipEntryPath converts an OS path to a forward-slash zip entry path.
func zipEntryPath(p string) string {
	return path.Clean(filepath.ToSlash(p))
}

func (b *Builder) writeEntry(entryName string, data []byte) error {
	w, err := b.zipWriter.Create(entryName)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

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
// Per-type add methods (mirrors ArchiveItemFactory / *Item classes)
// ---------------------------------------------------------------------------

// archiveDir returns the zip directory for a file — mirrors the TS archivePath
// logic used in GeneralItem / SourceItem / ModuleItem constructors.
func (b *Builder) archiveDir(file apw.File) string {
	if file.IsExtra {
		return b.opts.ExtraFileArchiveLocation
	}

	return filepath.Dir(file.Path)
}

// addWorkspaceItem places the .apw file at the zip root (no subdirectory).
// Mirrors WorkspaceItem.addToArchive().
func (b *Builder) addWorkspaceItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Base(file.Path))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	b.logVerbose("Added file: %s", entryName)

	return nil
}

// addGeneralItem places the file under its directory path inside the zip.
// Mirrors GeneralItem.addToArchive().
func (b *Builder) addGeneralItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Join(b.archiveDir(file), filepath.Base(file.Path)))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	b.logVerbose("Added file: %s", file.Path)

	return nil
}

// addSourceItem places the source file and, when requested, its compiled .tkn.
// Mirrors SourceItem.addToArchive().
func (b *Builder) addSourceItem(file apw.File) error {
	dir := b.archiveDir(file)
	entryName := zipEntryPath(filepath.Join(dir, filepath.Base(file.Path)))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	b.logVerbose("Added file: %s", file.Path)

	if !b.opts.IncludeCompiledSourceFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeSource]) +
		apw.AmxCompiledExtensions[apw.FileTypeSource]
	compiledEntry := zipEntryPath(filepath.Join(dir, filepath.Base(compiledPath)))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		b.logVerbose("Compiled file not found, skipping: %s", compiledPath)
	} else {
		b.logVerbose("Added file: %s", compiledPath)
	}

	return nil
}

// addModuleItem places the module source and, when requested, its compiled .tko.
// Mirrors ModuleItem.addToArchive().
func (b *Builder) addModuleItem(file apw.File) error {
	dir := b.archiveDir(file)
	entryName := zipEntryPath(filepath.Join(dir, filepath.Base(file.Path)))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	b.logVerbose("Added file: %s", file.Path)

	if !b.opts.IncludeCompiledModuleFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeModule]) +
		apw.AmxCompiledExtensions[apw.FileTypeModule]
	compiledEntry := zipEntryPath(filepath.Join(dir, filepath.Base(compiledPath)))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		b.logVerbose("Compiled file not found, skipping: %s", compiledPath)
	} else {
		b.logVerbose("Added file: %s", compiledPath)
	}

	return nil
}

// addEnvItem writes an in-memory .env file into the extra-file archive location.
// Mirrors EnvItem.addToArchive().
func (b *Builder) addEnvItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Join(b.opts.ExtraFileArchiveLocation, file.Path))

	if err := b.writeEntry(entryName, []byte(file.Content)); err != nil {
		return err
	}

	b.logVerbose("Added file: %s", entryName)

	return nil
}

// addFileToArchive dispatches to the correct per-type method.
// Mirrors ArchiveItemFactory.create().
func (b *Builder) addFileToArchive(file apw.File) error {
	switch file.Type {
	case apw.FileTypeWorkspace:
		return b.addWorkspaceItem(file)
	case apw.FileTypeModule:
		return b.addModuleItem(file)
	case apw.FileTypeSource, apw.FileTypeMasterSrc:
		return b.addSourceItem(file)
	case "Env":
		return b.addEnvItem(file)
	default:
		return b.addGeneralItem(file)
	}
}

// ---------------------------------------------------------------------------
// Workspace files pass
// ---------------------------------------------------------------------------

func (b *Builder) addWorkspaceFiles() error {
	for _, file := range b.apw.AllFiles() {
		if err := b.addFileToArchive(file); err != nil {
			fmt.Printf("Warning: could not add %s: %v\n", file.Path, err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Extra-files pass (files not in the workspace)
// ---------------------------------------------------------------------------

func (b *Builder) addEnvFile() {
	masterSrcPaths := b.apw.MasterSrcPath()
	if len(masterSrcPaths) == 0 {
		return
	}

	b.logVerbose("Adding env file to the archive...")

	masterSrcRelPath := filepath.Join("..", filepath.Base(masterSrcPaths[0]))

	file := apw.File{
		Type:    "Env",
		Path:    ".env",
		Exists:  true,
		IsExtra: true,
		Content: fmt.Sprintf("SOURCE_DIRECTORY_RELATIVE_PATH=%s", filepath.ToSlash(masterSrcRelPath)),
	}

	if err := b.addEnvItem(file); err != nil {
		fmt.Printf("Warning: could not add .env: %v\n", err)
	}
}

func (b *Builder) addSymlinkScripts() {
	b.logVerbose("Adding symlink scripts for extra files to the archive...")

	scripts := map[string][]byte{
		"symlink.bat": symlinkBat,
		"symlink.ps1": symlinkPS1,
	}

	for name, content := range scripts {
		entryName := zipEntryPath(filepath.Join(b.opts.ExtraFileArchiveLocation, name))

		if err := b.writeEntry(entryName, content); err != nil {
			fmt.Printf("Warning: could not add script %s: %v\n", name, err)
		} else {
			b.logVerbose("Added file: %s", entryName)
		}
	}
}

// getExtraFilesOnDisk walks the given locations and collects files that are of
// interest for extra-file matching (.axs, .axi, .jar, .xdd).
// Mirrors ArchiveBuilder.getExtraFilesOnDisk().
func (b *Builder) getExtraFilesOnDisk(locations []string) {
	b.logVerbose("Searching known locations for extra files...")

	for _, location := range locations {
		b.logVerbose("--> Searching %s", location)

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
	b.logVerbose("Found %d extra files referenced...", len(refs))

	for _, r := range refs {
		b.logVerbose("--> %s", r)
	}
}

func (b *Builder) isIgnored(filePath string) bool {
	base := filepath.Base(filePath)

	for _, ignored := range b.opts.IgnoredFiles {
		if base == ignored {
			return true
		}
	}

	return false
}

// searchForExtraFiles tries to locate each referenced ID on disk and records
// found files for later addition to the archive.
// Mirrors ArchiveBuilder.searchForExtraFiles().
func (b *Builder) searchForExtraFiles(refs []string) error {
	var newLocated []string

	for _, ref := range refs {
		var found string

		for _, diskFile := range b.extraFilesOnDisk {
			if strings.Contains(diskFile, ref) {
				found = diskFile
				break
			}
		}

		if found == "" {
			b.logVerbose("Could not find %s", ref)
			continue
		}

		if b.isIgnored(found) {
			b.logVerbose("Ignoring %s as per config", filepath.Base(found))
			continue
		}

		b.logVerbose("Found %s: %s", ref, found)
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

		b.logVerbose("Searching %s for references...", file)

		refs, err := b.apw.GetExtraFileReferencesFromFile(file)
		if err != nil {
			continue
		}

		if len(refs) == 0 {
			b.logVerbose("--> No references found")
			continue
		}

		var newRefs []string

		for _, r := range refs {
			alreadyKnown := false

			for _, existing := range b.extraFileReferences {
				if existing == r {
					alreadyKnown = true
					break
				}
			}

			if !alreadyKnown {
				newRefs = append(newRefs, r)
			}
		}

		if len(newRefs) == 0 {
			b.logVerbose("--> No new references found")
			continue
		}

		b.extraFileReferences = append(b.extraFileReferences, newRefs...)
		b.displayExtraFileReferences(newRefs)

		if err := b.searchForExtraFiles(newRefs); err != nil {
			return err
		}
	}

	return nil
}

// addExtraFiles orchestrates the discovery and addition of files that are
// referenced in the workspace but not listed as workspace members.
// Mirrors ArchiveBuilder.addExtraFiles().
func (b *Builder) addExtraFiles() error {
	if !b.opts.IncludeFilesNotInWorkspace {
		return nil
	}

	b.logVerbose("Searching for extra files that are not part of the workspace...")

	refs, err := b.apw.GetExtraFileReferences()
	if err != nil {
		return err
	}

	b.extraFileReferences = append(b.extraFileReferences, refs...)

	if len(b.extraFileReferences) == 0 {
		b.logVerbose("No extra file references found in the workspace file")
		return nil
	}

	searchLocations := append(
		append([]string{}, b.opts.ExtraFileSearchLocations...),
		filepath.Dir(b.apw.FilePath()),
	)

	b.displayExtraFileReferences(b.extraFileReferences)
	b.getExtraFilesOnDisk(searchLocations)

	if err := b.searchForExtraFiles(b.extraFileReferences); err != nil {
		return err
	}

	if len(b.locatedExtraRefs) == 0 {
		return nil
	}

	for _, ref := range b.locatedExtraRefs {
		file := apw.File{
			Type:    apw.GetFileType(ref),
			Path:    ref,
			Exists:  true,
			IsExtra: true,
		}

		if err := b.addFileToArchive(file); err != nil {
			fmt.Printf("Warning: could not add %s: %v\n", ref, err)
		}
	}

	b.addEnvFile()
	b.addSymlinkScripts()

	return nil
}

// ---------------------------------------------------------------------------
// Display helper
// ---------------------------------------------------------------------------

func (b *Builder) displayZippedFiles() {
	for _, entry := range b.entries {
		fmt.Printf("--> %s\n", entry)
	}
}
