#Region ;**** Directives created by AutoIt3Wrapper_GUI ****
#AutoIt3Wrapper_Res_Description=ColoTransfer transfers colors from Fusion360 to PrusaSlicer
#AutoIt3Wrapper_Res_Fileversion=0.1.0.41
#AutoIt3Wrapper_Res_ProductVersion=1.0
#AutoIt3Wrapper_Res_Fileversion_AutoIncrement=y
#AutoIt3Wrapper_Res_ProductName=ColorTransfer
#AutoIt3Wrapper_Res_CompanyName=schiffer-soft
#AutoIt3Wrapper_Res_LegalCopyright=schiffer-soft & MaxKlaxx
#AutoIt3Wrapper_Res_Language=1031
#EndRegion ;**** Directives created by AutoIt3Wrapper_GUI ****

#include <Array.au3>
#include <WinAPIGdi.au3>
#include <File.au3>
#include <buttonconstants.au3>
#include <guiconstantsex.au3>
#include <GDIPlus.au3>
#include <SendMessage.au3>
#include <WinAPISysWin.au3>
#include "Forms\colortransfer.isf"

Opt("GuiOnEventMode", 1)
_GDIPlus_Startup()

DirCreate(@AppDataDir & "\colortransfer")
FileInstall("colortransfer.exe", @AppDataDir & "\colortransfer\colortransfer.exe", 1 )
FileInstall("colortransfer.jpg", @AppDataDir & "\colortransfer\colortransfer.jpg" )

Global Const $ProgName = "ColorTransfer"
Global $PathToInput3mf, $PathToOutput3mf, $PathToInput3mf, $sOutputArray
Global $sExePath = @AppDataDir & "\colortransfer\colortransfer.exe"
Global $FileLoaded =  False 

