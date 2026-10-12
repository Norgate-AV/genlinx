package studio

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Preferences is the root element of a NetLinx Studio .epx preferences file.
type Preferences struct {
	XMLName                 xml.Name                `xml:"Preferences"`
	EditorSettings          EditorSettings          `xml:"EditorSettings"`
	WorkspaceSettings       WorkspaceSettings       `xml:"WorkspaceSettings"`
	TerminalSettings        TerminalSettings        `xml:"TerminalSettings"`
	GeneralSettings         GeneralSettings         `xml:"GeneralSettings"`
	NetlinxCompilerSettings NetlinxCompilerSettings `xml:"NetlinxCompilerSettings"`
	FileTransferSettings    FileTransferSettings    `xml:"FileTransferSettings"`
	DiagnosticsSettings     DiagnosticsSettings     `xml:"DiagnosticsSettings"`
	TCPIPHistory            TCPIPHistory            `xml:"TCPIPHistory"`
}

// ---------------------------------------------------------------------------
// EditorSettings
// ---------------------------------------------------------------------------

// EditorSettings contains all editor-related preferences.
type EditorSettings struct {
	ClipboardTextBufferMaxItems int     `xml:"ClipboardTextBufferMaxItems"`
	ClipboardTextBufferMaxWidth int     `xml:"ClipboardTextBufferMaxWidth"`
	AutoIndentEnabled           int     `xml:"AutoIndentEnabled"`
	ShowLineNumbersEnabled      int     `xml:"ShowLineNumbersEnabled"`
	CodeFoldingEnabled          int     `xml:"CodeFoldingEnabled"`
	AutoSuggestEnabled          int     `xml:"AutoSuggestEnabled"`
	IndentGuidesEnabled         int     `xml:"IndentGuidesEnabled"`
	CallTipsEnabled             int     `xml:"CallTipsEnabled"`
	UTF8FormatEnabled           int     `xml:"UTF8FormatEnabled"`
	IndentationWidth            int     `xml:"IndentationWidth"`
	TabSpaces                   int     `xml:"TabSpaces"`
	FontPointSize               int     `xml:"FontPointSize"`
	FontName                    string  `xml:"FontName"`
	OnlyFixedPitchFontsEnabled  int     `xml:"OnlyFixedPitchFontsEnabled"`
	PrinterColorMode            int     `xml:"PrinterColorMode"`
	TrimLinesRight              int     `xml:"TrimLinesRight"`
	AutoSuggestInComments       int     `xml:"AutoSuggestInComments"`
	ColumnEdgeMarker            int     `xml:"ColumnEdgeMarker"`
	ColumnEdgeMarkerValue       int     `xml:"ColumnEdgeMarkerValue"`
	RetainBookmarksOnFileClose  int     `xml:"RetainBookmarksOnFileClose"`
	HighlightMatchingBraces     int     `xml:"HighlightMatchingBraces"`
	UseTabEnabled               int     `xml:"UseTabEnabled"`
	Styles                      []Style `xml:"Styles>Style"`
}

// Style represents a single editor colour/font style entry.
type Style struct {
	StyleKey        string `xml:"StyleKey"`
	StyleName       string `xml:"StyleName"`
	StyleType       int    `xml:"StyleType"`
	Bold            int    `xml:"Bold"`
	Italic          int    `xml:"Italic"`
	Underline       int    `xml:"Underline"`
	BackgroundColor int    `xml:"BackgroundColor"`
	ForegroundColor int    `xml:"ForegroundColor"`
	Internal        int    `xml:"Internal"`
}

// ---------------------------------------------------------------------------
// WorkspaceSettings
// ---------------------------------------------------------------------------

// WorkspaceSettings holds workspace-related preferences.
type WorkspaceSettings struct {
	AutoDateStampFileEnabled     int `xml:"AutoDateStampFileEnabled"`
	RestoreWorkspaceEnabled      int `xml:"RestoreWorkspaceEnabled"`
	AutoSaveEnabled              int `xml:"AutoSaveEnabled"`
	SaveBeforeCompileEnabled     int `xml:"SaveBeforeCompileEnabled"`
	EnableCommConfigInSystemID   int `xml:"EnableCommConfigInSystemID"`
	AutoSaveDelay                int `xml:"AutoSaveDelay"`
	EnableWindowTabs             int `xml:"EnableWindowTabs"`
	EnableIconsInTabs            int `xml:"EnableIconsInTabs"`
	EnableCloseWindowInTabs      int `xml:"EnableCloseWindowInTabs"`
	ShowIRTab                    int `xml:"ShowIRTab"`
	ShowZeroConfigTab            int `xml:"ShowZeroConfigTab"`
	PromptToCloseAssociatedFiles int `xml:"PromptToCloseAssociatedFiles"`
	PromptToCloseRemovedFile     int `xml:"PromptToCloseRemovedFile"`
}

