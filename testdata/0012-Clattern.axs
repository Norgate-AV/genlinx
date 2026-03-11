PROGRAM_NAME='PR-MB-0012-Clattern-DVX-Main'

(***********************************************************)
#DEFINE __MAIN__
#DEFINE USING_RMS
#include 'NAVFoundation.Core.axi'
#include 'NAVFoundation.ArrayUtils.axi'
#include 'NAVFoundation.Math.axi'
#include 'NAVFoundation.UIUtils.axi'
#include 'NAVFoundation.Enova.axi'
#include 'KingstonStandard.axi'

#DEFINE USING_STANDARD_DSP_LECTERN_MIC
#DEFINE USING_STANDARD_DSP_PATCH_MIC
#DEFINE USING_STANDARD_DSP_ROAMING_MIC
#DEFINE USING_STANDARD_DSP_WIRELESS_MIC_1
#DEFINE USING_STANDARD_DSP_WIRELESS_MIC_2
#DEFINE USING_STANDARD_DSP_WIRELESS_MIC_3
#DEFINE USING_STANDARD_DSP_WIRELESS_MIC_4
#DEFINE USING_STANDARD_DSP_CEILING_MICS_1
#DEFINE INHIBIT_STANDARD_DSP_STATE_UI_PROGRAM
#DEFINE INHIBIT_STANDARD_DSP_LEVEL_UI_PROGRAM
#DEFINE USING_DSP_ONLINE_EVENT_CALLBACK
#DEFINE USING_DSP_DATA_INITIALIZED_EVENT_CALLBACK
#DEFINE USING_DSP_STATE_CHANNEL_EVENT_CALLBACK
#include 'StandardDspConfig.axi'

// #DEFINE USING_SESSION_MANAGEMENT_SESSION_END_EVENT_CALLBACK
// #include 'SessionManagement.axi'

#DEFINE USING_CEILING_MIC_1
#include 'CeilingMicManagement.axi'

#DEFINE USING_STANDARD_CAMERA_1
#include 'StandardCameraConfig.axi'

#DEFINE USING_EVENT_SCHEDULER_EVENT_CALLBACK
#include 'EventScheduler.axi'

#DEFINE USING_BLURAY_OBJECT_ONLINE_EVENT_CALLBACK
#include 'StandardBlurayConfig.axi'

