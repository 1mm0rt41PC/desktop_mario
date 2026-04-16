package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Block / sprite size in screen pixels (16 logical px × pixelScale).
const blk = 16 * pixelScale // 48

// ── Game struct ───────────────────────────────────────────────────────────────

// Game implements ebiten.Game.
type Game struct {
	W, H float64 // screen dimensions

	tick  int
	score int

	// Mario state
	mwx, my      float64
	mvx, mvy     float64
	onGround     bool
	facingRight  bool
	isBig        bool
	jumping      bool
	dead         bool
	deadTimer    int
	lastSafeWX   float64
	invincible   int
	shrinkTimer  int
	stompGrace   int
	fireCooldown int

	// Camera
	cam float64

	// World objects
	bricks    []brick
	qblocks   []qblock
	coins     []coin
	enemies   []enemy
	pipes     []pipe
	gaps      []gap
	mushrooms []mushroom
	fireballs []fireball
	popups    []scorePopup
	clouds    []cloud
	genX      float64

	// Ground tile positions (looping)
	groundTileXs []float64

	// Input state
	keys [ebiten.KeyMax]bool

	// Sprites (built once)
	spr spriteImages

	// Platform setup: transparency applied on first frame
	platformReady bool

	// Toggle channel (from platform hotkey goroutine)
	toggleCh <-chan struct{}

	visible bool
}

func newGame(w, h float64, toggleCh <-chan struct{}, gameNow bool) *Game {
	g := &Game{
		W:           w,
		H:           h,
		facingRight: true,
		visible:     gameNow,
		toggleCh:    toggleCh,
	}
	g.spr = buildSprites()

	g.mwx = 200
	g.my = g.groundY()
	g.lastSafeWX = g.mwx
	g.invincible = 60

	// Ground tiles: enough to fill width + 3 extra
	n := int(w/blk) + 4
	g.groundTileXs = make([]float64, n)
	for i := range g.groundTileXs {
		g.groundTileXs[i] = float64(i * blk)
	}

	// Clouds
	for i := 0; i < 5; i++ {
		cw := float64(rand.Intn(51) + 70)
		ch := float64(rand.Intn(15) + 26)
		g.clouds = append(g.clouds, cloud{
			x:     float64(rand.Intn(int(w))),
			y:     float64(rand.Intn(int(h/3)+1) + 30),
			speed: rand.Float64()*0.2 + 0.2,
			w:     cw,
			h:     ch,
		})
	}

	// Generate first chunk
	g._generate(600, int(w)+600)
	return g
}

func (g *Game) groundY() float64 {
	return g.H - 48 - blk
}

// ── ebiten.Game interface ─────────────────────────────────────────────────────

func (g *Game) Layout(ow, oh int) (int, int) { return ow, oh }

func (g *Game) Update() error {
	// Apply Win32 transparency on the very first frame (window now exists).
	if !g.platformReady {
		applyTransparency()
		if !g.visible {
			platformHideWindow()
		}
		g.platformReady = true
	}

	// Handle toggle from hotkey goroutine.
	select {
	case <-g.toggleCh:
		g.visible = !g.visible
		if g.visible {
			platformShowWindow()
		} else {
			platformHideWindow()
		}
	default:
	}

	if !g.visible {
		return nil
	}

	g._update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.visible {
		return
	}
	// Black background; pure-black pixels become transparent on Windows
	// via the WS_EX_LAYERED / LWA_COLORKEY mechanism.
	screen.Fill(color.RGBA{0, 0, 0, 255})
	g._draw(screen)
}