// ---------------------------------------------------------------------------
// TerminalSettings
// ---------------------------------------------------------------------------

// TerminalSettings holds terminal and Telnet window preferences.
type TerminalSettings struct {
	TerminalActiveBackgroundClr   int    `xml:"TerminalActiveBackgroundClr"`
	TerminalActiveTextClr         int    `xml:"TerminalActiveTextClr"`
	TerminalInactiveBackgroundClr int    `xml:"TerminalInactiveBackgroundClr"`
	TerminalInactiveTextClr       int    `xml:"TerminalInactiveTextClr"`
	AddLineFeed                   int    `xml:"AddLineFeed"`
	CtrlCharsToPort               int    `xml:"CtrlCharsToPort"`
	TerminalWindowFontName        string `xml:"TerminalWindowFontName"`
	TerminalWindowPointSize       int    `xml:"TerminalWindowPointSize"`
	TelnetWindowFontName          string `xml:"TelnetWindowFontName"`
	TelnetWindowPointSize         int    `xml:"TelnetWindowPointSize"`
	TelnetDefaultPGM              string `xml:"TelnetDefaultPGM"`
}

// ---------------------------------------------------------------------------
// GeneralSettings
// ---------------------------------------------------------------------------

// GeneralSettings holds general application and window colour preferences.
type GeneralSettings struct {
	MRUListSizxeTesting        int    `xml:"MRUListSizxeTesting"`
	OLTShowPortCounts          int    `xml:"OLTShowPortCounts"`
	OLTShowNexusDevices        int    `xml:"OLTShowNexusDevices"`
	OnlineTreeFontName         string `xml:"OnlineTreeFontName"`
	OnlineTreePointSize        int    `xml:"OnlineTreePointSize"`
	ZeroConfigFontName         string `xml:"ZeroConfigFontName"`
	ZeroConfigPointSize        int    `xml:"ZeroConfigPointSize"`
	StatusWindowFontName       string `xml:"StatusWindowFontName"`
	StatusWindowPointSize      int    `xml:"StatusWindowPointSize"`
	AsyncWindowFontName        string `xml:"AsyncWindowFontName"`
	AsyncWindowPointSize       int    `xml:"AsyncWindowPointSize"`
	DiagsWindowFontName        string `xml:"DiagsWindowFontName"`
	DiagsWindowPointSize       int    `xml:"DiagsWindowPointSize"`
	ApplicationLook            int    `xml:"ApplicationLook"`
	GCOStatusFGC               int    `xml:"GCO_STATUS_FGC"`
	GCOFindInFilesFGC          int    `xml:"GCO_FINDINFILES_FGC"`
	GCOFindIRLFGC              int    `xml:"GCO_FINDIRL_FGC"`
	GCOFileTransferStatusFGC   int    `xml:"GCO_FILETRANSFERSTATUS_FGC"`
	GCOAsyncNotificationsFGC   int    `xml:"GCO_ASYNCNOTIFICATIONS_FGC"`
	GCODiagnosticsFGC          int    `xml:"GCO_DIAGNOSTICS_FGC"`
	GCOStatusBGC               int    `xml:"GCO_STATUS_BGC"`
	GCOFindInFilesBGC          int    `xml:"GCO_FINDINFILES_BGC"`
	GCOFindIRLBGC              int    `xml:"GCO_FINDIRL_BGC"`
	GCOFileTransferStatusBGC   int    `xml:"GCO_FILETRANSFERSTATUS_BGC"`
	GCOAsyncNotificationsBGC   int    `xml:"GCO_ASYNCNOTIFICATIONS_BGC"`
	GCODiagnosticsBGC          int    `xml:"GCO_DIAGNOSTICS_BGC"`
	GCOWorkspaceProjectTreeBGC int    `xml:"GCO_WORKSPACE_PROJECT_TREE_BGC"`
	GCOWorkspaceProjectTreeFGC int    `xml:"GCO_WORKSPACE_PROJECT_TREE_FGC"`
	GCOOnlineTreeBGC           int    `xml:"GCO_ONLINE_TREE_BGC"`
	GCOOnlineTreeFGC           int    `xml:"GCO_ONLINE_TREE_FGC"`
	GCOWatchWindowBGC          int    `xml:"GCO_WATCH_WINDOW_BGC"`
	GCOWatchWindowFGC          int    `xml:"GCO_WATCH_WINDOW_FGC"`
	GCOTelnetWindowBGC         int    `xml:"GCO_TELNET_WINDOW_BGC"`
	GCOTelnetWindowFGC         int    `xml:"GCO_TELNET_WINDOW_FGC"`
	GCOZeroConfigWindowBGC     int    `xml:"GCO_ZEROCONFIG_WINDOW_BGC"`
	GCOZeroConfigWindowFGC     int    `xml:"GCO_ZEROCONFIG_WINDOW_FGC"`
}

