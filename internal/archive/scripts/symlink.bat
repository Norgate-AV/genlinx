@ECHO OFF

SET ScriptDirectory=%~dp0
SET ScriptName=%~n0
SET ScriptPath="%ScriptDirectory%%ScriptName%.ps1"

PowerShell -NoProfile -ExecutionPolicy Bypass -Command "& {Start-Process PowerShell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File ""%ScriptPath%""' -Verb RunAs}"
