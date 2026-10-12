package apw

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
)

// apwDoctype is the standard APW DOCTYPE declaration, matching the one
// produced by NetLinx Studio verbatim. Parse requires the string
// "<!DOCTYPE Workspace [" to be present; the declarations inside are
// informational for XML validators.
const apwDoctype = `<!DOCTYPE Workspace [

    <!-- Common Elements -->
    <!ELEMENT Identifier (#PCDATA)>
    <!ELEMENT Comments (#PCDATA)>
    <!ELEMENT MasterDirectory (#PCDATA)>
    <!ELEMENT CreationDate (#PCDATA)>

    <!-- Workspace Elements-->
    <!ELEMENT Workspace ( Identifier,
               CreateVersion,
               PJS_File?,
               PJS_ConvertDate?,
               PJS_CreateDate?,
               Comments?,
               Project*)>

    <!ATTLIST Workspace CurrentVersion CDATA #REQUIRED>

    <!ELEMENT CreateVersion (#PCDATA)>


    <!-- Conversion data only: only included for files converted from the old .pjs database -->
    <!ELEMENT PJS_File (#PCDATA)>
    <!ELEMENT PJS_ConvertDate (#PCDATA)>
    <!ELEMENT PJS_CreateDate (#PCDATA)>

    <!ELEMENT Project ( Identifier,
               Designer?,
               DealerID?,
               SalesOrder?,
               PurchaseOrder?,
               Comments?,
               System*)>

    <!ELEMENT Designer (#PCDATA)>
    <!ELEMENT DealerID (#PCDATA)>
    <!ELEMENT SalesOrder (#PCDATA)>
    <!ELEMENT PurchaseOrder (#PCDATA)>


    <!ELEMENT System (  Identifier,
                 SysID,
                 TransTCPIP?,
                 TransSerial?,
                 TransTCPIPEx?,
                 TransSerialEx?,
                 TransUSBEx?,
                 TransVNMEx?,
                 VirtualNetLinxMasterFlag?,
                 VNMSystemID?,
                 VNMIPAddress?,
                 VNMMaskAddress?,
                 UserName?,
                 Password?,
                 Comments?,
                 File*)>

    <!ATTLIST System
        IsActive (true | false) "false"
        Platform (Axcess | Netlinx) "Axcess"
        Transport (Serial | Modem | TCPIP) "Serial"
        TransportEx (Serial | USB | TCPIP | VNM) "Serial">

    <!ELEMENT SysID (#PCDATA)>
    <!ELEMENT TransSerial (#PCDATA)>
    <!ELEMENT TransTCPIP (#PCDATA)>
    <!ELEMENT TransTCPIPEx (#PCDATA)>
    <!ELEMENT TransSerialEx (#PCDATA)>
    <!ELEMENT TransUSBEx (#PCDATA)>
    <!ELEMENT TransVNMEx (#PCDATA)>
    <!ELEMENT VNMSystemID (#PCDATA)>
    <!ELEMENT VNMIPAddress (#PCDATA)>
    <!ELEMENT VNMMaskAddress (#PCDATA)>
    <!ELEMENT VirtualNetLinxMasterFlag (#PCDATA)>
    <!ELEMENT UserName (#PCDATA)>
    <!ELEMENT Password (#PCDATA)>


    <!ELEMENT File ( Identifier,
               FilePathName,
               Comments?,
               MasterDirectory?,
               DeviceMap*,
               IRDB*)>

    <!ATTLIST File
        Type (Source | MasterSrc | Include | Module | AXB | IR | TPD | TP4 | TP5 | KPD | TKO | AMX_IR_DB | IRN_DB | Other | DUET | TOK | TKN | KPB | XDD ) "Other"
        CompileType (Axcess | Netlinx | None) "None">

    <!ELEMENT FilePathName (#PCDATA)>

    <!ELEMENT DeviceMap (DevName)>
    <!ATTLIST DeviceMap DevAddr CDATA #REQUIRED>

    <!ELEMENT DevName (#PCDATA)>

    <!ELEMENT IRDB (Property,
                 DOSName,
                 UserDBPathName,
                 Notes)>
    <!ATTLIST IRDB DBKey CDATA #REQUIRED>

    <!ELEMENT Property (#PCDATA)>
    <!ELEMENT DOSName (#PCDATA)>
    <!ELEMENT UserDBPathName (#PCDATA)>
    <!ELEMENT Notes (#PCDATA)>
]>`

// Marshal serializes ws to a complete .apw document, including the XML
// declaration and DOCTYPE block that Parse requires.
// The result is round-trip safe: Marshal → os.WriteFile → os.ReadFile → Parse
// returns an equivalent workspace.
func Marshal(ws *Workspace) ([]byte, error) {
	if ws == nil {
		return nil, fmt.Errorf("cannot marshal nil workspace")
	}

	body, err := xml.MarshalIndent(ws, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workspace: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header) // "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	buf.WriteString(apwDoctype)
	buf.WriteByte('\n')
	buf.Write(body)
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

// compileTypeFor returns the CompileType appropriate for t.
// Source-code file types (Source, MasterSrc, Include, Module) compile with
// NetLinx; all other file types (binaries, touch panels, IR etc.) use None.
func compileTypeFor(t FileType) FileCompileType {
	switch t {
	case FileTypeSource, FileTypeMasterSrc, FileTypeInclude, FileTypeModule:
		return FileCompileTypeNetLinx
	default:
		return FileCompileTypeNone
	}
}

// Bytes serializes the APW to a complete .apw document.
// It is a thin wrapper around Marshal using the APW's internal workspace.
func (a *APW) Bytes() ([]byte, error) {
	return Marshal(a.ws)
}

// Write serializes the APW and writes it to path, creating or overwriting the
// file with mode 0644.
func (a *APW) Write(path string) error {
	data, err := a.Bytes()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