// ---------------------------------------------------------------------------
// NetlinxCompilerSettings
// ---------------------------------------------------------------------------

// NetlinxCompilerSettings holds NetLinx compiler preferences, including
// dynamically numbered library, include, and module directory lists
// (NLXLibraryDir000, NLXIncludeDir000, NLXModuleDir000, etc.).
type NetlinxCompilerSettings struct {
	BuildWithSource               int
	BuildWithDebugInfo            int
	PasswordProtect               int
	EnableWCPreprocessor          int
	Password                      string
	ShowDebugWindowOnSessionClose int
	ShowMainAXSOnSessionStart     int
	LibraryDirs                   []string // NLXLibraryDir000, NLXLibraryDir001, ...
	IncludeDirs                   []string // NLXIncludeDir000, NLXIncludeDir001, ...
	ModuleDirs                    []string // NLXModuleDir000, NLXModuleDir001, ...
}

// UnmarshalXML implements xml.Unmarshaler. It decodes fixed fields and
// accumulates the dynamically numbered NLXLibraryDir*, NLXIncludeDir*, and
// NLXModuleDir* elements into their respective slices in index order.
func (s *NetlinxCompilerSettings) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			var val string
			if err := d.DecodeElement(&val, &t); err != nil {
				return err
			}

			name := t.Name.Local
			switch {
			case name == "BuildWithSource":
				s.BuildWithSource, _ = strconv.Atoi(val)
			case name == "BuildWithDebugInfo":
				s.BuildWithDebugInfo, _ = strconv.Atoi(val)
			case name == "PasswordProtect":
				s.PasswordProtect, _ = strconv.Atoi(val)
			case name == "EnableWCPreprocessor":
				s.EnableWCPreprocessor, _ = strconv.Atoi(val)
			case name == "Password":
				s.Password = val
			case name == "ShowDebugWindowOnSessionClose":
				s.ShowDebugWindowOnSessionClose, _ = strconv.Atoi(val)
			case name == "ShowMainAXSOnSessionStart":
				s.ShowMainAXSOnSessionStart, _ = strconv.Atoi(val)
			case strings.HasPrefix(name, "NLXLibraryDir"):
				s.LibraryDirs = append(s.LibraryDirs, val)
			case strings.HasPrefix(name, "NLXIncludeDir"):
				s.IncludeDirs = append(s.IncludeDirs, val)
			case strings.HasPrefix(name, "NLXModuleDir"):
				s.ModuleDirs = append(s.ModuleDirs, val)
			}

		case xml.EndElement:
			return nil
		}
	}
}

// MarshalXML implements xml.Marshaler. It writes fixed fields first, then the
// dynamically numbered directory lists in the order NLXLibraryDir, NLXIncludeDir,
// NLXModuleDir, each zero-padded to three digits.
func (s NetlinxCompilerSettings) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// enc encodes a single element, short-circuiting on first error.
	var encErr error
	enc := func(name string, v any) {
		if encErr != nil {
			return
		}
		encErr = e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: name}})
	}

	enc("BuildWithSource", s.BuildWithSource)
	enc("BuildWithDebugInfo", s.BuildWithDebugInfo)
	enc("PasswordProtect", s.PasswordProtect)
	enc("EnableWCPreprocessor", s.EnableWCPreprocessor)
	enc("Password", s.Password)
	enc("ShowDebugWindowOnSessionClose", s.ShowDebugWindowOnSessionClose)
	enc("ShowMainAXSOnSessionStart", s.ShowMainAXSOnSessionStart)

	for i, dir := range s.LibraryDirs {
		enc(fmt.Sprintf("NLXLibraryDir%03d", i), dir)
	}
	for i, dir := range s.IncludeDirs {
		enc(fmt.Sprintf("NLXIncludeDir%03d", i), dir)
	}
	for i, dir := range s.ModuleDirs {
		enc(fmt.Sprintf("NLXModuleDir%03d", i), dir)
	}

	if encErr != nil {
		return encErr
	}

	return e.EncodeToken(start.End())
}

// ---------------------------------------------------------------------------
// FileTransferSettings
// ---------------------------------------------------------------------------

