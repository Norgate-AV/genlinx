# NLRC V3.1 - NetLinx Compiler Console Program - September 2015

**NLRC.EXE** is a console application to launch the NetLinx Compiler. There are now three ways to invoke this application via the command line. Below are shown how to invoke the three methods of the application from the command line. The **NLRC.EXE** resides in the `C:\Program Files (x86)\Common Files\AMXShare\COM` directory.

## Single AXS File to Compile Option

This is how the original V1.0 version of the program was invoked.

To Invoke:

```pwsh
NLRC.EXE "C:\AMXProjects\main.axs" "-xC:\AMXProjects\OutputLogFile.log"
```

`C:\AMXProjects\main.axs` - File name to compile. Must be in double quotes and be a fully qualified file name.

`-xC:\AMXProjects\OutputLogFile.log` - Optional. If NOT present, then all output goes to the command window. If present, the file name must be in double quotes and where `-x` must be one of the following options:

- `A` = Append status to the output file. If file does not exist, then the program will create a new one.
- `N` = Create a new output file. Overwrites if the file already exists.

Examples:

- `"-AC:\AMXProjects\OutputLogFile.log"` - Append to the output file.
- `"-NC:\AMXProjects\OutputLogFile.log"` - Create a new output file.

The NetLinx Compiler's Include, Library and Module directory path definitions along with compiler options that are defined with the NetLinx Studio application will be used by this option of the application.

## Configuration File to Compile Multiple AXS Files Option

You can pass in a configuration file name (`.CFG`) that will describe which files to compile and set various compiler options during the compile.

To Invoke:

```pwsh
NLRC -CFG"C:\AMX Projects\ACMECorp.cfg"
```

Where `-CFG` specifies the Configuration File name for the program (quotes are needed for the fully qualified file name).

The **NLRC Configuration File Layout** section of this document describes the various options that are available within the file.

The **Example CFG Files** Section of this document has two CFG files examples to review.

## Configuration File and Additional AXS File to Compile Option

You can pass in a configuration file name that will describe which files to compile along with a single AXS file to compile in addition to what is specified in the CFG file. You are NOT required to specify any AXS files within the CFG file if you want to use the CFG file as a type of "Compile Settings" file or compile library type AXS files using the single AXS file as the main program line and the CFG file as your "library" of AXS files.

Additionally, you may specify in what sequence the single AXS file should be compiled. By default, the single specified AXS file will be compiled **LAST** after the files listed in the CFG are compiled.

If you want the file to be compiled first, you must precede the file name with `-F` parameter.

To Invoke -- Compile the AXS file **LAST**:

```pwsh
NLRC "C:\AMX Projects\main.axs" -CFG"C:\AMX Projects\ACMECorp.cfg"
```

To Invoke -- Compile the AXS file **FIRST**:

```pwsh
NLRC -F"C:\AMX Projects\module.axs" -CFG"C:\AMX Projects\ACMECorp.cfg"
```

Where `-CFG` specifies the Configuration File name for the program (quotes are needed for the file name).

## NLRC CFG File Layout

The Configuration File used by the **NLRC.EXE** program has various types of options to instruction the program what files to compile, what NetLinx compiler options to use, additional path statements and what type of logging to produce.

A semi-colon in column 1 of the file indicates a comment line. Use this to document your configuration file.

Example:

```pwsh
;------------------------------------------------------------------------------
; The NLRC.EXE Configuration File
;
; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies how to
; invoke the NetLinx Compiler via a command line.
;------------------------------------------------------------------------------
```

### MainAXSRootDirectory Key

The first key in the CFG file is the `MainAXSRootDirectory` key. This key informs the NLRC program how to construct the paths for the `AXSFile` and `AdditionalxxxxxPath` keys that can be listed in the CFG file.

Example: `MainAXSRootDirectory=C:\AMXProjects\ACME_Company`

Or you can specify to the use the relative path of the CFG that was invoked:

Example: `MainAXSRootDirectory=-R`

So if the **NLRC.EXE** program was invoked using the `C:\AMXProject\SmallBoardRoom\main.cfg` file, then with the AXS Root Directory value would be set to `C:\AMXProject\SmallBoardRoom` for the duration of the compile process.

### AXSFile Key

The `AXSFile` key is used to specify the AXS files to compile. Order of compile is the order the files are listed. Basically there is no limited to the number of files you can list. You have two options on how to use this key.

1. Specify the fully qualified file name (no quotes are needed).

    Example: `AXSFile=C:\AMXProjects\ACME_Company\QuantumData.axs`

2. Specify a relative path to the file. The `MainAXSRootDirectory` value will be used to construct the full file name.

    Example:

    ```pwsh
    AXSFile=Modules\TV.axs
    AXSFile=DVD.axs
    ```

    So with `MainAXSRootDirectoryKey` was set to `MainAXSRootDirectory=C:\AMXProjects\ACME_Company`, then the files to compile would be `C:\AMXProjects\ACME_Company\Modules\TV.axs` and `C:\AMXProjects\ACME_Company\DVD.axs`.

