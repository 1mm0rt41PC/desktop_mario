// Desktop Mario – A stress-relief mini-game that lives on your desktop.
//
//	Ctrl+Alt+M  – show/hide the game (Windows only)
//	Arrow keys  – move
//	Space       – jump (hold longer = higher)
//	Shift       – run faster
//	F / Z       – throw fireball (Big Mario only)
//	ESC         – hide
//
// Build: go build -o desktop_mario .
// No installer required; produces a single executable.
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
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

	// Bind ESC key: hide (not quit) — same as the Python version.
	// We handle this inside Game.Update via the keys map.

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