// FileTransferSettings holds file transfer preferences.
type FileTransferSettings struct {
	TKNReboot                         int `xml:"TKNReboot"`
	XDDReboot                         int `xml:"XDDReboot"`
	JARReboot                         int `xml:"JARReboot"`
	TP4SmartTransfer                  int `xml:"TP4SmartTransfer"`
	TP5SmartTransfer                  int `xml:"TP5SmartTransfer"`
	AutoSendSRC                       int `xml:"AutoSendSRC"`
	ReportNumberofItemsLoadedFromList int `xml:"ReportNumberofItemsLoadedFromList"`
	HttpServerPort                    int `xml:"HttpServerPort"`
	HttpServerSelected                int `xml:"HttpServerSelected"`
}

// ---------------------------------------------------------------------------
// DiagnosticsSettings
// ---------------------------------------------------------------------------

// DiagnosticsSettings holds diagnostics and notifications window preferences.
type DiagnosticsSettings struct {
	DisplayLineNumbersNotification      int `xml:"DisplayLineNumbersNotification"`
	DisplayLineNumbersDiagnostics       int `xml:"DisplayLineNumbersDiagnostics"`
	DisplayTimeStampDiagnostics         int `xml:"DisplayTimeStampDiagnostics"`
	DisplayTimeStampNotification        int `xml:"DisplayTimeStampNotification"`
	DisplayTimeMillisecondsDiagnostics  int `xml:"DisplayTimeMillisecondsDiagnostics"`
	DisplayTimeMillisecondsNotification int `xml:"DisplayTimeMillisecondsNotification"`
	DisplayMessagesNotification         int `xml:"DisplayMessagesNotification"`
	DisplayMessagesDiagnostics          int `xml:"DisplayMessagesDiagnostics"`
	LinesToRead                         int `xml:"LinesToRead"`
	BufferFileSize                      int `xml:"BufferFileSize"`
	DisplayDateStampNotification        int `xml:"DisplayDateStampNotification"`
	DisplayDateStampDiagnostics         int `xml:"DisplayDateStampDiagnostics"`
}

// ---------------------------------------------------------------------------
// TCPIPHistory
// ---------------------------------------------------------------------------

// TCPIPEntry represents a single TCP/IP connection history entry. The raw
// format stored in the EPX file is a pipe-delimited string:
//
//	host|port|pingTest|name|username|password
//
// where pingTest is 1 (true) or 0 (false), and username/password are stored
// as-is (NetLinx Studio uses base64 encoding internally).
type TCPIPEntry struct {
	Host     string
	Port     int
	PingTest bool
	Name     string
	Username string
	Password string
}

func parseTCPIPEntry(raw string) TCPIPEntry {
	const fieldCount = 6
	parts := strings.SplitN(raw, "|", fieldCount)
	for len(parts) < fieldCount {
		parts = append(parts, "")
	}

	port, _ := strconv.Atoi(parts[1])
	ping, _ := strconv.Atoi(parts[2])

	return TCPIPEntry{
		Host:     parts[0],
		Port:     port,
		PingTest: ping != 0,
		Name:     parts[3],
		Username: parts[4],
		Password: parts[5],
	}
}

func (e TCPIPEntry) encode() string {
	pingInt := 0
	if e.PingTest {
		pingInt = 1
	}

	return fmt.Sprintf("%s|%d|%d|%s|%s|%s", e.Host, e.Port, pingInt, e.Name, e.Username, e.Password)
}

// TCPIPHistory holds the ordered list of TCP/IP connection history entries.
// The entries are indexed as TCPIPHistoryEX0, TCPIPHistoryEX1, etc. in the
// EPX file; index order is preserved on both parse and write.
type TCPIPHistory struct {
	Entries []TCPIPEntry
}

// UnmarshalXML implements xml.Unmarshaler. It reads each TCPIPHistoryEX{n}
// element in document order and parses its pipe-delimited value.
func (h *TCPIPHistory) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if strings.HasPrefix(t.Name.Local, "TCPIPHistoryEX") {
				var raw string
				if err := d.DecodeElement(&raw, &t); err != nil {
					return err
				}

				h.Entries = append(h.Entries, parseTCPIPEntry(raw))
			} else {
				// Skip any unrecognised child elements.
				if err := d.Skip(); err != nil {
					return err
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}

// MarshalXML implements xml.Marshaler. It writes each entry as
// TCPIPHistoryEX{n} with a pipe-delimited value.
func (h TCPIPHistory) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for i, entry := range h.Entries {
		name := fmt.Sprintf("TCPIPHistoryEX%d", i)
		if err := e.EncodeElement(entry.encode(), xml.StartElement{Name: xml.Name{Local: name}}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}
