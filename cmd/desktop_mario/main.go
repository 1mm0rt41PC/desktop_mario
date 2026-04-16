// Desktop Mario – A stress-relief mini-game that lives on your desktop.
//
//	Ctrl+Alt+M  – show/hide the game (Windows only)
//	Arrow keys  – move
//	Space       – jump (hold longer = higher)
//	Shift       – run faster
//	F / Z       – throw fireball (Big Mario only)
//	ESC         – hide
//	--game-now  – start the game immediately visible (skip hotkey trigger)
//
// Build: go build -ldflags="-H=windowsgui" -o desktop_mario.exe ./cmd/desktop_mario
package main

import (
	"flag"
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

	gameNow := flag.Bool("game-now", false, "start the game immediately visible without requiring the Ctrl+Alt+M hotkey")
	flag.Parse()

	// Full screen dimensions for window and game coordinate system.
	sw, sh := ebiten.ScreenSizeInFullscreen()
	W := float64(sw)
	H := float64(sh)

	// Work area excludes the taskbar — compute taskbar height so the ground
	// can be positioned above it. Fall back to 0 if unavailable.
	_, wh := getWorkArea()
	taskbarH := 0
	if wh > 0 && wh < sh {
		taskbarH = sh - wh
	}
	log.Printf("screen=%dx%d workH=%d taskbarH=%d", sw, sh, wh, taskbarH)

	// Configure the Ebiten window: full-screen borderless, always-on-top.
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Desktop Mario")
	ebiten.SetWindowMousePassthrough(true)
	// 60 TPS for smooth gameplay.
	ebiten.SetTPS(60)

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

	g := newGame(W, H, taskbarH, toggleCh, *gameNow)

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
		SkipTaskbar:       true,
	}
	if err := ebiten.RunGameWithOptions(g, opts); err != nil {
		log.Fatal(err)
	}
}
