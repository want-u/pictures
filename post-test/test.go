package main

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

func main() {
	err := openAmlogicBurningTool()
	if err != nil {
		fmt.Println("open tool error : ", err)
		return
	}

	// 延时等待软件启动
	time.Sleep(time.Second)

	processName := "OpenVPNConnect.exe"
	fpid, err := findProcess(processName)
	fmt.Printf("fpid: %v\n", fpid)
	if err != nil {
		fmt.Println("find process error : ", err)
		return
	}

	activeWin := robotgo.GetActive()
	defer robotgo.SetActive(activeWin)
	if err := robotgo.ActivePID(fpid); err != nil {
		fmt.Println("active pid error : ", err)
		return
	}
	fmt.Println("robotgo.ActivePID( ok")

	hwnd := robotgo.GetHWND()
	fmt.Printf("hwnd: %v\n", hwnd)
	if !win.EnumChildWindows(hwnd, callback(EnumMainTVWindowCn), 0) {
		fmt.Println("start.. ")
		fmt.Printf("startButton: %v\n", startButton)
		win.SendMessage(startButton, win.WM_LBUTTONDOWN, 0, 0)
		win.SendMessage(startButton, win.WM_LBUTTONUP, 0, 0)
	}

	// 关闭这个窗口
	fmt.Println(robotgo.Kill(fpid))
}

func windowText(hchildWnd win.HWND) string {
	textLength := win.SendMessage(hchildWnd, win.WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, textLength+1)
	win.SendMessage(hchildWnd, win.WM_GETTEXT, uintptr(textLength+1), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

var startButton win.HWND

func EnumMainTVWindowCn(hWnd win.HWND, lParam uintptr) uintptr {
	ret := windowText(hWnd)
	fmt.Printf("ret: %v\n", ret)
	if ret == "Chrome Legacy Window" {
		fmt.Printf("ret: %v\n", ret)
		startButton = hWnd
		return 0
	}

	return 1
}

func callback(fn interface{}) uintptr {
	return syscall.NewCallback(fn)
}

func openAmlogicBurningTool() error {
	command := "C:\\Program Files\\OpenVPN Connect\\OpenVPNConnect.exe"
	cmd := exec.Command("cmd.exe", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func findProcess(name string) (int32, error) {
	processes, err := robotgo.Process()
	if err != nil {
		return 0, err
	}
	for _, process := range processes {
		if process.Name == name {
			return process.Pid, nil
		}
	}
	return 0, errors.New("process not found")
}
