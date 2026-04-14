//go:build windows

package main

import (
	"syscall"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procFindWindowW             = user32.NewProc("FindWindowW")
	procGetWindowLongW          = user32.NewProc("GetWindowLongW")
	procSetWindowLongW          = user32.NewProc("SetWindowLongW")
	procSetLayeredWindowAttribs = user32.NewProc("SetLayeredWindowAttributes")
	procShowWindow              = user32.NewProc("ShowWindow")
	procRegisterHotKey          = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey        = user32.NewProc("UnregisterHotKey")
	procGetMessageW             = user32.NewProc("GetMessageW")
)

// gwlExStyleUint is GWL_EXSTYLE (-20) represented as a uint32 value
// suitable for passing to GetWindowLongW/SetWindowLongW via syscall.Call.
// -20 as int32 = 0xFFFFFFEC; Windows reads it as a signed index.
const gwlExStyleUint = uintptr(0xFFFFFFEC)

const (
	wsExLayered    = uintptr(0x00080000)
	lwaColorKey    = uintptr(0x00000001)
	swShow         = uintptr(5)
	swHide         = uintptr(0)
	wmHotkey       = uint32(0x0312)
	modAlt         = uintptr(0x0001)
	modControl     = uintptr(0x0002)
	vkM            = uintptr(0x4D)
	hotkeyID       = uintptr(1)
	transparentRGB = uintptr(0x00000000) // pure black (COLORREF 0x00BBGGRR)
)

// windowTitle must match the title set by ebiten.SetWindowTitle in main.go.
const windowTitle = "Desktop Mario"

// getHWND locates the game window using FindWindowW and the known title.
func getHWND() uintptr {
	p, err := syscall.UTF16PtrFromString(windowTitle)
	if err != nil {
		return 0
	}
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(p)))
	return hwnd
}

// applyTransparency makes the window layered with pure black as the color key.
// Must be called after the window is created (inside Update/Draw, not before RunGame).
func applyTransparency() {
	hwnd := getHWND()
	if hwnd == 0 {
		return
	}
	exStyle, _, _ := procGetWindowLongW.Call(hwnd, gwlExStyleUint)
	procSetWindowLongW.Call(hwnd, gwlExStyleUint, exStyle|wsExLayered)
	// crKey=0 (black), bAlpha=255, dwFlags=LWA_COLORKEY
	procSetLayeredWindowAttribs.Call(hwnd, transparentRGB, 255, lwaColorKey)
}

// platformShowWindow shows the game window.
func platformShowWindow() {
	hwnd := getHWND()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, swShow)
		return
	}
	ebiten.RestoreWindow()
}

// platformHideWindow hides the game window without destroying it.
func platformHideWindow() {
	hwnd := getHWND()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, swHide)
		return
	}
	ebiten.MinimizeWindow()
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
func setupPlatform(toggleFn func()) {
	go func() {
		procUnregisterHotKey.Call(0, hotkeyID)
		ret, _, _ := procRegisterHotKey.Call(0, hotkeyID, modAlt|modControl, vkM)
		if ret == 0 {
			return // could not register (another instance may own it) — silently skip
		}
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