GUICtrlSetData($h_label_version, "v "& FileGetVersion(@ScriptDir & "\" & @ScriptName))
GUICtrlSetImage($h_image, @AppDataDir & "\colortransfer\colortransfer.jpg")

_toCMD('please select inputfile.')
GUISetState(@SW_SHOW)

While 1
	Sleep(10)
WEnd 

Func _exit()
	Exit 
EndFunc

Func _button_inputfile()
	GUICtrlSetData($h_console, "")
	_toCMD('please select inputfile.')
		
	$PathToInput3mf =  FileOpenDialog($ProgName, @ScriptDir & '\',"3mf (*.3mf)", 1 +2)
	If @error Then 
		GUICtrlSetData($h_inputfile, "Error/abort inputfile.")
		_toCMD('Error/abort inputfile.')
		$PathToInput3mf =  ""
		$FileLoaded =  False 
		GUICtrlSetState($h_button_ConvertSave, $GUI_Disable )
		Return 
	EndIf
	$FileLoaded =  True 
	GUICtrlSetData($h_inputfile, $PathToInput3mf)
	_toCMD('load: "'&$PathToInput3mf & '"')
	_toCMD('open and parse the file.')
	Sleep(10)
	_preload3mf($PathToInput3mf)
	
EndFunc


Func _preload3mf($PathToInput3mf)
	
	GUICtrlSetBkColor($h_lable_c1, 0x555555)
	GUICtrlSetBkColor($h_lable_c2, 0x555555)
	GUICtrlSetBkColor($h_lable_c3, 0x555555)
	GUICtrlSetBkColor($h_lable_c4, 0x555555)
	GUICtrlSetBkColor($h_lable_c5, 0x555555)
	
	GUICtrlSetBkColor($h_lable_e1, 0x555555)
	GUICtrlSetBkColor($h_lable_e2, 0x555555)
	GUICtrlSetBkColor($h_lable_e3, 0x555555)
	GUICtrlSetBkColor($h_lable_e4, 0x555555)
	GUICtrlSetBkColor($h_lable_e5, 0x555555)
	

	Local $iPID = Run($sExePath & ' "' & $PathToInput3mf & '"', @WorkingDir, @SW_HIDE, $STDOUT_CHILD)

	Local $sOutput = ""
	While True
		$sLine = StdOutRead($iPID)
		If @error Then ExitLoop
		$sOutput &= $sLine
	WEnd
	ConsoleWrite($sOutput & @CRLF)
	$sOutputArray =  StringSplit($sOutput,Chr(10) )
;~ 	_ArrayDisplay($sOutputArray)
	
	If StringInStr($sOutput, "Error") Or Not StringInStr($sOutput, "#") Or UBound($sOutputArray) < 7  Then 
		_toCMD('Return: ' & $sOutput)
		GUICtrlSetState($h_button_ConvertSave, $GUI_DISABLE )
		Return 
		
	EndIf
	GUICtrlSetState($h_button_ConvertSave, $GUI_ENABLE )
	_toCMD('Return: Found colors ' & _ArrayToString($sOutputArray, ",", 1, 5))
		
	GUICtrlSetBkColor($h_lable_c1, "0x" & StringTrimLeft($sOutputArray[1], 1))
	GUICtrlSetBkColor($h_lable_c2, "0x" & StringTrimLeft($sOutputArray[2], 1))
	GUICtrlSetBkColor($h_lable_c3, "0x" & StringTrimLeft($sOutputArray[3], 1))
	GUICtrlSetBkColor($h_lable_c4, "0x" & StringTrimLeft($sOutputArray[4], 1))
	GUICtrlSetBkColor($h_lable_c5, "0x" & StringTrimLeft($sOutputArray[5], 1))
	
	GUICtrlSetBkColor($h_lable_e1, "0x" & StringTrimLeft($sOutputArray[1], 1))
	GUICtrlSetBkColor($h_lable_e2, "0x" & StringTrimLeft($sOutputArray[2], 1))
	GUICtrlSetBkColor($h_lable_e3, "0x" & StringTrimLeft($sOutputArray[3], 1))
	GUICtrlSetBkColor($h_lable_e4, "0x" & StringTrimLeft($sOutputArray[4], 1))
	GUICtrlSetBkColor($h_lable_e5, "0x" & StringTrimLeft($sOutputArray[5], 1))
	
	;Vorschaubild setzen
	_toCMD('Loading Fusion 360 preview image.')
	$Bmp_Logo = _GDIPlus_BitmapCreateFromMemory(_Base64StringEx($sOutputArray[$sOutputArray[0] -2]), True)
	$hBitmap = _GDIPlus_BitmapCreateFromHBITMAP($Bmp_Logo)
	$hBitmap2 = _GDIPlus_ImageResize($hBitmap,144,144)
	$hHBitmapM=_GDIPlus_BitmapCreateHBITMAPFromBitmap($hBitmap2)
	_WinAPI_DeleteObject(GUICtrlSendMsg($h_image, $STM_SETIMAGE, $IMAGE_BITMAP, $hHBitmapM))
		
	;Suchen ob die Farbkombination schon genutzt wurde
	$ret =  _GetNewColors(_ArrayToString($sOutputArray, ";", 1, 5), @AppDataDir & "\colortransfer\remembercolors.dat")
	ConsoleWrite("Farben bereits benutzt: " & $ret &@CRLF )
	If $ret <> "" Then 
		_toCMD('Remenber this colors. Set your last selection.')
		$aNewColors = StringSplit($ret, ";")
		GUICtrlSetBkColor($h_lable_e1, "0x" & StringTrimLeft($aNewColors[1], 1))
		GUICtrlSetBkColor($h_lable_e2, "0x" & StringTrimLeft($aNewColors[2], 1))
		GUICtrlSetBkColor($h_lable_e3, "0x" & StringTrimLeft($aNewColors[3], 1))
		GUICtrlSetBkColor($h_lable_e4, "0x" & StringTrimLeft($aNewColors[4], 1))
		GUICtrlSetBkColor($h_lable_e5, "0x" & StringTrimLeft($aNewColors[5], 1))
			
	EndIf
EndFunc

Func _ConvertSave()
	
	$split = StringSplit($PathToInput3mf, "/")
	$split = StringSplit($split[$split[0]], ".", 2)
	$PathToOutput3mf =  FileSaveDialog( $ProgName, @ScriptDir & '\',"3mf (*.3mf)",2 +16,  $split[0] &"_out."& $split[1])
	If @error Then 
		_toCMD('Error/abort outputfile.')
		Return 
	EndIf
	WinActivate($colortransfer)
	_toCMD('Export Path: ' &$PathToOutput3mf)
	
	$ColorsOld =  "#"&  hex(GUICtrlGetBkColor($h_lable_c1), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_c2), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_c3), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_c4), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_c5), 6)
						
	$ColorsNew =   "#"& hex(GUICtrlGetBkColor($h_lable_e1), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_e2), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_e3), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_e4), 6) &";#"& _
						hex(GUICtrlGetBkColor($h_lable_e5), 6)
						
	$ColorOldAndNew =  $ColorsOld &"|"& $ColorsNew
	
	_SaveOrUpdateColors($ColorsOld, $ColorsNew, @AppDataDir & "\colortransfer\remembercolors.dat")
