#NoEnv ; Recommended for performance and compatibility with future AutoHotkey releases.
; #Warn  ; Enable warnings to assist with detecting common errors.
SendMode Input ; Recommended for new scripts due to its superior speed and reliability.
SetWorkingDir %A_ScriptDir% ; Ensures a consistent starting directory.

; 打开一个网址
;=========================================================
!s::Run, https://www.baidu.com/

; 最钟爱代码之音量随心所欲
;=========================================================
~lbutton & enter:: ;鼠标放在任务栏，滚动滚轮实现音量的加减
exitapp 
~WheelUp:: 
    if (existclass("ahk_class Shell_TrayWnd")=1) 
        Send,{Volume_Up} 
Return 
~WheelDown:: 
    if (existclass("ahk_class Shell_TrayWnd")=1) 
        Send,{Volume_Down} 
Return 
~MButton:: 
    if (existclass("ahk_class Shell_TrayWnd")=1) 
        Send,{Volume_Mute} 
Return 

Existclass(class) 
{ 
    MouseGetPos,,,win 
    WinGet,winid,id,%class% 
    if win = %winid% 
        Return,1 
    Else 
        Return,0 
}

; 最常用功能之打开 - Edge
;=========================================================
!q:: 
    IfWinNotExist ahk_class Chrome_WidgetWin_1
    {
        run "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Microsoft Edge.lnk",,max
        WinActivate
    }
    Else IfWinNotActive ahk_class Chrome_WidgetWin_1
    {
        WinActivate
        WinMaximize,A
    }
    Else
    {
        WinMinimize
    }
Return

; 最常用功能之打开 - MobaXterm
;=========================================================
!q:: 
    IfWinNotExist ahk_class TMobaXtermForm
    {
        run "C:\Users\luoxian\Documents\我的应用\MobaXterm.lnk",,max
        WinActivate
    }
    Else IfWinNotActive ahk_class TMobaXtermForm
    {
        WinActivate
        WinMaximize,A
    }
    Else
    {
        WinMinimize
    }
Return