### Logging Keys

There are a three logging keys that can be specified to log the compiler’s status output.

#### OutputLogFile Key

If the `OutputLogFile` key is present, then you must specify a fully qualified file name for the log file (no quotes are needed). If not present, then no logging is sent to a file.

Example: `OutputLogFile=C:\AMXProjects\ACME_Company\mainCompile.log`

#### OutputLogFileOption Key

Used in conjunction with the `OutputLogFile` key. Two options are available (default is `N`):

- `A` = Append status to the output file. If file does not exist, then the program will create a new one.
- `N` = Create a new output file. Overwrites if the file already exists.

Example: `OutputLogFileOption=N`

#### OutputLogConsoleOption Key

Send the logging to the console window that invoked the **NLRC.EXE** program. Two option are available (default is `N`):

- `Y` = Send log info to the console window.
- `N` = Do no send log info to the console window.

Example: `OutputLogConsoleOption=N`

### Optional Build Keys

There are three optional build keys that can be specified within the CFG file. If these keys are not specified, then the options set via the NetLinx Studio application will be use.

#### BuildWithDebugInformation Key

Compile the AXS files with NetLinx Debug information. Two options are available (`Y` and `N`).

Example: `BuildWithDebugInformation=Y`

#### BuildWithSource Key

Compile the AXS files with NetLinx source code information. Two options are available (`Y` and `N`).

Example: `BuildWithSource=Y`

#### BuildWithWC Key

Enable the `_WC` PreProcessor (wide-character) Compiler for the AXS files. Two options are available (`Y` and `N`).

Example: `BuildWithWC=Y`

### Additional Path Keys

There are three types of additional path keys that can be set within the CFG file: Include, Library and Module. These paths will be used in conjunction with the paths set in the NetLinx Studio application. You can specify up to 50 additional paths for each type (one directory per key up to 50 keys per type). No quotes are needed for the directory names. As with the `AXSFile` key, you can use fully qualified paths or relative paths depending on how you set the `MainAXSRootDirectory` key.

Example: Having the `MainAXSRootDirectory=-R` key set. No leading slash is required for the paths below:

```pwsh
AdditionalIncludePath=Small_Room_Includes
AdditionalIncludePath=GenUtility_Includes

AdditionalModulePath=Small_Room_Modules\TypeX
AdditionalModulePath=My_Duet_Modules\TypeX

AdditionalLibraryPath=General_AMX_Libraries\Network
```

If the **NLRC.EXE** program was invoked using the `C:\AMXProject\ACME_Company\SmallBoardRoom\main.cfg` file, then paths would be expanded accordingly:

```pwsh
C:\AMXProject\ACME_Company\SmallBoardRoom\Small_Room_Includes
C:\AMXProject\ACME_Company\SmallBoardRoom\GenUtility_Includes

C:\AMXProject\ACME_Company\SmallBoardRoom\Small_Room_Modules\TypeX

C:\AMXProject\ACME_Company\SmallBoardRoom\My_Duet_Modules\TypeX
C:\AMXProject\ACME_Company\SmallBoardRoom\General_AMX_Libraries\Network
```

Example: Have the `MainAXSRootDirectory=C:\AMXProject\ACME_Company` key set, then you will need fully qualified path names (no quotes are need):

```pwsh
AdditionalIncludePath=C:\AMXProject\Small_Room_Includes
AdditionalIncludePath=C:\AMXProject\GenUtility_Includes

AdditionalModulePath=C:\AMXProject\Small_Room_Modules\TypeX
AdditionalModulePath=C:\AMXProject\My_Duet_Modules\TypeX

AdditionalLibraryPath=C:\AMXProject\General_AMX_Libraries\Network
```

## Example CFG Files

Below are example CFG files and descriptions of the various options within the CFG file.

### Example 1

