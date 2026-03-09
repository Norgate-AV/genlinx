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
	logRed   = color.New(color.FgRed)
)

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

	outputFile := fmt.Sprintf("%s.%s", b.apw.ID(), b.opts.OutputFileSuffix)

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	b.zipWriter = zip.NewWriter(f)

	if err := b.addWorkspaceFiles(); err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return err
	}

	if err := b.addExtraFiles(); err != nil {
		_ = b.zipWriter.Close()
		_ = f.Close()
		return err
	}

	// Marshal the workspace with rewritten flat paths and write it to the archive.
	apwData, err := b.apw.Bytes()
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

	fmt.Println(color.New(color.FgGreen, color.Bold).Sprintf("Created archive: %s", outputFile))

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
	fmt.Fprintln(os.Stderr, color.YellowString("Warning: "+format, args...))
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

	b.logVerbose(logCyan, "Added file: %s", entryName)

	return nil
}

// addSourceItem adds the source file and, when requested, its compiled .tkn.
// Both are placed at the archive root.
func (b *Builder) addSourceItem(file apw.File) error {
	entryName := zipEntryPath(filepath.Base(file.Path))

	if err := b.addDiskFile(file.Path, entryName); err != nil {
		return err
	}

	b.logVerbose(logCyan, "Added file: %s", entryName)

	if !b.opts.IncludeCompiledSourceFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeSource]) +
		apw.AmxCompiledExtensions[apw.FileTypeSource]
	compiledEntry := zipEntryPath(filepath.Base(compiledPath))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		warnf("compiled file not found, skipping: %s", compiledPath)
	} else {
		b.logVerbose(logCyan, "Added file: %s", compiledEntry)
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

	b.logVerbose(logCyan, "Added file: %s", entryName)

	if !b.opts.IncludeCompiledModuleFiles {
		return nil
	}

	compiledPath := strings.TrimSuffix(file.Path, apw.AmxExtensions[apw.FileTypeModule]) +
		apw.AmxCompiledExtensions[apw.FileTypeModule]
	compiledEntry := zipEntryPath(filepath.Base(compiledPath))

	if err := b.addDiskFile(compiledPath, compiledEntry); err != nil {
		warnf("compiled file not found, skipping: %s", compiledPath)
	} else {
		b.logVerbose(logCyan, "Added file: %s", compiledEntry)
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
func (b *Builder) addWorkspaceFiles() error {
	workspaceDir := filepath.Dir(b.apw.FilePath())

	for _, proj := range b.apw.Workspace().Projects {
		for _, sys := range proj.Systems {
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
			b.logVerbose(logRed, "Could not find %s", ref)
			continue
		}

		if b.isIgnored(found) {
			b.logVerbose(logBlue, "Ignoring %s as per config", filepath.Base(found))
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

		b.logVerbose(logBlue, "Searching %s for references...", file)

		refs, err := b.apw.GetExtraFileReferencesFromFile(file)
		if err != nil {
			continue
		}

		if len(refs) == 0 {
			b.logVerbose(logCyan, "--> No references found")
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
			b.logVerbose(logCyan, "--> No new references found")
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

func (b *Builder) addExtraFiles() error {
	if !b.opts.IncludeFilesNotInWorkspace {
		return nil
	}

	b.logVerbose(logBlue, "Searching for extra files that are not part of the workspace...")

	refs, err := b.apw.GetExtraFileReferences()
	if err != nil {
		return err
	}

	b.extraFileReferences = append(b.extraFileReferences, refs...)

	if len(b.extraFileReferences) == 0 {
		b.logVerbose(logBlue, "No extra file references found in the workspace file")
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
