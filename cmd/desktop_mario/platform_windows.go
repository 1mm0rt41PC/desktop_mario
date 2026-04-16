//go:build windows

package main

import (
	"log"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procFindWindowW           = user32.NewProc("FindWindowW")
	procShowWindow            = user32.NewProc("ShowWindow")
	procRegisterHotKey        = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey      = user32.NewProc("UnregisterHotKey")
	procGetMessageW           = user32.NewProc("GetMessageW")
	procSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

const (
	swShow         = uintptr(5)
	swHide         = uintptr(0)
	wmHotkey       = uint32(0x0312)
	modAlt         = uintptr(0x0001)
	modControl     = uintptr(0x0002)
	vkM            = uintptr(0x4D)
	hotkeyID       = uintptr(1)
	spiGetWorkArea = uintptr(0x0030)
)

// windowTitle must match the title set by ebiten.SetWindowTitle in main.go.
const windowTitle = "Desktop Mario"

// winRECT mirrors the Win32 RECT structure.
type winRECT struct{ Left, Top, Right, Bottom int32 }

// getWorkArea returns the usable screen dimensions excluding the taskbar.
func getWorkArea() (int, int) {
	var r winRECT
	procSystemParametersInfoW.Call(spiGetWorkArea, 0, uintptr(unsafe.Pointer(&r)), 0)
	w := int(r.Right - r.Left)
	h := int(r.Bottom - r.Top)
	if w <= 0 || h <= 0 {
		return 0, 0 // caller falls back to full screen size
	}
	return w, h
}

// getHWND locates the game window using FindWindowW and the known title.
// Returns 0 if not found yet.
func getHWND() uintptr {
	p, err := syscall.UTF16PtrFromString(windowTitle)
	if err != nil {
		return 0
	}
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(p)))
	return hwnd
}

// applyTransparency is a no-op here — transparency is handled by ebiten's
// RunGameWithOptions{ScreenTransparent: true}. Returns true once the window
// exists (HWND is valid), which signals platformReady.
func applyTransparency() bool {
	return getHWND() != 0
}

// platformShowWindow shows the game window.
func platformShowWindow() {
	if hwnd := getHWND(); hwnd != 0 {
		procShowWindow.Call(hwnd, swShow)
	}
}

// platformHideWindow hides the game window without destroying it.
func platformHideWindow() {
	if hwnd := getHWND(); hwnd != 0 {
		procShowWindow.Call(hwnd, swHide)
	}
}

// winMSG is the Windows MSG structure.
type winMSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      [2]int32
}

// setupPlatform registers the global Ctrl+Alt+M hotkey and fires toggleFn when pressed.
// LockOSThread is required so that RegisterHotKey and GetMessageW always run
// on the same OS thread.
func setupPlatform(toggleFn func()) {
	go func() {
		runtime.LockOSThread()
		procUnregisterHotKey.Call(0, hotkeyID)
		ret, _, err := procRegisterHotKey.Call(0, hotkeyID, modAlt|modControl, vkM)
		if ret == 0 {
			log.Printf("RegisterHotKey Ctrl+Alt+M failed: %v (another instance may own it)", err)
			return
		}
		log.Println("RegisterHotKey Ctrl+Alt+M registered successfully")
		var msg winMSG
		for {
			r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if r == 0 || r == ^uintptr(0) { // WM_QUIT or error
				break
			}
			if msg.Message == wmHotkey {
				toggleFn()
			}
		}
		procUnregisterHotKey.Call(0, hotkeyID)
	}()
}