```pwsh
;------------------------------------------------------------------------------
; The NLRCExample_1.cfg Configuration File
;
; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies
; how to invoke the NetLinx Compiler with a configuration file via a
; command console window.
;
;   > NLRC -CFG"C:\AMX Projects\NLRCExample_1.cfg"
;
;------------------------------------------------------------------------------

;------------------------------------------------------------------------------
;
;  Main AXS Root Directory Reference
;
;------------------------------------------------------------------------------
MainAXSRootDirectory=C:\AMX Projects\ACME Corporation

;------------------------------------------------------------------------------
;
; AXS files when specifying the MainAXSRootDirectory key above. You can have more
; than one, order of the compile is as written.
;
;------------------------------------------------------------------------------
AXSFile=Modules\TV.axs
AXSFile=Modules\DVD.axs
AXSFile=General Utility\Network\Network.axs

;------------------------------------------------------------------------------
;
; Output Log File and Log File Options.
;
; OutputLogFile=        <--: Output log file name
;
;    Fully qualified file name (no quotes are needed)
;    If no OutputLogFile key present, then by default, log to the console
;    window.  Unless the OptionLogConsoleOptions= is specified (see below).
;
; OutputLogFileOption=  <--: Output log file option
;
;    A = Append status to the output file. If file does not exist,
;        then the program will create a new one.
;    N = Create a new output file. Overwrites if the file already exists.
;
;  If no OutputLogFileOption key present, then the default is N.
;
; OutputLogConsoleOption= <--: Output Log to the Console
;
;    Y = Send log info to the console.
;    N = Do no send log info to the console.
;------------------------------------------------------------------------------
OutputLogFile=C:\AMXProjects\Example1_Compile.log
OutputLogFileOption=N
OutputLogConsoleOption=N

;------------------------------------------------------------------------------
;
; NetLinx Compiler Option Overrides
;
;   Ability to override the NetLinx Studio Compiler options that are defined
;   within NetLinx Studio.
;
;   Y = Yes   N = No
;
; Comment these options out if you want to use the NetLinx Studio options.
;------------------------------------------------------------------------------
BuildWithDebugInformation=Y
BuildWithSource=Y
BuildWithWC=Y

;------------------------------------------------------------------------------
; Additional Paths
;
; If you need to specify additional paths for the NetLinx compiler, you can add
; the following keys:
;
;    AdditionalIncludePath=
;    AdditionalLibraryPath=
;    AdditionalModulePath=
;
; You can specify upto 50 additional paths for each type (one directory per
; key upto 50 keys per type).  No quotes are needed for the directory names.
;------------------------------------------------------------------------------
AdditionalIncludePath=C:\AMXProjects\Includes
AdditionalModulePath=C:\AMXProjects\Small Room Modules
AdditionalLibraryPath=C:\AMXProjects\Small Room Library
```

### Example 2

```pwsh
;------------------------------------------------------------------------------
; The NLRCExample_2.cfg Configuration File
;
; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies
; how to invoke the NetLinx Compiler with a configuration file via a
; command console window.
;
; To Invoke:
;
;   > NLRC -CFG"C:\AMX Projects\NLRCExample_2.cfg"
;------------------------------------------------------------------------------

;------------------------------------------------------------------------------
;
;  Main AXS Root Directory Reference
;
;  MainAXSRootDirectory=-R ---> Use the relative path of the CFG file.
;
;  With the Invoke statement above, then MainAXRootDirectory=C:\AMX Projects
;
;------------------------------------------------------------------------------
MainAXSRootDirectory=-R

;------------------------------------------------------------------------------
;
; AXS files to compile with fully qualified paths. You can have more than one,
; order of the compile is as written (no quotes are needed for the file names).
;
;------------------------------------------------------------------------------
AXSFile=C:\AMXProjects\MainBoardRoom\QuantumData.axs
AXSFile=C:\AMXProjects\MainBoardRoom\Projector.axs
AXSFile=C:\AMXProjects\MainBoardRoom\main.axs

;------------------------------------------------------------------------------
;
; Output Log File and Log File Options.
;
; OutputLogFile=        <--: Output log file name
;
;    Fully qualified file name (no quotes are needed)
;    If no OutputLogFile key present, then by default, log to the console
;    window.  Unless the OptionLogConsoleOptions= is specified (see below).
;
; OutputLogFileOption=  <--: Output log file option
;
;    A = Append status to the output file. If file does not exist,
;        then the program will create a new one.
;    N = Create a new output file. Overwrites if the file already exists.
;
;  If no OutputLogFileOption key present, then the default is N.
;
; OutputLogConsoleOption= <--: Output Log to the Console
;
;    Y = Send log info to the console.
;    N = Do no send log info to the console.
;------------------------------------------------------------------------------
OutputLogFile=C:\AMXProjects\Example2_Compile.log
OutputLogFileOption=N
OutputLogConsoleOption=N

;------------------------------------------------------------------------------
;
; NetLinx Compiler Option Overrides
;
;   Ability to override the NetLinx Studio Compiler options that are defined
;   within NetLinx Studio.
;
;   Y = Yes   N = No
;
; Comment these options out if you want to use the NetLinx Studio options.
;------------------------------------------------------------------------------
BuildWithDebugInformation=Y
BuildWithSource=Y

;------------------------------------------------------------------------------
; Additional Paths
;
; If you need to specify additional paths for the NetLinx compiler, you can add
; the following keys:
;
;    AdditionalIncludePath=
;    AdditionalLibraryPath=
;    AdditionalModulePath=
;
; You can specify up to 50 additional paths for each type (one directory per
; key up to 50 keys per type).  No quotes are needed for the directory names.
;
; With the MainAXSRootDirectory=-R key defined above.
;
;------------------------------------------------------------------------------
AdditionalIncludePath=Small Room Includes
AdditionalIncludePath=GenUtility Includes

AdditionalModulePath=Small Room Modules
AdditionalModulePath=My Duet Modules

AdditionalLibraryPath=General AMX Libraries\Network
```
