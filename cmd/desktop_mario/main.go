// Desktop Mario – A stress-relief mini-game that lives on your desktop.
//
//	Ctrl+Alt+M  – show/hide the game (Windows only)
//	Arrow keys  – move
//	Space       – jump (hold longer = higher)
//	Shift       – run faster
//	F / Z       – throw fireball (Big Mario only)
//	ESC         – hide
//
// Build: go build -o desktop_mario ./cmd/desktop_mario
// No installer required; produces a single executable.
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// setupLogging redirects log output to a file next to the executable so that
// messages are visible even when the OS console window is hidden (Windows).
func setupLogging() {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	logPath := filepath.Join(filepath.Dir(ex), "desktop_mario.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("=== Desktop Mario started ===")
}

func main() {
	setupLogging()

	// Determine screen size before opening the window.
	sw, sh := ebiten.ScreenSizeInFullscreen()
	W := float64(sw)
	H := float64(sh)

	// Configure the Ebiten window: full-screen borderless, always-on-top.
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Desktop Mario")
	// 30 TPS matches the original Python ~30 fps loop.
	ebiten.SetTPS(30)

	// Channel for the global hotkey → toggle visibility.
	toggleCh := make(chan struct{}, 1)
	toggleFn := func() {
		select {
		case toggleCh <- struct{}{}:
		default:
		}
	}

	// Start platform-specific setup (Win32 hotkey registration, etc.).
	setupPlatform(toggleFn)

	g := newGame(W, H, toggleCh)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