/*
 _   _                       _          ___     __
| \ | | ___  _ __ __ _  __ _| |_ ___   / \ \   / /
|  \| |/ _ \| '__/ _` |/ _` | __/ _ \ / _ \ \ / /
| |\  | (_) | | | (_| | (_| | ||  __// ___ \ V /
|_| \_|\___/|_|  \__, |\__,_|\__\___/_/   \_\_/
                 |___/

MIT License

Copyright (c) 2023 Norgate AV Services Limited

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

(***********************************************************)
(*          DEVICE NUMBER DEFINITIONS GO BELOW             *)
(***********************************************************)
DEFINE_DEVICE
/////////////////////////////////////////////////////////////
// IP Devices
/////////////////////////////////////////////////////////////
dvDocCam                            =   0:3:0       // Wolfvision VZ-8
dvSwitcher_1                        =   0:5:0       // Lightware UVX-4x2
dvSwitcher_2                        =   0:6:0       // Lightware MMX2-4x1
dvScaler                            =   0:7:0       // Extron DVS605


/////////////////////////////////////////////////////////////
// RS232 Devices
/////////////////////////////////////////////////////////////
dvDisplay_1                         =    5001:1:0    // Panasonic Projector
dvLighting                          =    5001:2:0    // Lighting
dvSSP                               =    6001:1:0    // Extron SSP-200



/////////////////////////////////////////////////////////////
// Relays
/////////////////////////////////////////////////////////////
dvRelays                            =    5001:21:0


/////////////////////////////////////////////////////////////
// IR/One-Way RS232 Devices
/////////////////////////////////////////////////////////////
// None


/////////////////////////////////////////////////////////////
// I/O
/////////////////////////////////////////////////////////////
dvIO                                =    5001:22:0


/////////////////////////////////////////////////////////////
// Touch Panels
/////////////////////////////////////////////////////////////
dvTP_Main                           =    10001:1:0
dvTP_SSP_Fader_Program              =    10001:2:0
dvTP_Lighting                       =    10001:12:0
dvTP_Display                        =    10001:13:0
dvTP_Doc_Cam                        =    10001:16:0


/////////////////////////////////////////////////////////////
// Keypads
/////////////////////////////////////////////////////////////
// None


/////////////////////////////////////////////////////////////
// Netlinx Devices
/////////////////////////////////////////////////////////////
dvMatrix_Port_1                     =    5002:1:0
dvMatrix_Port_2                     =    5002:2:0
dvMatrix_Port_3                     =    5002:3:0
dvMatrix_Port_4                     =    5002:4:0
dvMatrix_Port_5                     =    5002:5:0
dvMatrix_Port_6                     =    5002:6:0
dvMatrix_Port_7                     =    5002:7:0
dvMatrix_Port_8                     =    5002:8:0
dvMatrix_Port_9                     =    5002:9:0
dvMatrix_Port_10                    =    5002:10:0
dvMatrix_Port_11                    =    5002:11:0
dvMatrix_Port_12                    =    5002:12:0
dvMatrix_Port_13                    =    5002:13:0
dvMatrix_Port_14                    =    5002:14:0


/////////////////////////////////////////////////////////////
// Virtual Devices
/////////////////////////////////////////////////////////////
vdvConfigFromFile                   =   33001:1:0

vdvMatrix                           =   33201:1:0
vdvDisplay_1                        =   33202:1:0
vdvLighting                         =   33203:1:0
vdvScaler                           =   33205:1:0
vdvSSP                              =   33206:1:0

vdvDocCam                        =   33294:1:0

vdvSwitcher_1                    =   33295:1:0
vdvSwitcher_2                    =   33296:1:0


(***********************************************************)
(*               CONSTANT DEFINITIONS GO BELOW             *)
(***********************************************************)
DEFINE_CONSTANT
/////////////////////////////////////////////////////////////
// Room Defs
/////////////////////////////////////////////////////////////
constant sinteger ROOM_START_UP_SHUT_DOWN_TIMES[]    = { 90, 10 }


/////////////////////////////////////////////////////////////
// Timeline Defs
/////////////////////////////////////////////////////////////
constant long TL_START_UP       = 1
constant long TL_SHUT_DOWN      = 2

constant long TL_POWER_CYCLE_INTERVAL[]    = { 100 }


/////////////////////////////////////////////////////////////
// IO Defs
/////////////////////////////////////////////////////////////
constant integer IO_FIRE_ALARM    = 1


/////////////////////////////////////////////////////////////
// System Mode Defs
/////////////////////////////////////////////////////////////
constant integer SYSTEM_MODE_SINGLE     = 1
constant integer SYSTEM_MODE_DUAL       = 2


/////////////////////////////////////////////////////////////
// Dual Mode Defs
/////////////////////////////////////////////////////////////
constant integer DUAL_MODE_1    = 1
constant integer DUAL_MODE_2    = 2
constant integer DUAL_MODE_3    = 3
constant integer DUAL_MODE_4    = 4
constant integer DUAL_MODE_5    = 5
constant integer DUAL_MODE_6    = 6
constant integer DUAL_MODE_7    = 7
constant integer DUAL_MODE_8    = 8
constant integer DUAL_MODE_9    = 9
constant integer DUAL_MODE_10    = 10

constant integer DEFAULT_DUAL_MODE    = DUAL_MODE_9


/////////////////////////////////////////////////////////////
// Source Defs
/////////////////////////////////////////////////////////////
constant integer SOURCE_LAPTOP      = 1
constant integer SOURCE_DOC_CAM     = 2
constant integer SOURCE_PC          = 3
constant integer SOURCE_BLURAY      = 4
constant integer SOURCE_WIRELESS    = 5
constant integer SOURCE_SCALER      = 6

constant char SOURCE_NAME[][NAV_MAX_CHARS]    =    {
                                                    'Laptop',
                                                    'Visualiser',
                                                    'PC',
                                                    'Blu-ray',
                                                    'BYOD'
                                                }

constant integer SOURCE_DEFAULT = SOURCE_PC

constant integer SOURCE_LAPTOP_1    = 1
constant integer SOURCE_LAPTOP_2    = 2
constant integer SOURCE_LAPTOP_3    = 3

constant integer SOURCE_LAPTOP_DEFAULT = SOURCE_LAPTOP_1


/////////////////////////////////////////////////////////////
// Display Defs
/////////////////////////////////////////////////////////////
constant integer NUMBER_OF_DISPLAYS    = 1
constant integer DISPLAY_1 = 1

constant dev DVA_DISPLAYS[] =   {
                                    vdvDisplay_1
                                }

constant char DISPLAY_INPUT_FOR_SOURCE[][NAV_MAX_CHARS]    =    {
                                                                    'HDMI,1',
                                                                    'HDMI,1',
                                                                    'HDMI,1',
                                                                    'HDMI,1',
                                                                    'HDMI,1',
                                                                    'HDMI,1'
                                                                }


/////////////////////////////////////////////////////////////
// SSP Defs
/////////////////////////////////////////////////////////////
constant integer SSP_INPUT_FOR_SOURCE[] =   {
                                                01,
                                                01,
                                                01,
                                                01,
                                                01,
                                                01
                                            }


/////////////////////////////////////////////////////////////
// Switcher Defs
/////////////////////////////////////////////////////////////
constant integer NUMBER_OF_SWITCHERS    = 2
constant integer SWITCHER_1 = 1
constant integer SWITCHER_2 = 2

constant dev DVA_SWITCHER[] =   {
                                    vdvSwitcher_1,
                                    vdvSwitcher_2
                                }

constant integer SWITCHER_INPUT_FOR_SOURCE[][]  =   {
                                                        { 00, 00, 04, 00, 00, 00 },
                                                        { 00, 01, 00, 03, 02, 00 }
                                                    }

constant integer SWITCHER_INPUT_FOR_SOURCE_LAPTOP[][]   =   {
                                                                { 03, 02, 00 },
                                                                { 00, 00, 00 }
                                                            }

constant integer SWITCHER_INPUT_MATRIX      = 01

constant integer SWITCHER_OUTPUT_MAIN       = 1
constant integer SWITCHER_OUTPUT_MONITOR    = 2
constant integer SWITCHER_OUTPUTS[][]   =   {
                                                { 01, 02 },
                                                { 01 }
                                            }


/////////////////////////////////////////////////////////////
// Matrix Defs
/////////////////////////////////////////////////////////////
constant dev DVA_MATRIX[]    =  {
                                    dvMatrix_Port_1,
                                    dvMatrix_Port_2,
                                    dvMatrix_Port_3,
                                    dvMatrix_Port_4,
                                    dvMatrix_Port_5,
                                    dvMatrix_Port_6,
                                    dvMatrix_Port_7,
                                    dvMatrix_Port_8,
                                    dvMatrix_Port_9,
                                    dvMatrix_Port_10,
                                    dvMatrix_Port_11,
                                    dvMatrix_Port_12,
                                    dvMatrix_Port_13,
                                    dvMatrix_Port_14
                                }

constant integer MATRIX_INPUT_FOR_SOURCE[][]    =   {
                                                        { 01, 02, 01, 02, 02, 03 }, // Video
                                                        { 01, 02, 01, 02, 02, 00 }  // Audio
                                                    }

constant integer MATRIX_INPUT_FOR_SOURCE_LAPTOP[][]    =    {
                                                                { 01, 01, 01 },
                                                                { 01, 01, 01 }
                                                            }


constant integer MATRIX_VID_OUTPUT_SCALER_INPUT_1    = 1
constant integer MATRIX_VID_OUTPUT_SCALER_INPUT_2    = 2
constant integer MATRIX_VID_OUTPUT_DISPLAY    = 3
constant integer MATRIX_VID_OUTPUT_SSP    = 4
constant integer MATRIX_AUD_OUTPUT_DSP    = 1
constant integer MATRIX_AUD_OUTPUT_SSP    = 2
constant integer MATRIX_AUD_OUTPUT_LOOPBACK    = 4
constant integer MATRIX_OUTPUTS[][]    =    {
                                                { 01, 02, 03, 04 },
                                                { 01, 02, 03, 04 }
                                            }

constant integer MATRIX_MONITOR_SOURCE_FOR_SOURCE[]    = { 07, 07, 07, 08, 07 }

constant integer MATRIX_VIDEO_OUT_MUTE_CHANNEL = 210


/////////////////////////////////////////////////////////////
// Scaler Defs
/////////////////////////////////////////////////////////////
constant integer SCALER_WINDOW_1    = 1
constant integer SCALER_WINDOW_2    = 2

constant integer SCALER_DUAL_WINDOW_PRESET_1    = 1
constant integer SCALER_DUAL_WINDOW_PRESET_2    = 2
constant integer SCALER_DUAL_WINDOW_PRESET_3    = 3
constant integer SCALER_DUAL_WINDOW_PRESET_4    = 4
constant integer SCALER_DUAL_WINDOW_PRESET_5    = 5
constant integer SCALER_DUAL_WINDOW_PRESET_6    = 6
constant integer SCALER_DUAL_WINDOW_PRESET_7    = 7
constant integer SCALER_DUAL_WINDOW_PRESET_8    = 8
constant integer SCALER_DUAL_WINDOW_PRESET_9    = 9
constant integer SCALER_DUAL_WINDOW_PRESET_10    = 10

constant integer SCALER_PRESET_FOR_DUAL_MODE_MAP[]    = { 01, 02, 03, 04, 05, 06, 07, 08, 09, 10 }

constant integer SCALER_MAIN_INPUT  = 3
constant integer SCALER_PIP_INPUT   = 4

constant char SCALER_OUTPUT_RATE[]  = '1080p/60'


/////////////////////////////////////////////////////////////
// Relay Defs
/////////////////////////////////////////////////////////////
constant integer RELAY_HOIST_UP     = 1
constant integer RELAY_HOIST_DOWN   = 2


/////////////////////////////////////////////////////////////
// TP Defs
/////////////////////////////////////////////////////////////
constant dev DVA_TP_MAIN[]                  = { dvTP_Main }
constant dev DVA_TP_SSP_FADER_PROGRAM[]     = { dvTP_SSP_Fader_Program }
constant dev DVA_TP_LIGHTING[]              = { dvTP_Lighting }
constant dev DVA_TP_DISPLAY[]               = { dvTP_Display }
constant dev DVA_TP_DOC_CAM[]               = { dvTP_Doc_Cam }


/////////////////////////////////////////////////////////////
// Page Defs
/////////////////////////////////////////////////////////////
constant integer PAGE_LOGO              = 1
constant integer PAGE_MODE              = 2
constant integer PAGE_MAIN              = 3
constant integer PAGE_STARTING_UP       = 4
constant integer PAGE_SHUTTING_DOWN     = 5
constant integer PAGE_DUAL              = 6
constant char PAGE_NAMES[][NAV_MAX_CHARS]   =   {
                                                    'Logo',
                                                    'Mode',
                                                    'Main',
                                                    'Starting Up',
                                                    'Shutting Down',
                                                    'Dual'
                                                }


/////////////////////////////////////////////////////////////
// Popup Defs
/////////////////////////////////////////////////////////////
constant integer POPUP_LAPTOP       = 1
constant integer POPUP_DOC_CAM      = 2
constant integer POPUP_PC           = 3
constant integer POPUP_BLURAY       = 4
constant integer POPUP_WIRELESS     = 5
constant char POPUP_NAMES[][NAV_MAX_CHARS]  =   {
                                                    'Sources - Laptop',
                                                    'Sources - Doc Cam',
                                                    'Sources - PC',
                                                    'Sources - Bluray',
                                                    'Sources - Wireless'
                                                }


/////////////////////////////////////////////////////////////
// Button Defs
/////////////////////////////////////////////////////////////
constant integer BUTTON_TOUCH_TO_START      = 1
constant integer BUTTON_EXIT                = 2
constant integer BUTTON_SHUT_DOWN_OK        = 3
constant integer BUTTON_SHUT_DOWN_CANCEL    = 4

constant integer BUTTON_SOURCE_LAPTOP       = 31
constant integer BUTTON_SOURCE_DOC_CAM      = 32
constant integer BUTTON_SOURCE_PC           = 33
constant integer BUTTON_SOURCE_BLURAY       = 34
constant integer BUTTON_SOURCE_WIRELESS     = 35
constant integer BUTTON_SOURCES[]   =   {
                                            BUTTON_SOURCE_LAPTOP,
                                            BUTTON_SOURCE_DOC_CAM,
                                            BUTTON_SOURCE_PC,
                                            BUTTON_SOURCE_BLURAY,
                                            BUTTON_SOURCE_WIRELESS
                                        }

constant integer BUTTON_SCALER_WINDOW_1     = 41
constant integer BUTTON_SCALER_WINDOW_2     = 42
constant integer BUTTON_SCALER_WINDOWS[]    =   {
                                                    BUTTON_SCALER_WINDOW_1,
                                                    BUTTON_SCALER_WINDOW_2
                                                }

constant integer BUTTON_SYSTEM_MODE_SINGLE  = 51
constant integer BUTTON_SYSTEM_MODE_DUAL    = 52
constant integer BUTTON_SYSTEM_MODES[]      =   {
                                                    BUTTON_SYSTEM_MODE_SINGLE,
                                                    BUTTON_SYSTEM_MODE_DUAL
                                                }

constant integer BUTTON_AUDIO_WINDOW_1      = 61
constant integer BUTTON_AUDIO_WINDOW_2      = 62
constant integer BUTTON_AUDIO_WINDOWS[]     =   {
                                                    BUTTON_AUDIO_WINDOW_1,
                                                    BUTTON_AUDIO_WINDOW_2
                                                }

constant integer BUTTON_DUAL_MODE_1     = 71
constant integer BUTTON_DUAL_MODE_2     = 72
constant integer BUTTON_DUAL_MODE_3     = 73
constant integer BUTTON_DUAL_MODE_4     = 74
constant integer BUTTON_DUAL_MODE_5     = 75
constant integer BUTTON_DUAL_MODE_6     = 76
constant integer BUTTON_DUAL_MODE_7     = 77
constant integer BUTTON_DUAL_MODE_8     = 78
constant integer BUTTON_DUAL_MODE_9     = 79
constant integer BUTTON_DUAL_MODE_10    = 80
constant integer BUTTON_DUAL_MODES[]    =   {
                                                BUTTON_DUAL_MODE_1,
                                                BUTTON_DUAL_MODE_2,
                                                BUTTON_DUAL_MODE_3,
                                                BUTTON_DUAL_MODE_4,
                                                BUTTON_DUAL_MODE_5,
                                                BUTTON_DUAL_MODE_6,
                                                BUTTON_DUAL_MODE_7,
                                                BUTTON_DUAL_MODE_8,
                                                BUTTON_DUAL_MODE_9,
                                                BUTTON_DUAL_MODE_10
                                            }

constant integer BUTTON_RESET_AUDIO    = 81

constant integer BUTTON_SOURCE_LAPTOP_HDMI      = 91
constant integer BUTTON_SOURCE_LAPTOP_VGA       = 92
constant integer BUTTON_SOURCE_LAPTOP_WEPRESENT = 93
constant integer BUTTON_SOURCE_LAPTOPS[]        =   {
                                                        BUTTON_SOURCE_LAPTOP_HDMI,
                                                        BUTTON_SOURCE_LAPTOP_VGA,
                                                        BUTTON_SOURCE_LAPTOP_WEPRESENT
                                                    }

constant integer BUTTON_HOIST_UP    = 101
constant integer BUTTON_HOIST_DOWN  = 102
constant integer BUTTON_HOIST[]     =   {
                                            BUTTON_HOIST_UP,
                                            BUTTON_HOIST_DOWN
                                        }

constant integer BUTTON_DUAL_MODE   = 111

constant integer BUTTON_AV_MUTE     = 131

(***********************************************************)
(*              DATA TYPE DEFINITIONS GO BELOW             *)
(***********************************************************)
DEFINE_TYPE

struct _LightingConfig {
    char AreaNumber[NAV_MAX_CHARS]
    char IntegrationId[NAV_MAX_CHARS]
}

struct _RoomConfig {
    _StandardConfig StandardConfig
    _LightingConfig LightingConfig
    char DocCamIpAddress[NAV_MAX_CHARS]
    char SwitcherIpAddress[NUMBER_OF_SWITCHERS][NAV_MAX_CHARS]
    char ScalerIpAddress[NAV_MAX_CHARS]
    char BlurayIpAddress[NAV_MAX_CHARS]
    char DspIpAddress[NAV_MAX_CHARS]
    char CameraIpAddress[MAX_CAMERAS][NAV_MAX_CHARS]
    char CeilingMicIpAddress[MAX_CEILING_MICS][NAV_MAX_CHARS]
}

(***********************************************************)
(*               VARIABLE DEFINITIONS GO BELOW             *)
(***********************************************************)
DEFINE_VARIABLE
/////////////////////////////////////////////////////////////
// Config Data
/////////////////////////////////////////////////////////////
volatile _RoomConfig roomConfig


/////////////////////////////////////////////////////////////
// Fire Alarm Data
/////////////////////////////////////////////////////////////
volatile integer fireAlarmState


/////////////////////////////////////////////////////////////
// Page Data
/////////////////////////////////////////////////////////////
volatile integer requiredPage = PAGE_LOGO


/////////////////////////////////////////////////////////////
// Popup Data
/////////////////////////////////////////////////////////////
volatile integer requiredPopup


/////////////////////////////////////////////////////////////
// System Mode Data
/////////////////////////////////////////////////////////////
volatile integer systemMode
volatile integer dualMode


/////////////////////////////////////////////////////////////
// Source Data
/////////////////////////////////////////////////////////////
volatile integer selectedSource
volatile integer selectedSourceLaptop
volatile integer currentSourceSingle
volatile integer currentSourceSingleLaptop
volatile integer currentSourceDual[]    = { 0, 0 }
volatile integer currentSourceDualLaptop[]    = { 0, 0 }


/////////////////////////////////////////////////////////////
// Scaler Window Data
/////////////////////////////////////////////////////////////
volatile integer selectedScalerWindow


/////////////////////////////////////////////////////////////
// Includes
/////////////////////////////////////////////////////////////
#IF_DEFINED USING_RMS
#include 'RMSMainCommon.axi'
#include 'RMSTouchPanelCommon.axi'
#include 'RMSMainPR-MB-0012-Clattern.axi'
#END_IF


(***********************************************************)
(*               LATCHING DEFINITIONS GO BELOW             *)
(***********************************************************)
DEFINE_LATCHING

(***********************************************************)
(*       MUTUALLY EXCLUSIVE DEFINITIONS GO BELOW           *)
(***********************************************************)
DEFINE_MUTUALLY_EXCLUSIVE

(***********************************************************)
(*        SUBROUTINE/FUNCTION DEFINITIONS GO BELOW         *)
(***********************************************************)
(* EXAMPLE: DEFINE_FUNCTION <RETURN_TYPE> <NAME> (<PARAMETERS>) *)
(* EXAMPLE: DEFINE_CALL '<NAME>' (<PARAMETERS>) *)

define_function PanelRefresh() {
    NAVPopupShowArray(DVA_TP_MAIN, 'Header', PAGE_NAMES[PAGE_MAIN])
    NAVPopupKillArray(DVA_TP_MAIN, 'Dialogs - Audio')
    NAVPageArray(DVA_TP_MAIN, PAGE_NAMES[requiredPage])

    switch (requiredPage) {
        case PAGE_LOGO: {
            NAVPopupShowArray(DVA_TP_MAIN, 'Sources - Off', PAGE_NAMES[PAGE_MAIN])
        }
        case PAGE_MAIN: {
            ShowDocCamButtons()
            ShowDualModeButtons((systemMode == SYSTEM_MODE_DUAL))

            if (selectedSource) {
                NAVPopupShowArray(DVA_TP_MAIN, POPUP_NAMES[requiredPopup], PAGE_NAMES[requiredPage])
            }
            else {
                NAVPopupShowArray(DVA_TP_MAIN, 'Sources - Off', PAGE_NAMES[requiredPage])
            }
        }
    }

    if (fireAlarmState) {
        NAVPopupShowArray(DVA_TP_MAIN, 'FIRE', '')
    }
    else {
        NAVPopupKillArray(DVA_TP_MAIN, 'FIRE')
    }
}


define_function ShowDualModeButtons(char state) {
    stack_var integer x

    for (x = 1; x <= length_array(BUTTON_SCALER_WINDOWS); x++) {
        NAVShowButtonArray(DVA_TP_MAIN, BUTTON_SCALER_WINDOWS[x], state)
    }

    for (x = 1; x <= length_array(BUTTON_AUDIO_WINDOWS); x++) {
        NAVShowButtonArray(DVA_TP_MAIN, BUTTON_AUDIO_WINDOWS[x], state)
    }

    NAVShowButtonArray(DVA_TP_MAIN, BUTTON_DUAL_MODE, state)
}


define_function PanelReset() {
    if (!NAVDeviceIsOnline(DVA_TP_MAIN[1])) {
        return
    }

    NAVTextArray(DVA_TP_MAIN, 1, '0', roomConfig.StandardConfig.RoomName)
    NAVTextArray(DVA_TP_MAIN, 100, '0', "__FILE__, ' (Compiled on ',__DATE__,' at ',__TIME__,')'")

    NAVPopupsClearArray(DVA_TP_MAIN)

    NAVDoubleBeepArray(DVA_TP_MAIN)
    PanelRefresh()
}


define_function SetupScaler() {
    NAVSwitch(vdvScaler,
                SCALER_MAIN_INPUT,
                0,
                NAV_SWITCH_LEVEL_ALL)

    NAVCommand(vdvScaler, "'PIP-', itoa(SCALER_PIP_INPUT)")
    NAVCommand(vdvScaler, "'OUTPUT_RATE-', SCALER_OUTPUT_RATE")
}


define_function SelectSource(integer source) {
    selectedSource = source
    requiredPopup = source

    #WARN 'Not sure what is happening with routing to the monitor'
    // NAVSwitch(vdvMatrix,
    //             MATRIX_MONITOR_SOURCE_FOR_SOURCE[source],
    //             MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_MONITOR],
    //             NAV_SWITCH_LEVEL_VID)

    SetupScaler()

    switch (source) {
        case SOURCE_LAPTOP: {
            if (!selectedSourceLaptop) {
                SelectSourceLaptop(SOURCE_LAPTOP_DEFAULT)
            }
            else {
                SelectSourceLaptop(selectedSourceLaptop)
            }
        }
        default: {
            switch (systemMode) {
                case SYSTEM_MODE_SINGLE: {
                    SendSource(source, 0)
                }
                case SYSTEM_MODE_DUAL: {
                    SendSource(source, selectedScalerWindow)
                }
                default: {
                    SelectSystemMode(SYSTEM_MODE_SINGLE)
                }
            }
        }
    }

    PanelRefresh()
}


define_function SetDisplayPower(integer state) {
    switch (state) {
        case true: {
            NAVTimelineStart(TL_START_UP,
                                TL_POWER_CYCLE_INTERVAL,
                                TIMELINE_ABSOLUTE,
                                TIMELINE_REPEAT)
        }
        case false: {
            NAVTimelineStart(TL_SHUT_DOWN,
                                TL_POWER_CYCLE_INTERVAL,
                                TIMELINE_ABSOLUTE,
                                TIMELINE_REPEAT)
        }
    }
}


define_function SetDisplayInput(integer source) {
    NAVInput(vdvDisplay_1, DISPLAY_INPUT_FOR_SOURCE[source])
}


define_function SendSource(integer source, integer output) {
    if (!NAVGetPower(vdvDisplay_1)) {
        SetDisplayPower(true)
    }

    if (source == SOURCE_DOC_CAM) {
        pulse[vdvDocCam, PWR_ON]
    }

    SetDisplayInput(source)

    NAVSwitch(vdvSSP,
                SSP_INPUT_FOR_SOURCE[source],
                0,
                NAV_SWITCH_LEVEL_AUD)

    switch (output) {
        case SCALER_WINDOW_1:
        case SCALER_WINDOW_2: {
            stack_var integer previousSource

            previousSource = currentSourceDual[output]
            currentSourceDual[output] = source

            RouteVideoSource(MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_VID][SOURCE_SCALER],
                                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_DISPLAY])

            #IF_DEFINED USING_RMS
            RmsSourceUsageDeactivateSource(previousSource)
            if (currentSourceDual[output] != currentSourceDual[GetOppositeWindow(output)]) {
                RmsSourceUsageActivateSource(source)
            }
            #END_IF

            switch (source) {
                case SOURCE_LAPTOP: {
                    SendSourceLaptop(selectedSourceLaptop, output)
                }
                default: {
                    stack_var integer x

                    for (x = 1; x <= length_array(DVA_SWITCHER); x++) {
                        NAVSwitch(DVA_SWITCHER[x],
                                    SWITCHER_INPUT_FOR_SOURCE[x][source],
                                    0,
                                    NAV_SWITCH_LEVEL_ALL)
                    }

                    RouteVideoSource(MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_VID][source],
                                        MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][output])
                    RouteAudioSource(MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_AUD][source])
                }
            }
        }
        default: {
            currentSourceSingle = source

            #IF_DEFINED USING_RMS
            RmsSourceUsageDeactivateAllSources()
            RmsSourceUsageActivateSource(source)
            #END_IF

            switch (source) {
                case SOURCE_LAPTOP: {
                    SendSourceLaptop(selectedSourceLaptop, 0)
                }
                default: {
                    stack_var integer x

                    for (x = 1; x <= length_array(DVA_SWITCHER); x++) {
                        NAVSwitch(DVA_SWITCHER[x],
                                    SWITCHER_INPUT_FOR_SOURCE[x][source],
                                    0,
                                    NAV_SWITCH_LEVEL_ALL)
                    }

                    RouteVideoSource(MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_VID][source],
                                        MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_DISPLAY])
                    RouteAudioSource(MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_AUD][source])
                }
            }
        }
    }
}


define_function integer GetOppositeWindow(integer window) {
    switch (window) {
        case SCALER_WINDOW_1: { return SCALER_WINDOW_2 }
        case SCALER_WINDOW_2: { return SCALER_WINDOW_1 }
    }
}


define_function SelectSourceLaptop(integer source) {
    selectedSourceLaptop = source

    switch (systemMode) {
        case SYSTEM_MODE_SINGLE: {
            SendSource(selectedSource, 0)
        }
        case SYSTEM_MODE_DUAL: {
            SendSource(selectedSource, selectedScalerWindow)
        }
    }
}


define_function SendSourceLaptop(integer source, integer output) {
    switch (output) {
        case SCALER_WINDOW_1:
        case SCALER_WINDOW_2: {
            stack_var integer x

            currentSourceDualLaptop[output] = source

            for (x = 1; x <= length_array(DVA_SWITCHER); x++) {
                NAVSwitch(DVA_SWITCHER[x],
                            SWITCHER_INPUT_FOR_SOURCE_LAPTOP[x][source],
                            0,
                            NAV_SWITCH_LEVEL_ALL)
            }

            RouteVideoSource(MATRIX_INPUT_FOR_SOURCE_LAPTOP[NAV_SWITCH_LEVEL_VID][source],
                                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][output])
        }
        default: {
            stack_var integer x

            currentSourceSingleLaptop = source

            for (x = 1; x <= length_array(DVA_SWITCHER); x++) {
                NAVSwitch(DVA_SWITCHER[x],
                            SWITCHER_INPUT_FOR_SOURCE_LAPTOP[x][source],
                            0,
                            NAV_SWITCH_LEVEL_ALL)
            }

            RouteVideoSource(MATRIX_INPUT_FOR_SOURCE_LAPTOP[NAV_SWITCH_LEVEL_VID][source],
                                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_DISPLAY])
        }
    }

    RouteAudioSource(MATRIX_INPUT_FOR_SOURCE_LAPTOP[NAV_SWITCH_LEVEL_AUD][source])
}


define_function RouteVideoSource(integer source, integer output) {
    NAVSwitch(vdvMatrix,
                source,
                output,
                NAV_SWITCH_LEVEL_VID)
}


define_function RouteAudioSource(integer source) {
    NAVSwitch(vdvMatrix,
                source,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_SSP],
                NAV_SWITCH_LEVEL_AUD)
    NAVSwitch(vdvMatrix,
                source,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_DSP],
                NAV_SWITCH_LEVEL_AUD)
    NAVSwitch(vdvMatrix,
                source,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_LOOPBACK],
                NAV_SWITCH_LEVEL_AUD)
}


define_function RouteAudioFromWindow(integer window) {
    if (!currentSourceDual[window]) {
        return
    }

    switch (currentSourceDual[window]) {
        case SOURCE_LAPTOP: {
            NAVSwitch(vdvMatrix,
                        MATRIX_INPUT_FOR_SOURCE_LAPTOP[NAV_SWITCH_LEVEL_AUD][currentSourceDualLaptop[window]],
                        MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_DSP],
                        NAV_SWITCH_LEVEL_AUD)
        }
        default: {
            NAVSwitch(vdvMatrix,
                        MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_AUD][currentSourceDual[window]],
                        MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_DSP],
                        NAV_SWITCH_LEVEL_AUD)
        }
    }
}


define_function ClearDualOutputs() {
    stack_var integer x

    for (x = 1; x <= length_array(currentSourceDual); x++) {
        UnRouteSourceFromDualOutput(x)
    }
}


define_function UnRouteSourceFromDualOutput(integer source) {
    currentSourceDual[source] = 0
    currentSourceDualLaptop[source] = 0

    NAVSwitch(vdvMatrix,
                0,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][source],
                NAV_SWITCH_LEVEL_VID)
}


define_function ShutDown() {
    selectedSource = 0
    selectedSourceLaptop = 0

    currentSourceSingle = 0
    currentSourceSingleLaptop = 0

    currentSourceDual[1] = 0
    currentSourceDual[2] = 0

    currentSourceDualLaptop[1] = 0
    currentSourceDualLaptop[2] = 0

    requiredPopup = 0

    systemMode = 0
    dualMode = DEFAULT_DUAL_MODE

    ClearDualOutputs()

    pulse[vdvDocCam, PWR_OFF]

    NAVSwitch(vdvMatrix,
                MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_VID][SOURCE_PC],
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_DISPLAY],
                NAV_SWITCH_LEVEL_VID)
    NAVSwitch(vdvMatrix,
                0,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_DISPLAY],
                NAV_SWITCH_LEVEL_VID)
    NAVSwitch(vdvMatrix,
                0,
                MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_AUD][MATRIX_AUD_OUTPUT_DSP],
                NAV_SWITCH_LEVEL_AUD)

    #IF_DEFINED USING_RMS
    RmsSourceUsageDeactivateAllSources()
    RmsSystemPowerOff()
    #END_IF

    if (NAVGetPower(vdvDisplay_1)) {
        SetDisplayPower(false)
    }
    else {
        requiredPage = PAGE_LOGO
        PanelRefresh()
    }
}


define_function FireAlarm(integer state) {
    cancel_all_wait

    wait 20 {
        fireAlarmState = !state

        pulse[vdvDSP_Preset, (fireAlarmState % 2) + 1]

        if (fireAlarmState) {
            NAVCommand(vdvSSP, 'MUTE-ON')
            pulse[vdvLighting, 4]
        }
        else {
            NAVCommand(vdvSSP, 'MUTE-OFF')
        }

        if (NAVEnovaIsReady(DVA_MATRIX)) {
            PanelRefresh()
        }
    }
}


define_function SetSystemMode(integer mode) {
    systemMode = mode
}


define_function SelectSystemMode(integer mode) {
    SetSystemMode(mode)

    switch (mode) {
        case SYSTEM_MODE_SINGLE: {
            if (currentSourceDual[SCALER_WINDOW_1]) {
                SelectSource(currentSourceDual[SCALER_WINDOW_1])
            }

            dualMode = 0
            ClearDualOutputs()
        }
        case SYSTEM_MODE_DUAL: {
            if (!dualMode) {
                SelectDualMode(DEFAULT_DUAL_MODE)
            }
        }
    }

    PanelRefresh()
}


define_function SelectDualMode(integer mode) {
    dualMode = mode

    pulse[vdvScaler, SCALER_PRESET_FOR_DUAL_MODE_MAP[mode]]

    if (!selectedScalerWindow) {
        SelectScalerWindow(SCALER_WINDOW_1)
    }

    if (currentSourceSingle) {
        SelectSource(currentSourceSingle)
        SendSource(currentSourceSingle, SCALER_WINDOW_1)
        currentSourceSingle = 0
        currentSourceSingleLaptop = 0
    }
    else {
        if (!currentSourceDual[SCALER_WINDOW_1] && !currentSourceDual[SCALER_WINDOW_2]) {
            SelectSource(SOURCE_PC)
            SendSource(SOURCE_PC, SCALER_WINDOW_1)
            SendSource(SOURCE_DOC_CAM, SCALER_WINDOW_2)
        }
    }
}


define_function SelectScalerWindow(integer window) {
    selectedScalerWindow = window

    if (currentSourceDual[window]) {
        selectedSource = currentSourceDual[window]
        requiredPopup = currentSourceDual[window]
    }
    else {
        selectedSource = 0
        requiredPopup = 0
    }

    PanelRefresh()
}


define_function StopTimeline(long id) {
    switch (id) {
        case TL_START_UP: {
            requiredPage = PAGE_MAIN
        }
        case TL_SHUT_DOWN: {
            requiredPage = PAGE_LOGO
        }
    }

    PanelRefresh()
    NAVTimelineStop(id)

    NAVSendLevelArray(DVA_TP_MAIN, 1, 0)
}


define_function ResetAudio() {
    local_var integer x

    if (![vdvDSP_Main_Comm, DATA_INITIALIZED] || ![vdvSSP, DATA_INITIALIZED]) {
        return
    }

    pulse[vdvDSP_Preset, DSP_PRESET_AUDIO_RESET]

    NAVCommand(vdvSSP, "'VOLUME-ABS,75'")

    wait 50 {
        for (x = 1; x <= length_array(DVA_DSP_LEVEL_OBJECTS); x++) {
            NAVCommand(DVA_DSP_LEVEL_OBJECTS[x], 'VOLUME-HALF')
        }

        for (x = 1; x <= length_array(DVA_DSP_STATE_OBJECTS); x++) {
            NAVCommand(DVA_DSP_STATE_OBJECTS[x],
                    "'MUTE-', upper_string(NAVBooleanToOnOffString(DSP_STATE_OBJECT_DEFAULT_STATE[x]))")
        }
    }

    NAVEnovaSetVideoMuteStateAll(DVA_MATRIX, false)
}


define_function integer HasControllableDocCam() {
    return length_array(roomConfig.DocCamIpAddress) > 0
}


define_function ShowDocCamButtons() {
    NAVShowButtonArray(DVA_TP_DOC_CAM, ZOOM_IN, HasControllableDocCam())
    NAVShowButtonArray(DVA_TP_DOC_CAM, ZOOM_OUT, HasControllableDocCam())
    NAVShowButtonArray(DVA_TP_DOC_CAM, FOCUS_NEAR, HasControllableDocCam())
    NAVShowButtonArray(DVA_TP_DOC_CAM, FOCUS_FAR, HasControllableDocCam())
    NAVShowButtonArray(DVA_TP_DOC_CAM, AUTO_FOCUS, HasControllableDocCam())
}


#IF_DEFINED USING_DSP_ONLINE_EVENT_CALLBACK
define_function DspOnlineEventCallback(tdata data) {
    NAVCommand(data.device, "'PROPERTY-PASSWORD,', DEFAULT_PASSWORD")

    if (!length_array(roomConfig.DspIpAddress)) {
        return
    }

    NAVCommand(data.device, "'PROPERTY-IP_ADDRESS,', roomConfig.DspIpAddress")
}
#END_IF


#IF_DEFINED USING_DSP_DATA_INITIALIZED_EVENT_CALLBACK
define_function DspDataInitializedEventCallback(tchannel channel, char state) {
    if (!state) {
        return
    }

    wait 20 {
        ResetAudio()
    }
}
#END_IF


#IF_DEFINED USING_SESSION_MANAGEMENT_SESSION_END_EVENT_CALLBACK
define_function SessionManagementSessionEndEventCallback() {
    if (roomConfig.StandardConfig.SessionInhibit) {
        return
    }

    ShutDown()
}
#END_IF


#IF_DEFINED USING_DSP_STATE_CHANNEL_EVENT_CALLBACK
define_function DspStateChannelEventCallback(tchannel channel, char state) {
    stack_var integer object

    if (channel.channel != VOL_MUTE_FB) {
        return
    }

    object = NAVFindInArrayDevice(DVA_DSP_STATE_OBJECTS, channel.device)

    switch (DSP_OBJECT_NAME[object]) {
        case 'Ceiling Mics 1': {
            CeilingMicsSetMicMute(state)
        }
    }
}
#END_IF


#IF_DEFINED USING_EVENT_SCHEDULER_EVENT_CALLBACK
define_function EventSchedulerEventCallback(char event[]) {
    switch (event) {
        case 'ShutDown': {
            if (roomConfig.StandardConfig.ShutDownInhibit) {
                return
            }

            ShutDown()
        }
    }
}
#END_IF


#IF_DEFINED USING_BLURAY_OBJECT_ONLINE_EVENT_CALLBACK
define_function BlurayObjectOnlineEventCallback(tdata data) {
    if (!length_array(roomConfig.BlurayIpAddress)) {
        return
    }

    NAVCommand(data.device, "'PROPERTY-IP_ADDRESS,', roomConfig.BlurayIpAddress")
}
#END_IF


(***********************************************************)
(*                STARTUP CODE GOES BELOW                  *)
(***********************************************************)
DEFINE_START {
    StandardConfigInit(roomConfig.StandardConfig)
}

/////////////////////////////////////////////////////////////
// Module Defs
/////////////////////////////////////////////////////////////
define_module 'mEnovaDVX' MatrixComm(vdvMatrix, DVA_MATRIX[1])

define_module 'mPanasonicProjector' DisplayComm(vdvDisplay_1, dvDisplay_1)
define_module 'mGenericProjectorUIArray' DisplayUIComm(DVA_TP_DISPLAY, vdvDisplay_1)

define_module 'mLutronQuantum' LightingComm(vdvLighting, dvLighting)
define_module 'mLutronQuantumUIArray' LightingUIComm(DVA_TP_LIGHTING, vdvLighting)

define_module 'mExtronDVS605' ScalerComm(vdvScaler, dvScaler)

define_module 'mExtronSSP' SSPComm(vdvSSP, dvSSP)
define_module 'mExtronSSPLevelUIArray' SSPFaderProgramUIComm(DVA_TP_SSP_FADER_PROGRAM, vdvSSP)

define_module 'mWolfVisionVisualizer' DocCamComm(vdvDocCam, dvDocCam)
define_module 'mGenericDocCamUIArray' DocCamUIComm(DVA_TP_DOC_CAM, vdvDocCam)

define_module 'mLightwareLW3' Switcher1Comm(vdvSwitcher_1, dvSwitcher_1)
define_module 'mLightwareLW3' Switcher2Comm(vdvSwitcher_2, dvSwitcher_2)

define_module 'mConfigFromFile' ConfigFromFileComm(vdvConfigFromFile)


(***********************************************************)
(*                THE EVENTS GO BELOW                      *)
(***********************************************************)
DEFINE_EVENT

data_event[vdvConfigFromFile] {
    online: {
        NAVCommand(data.device, 'GET_TEXT')
    }
    string: {
        stack_var _NAVSnapiMessage message

        NAVParseSnapiMessage(data.text, message)

        switch (message.Header) {
            case 'LINE': {
                if (!length_array(message.Parameter[2])) {
                    break
                }

                switch (atoi(message.Parameter[1])) {
                    case 1: {
                        roomConfig.StandardConfig.RoomName = message.Parameter[2]

                        #IF_DEFINED USING_RMS
                        if (length_array(roomConfig.StandardConfig.RoomName)) {
                            roomConfig.StandardConfig.RmsConnection.Name = roomConfig.StandardConfig.RoomName

                            // NAVRmsConnectionCopy(roomConfig.StandardConfig.RmsConnection, rmsClient.Connection)
                            NAVRmsAdapterConnectionUpdate(vdvRMS, roomConfig.StandardConfig.RmsConnection)
                        }
                        #END_IF
                    }

                    #IF_DEFINED USING_RMS
                    case 2: {
                        roomConfig.StandardConfig.RmsConnection.Url = lower_string(message.Parameter[2])

                        // NAVRmsConnectionCopy(roomConfig.StandardConfig.RmsConnection, rmsClient.Connection)
                        NAVRmsAdapterConnectionUpdate(vdvRMS, roomConfig.StandardConfig.RmsConnection)
                    }
                    case 3: {
                        roomConfig.StandardConfig.RmsConnection.Password = message.Parameter[2]

                        // NAVRmsConnectionCopy(roomConfig.StandardConfig.RmsConnection, rmsClient.Connection)
                        NAVRmsAdapterConnectionUpdate(vdvRMS, roomConfig.StandardConfig.RmsConnection)
                    }
                    case 4: {
                        roomConfig.StandardConfig.RmsConnection.Enabled = lower_string(message.Parameter[2])

                        // NAVRmsConnectionCopy(roomConfig.StandardConfig.RmsConnection, rmsClient.Connection)
                        NAVRmsAdapterConnectionUpdate(vdvRMS, roomConfig.StandardConfig.RmsConnection)
                    }
                    #END_IF

                    case 5: {
                        roomConfig.StandardConfig.ShutDownTime = message.Parameter[2]

                        #IF_DEFINED __EVENT_SCHEDULER__
                        NAVCommand(vdvEventScheduler, "'UPDATE_EVENT-ShutDown,', roomConfig.StandardConfig.ShutDownTime")
                        #END_IF

                        #IF_DEFINED __SESSION_MANAGEMENT__
                        NAVCommand(vdvSessionManager, "'PROPERTY-SESSION_EDIT_LIMIT,', roomConfig.StandardConfig.ShutDownTime")
                        #END_IF
                    }
                    case 6: {
                        roomConfig.StandardConfig.ShutDownInhibit = (atoi(message.Parameter[2]) == true)
                    }
                    case 7: {
                        roomConfig.StandardConfig.SessionDuration = message.Parameter[2]

                        #IF_DEFINED __SESSION_MANAGEMENT__
                        NAVCommand(vdvSessionManager, "'PROPERTY-DEFAULT_SESSION_DURATION,', roomConfig.StandardConfig.SessionDuration")
                        #END_IF
                    }
                    case 8: {
                        roomConfig.StandardConfig.SessionInhibit = (atoi(message.Parameter[2]) == true)
                    }
                    case 9: {
                        roomConfig.DocCamIpAddress = message.Parameter[2]
                        NAVCommand(vdvDocCam, "'PROPERTY-IP_ADDRESS,', roomConfig.DocCamIpAddress")
                    }
                    case 10: {
                        roomConfig.DspIpAddress = message.Parameter[2]
                        NAVCommand(vdvDSP_Main_Comm, "'PROPERTY-IP_ADDRESS,', roomConfig.DspIpAddress")
                    }
                    case 11: {
                        roomConfig.SwitcherIpAddress[1] = message.Parameter[2]
                        NAVCommand(DVA_SWITCHER[1], "'PROPERTY-IP_ADDRESS,', roomConfig.SwitcherIpAddress[1]")
                    }
                    case 12: {
                        roomConfig.SwitcherIpAddress[2] = message.Parameter[2]
                        NAVCommand(DVA_SWITCHER[2], "'PROPERTY-IP_ADDRESS,', roomConfig.SwitcherIpAddress[2]")
                    }
                    case 13: {
                        roomConfig.BlurayIpAddress = message.Parameter[2]
                        NAVCommand(vdvBluray, "'PROPERTY-IP_ADDRESS,', roomConfig.BlurayIpAddress")
                    }
                    case 14: {
                        roomConfig.ScalerIpAddress = message.Parameter[2]
                        NAVCommand(vdvScaler, "'PROPERTY-IP_ADDRESS,', roomConfig.ScalerIpAddress")
                    }
                    case 15: {
                        roomConfig.LightingConfig.AreaNumber = message.Parameter[2]
                        NAVCommand(vdvLighting, "'PROPERTY-AREA,', roomConfig.LightingConfig.AreaNumber")
                    }
                    case 16: {
                        roomConfig.LightingConfig.IntegrationId = message.Parameter[2]
                        NAVCommand(vdvLighting, "'PROPERTY-ID,', roomConfig.LightingConfig.IntegrationId")
                    }
                    case 17: {
                        roomConfig.CameraIpAddress[CAMERA_1] = message.Parameter[2]
                        NAVCommand(vdvCamera_1, "'PROPERTY-IP_ADDRESS,', roomConfig.CameraIpAddress[CAMERA_1]")
                    }
                }
            }
            case 'DONE': {
                PanelReset()
            }
        }
    }
}


data_event[DVA_TP_MAIN] {
    online: {
        if (NAVEnovaIsReady(DVA_MATRIX)) {
            PanelReset()
        }
        else {
            NAVPage(data.device, 'Init')
        }
    }
}


data_event[vdvDocCam] {
    online: {
        if (length_array(roomConfig.DocCamIpAddress)) {
            NAVCommand(data.device, "'PROPERTY-IP_ADDRESS,', roomConfig.DocCamIpAddress")
        }
    }
}


data_event[DVA_SWITCHER] {
    online: {
        stack_var integer switcher

        switcher = get_last(DVA_SWITCHER)

        NAVCommand(data.device, "'PROPERTY-PASSWORD,', DEFAULT_PASSWORD")

        if (length_array(roomConfig.SwitcherIpAddress[switcher])) {
            NAVCommand(data.device, "'PROPERTY-IP_ADDRESS,', roomConfig.SwitcherIpAddress[switcher]")
        }
    }
    string: {}
}


data_event[vdvScaler] {
    online: {
        NAVCommand(data.device, "'PROPERTY-PASSWORD,', DEFAULT_PASSWORD")

        if (length_array(roomConfig.ScalerIpAddress)) {
            NAVCommand(data.device, "'PROPERTY-IP_ADDRESS,', roomConfig.ScalerIpAddress")
        }
    }
}


channel_event[vdvScaler, DATA_INITIALIZED] {
    on: {
        SetupScaler()
    }
}


data_event[vdvSSP] {
    online: {
        NAVCommand(data.device, "'PROPERTY-PASSWORD,', DEFAULT_PASSWORD")
    }
}


channel_event[vdvSSP, DATA_INITIALIZED] {
    on: {
        wait 20 {
            ResetAudio()
        }
    }
}


data_event[DVA_MATRIX] {
    online: {
        if (data.device == DVA_MATRIX[length_array(DVA_MATRIX)]) {
            #WARN 'Not sure what is happening with routing to the monitor'
            // NAVSwitch(vdvMatrix,
            //             MATRIX_INPUT_FOR_SOURCE[NAV_SWITCH_LEVEL_VID][SOURCE_PC],
            //             MATRIX_OUTPUTS[NAV_SWITCH_LEVEL_VID][MATRIX_VID_OUTPUT_MONITOR],
            //             NAV_SWITCH_LEVEL_VID)

            PanelReset()
        }
    }
}


data_event[dvIO] {
    online: {
        FireAlarm([data.device, IO_FIRE_ALARM])
    }
}


data_event[vdvLighting] {
    online: {
        if (length_array(roomConfig.LightingConfig.AreaNumber)) {
            NAVCommand(data.device, "'PROPERTY-AREA,', roomConfig.LightingConfig.AreaNumber")
        }

        if (length_array(roomConfig.LightingConfig.IntegrationId)) {
            NAVCommand(data.device, "'PROPERTY-ID,', roomConfig.LightingConfig.IntegrationId")
        }
    }
    string: {}
}


button_event[DVA_TP_MAIN, BUTTON_TOUCH_TO_START] {
    push: {
        #IF_DEFINED __SESSION_MANAGEMENT__
        NAVCommand(vdvSessionManager, 'SESSION-START')
        #END_IF

        SetSystemMode(SYSTEM_MODE_SINGLE)

        requiredPage = PAGE_MAIN
        PanelRefresh()

        SelectSource(SOURCE_DEFAULT)
    }
}


button_event[DVA_TP_MAIN, BUTTON_SYSTEM_MODES] {
    push: {
        SelectSystemMode(get_last(BUTTON_SYSTEM_MODES))
    }
}


button_event[DVA_TP_MAIN, BUTTON_EXIT] {
    push: {
        if (NAVGetPower(vdvDisplay_1)) {
            NAVPopupShowArray(DVA_TP_MAIN, 'Dialogs - Shut Down', PAGE_NAMES[PAGE_MAIN])
        }
        else {
            ShutDown()
        }
    }
}


button_event[DVA_TP_MAIN, BUTTON_SHUT_DOWN_CANCEL] {
    push: {
        requiredPage = PAGE_MAIN
        PanelRefresh()
    }
}


button_event[DVA_TP_MAIN, BUTTON_SHUT_DOWN_OK] {
    push: {
        #IF_DEFINED __SESSION_MANAGEMENT__
        NAVCommand(vdvSessionManager, 'SESSION-END_EARLY')
        #END_IF

        ShutDown()
    }
}


button_event[DVA_TP_MAIN, BUTTON_SOURCES] {
    push: {
        SelectSource(get_last(BUTTON_SOURCES))
    }
}


button_event[DVA_TP_MAIN, BUTTON_SCALER_WINDOWS] {
    push: {
        SelectScalerWindow(get_last(BUTTON_SCALER_WINDOWS))
    }
}


button_event[DVA_TP_MAIN, BUTTON_AUDIO_WINDOWS] {
    push: {
        RouteAudioFromWindow(get_last(BUTTON_AUDIO_WINDOWS))
    }
}


button_event[DVA_TP_MAIN, BUTTON_DUAL_MODES] {
    push: {
        SelectDualMode(get_last(BUTTON_DUAL_MODES))
    }
}


button_event[DVA_TP_MAIN, BUTTON_SOURCE_LAPTOPS] {
    push: {
        SelectSourceLaptop(get_last(BUTTON_SOURCE_LAPTOPS))
    }
}


button_event[DVA_TP_MAIN, BUTTON_RESET_AUDIO] {
    push: {
        ResetAudio()
    }
}


channel_event[dvIO, IO_FIRE_ALARM] {
    on: {
        FireAlarm([channel.device, channel.channel])
    }
    off: {
        FireAlarm([channel.device, channel.channel])
    }
}


timeline_event[TL_START_UP]
timeline_event[TL_SHUT_DOWN] {
    select {
        active (!timeline.repetition): {
            switch (timeline.id) {
                case TL_START_UP: {
                    requiredPage = PAGE_STARTING_UP
                }
                case TL_SHUT_DOWN: {
                    requiredPage = PAGE_SHUTTING_DOWN
                    pulse[vdvDisplay_1, PWR_OFF]
                }
            }

            ResetAudio()
            PanelRefresh()
        }
        active (timeline.repetition == (ROOM_START_UP_SHUT_DOWN_TIMES[timeline.id] * 10)): {
            StopTimeline(timeline.id)
        }
        active (true): {
            NAVSendLevelArray(DVA_TP_MAIN,
                                1,
                                type_cast(NAVScaleValue(type_cast(timeline.repetition),
                                                        (ROOM_START_UP_SHUT_DOWN_TIMES[timeline.id] * 10),
                                                        255,
                                                        0)))
        }
    }
}


channel_event[vdvDisplay_1, POWER_FB] {
    off: {
        StopTimeline(TL_SHUT_DOWN)

        #IF_DEFINED USING_RMS
        RmsSystemPowerOff()
        #END_IF
    }
    on: {
        StopTimeline(TL_START_UP)

        #IF_DEFINED USING_RMS
        RmsSystemPowerOn()
        #END_IF
    }
}


button_event[DVA_TP_MAIN, BUTTON_HOIST] {
    push: {
        pulse[dvRelays, get_last(BUTTON_HOIST)]
    }
}


button_event[DVA_TP_MAIN, BUTTON_AV_MUTE] {
    push: {
        [DVA_MATRIX[3], MATRIX_VIDEO_OUT_MUTE_CHANNEL]    = ![DVA_MATRIX[3], MATRIX_VIDEO_OUT_MUTE_CHANNEL]

        if ([DVA_MATRIX[3], MATRIX_VIDEO_OUT_MUTE_CHANNEL]) {
            NAVCommand(vdvSSP, 'MUTE-ON')
        }
        else {
            NAVCommand(vdvSSP, 'MUTE-OFF')
        }
    }
}


button_event[DVA_TP_SSP_FADER_PROGRAM, VOL_MUTE] {
    push: {
        if (![DVA_MATRIX[3], MATRIX_VIDEO_OUT_MUTE_CHANNEL]) {
            if (![vdvSSP, VOL_MUTE_FB]) {
                NAVCommand(vdvSSP, 'MUTE-ON')
            }
            else {
                NAVCommand(vdvSSP, 'MUTE-OFF')
            }
        }
    }
}


timeline_event[TL_NAV_FEEDBACK] {
    NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_DUAL_MODES, dualMode)
    NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_SYSTEM_MODES, systemMode)

    switch (systemMode) {
        case SYSTEM_MODE_SINGLE: {
            NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_SOURCES, selectedSource)
            NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_SOURCE_LAPTOPS, selectedSourceLaptop)
        }
        case SYSTEM_MODE_DUAL: {
            NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_SOURCES, currentSourceDual[selectedScalerWindow])
            NAVFeedbackWithDevArray(DVA_TP_MAIN, BUTTON_SOURCE_LAPTOPS, currentSourceDualLaptop[selectedScalerWindow])
        }
    }

    {
        stack_var integer x

        for (x = 1; x <= length_array(BUTTON_SCALER_WINDOWS); x++) {
            select {
                active (currentSourceDual[x]): {
                    [DVA_TP_MAIN, BUTTON_SCALER_WINDOWS[x]]    =    (selectedScalerWindow == x)
                }
                active (true): {
                    [DVA_TP_MAIN, BUTTON_SCALER_WINDOWS[x]]    =    (NAVBlinker)
                }
            }
        }
    }

    [DVA_TP_SSP_FADER_PROGRAM, VOL_MUTE]    = ([vdvSSP, VOL_MUTE_FB]
                                                && NAVBlinker)
    [DVA_TP_MAIN, BUTTON_AV_MUTE]    =    ([DVA_MATRIX[3], MATRIX_VIDEO_OUT_MUTE_CHANNEL]
                                            && [vdvSSP, VOL_MUTE_FB]
                                            && NAVBlinker)
}


(***********************************************************)
(*                     END OF PROGRAM                      *)
(*        DO NOT PUT ANY CODE BELOW THIS COMMENT           *)
(***********************************************************)