;~ 	$RememberColors =  IniRead(@AppDataDir & "\colortransfer\options.ini", "Remember", "Colors","")
;~ 	If StringInStr($RememberColors,$ColorsOld ) Then 
;~ 	IniWrite(@AppDataDir & "\colortransfer\options.ini", "Remember", "Colors",$ColorOldAndNew )
	
	ConsoleWrite($ColorsNew & @CRLF)
;~ 	_toCMD('New Colors: ' &$ColorsNew)
	
	Local $CMD =  $sExePath & ' "' & $PathToInput3mf & '" "' & $PathToOutput3mf & '" ' & '"'& $ColorsNew& '"'
	ConsoleWrite("CMD: " & $CMD & @CRLF)
	Local $iPID = Run($CMD, @WorkingDir, @SW_HIDE, $STDOUT_CHILD)
	Local $sOutput = ""
	While True
		$sLine = StdOutRead($iPID)
		If @error Then ExitLoop
		$sOutput &= $sLine
	WEnd
	$sOutput =  StringTrimRight($sOutput, 2)
	ConsoleWrite($sOutput & @CRLF)
	_toCMD('Return: ' &@CRLF &@CRLF &StringReplace($sOutput, Chr(10), @CRLF))
	
EndFunc 

Func button_ext1_d()
	label_swap($h_lable_e1, $h_lable_e2)
EndFunc 
Func button_ext2_d()
	label_swap($h_lable_e2, $h_lable_e3)
EndFunc 
Func button_ext3_d()
	label_swap($h_lable_e3, $h_lable_e4)
EndFunc 
Func button_ext4_d()
	label_swap($h_lable_e4, $h_lable_e5)
EndFunc 	
	
Func button_ext2_u()
	label_swap($h_lable_e1, $h_lable_e2)
EndFunc 
Func button_ext3_u()
	label_swap($h_lable_e2, $h_lable_e3)
EndFunc 
Func button_ext4_u()
	label_swap($h_lable_e3, $h_lable_e4)
EndFunc 
Func button_ext5_u()
	label_swap($h_lable_e4, $h_lable_e5)
EndFunc 	
	
Func label_swap($label1, $label2)
	$color_old =  "0x" &hex(GUICtrlGetBkColor($label1))
	GUICtrlSetBkColor($label1, "0x" &hex(GUICtrlGetBkColor($label2)))
	GUICtrlSetBkColor($label2, $color_old)
EndFunc

Func GUICtrlGetBkColor($hWnd)
    If Not IsHWnd($hWnd) Then
        $hWnd = GUICtrlGetHandle($hWnd)
    EndIf
    Local $hDC = _WinAPI_GetDC($hWnd)
    Local $iColor = _WinAPI_GetPixel($hDC, 0, 0)
    _WinAPI_ReleaseDC($hWnd, $hDC)
    Return $iColor
EndFunc   ;==>GUICtrlGetBkColor

Func _toCMD($sData)
    Local $SB_SCROLLCARET = 0xB7
    GUICtrlSetData($h_console, GUICtrlRead($h_console) & $sData & @CRLF)
    GUICtrlSendMsg($h_console, $SB_SCROLLCARET, 0, 0)
    ControlFocus($colortransfer, "", $h_console)
EndFunc

