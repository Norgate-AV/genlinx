# NetLinx Studio EPX ↔ Registry Mapping

Deep analysis mapping every EPX (`<Preferences>`) element to its corresponding Windows registry
key-value pair.

**Registry base path:** `HKCU\Software\AMX Corp.\NetLinx Studio\`  
**Compiler dirs base:** `HKLM\SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\`

---

## 1. `<EditorSettings>` → `Editor Preferences`

Registry key: `HKCU\...\Editor Preferences`

| EPX Element                  | Registry Value Name        | Type        | Notes                                 |
| ---------------------------- | -------------------------- | ----------- | ------------------------------------- |
| `AutoIndentEnabled`          | `AutoIdent`                | `REG_DWORD` | Renamed                               |
| `ShowLineNumbersEnabled`     | `ShowLineNumbers`          | `REG_DWORD` | Renamed                               |
| `CodeFoldingEnabled`         | `CodeFolding`              | `REG_DWORD` | Renamed                               |
| `AutoSuggestEnabled`         | `AutoSuggest`              | `REG_DWORD` | Renamed                               |
| `IndentGuidesEnabled`        | `IndentGuides`             | `REG_DWORD` | Renamed                               |
| `CallTipsEnabled`            | `CallTips`                 | `REG_DWORD` | Renamed                               |
| `UTF8FormatEnabled`          | `UTF8FormatEnabled`        | `REG_DWORD` | Identical                             |
| `IndentationWidth`           | `IndentWidth`              | `REG_DWORD` | Renamed                               |
| `TabSpaces`                  | `TabWidth`                 | `REG_DWORD` | Renamed                               |
| `FontPointSize`              | `PointSize`                | `REG_DWORD` | Renamed                               |
| `FontName`                   | `FontName`                 | `REG_SZ`    | Identical                             |
| `OnlyFixedPitchFontsEnabled` | `UseOnlyFixedPitchFonts`   | `REG_DWORD` | Renamed                               |
| `PrinterColorMode`           | `PrinterColorMode`         | `REG_DWORD` | Identical                             |
| `TrimLinesRight`             | `RightTrimLines`           | `REG_DWORD` | Renamed                               |
| `AutoSuggestInComments`      | `AutoSuggestInComments`    | `REG_DWORD` | Identical                             |
| `ColumnEdgeMarker`           | `ColumnEdgeMarker`         | `REG_DWORD` | Identical                             |
| `ColumnEdgeMarkerValue`      | `ColumnEdgeColumnNumber`   | `REG_DWORD` | Renamed                               |
| `RetainBookmarksOnFileClose` | `SaveBookmarksOnClose`     | `REG_DWORD` | Renamed                               |
| `HighlightMatchingBraces`    | `HightlightMatchingBraces` | `REG_DWORD` | Registry has typo: `Hight` not `High` |
| `UseTabEnabled`              | `UseTabs`                  | `REG_DWORD` | Renamed                               |

### Clipboard settings → `Text Clipboard History`

Registry key: `HKCU\...\Text Clipboard History`

These two fields live in a separate key — **not** `Editor Preferences`:

| EPX Element                   | Registry Value Name                   | Type        | Notes   |
| ----------------------------- | ------------------------------------- | ----------- | ------- |
| `ClipboardTextBufferMaxItems` | `Clipboard History Max Items`         | `REG_DWORD` | Renamed |
| `ClipboardTextBufferMaxWidth` | `Clipboard History Max Display Width` | `REG_DWORD` | Renamed |

---

## 2. `<EditorSettings><Styles>` → `AXSViewColorPreferences` / `TextViewColorPreferences`

Each `<Style>` element maps to a group of registry values spread across **two** registry keys selected
by the `<StyleKey>` field. The `<StyleType>` integer becomes the numeric suffix on each registry value
name.

Registry keys:

- `HKCU\...\AXSViewColorPreferences`
- `HKCU\...\TextViewColorPreferences`

### Mapping pattern

| EPX Field           | Registry Value Name Pattern                 | Type        |
| ------------------- | ------------------------------------------- | ----------- |
| `StyleKey`          | Determines which registry key               | —           |
| `StyleType`         | Numeric suffix `N` in all value names below | —           |
| `StyleName`         | `StyleName-N`                               | `REG_SZ`    |
| `StyleType` (value) | `StyleType-N`                               | `REG_DWORD` |
| `Bold`              | `Bold-N`                                    | `REG_DWORD` |
| `Italic`            | `Italic-N`                                  | `REG_DWORD` |
| `Underline`         | `Underline-N`                               | `REG_DWORD` |
| `BackgroundColor`   | `BackgroundColor-N`                         | `REG_DWORD` |
| `ForegroundColor`   | `ForegroundColor-N`                         | `REG_DWORD` |
| `Internal`          | `Internal-N`                                | `REG_DWORD` |

### Style index reference — `AXSViewColorPreferences`

| StyleType (N) | StyleName             |
| ------------- | --------------------- |
| 0             | Space                 |
| 1             | Comment               |
| 2             | InComment1            |
| 3             | InComment2            |
| 4             | InComment3            |
| 5             | InComment4            |
| 6             | InComment5            |
| 7             | Reserved Word         |
| 8             | String                |
| 9             | InString              |
| 10            | Operator              |
| 11            | Number                |
| 12            | Float                 |
| 13            | Device                |
| 14            | Constant              |
| 15            | Variable              |
| 16            | Type                  |
| 17            | Default Text          |
| 18            | Selected Text         |
| 19            | Caret                 |
| 20            | Function Names        |
| 21            | Stack/Param Variables |
| 22            | Code-Fold Margin      |
| 23            | Code-Fold Plus        |
| 24            | Code-Fold Minus       |
| 26            | Matching Braces       |
| 32            | Window Background     |
| 33            | Margin                |

### Style index reference — `TextViewColorPreferences`

| StyleType (N) | StyleName         |
| ------------- | ----------------- |
| 17            | Default Text      |
| 18            | Selected Text     |
| 19            | Caret             |
| 32            | Window Background |
| 33            | Margin            |

---

## 3. `<WorkspaceSettings>` → `General Options`

Registry key: `HKCU\...\General Options`

| EPX Element                    | Registry Value Name                       | Type        | Notes                         |
| ------------------------------ | ----------------------------------------- | ----------- | ----------------------------- |
| `AutoDateStampFileEnabled`     | `Enable Auto Time Stamp Source File`      | `REG_DWORD` | Renamed                       |
| `RestoreWorkspaceEnabled`      | `Restore Workspace`                       | `REG_DWORD` | Renamed                       |
| `AutoSaveEnabled`              | `Enable AutoSave`                         | `REG_DWORD` | Renamed                       |
| `SaveBeforeCompileEnabled`     | `Enable Save Before Compile`              | `REG_DWORD` | Renamed                       |
| `EnableCommConfigInSystemID`   | `Display CommConfig In System Identifier` | `REG_DWORD` | Renamed                       |
| `AutoSaveDelay`                | `AutoSave Delay`                          | `REG_DWORD` | Renamed                       |
| `EnableWindowTabs`             | `Enable Window Tabs`                      | `REG_DWORD` | Renamed                       |
| `EnableIconsInTabs`            | `Enable Icons In Tabs`                    | `REG_DWORD` | Renamed                       |
| `EnableCloseWindowInTabs`      | `Enable Close Button In Tabs`             | `REG_DWORD` | Renamed                       |
| `ShowIRTab`                    | `Workspace-ShowIRTab`                     | `REG_DWORD` | Renamed                       |
| `ShowZeroConfigTab`            | `Workspace-ShowZeroConfigTab`             | `REG_DWORD` | Renamed                       |
| `PromptToCloseAssociatedFiles` | `Close Associated Workspace Files`        | `REG_DWORD` | Renamed; values: 1=Yes, 2=Ask |
| `PromptToCloseRemovedFile`     | `Close Removed File`                      | `REG_DWORD` | Renamed; values: 1=Yes, 2=Ask |

---

## 4. `<TerminalSettings>` → `General Options`

Registry key: `HKCU\...\General Options`

| EPX Element                     | Registry Value Name           | Type        | Notes                                                                   |
| ------------------------------- | ----------------------------- | ----------- | ----------------------------------------------------------------------- |
| `TerminalActiveBackgroundClr`   | `Terminal Bckgrnd`            | `REG_DWORD` | Renamed; COLORREF                                                       |
| `TerminalActiveTextClr`         | `Terminal Text`               | `REG_DWORD` | Renamed; COLORREF                                                       |
| `TerminalInactiveBackgroundClr` | `Terminal Inactive Bckgrnd`   | `REG_DWORD` | Renamed; COLORREF                                                       |
| `TerminalInactiveTextClr`       | `Terminal Inactive Text`      | `REG_DWORD` | Renamed; COLORREF                                                       |
| `AddLineFeed`                   | `Terminal Add LF`             | `REG_DWORD` | Renamed                                                                 |
| `CtrlCharsToPort`               | `Terminal Ctrl Chars to Port` | `REG_DWORD` | Renamed                                                                 |
| `TerminalWindowFontName`        | `TerminalWindow-FontNameSize` | `REG_SZ`    | EPX splits font/size; registry stores name only in the FontNameSize key |
| `TerminalWindowPointSize`       | `TerminalWindow-PointSize`    | `REG_DWORD` | —                                                                       |
| `TelnetWindowFontName`          | `TelnetWindow-FontNameSize`   | `REG_SZ`    | —                                                                       |
| `TelnetWindowPointSize`         | `TelnetWindow-PointSize`      | `REG_DWORD` | —                                                                       |
| `TelnetDefaultPGM`              | `Default TELNET Program`      | `REG_SZ`    | Renamed                                                                 |

---

## 5. `<GeneralSettings>` → `General Options`

Registry key: `HKCU\...\General Options`

| EPX Element                      | Registry Value Name                  | Type        | Notes                                                |
| -------------------------------- | ------------------------------------ | ----------- | ---------------------------------------------------- |
| `MRUListSizxeTesting`            | `MaxMRUList`                         | `REG_DWORD` | EPX element name appears to contain a typo (`Sizxe`) |
| `OLTShowPortCounts`              | `OnlineTree-ShowOnlyPortCounts`      | `REG_DWORD` | Renamed                                              |
| `OLTShowNexusDevices`            | `OnlineTree-ShowNexusDevices`        | `REG_DWORD` | Renamed                                              |
| `OnlineTreeFontName`             | `OnlineTree-FontNameSize`            | `REG_SZ`    | —                                                    |
| `OnlineTreePointSize`            | `OnlineTree-PointSize`               | `REG_DWORD` | —                                                    |
| `ZeroConfigFontName`             | `ZeroConfig-FontNameSize`            | `REG_SZ`    | —                                                    |
| `ZeroConfigPointSize`            | `ZeroConfig-PointSize`               | `REG_DWORD` | —                                                    |
| `StatusWindowFontName`           | `StatusWindow-FontNameSize`          | `REG_SZ`    | —                                                    |
| `StatusWindowPointSize`          | `StatusWindow-PointSize`             | `REG_DWORD` | —                                                    |
| `AsyncWindowFontName`            | `AsyncWindow-FontNameSize`           | `REG_SZ`    | —                                                    |
| `AsyncWindowPointSize`           | `AsyncWindow-PointSize`              | `REG_DWORD` | —                                                    |
| `DiagsWindowFontName`            | `DiagsWindow-FontNameSize`           | `REG_SZ`    | —                                                    |
| `DiagsWindowPointSize`           | `DiagsWindow-PointSize`              | `REG_DWORD` | —                                                    |
| `ApplicationLook`                | `ApplicationLook`                    | `REG_DWORD` | Identical                                            |
| `GCO_STATUS_FGC`                 | `OutputBar Foreground Color-0`       | `REG_DWORD` | Status bar (index 0); COLORREF                       |
| `GCO_FINDINFILES_FGC`            | `OutputBar Foreground Color-1`       | `REG_DWORD` | Find in Files (index 1); COLORREF                    |
| `GCO_FINDIRL_FGC`                | `OutputBar Foreground Color-2`       | `REG_DWORD` | Find IR files (index 2); COLORREF                    |
| `GCO_FILETRANSFERSTATUS_FGC`     | `OutputBar Foreground Color-3`       | `REG_DWORD` | File Transfer Status (index 3); COLORREF             |
| `GCO_ASYNCNOTIFICATIONS_FGC`     | `OutputBar Foreground Color-4`       | `REG_DWORD` | Async Notifications (index 4); COLORREF              |
| `GCO_DIAGNOSTICS_FGC`            | `OutputBar Foreground Color-5`       | `REG_DWORD` | Diagnostics (index 5); COLORREF                      |
| `GCO_STATUS_BGC`                 | `OutputBar Background Color-0`       | `REG_DWORD` | COLORREF                                             |
| `GCO_FINDINFILES_BGC`            | `OutputBar Background Color-1`       | `REG_DWORD` | COLORREF                                             |
| `GCO_FINDIRL_BGC`                | `OutputBar Background Color-2`       | `REG_DWORD` | COLORREF                                             |
| `GCO_FILETRANSFERSTATUS_BGC`     | `OutputBar Background Color-3`       | `REG_DWORD` | COLORREF                                             |
| `GCO_ASYNCNOTIFICATIONS_BGC`     | `OutputBar Background Color-4`       | `REG_DWORD` | COLORREF                                             |
| `GCO_DIAGNOSTICS_BGC`            | `OutputBar Background Color-5`       | `REG_DWORD` | COLORREF                                             |
| `GCO_WORKSPACE_PROJECT_TREE_BGC` | `Workspace Background Color`         | `REG_DWORD` | COLORREF                                             |
| `GCO_WORKSPACE_PROJECT_TREE_FGC` | `Workspace Foreground Color`         | `REG_DWORD` | COLORREF                                             |
| `GCO_ONLINE_TREE_BGC`            | `OnLine Tree Background Color`       | `REG_DWORD` | COLORREF                                             |
| `GCO_ONLINE_TREE_FGC`            | `OnLine Tree Foreground Color`       | `REG_DWORD` | COLORREF                                             |
| `GCO_WATCH_WINDOW_BGC`           | `Watch Window Background Color`      | `REG_DWORD` | COLORREF                                             |
| `GCO_WATCH_WINDOW_FGC`           | `Watch Window Foreground Color`      | `REG_DWORD` | COLORREF                                             |
| `GCO_TELNET_WINDOW_BGC`          | `Telnet Window Background Color`     | `REG_DWORD` | COLORREF                                             |
| `GCO_TELNET_WINDOW_FGC`          | `Telnet Window Foreground Color`     | `REG_DWORD` | COLORREF                                             |
| `GCO_ZEROCONFIG_WINDOW_BGC`      | `ZeroConfig Window Background Color` | `REG_DWORD` | COLORREF                                             |
| `GCO_ZEROCONFIG_WINDOW_FGC`      | `ZeroConfig Window Foreground Color` | `REG_DWORD` | COLORREF                                             |

### Registry-only — not exported to EPX

| Registry Value Name              | Type        | Notes                                                 |
| -------------------------------- | ----------- | ----------------------------------------------------- |
| `DebugWindow-FontNameSize`       | `REG_SZ`    | Debug output window font name                         |
| `DebugWindow-PointSize`          | `REG_DWORD` | Debug output window font size                         |
| `FileTxWindow-FontNameSize`      | `REG_SZ`    | File Transfer window font name                        |
| `FileTxWindow-PointSize`         | `REG_DWORD` | File Transfer window font size                        |
| `FindInFilesWindow-FontNameSize` | `REG_SZ`    | Find in Files window font name                        |
| `FindInFilesWindow-PointSize`    | `REG_DWORD` | Find in Files window font size                        |
| `FindIrFilesWindow-FontNameSize` | `REG_SZ`    | Find IR Files window font name                        |
| `FindIrFilesWindow-PointSize`    | `REG_DWORD` | Find IR Files window font size                        |
| `Extract Src Dest Dir`           | `REG_SZ`    | Runtime/session state — last extract destination path |
| `Extract Src Dir`                | `REG_SZ`    | Runtime/session state — last extract source path      |

---

## 6. `<NetlinxCompilerSettings>` → Multiple keys

### Compiler options — `HKCU\...\NLXCompiler_Options`

| EPX Element                     | Registry Value Name  | Type        | Notes     |
| ------------------------------- | -------------------- | ----------- | --------- |
| `BuildWithSource`               | `BuildWithSource`    | `REG_DWORD` | Identical |
| `BuildWithDebugInfo`            | `BuildWithDebugInfo` | `REG_DWORD` | Identical |
| `PasswordProtect`               | `PasswordProtect`    | `REG_DWORD` | Identical |
| `EnableWCPreprocessor`          | `EnableWC`           | `REG_DWORD` | Renamed   |
| `Password`                      | `Password`           | `REG_SZ`    | Identical |
| `ShowDebugWindowOnSessionClose` | `ShowDebugWindow`    | `REG_DWORD` | Renamed   |
| `ShowMainAXSOnSessionStart`     | `ShowMainAXS`        | `REG_DWORD` | Renamed   |

### Library directories — `HKLM\...\NLXCompiler_Libs`

| EPX Element        | Registry Value Name   | Type     | Notes          |
| ------------------ | --------------------- | -------- | -------------- |
| `NLXLibraryDir000` | `LibraryDirectory000` | `REG_SZ` | Renamed prefix |
| `NLXLibraryDir001` | `LibraryDirectory001` | `REG_SZ` | (if present)   |

### Include directories — `HKLM\...\NLXCompiler_Includes`

| EPX Element        | Registry Value Name   | Type     | Notes          |
| ------------------ | --------------------- | -------- | -------------- |
| `NLXIncludeDir000` | `IncludeDirectory000` | `REG_SZ` | Renamed prefix |
| `NLXIncludeDir001` | `IncludeDirectory001` | `REG_SZ` | (if present)   |

### Module directories — `HKLM\...\NLXCompiler_Modules`

| EPX Element       | Registry Value Name  | Type     | Notes          |
| ----------------- | -------------------- | -------- | -------------- |
| `NLXModuleDir000` | `ModuleDirectory000` | `REG_SZ` | Renamed prefix |
| `NLXModuleDir001` | `ModuleDirectory001` | `REG_SZ` | —              |
| `NLXModuleDir002` | `ModuleDirectory002` | `REG_SZ` | —              |

> **Note:** Registry locations are under `HKLM` (machine-wide), whereas the compiler _options_ live under `HKCU` (per user). The EPX stores both in the same `<NetlinxCompilerSettings>` element.

---

## 7. `<FileTransferSettings>` → `Batch Transfer User Options`

Registry key: `HKCU\...\Batch Transfer User Options`

| EPX Element                         | Registry Value Name                       | Type        | Notes                   |
| ----------------------------------- | ----------------------------------------- | ----------- | ----------------------- |
| `TKNReboot`                         | `TKN Reboot`                              | `REG_DWORD` | Renamed                 |
| `XDDReboot`                         | `XDD Reboot`                              | `REG_DWORD` | Renamed                 |
| `JARReboot`                         | `JAR Reboot`                              | `REG_DWORD` | Renamed                 |
| `TP4SmartTransfer`                  | `TP4 Smart Transfer`                      | `REG_DWORD` | Renamed                 |
| `TP5SmartTransfer`                  | `TP5 Smart Transfer`                      | `REG_DWORD` | Renamed                 |
| `AutoSendSRC`                       | `Auto Send SRC`                           | `REG_DWORD` | Renamed                 |
| `ReportNumberofItemsLoadedFromList` | `Report Number of Items Loaded From List` | `REG_DWORD` | Renamed                 |
| `HttpServerPort`                    | see `KIT Files` below                     | `REG_DWORD` | Different key — see §7a |
| `HttpServerSelected`                | see `KIT Files` below                     | `REG_DWORD` | Different key — see §7a |

### Registry-only — not exported to EPX

| Registry Value Name               | Type        | Notes                                       |
| --------------------------------- | ----------- | ------------------------------------------- |
| `Send Empty SRC`                  | `REG_DWORD` | Not included in EPX export                  |
| `Remember Last Items Transferred` | `REG_DWORD` | Not included in EPX export                  |
| `Use Relative Path On Save`       | `REG_DWORD` | Not included in EPX export                  |
| `LastLoadListDir`                 | `REG_SZ`    | Runtime/session state — last directory used |

### 7a. `<FileTransferSettings>` (HTTP fields) → `KIT Files`

Registry key: `HKCU\...\KIT Files`

The HTTP server settings are stored in the **KIT Files** key (the firmware/KIT file transfer manager), not with batch transfer options:

| EPX Element          | Registry Value Name  | Type        | Notes     |
| -------------------- | -------------------- | ----------- | --------- |
| `HttpServerPort`     | `HttpServerPort`     | `REG_DWORD` | Identical |
| `HttpServerSelected` | `HttpServerSelected` | `REG_DWORD` | Identical |

### Registry-only — `KIT Files` (not in EPX)

| Registry Value Name        | Type        | Notes                                               |
| -------------------------- | ----------- | --------------------------------------------------- |
| `Directory`                | `REG_SZ`    | Last used firmware directory                        |
| `ColModified`              | `REG_DWORD` | UI column width: Modified date                      |
| `ColName`                  | `REG_DWORD` | UI column width: Name                               |
| `ColSize`                  | `REG_DWORD` | UI column width: Size                               |
| `HttpServerPCIPSelected`   | `REG_SZ`    | Last selected PC IP address for HTTP server         |
| `KITDirectoryHistory`      | `REG_SZ`    | Current directory                                   |
| `KITDirectoryHistory0`–`N` | `REG_SZ`    | MRU list of previously browsed firmware directories |

---

## 8. `<DiagnosticsSettings>` → `Diagnostic Preferences`

Registry key: `HKCU\...\Diagnostic Preferences`

> All element names are **identical** between EPX and registry.

| EPX Element / Registry Value Name     | Type        |
| ------------------------------------- | ----------- |
| `DisplayLineNumbersNotification`      | `REG_DWORD` |
| `DisplayLineNumbersDiagnostics`       | `REG_DWORD` |
| `DisplayTimeStampDiagnostics`         | `REG_DWORD` |
| `DisplayTimeStampNotification`        | `REG_DWORD` |
| `DisplayTimeMillisecondsDiagnostics`  | `REG_DWORD` |
| `DisplayTimeMillisecondsNotification` | `REG_DWORD` |
| `DisplayMessagesNotification`         | `REG_DWORD` |
| `DisplayMessagesDiagnostics`          | `REG_DWORD` |
| `LinesToRead`                         | `REG_DWORD` |
| `BufferFileSize`                      | `REG_DWORD` |
| `DisplayDateStampNotification`        | `REG_DWORD` |
| `DisplayDateStampDiagnostics`         | `REG_DWORD` |

---

## 9. `<TCPIPHistory>` → `RecentConnectionsHistory`

Registry key: `HKCU\...\RecentConnectionsHistory`

### Value format comparison

|              | Format                                                                              |
| ------------ | ----------------------------------------------------------------------------------- |
| **EPX**      | `<TCPIPHistoryEXN>host\|port\|pingTest\|name\|username\|password</TCPIPHistoryEXN>` |
| **Registry** | `Recent Connection HistoryN = "T-host\|port\|pingTest\|name\|username\|password"`   |

The registry prepends a connection-type prefix (`T-` TCP/IP, `S-` serial, `U-` USB/direct) that the EPX omits — the EPX `<TCPIPHistory>` section stores **TCP/IP entries only**.

### Key differences

| Aspect           | EPX `<TCPIPHistory>`            | Registry `RecentConnectionsHistory`               |
| ---------------- | ------------------------------- | ------------------------------------------------- |
| Connection types | TCP/IP only (`T-` entries)      | All types: `T-` TCP, `S-` serial, `U-` USB/direct |
| Entry ordering   | Sorted alphabetically by label  | Most recently used first                          |
| Index numbering  | Sequential `EX0`, `EX1`, ...    | Sequential `0`, `1`, ... but **different order**  |
| Type prefix      | Absent (protocol always TCP/IP) | `T-`, `S-`, `U-` prefix on host field             |
| Value type       | `REG_SZ`                        | `REG_SZ`                                          |

### Registry-only — not exported to EPX

| Registry Value Name                         | Type        | Notes                                    |
| ------------------------------------------- | ----------- | ---------------------------------------- |
| `ColWidth1206-0`                            | `REG_DWORD` | UI column width: Name column             |
| `ColWidth1206-1`                            | `REG_DWORD` | UI column width: Address column          |
| `ColWidth1206-2`                            | `REG_DWORD` | UI column width: Port column             |
| Serial entries (`S-COM1\|…`)                | `REG_SZ`    | Serial connection history not in EPX     |
| USB entries (`U-Model\|ID\|Serial\|MAC\|…`) | `REG_SZ`    | USB/direct connection history not in EPX |

---

## Summary — Settings Not Captured by EPX Export

The EPX export does **not** round-trip all registry state. The following categories are registry-only:

| Category                       | Registry Key                  | Values                                                                                              |
| ------------------------------ | ----------------------------- | --------------------------------------------------------------------------------------------------- |
| Debug window font              | `General Options`             | `DebugWindow-FontNameSize`, `DebugWindow-PointSize`                                                 |
| File Transfer window font      | `General Options`             | `FileTxWindow-FontNameSize`, `FileTxWindow-PointSize`                                               |
| Find in Files window font      | `General Options`             | `FindInFilesWindow-FontNameSize`, `FindInFilesWindow-PointSize`                                     |
| Find IR Files window font      | `General Options`             | `FindIrFilesWindow-FontNameSize`, `FindIrFilesWindow-PointSize`                                     |
| Last extract paths             | `General Options`             | `Extract Src Dir`, `Extract Src Dest Dir`                                                           |
| Batch transfer runtime state   | `Batch Transfer User Options` | `LastLoadListDir`, `Send Empty SRC`, `Remember Last Items Transferred`, `Use Relative Path On Save` |
| Clipboard item content         | `Text Clipboard History`      | `Item0`–`ItemN` (actual clipboard strings — runtime state)                                          |
| KIT file manager state         | `KIT Files`                   | `Directory`, `KITDirectoryHistory*`, `HttpServerPCIPSelected`, column widths                        |
| Connection list column widths  | `RecentConnectionsHistory`    | `ColWidth1206-0/1/2`                                                                                |
| Serial connections history     | `RecentConnectionsHistory`    | All `S-` prefixed entries                                                                           |
| USB/direct connections history | `RecentConnectionsHistory`    | All `U-` prefixed entries                                                                           |

---

## Notable Anomalies

| Issue                   | Detail                                                                                                                                                                                                                 |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Registry typo           | `HightlightMatchingBraces` (registry) vs `HighlightMatchingBraces` (EPX) — missing `l` in `Hight`                                                                                                                      |
| EPX element name typo   | `<MRUListSizxeTesting>` — appears to contain `Sizxe` instead of `Size`                                                                                                                                                 |
| HKCU vs HKLM split      | Compiler _options_ are user-scoped (`HKCU`); compiler _directory lists_ are machine-scoped (`HKLM WOW6432Node`) — EPX treats them as one unified section                                                               |
| Font key naming         | Registry stores `*-FontNameSize` (name only, despite the key name) alongside a separate `*-PointSize`; EPX uses `*FontName` + `*PointSize` — the `FontNameSize` suffix in the registry name is misleading              |
| HTTP server key split   | `HttpServerPort` and `HttpServerSelected` in EPX `<FileTransferSettings>` map to `HKCU\...\KIT Files`, **not** `Batch Transfer User Options` — they belong to the firmware transfer manager, not the batch file sender |
| EPX connection ordering | `<TCPIPHistory>` entries are sorted alphabetically, not by recency — entry indices have no correspondence to `Recent Connection HistoryN` indices                                                                      |
