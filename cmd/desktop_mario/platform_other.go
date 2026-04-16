//go:build !windows

package main

import "github.com/hajimehoshi/ebiten/v2"

// applyTransparency is a no-op on non-Windows platforms; always succeeds.
func applyTransparency() bool { return true }

// getWorkArea returns 0,0 on non-Windows; caller falls back to full screen size.
func getWorkArea() (int, int) { return 0, 0 }

// platformShowWindow restores the window from minimized state.
func platformShowWindow() { ebiten.RestoreWindow() }

// platformHideWindow minimizes the window (no system-level hide on non-Windows).
func platformHideWindow() { ebiten.MinimizeWindow() }

// setupPlatform is a no-op on non-Windows platforms (no global hotkey support).
func setupPlatform(_ func()) {}