Func _GetNewColors($sOld, $sFile)
    If Not FileExists($sFile) Then Return ""
    Local $aLines = FileReadToArray($sFile)
    If @error Then Return ""
;~ 	_ArrayDisplay($aLines)
    For $i = 0 To UBound($aLines) -1
        Local $aSplit = StringSplit($aLines[$i], "|" ) 
        If $aSplit[0] = 2 Then
            If $aSplit[1] = $sOld Then
                ; Gefunden -> neuen Farbstring zurückgeben
                Return $aSplit[2]
            EndIf
        EndIf
    Next
    Return ""
EndFunc

Func _SaveOrUpdateColors($sOld, $sNew, $sFile)
    ; 1) Existiert die Datei noch nicht...
    If Not FileExists($sFile) Then
        Local $hFile = FileOpen($sFile, 2) ; Modus=2 -> Schreiben/Überschreiben
        If $hFile <> -1 Then
            FileWriteLine($hFile, $sOld & "|" & $sNew)
            FileClose($hFile)
        EndIf
        Return
    EndIf

    ; 2) Datei existiert bereits..
    Local $aLines[1] 
	_FileReadToArray($sFile,$aLines)
    If @error Then Return
      
    Local $bFound = False
    ; 3) Jede Zeile durchgehen und schauen, ob die alten Farben vorhanden sind
    For $i = 1 To $aLines[0]
        Local $aSplit = StringSplit($aLines[$i], "|")
        If $aSplit[0] = 2 Then
            If $aSplit[1] = $sOld Then
                $aLines[$i] = $sOld & "|" & $sNew
                $bFound = True
                ExitLoop
            EndIf
        EndIf
    Next
	
;~ 	_ArrayDisplay($aLines)
    ; 4) Wurden die alten Farben nicht gefunden, wird ein neuer Eintrag angehängt.
    If Not $bFound Then
		_ArrayAdd($aLines,$sOld & "|" & $sNew, 0, "&")
		$aLines[0] =  UBound($aLines) -1
    EndIf
;~ 	_ArrayDisplay($aLines)
	
    ; 5) Das aktualisierte Array zurück in die Datei schreiben
    Local $hFile = FileOpen($sFile, 2) 
    If $hFile <> -1 Then
        For $i = 1 To $aLines[0]
            FileWriteLine($hFile, $aLines[$i])
        Next
        FileClose($hFile)
    EndIf
EndFunc

Func _currentversionlink()
	ShellExecute("https://github.com/schiffer-soft/ColorTransfer/releases")
EndFunc

Func _button_help()
	ShellExecute("https://github.com/schiffer-soft/ColorTransfer")
EndFunc

Func _Xmaxklaxx()
	ShellExecute("https://github.com/MaxKlaxxMiner")
EndFunc 
	
Func _Xschiffersoft()
	ShellExecute("https://x.com/schiffer_soft")
EndFunc 

Func _paypalme()
	ShellExecute("https://paypal.me/schiffersoft")
EndFunc 
	
Func _Base64StringEx($externImage)
	Local $bString = _WinAPI_Base64Decode($externImage)
	If @error Then Return SetError(1, 0, 0)
	$bString = Binary($bString)
	Return $bString
EndFunc   ;==>_Base64String

Func _WinAPI_Base64Decode($sB64String)
	Local $aCrypt = DllCall("Crypt32.dll", "bool", "CryptStringToBinaryA", "str", $sB64String, "dword", 0, "dword", 1, "ptr", 0, "dword*", 0, "ptr", 0, "ptr", 0)
	If @error Or Not $aCrypt[0] Then Return SetError(1, 0, "")
	Local $bBuffer = DllStructCreate("byte[" & $aCrypt[5] & "]")
	$aCrypt = DllCall("Crypt32.dll", "bool", "CryptStringToBinaryA", "str", $sB64String, "dword", 0, "dword", 1, "struct*", $bBuffer, "dword*", $aCrypt[5], "ptr", 0, "ptr", 0)
	If @error Or Not $aCrypt[0] Then Return SetError(2, 0, "")
	Return DllStructGetData($bBuffer, 1)
EndFunc   ;==>_WinAPI_Base64Decode