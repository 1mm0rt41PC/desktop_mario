package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Block / sprite size in screen pixels (16 logical px × pixelScale).
const blk = 16 * pixelScale // 48

// ── Data types ────────────────────────────────────────────────────────────────

type brick struct{ wx, y float64 }
type qblock struct {
	wx, y   float64
	hit     bool
	reward  string // "coin" | "mushroom"
}
type coin struct {
	wx, wy float64
	got    bool
	ft     int
}
type pipe struct{ wx, y, w, h, lipH float64 }
type gap struct{ startX, endX float64 }
type mushroom struct {
	wx, wy float64
	vx, vy float64
	active bool
}
type fireball struct {
	wx, wy float64
	vx, vy float64
}
type scorePopup struct {
	x, y float64
	text string
	life int
}
type cloud struct {
	x, y, speed, w, h float64
}

type enemyKind int

const (
	kindGoomba enemyKind = iota
	kindKoopa
	kindRedKoopa
	kindBobomb
)

type enemyState int

const (
	stateWalk enemyState = iota
	stateFlat
	stateShellStill
	stateShell
	stateFuse
)

type enemy struct {
	kind  enemyKind
	state enemyState
	wx, wy float64
	vx    float64
	timer int
}

// ── Sprite image cache ────────────────────────────────────────────────────────

type spriteImages struct {
	// Mario small: [stand_r, run1_r, run2_r, jump_r, stand_l, run1_l, run2_l, jump_l]
	marioSmall [8]*ebiten.Image
	// Mario big (same order)
	marioBig [8]*ebiten.Image
	// Enemies
	goomba   [3]*ebiten.Image // walk1, walk2, flat
	koopa    [3]*ebiten.Image // walk1_l, walk2_l, shell
	redKoopa [5]*ebiten.Image // walk1_l, walk2_l, shell, walk1_r, walk2_r
	bobomb   [3]*ebiten.Image // walk1, walk2, explode
	// Blocks / items
	groundBlock  *ebiten.Image
	brickBlock   *ebiten.Image
	qBlock       *ebiten.Image
	qBlockUsed   *ebiten.Image
	coin1        *ebiten.Image
	coin2        *ebiten.Image
	mushroomImg  *ebiten.Image
	fireball1Img *ebiten.Image
	fireball2Img *ebiten.Image
}

func buildSprites() spriteImages {
	var s spriteImages
	px := pixelScale

	// Small Mario (right-facing)
	s.marioSmall[0] = frameToImage(marioStand, px)
	s.marioSmall[1] = frameToImage(marioRun1, px)
	s.marioSmall[2] = frameToImage(marioRun2, px)
	s.marioSmall[3] = frameToImage(marioJump, px)
	// Small Mario (left-facing)
	s.marioSmall[4] = frameToImage(flipFrame(marioStand), px)
	s.marioSmall[5] = frameToImage(flipFrame(marioRun1), px)
	s.marioSmall[6] = frameToImage(flipFrame(marioRun2), px)
	s.marioSmall[7] = frameToImage(flipFrame(marioJump), px)
	// Big Mario
	s.marioBig[0] = frameToImage(bigMarioStand, px)
	s.marioBig[1] = frameToImage(bigMarioRun1, px)
	s.marioBig[2] = frameToImage(bigMarioRun2, px)
	s.marioBig[3] = frameToImage(bigMarioJump, px)
	s.marioBig[4] = frameToImage(flipFrame(bigMarioStand), px)
	s.marioBig[5] = frameToImage(flipFrame(bigMarioRun1), px)
	s.marioBig[6] = frameToImage(flipFrame(bigMarioRun2), px)
	s.marioBig[7] = frameToImage(flipFrame(bigMarioJump), px)
	// Enemies
	s.goomba[0] = frameToImage(goomba1, px)
	s.goomba[1] = frameToImage(goomba2, px)
	s.goomba[2] = frameToImage(goombaFlat, px)
	s.koopa[0] = frameToImage(koopaL1, px)
	s.koopa[1] = frameToImage(koopaL2, px)
	s.koopa[2] = frameToImage(shellSprite, px)
	s.redKoopa[0] = frameToImage(koopaL1, px)
	s.redKoopa[1] = frameToImage(koopaL2, px)
	s.redKoopa[2] = frameToImage(shellSprite, px)
	s.redKoopa[3] = frameToImage(flipFrame(koopaL1), px)
	s.redKoopa[4] = frameToImage(flipFrame(koopaL2), px)
	s.bobomb[0] = frameToImage(bobomb1, px)
	s.bobomb[1] = frameToImage(bobomb2, px)
	s.bobomb[2] = frameToImage(bobombExplode, px)
	// Blocks / items
	s.groundBlock = frameToImage(groundBlockFrame, px)
	s.brickBlock = frameToImage(brickFrame, px)
	s.qBlock = frameToImage(qblockFrame, px)
	s.qBlockUsed = frameToImage(qblockUsedFrame, px)
	s.coin1 = frameToImage(coin1Frame, px)
	s.coin2 = frameToImage(coin2Frame, px)
	s.mushroomImg = frameToImage(mushroomFrame, px)
	s.fireball1Img = frameToImage(fireball1, px)
	s.fireball2Img = frameToImage(fireball2, px)
	return s
}

// ── Game ──────────────────────────────────────────────────────────────────────

// Game implements ebiten.Game.
type Game struct {
	W, H float64 // screen dimensions

	tick  int
	score int

	// Mario state
	mwx, my     float64
	mvx, mvy    float64
	onGround    bool
	facingRight bool
	isBig       bool
	jumping     bool
	dead        bool
	deadTimer   int
	lastSafeWX  float64
	invincible  int
	shrinkTimer int
	stompGrace  int
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

func newGame(w, h float64, toggleCh <-chan struct{}) *Game {
	g := &Game{
		W:           w,
		H:           h,
		facingRight: true,
		visible:     true,
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

// ── Game logic (update) ───────────────────────────────────────────────────────

func (g *Game) _update() {
	// Populate key state from Ebiten input.
	for k := ebiten.Key(0); k < ebiten.KeyMax; k++ {
		g.keys[k] = ebiten.IsKeyPressed(k)
	}

	// ESC → hide window.
	if g.keys[ebiten.KeyEscape] {
		g.visible = false
		platformHideWindow()
		return
	}

	if g.dead {
		g.deadTimer--
		if g.deadTimer > 25 {
			g.my -= 8
		} else {
			g.my += 10
		}
		if g.deadTimer <= 0 {
			g._respawn()
		}
		return
	}

	g.tick++
	MH := float64(blk) // Mario hitbox height
	if g.isBig {
		MH = 32 * pixelScale
	}
	MT := g.my
	if g.isBig {
		MT = g.my - (MH - blk)
	}

	// Countdowns
	if g.stompGrace > 0 {
		g.stompGrace--
	}
	if g.invincible > 0 {
		g.invincible--
	}
	if g.shrinkTimer > 0 {
		g.shrinkTimer--
	}

	// Clouds
	for i := range g.clouds {
		g.clouds[i].x -= g.clouds[i].speed
		if g.clouds[i].x+g.clouds[i].w < 0 {
			g.clouds[i].x = g.W + float64(rand.Intn(200)+50)
			g.clouds[i].y = float64(rand.Intn(int(g.H/3)+1) + 30)
		}
	}

	// Input → velocity
	running := g.keys[ebiten.KeyShiftLeft] || g.keys[ebiten.KeyShiftRight]
	accel := 0.8
	maxSpeed := 5.0
	friction := 0.6
	if running {
		accel = 1.2
		maxSpeed = 8.0
	}

	if g.keys[ebiten.KeyArrowRight] {
		g.mvx = math.Min(g.mvx+accel, maxSpeed)
		g.facingRight = true
	} else if g.keys[ebiten.KeyArrowLeft] {
		g.mvx = math.Max(g.mvx-accel, -maxSpeed)
		g.facingRight = false
	} else {
		if math.Abs(g.mvx) < friction {
			g.mvx = 0
		} else if g.mvx > 0 {
			g.mvx -= friction
		} else {
			g.mvx += friction
		}
	}

	// Variable-height jump
	if g.keys[ebiten.KeySpace] && g.onGround {
		g.mvy = -16
		g.onGround = false
		g.jumping = true
	}
	if g.jumping && g.keys[ebiten.KeySpace] && g.mvy < 0 {
		g.mvy += 0.55
	} else {
		g.mvy += 1.2
		if g.mvy >= 0 {
			g.jumping = false
		}
	}
	if g.mvy > 18 {
		g.mvy = 18
	}

	// Fireball (Big Mario, press F or Z)
	if g.fireCooldown > 0 {
		g.fireCooldown--
	}
	if g.isBig && g.fireCooldown <= 0 && (g.keys[ebiten.KeyF] || g.keys[ebiten.KeyZ]) {
		dir := 1.0
		if !g.facingRight {
			dir = -1
		}
		fbwx := g.mwx + blk*dir
		if !g.facingRight {
			fbwx = g.mwx - 12
		}
		g.fireballs = append(g.fireballs, fireball{
			wx: fbwx, wy: g.my + float64(blk)/3,
			vx: 7.0 * dir, vy: -4.0,
		})
		g.fireCooldown = 12
	}

	// Move X → check solids
	g.mwx += g.mvx
	solids := g._allSolids()
	for _, s := range solids {
		if overlap(g.mwx, MT, float64(blk), MH, s[0], s[1], s[2], s[3]) {
			if g.mvx > 0 {
				g.mwx = s[0] - float64(blk)
			} else if g.mvx < 0 {
				g.mwx = s[0] + s[2]
			}
			g.mvx = 0
			break
		}
	}

	// Move Y → check solids
	g.my += g.mvy
	MT = g.my
	if g.isBig {
		MT = g.my - (MH - blk)
	}
	g.onGround = false
	for _, s := range solids {
		if overlap(g.mwx+2, MT, float64(blk)-4, MH, s[0], s[1], s[2], s[3]) {
			if g.mvy > 0 {
				g.my = s[1] - float64(blk)
				g.mvy = 0
				g.onGround = true
			} else if g.mvy < 0 {
				if g.isBig {
					g.my = s[1] + s[3] + (MH - blk)
				} else {
					g.my = s[1] + s[3]
				}
				g.mvy = 1
				// Bump ?-blocks
				for i := range g.qblocks {
					q := &g.qblocks[i]
					if !q.hit && q.wx == s[0] && q.y == s[1] {
						q.hit = true
						scrX := q.wx - g.cam
						if q.reward == "mushroom" && !g.isBig {
							g.mushrooms = append(g.mushrooms, mushroom{
								wx: q.wx, wy: q.y - float64(blk),
								vx: 2.0, active: true,
							})
						} else {
							g._addScore(100, scrX+float64(blk)/2, q.y-20)
						}
					}
				}
			}
			break
		}
	}

	// Ground floor (with gap check)
	if g.my >= g.groundY() {
		if !g._inGap(g.mwx, float64(blk)) {
			g.my = g.groundY()
			g.mvy = 0
			g.onGround = true
		}
	}

	if g.onGround && !g._inGap(g.mwx, float64(blk)) {
		g.lastSafeWX = g.mwx
	}

	// Fell off
	if g.my > g.H+100 {
		g._die()
		return
	}

	// Can't go left past camera
	if g.mwx < g.cam {
		g.mwx = g.cam
		g.mvx = 0
	}

	// Camera
	g.cam = g.mwx - g.W*0.3

	// Generate level ahead
	edge := g.cam + g.W + 400
	if edge > g.genX {
		g._generate(int(g.genX), int(edge))
	}

	// ── Coins ────────────────────────────────────────────────────────────────
	aliveCoins := g.coins[:0]
	for i := range g.coins {
		c := &g.coins[i]
		sx := c.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		if !c.got {
			if overlap(g.mwx, MT, float64(blk), MH, c.wx, c.wy, 8*pixelScale, float64(blk)) {
				c.got = true
				c.ft = 15
				g._addScore(100, sx, c.wy)
			}
		} else {
			c.ft--
			c.wy -= 5
			if c.ft <= 0 {
				continue
			}
		}
		aliveCoins = append(aliveCoins, *c)
	}
	g.coins = aliveCoins

	// ── Mushrooms ─────────────────────────────────────────────────────────────
	aliveMushrooms := g.mushrooms[:0]
	for i := range g.mushrooms {
		m := &g.mushrooms[i]
		sx := m.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		if m.active {
			m.vy += 1.0
			if m.vy > 10 {
				m.vy = 10
			}
			m.wx += m.vx
			m.wy += m.vy
			if m.wy >= g.groundY() {
				m.wy = g.groundY()
				m.vy = 0
			}
			for _, s := range solids {
				if overlap(m.wx, m.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
					if m.vx > 0 {
						m.wx = s[0] - float64(blk)
					} else {
						m.wx = s[0] + s[2]
					}
					m.vx *= -1
					break
				}
			}
			if overlap(g.mwx, MT, float64(blk), MH, m.wx, m.wy, float64(blk), float64(blk)) {
				m.active = false
				if !g.isBig {
					g.isBig = true
					g._addScore(1000, sx, m.wy-20)
				}
				continue
			}
		}
		aliveMushrooms = append(aliveMushrooms, *m)
	}
	g.mushrooms = aliveMushrooms

	// ── Fireballs ─────────────────────────────────────────────────────────────
	fbSz := float64(8 * pixelScale)
	aliveFireballs := g.fireballs[:0]
	for i := range g.fireballs {
		fb := &g.fireballs[i]
		fb.wy += fb.vy
		fb.wx += fb.vx
		fb.vy += 1.0
		if fb.wy >= g.groundY() {
			fb.wy = g.groundY()
			fb.vy = -8.0
		}
		hitSolid := false
		for _, s := range solids {
			if overlap(fb.wx, fb.wy, fbSz, fbSz, s[0], s[1], s[2], s[3]) {
				if fb.vy <= 0 || fb.wy+fbSz > s[1]+fbSz/2 {
					hitSolid = true
				} else {
					fb.wy = s[1] - fbSz
					fb.vy = -8.0
				}
				break
			}
		}
		if hitSolid {
			continue
		}
		fbsx := fb.wx - g.cam
		if fbsx < -float64(blk)*2 || fbsx > g.W+float64(blk)*2 || fb.wy > g.H+50 {
			continue
		}
		killed := false
		for j := range g.enemies {
			e := &g.enemies[j]
			if e.state == stateFlat {
				continue
			}
			if overlap(fb.wx, fb.wy, fbSz, fbSz, e.wx, e.wy, float64(blk), float64(blk)) {
				e.state = stateFlat
				e.timer = 18
				g._addScore(200, e.wx-g.cam, e.wy-20)
				killed = true
				break
			}
		}
		if killed {
			continue
		}
		aliveFireballs = append(aliveFireballs, *fb)
	}
	g.fireballs = aliveFireballs

	// ── Enemies ───────────────────────────────────────────────────────────────
	marioL, marioT := g.mwx, MT
	marioB := g.my + float64(blk)
	aliveEnemies := g.enemies[:0]
	for i := range g.enemies {
		e := &g.enemies[i]
		sx := e.wx - g.cam
		if sx < -float64(blk)*5 {
			continue
		}
		if sx > g.W+float64(blk)*3 {
			if e.state == stateWalk {
				e.wx += e.vx
			}
			aliveEnemies = append(aliveEnemies, *e)
			continue
		}

		switch e.state {
		case stateWalk:
			e.wx += e.vx
			for _, s := range solids {
				if overlap(e.wx, e.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
					if e.vx > 0 {
						e.wx = s[0] - float64(blk)
					} else {
						e.wx = s[0] + s[2]
					}
					e.vx *= -1
					break
				}
			}
			if e.kind == kindRedKoopa {
				aheadX := e.wx + float64(blk)
				if e.vx < 0 {
					aheadX = e.wx - 4
				}
				onSolid := false
				if !g._inGap(aheadX, 4) {
					if e.wy >= g.groundY() {
						onSolid = true
					} else {
						for _, s := range solids {
							if overlap(aheadX, e.wy+float64(blk), 4, 4, s[0], s[1], s[2], s[3]) {
								onSolid = true
								break
							}
						}
					}
				}
				if !onSolid {
					e.vx *= -1
				}
			}
			if g.invincible <= 0 && g.stompGrace <= 0 && g.shrinkTimer <= 0 &&
				overlap(marioL+4, marioT, float64(blk)-8, MH, e.wx, e.wy, float64(blk), float64(blk)) {
				if g.mvy > 0 && marioB < e.wy+float64(blk)*0.6 {
					switch e.kind {
					case kindGoomba:
						e.state = stateFlat
						e.timer = 18
						e.wy += float64(blk) / 2
					case kindBobomb:
						e.state = stateFuse
						e.timer = 90
						e.vx = 0
					default:
						e.state = stateShellStill
						e.vx = 0
						e.timer = 300
					}
					g.mvy = -10
					g.stompGrace = 25
					g._addScore(200, sx, e.wy-20)
				} else {
					g._takeHit()
					aliveEnemies = append(aliveEnemies, *e)
					for j := i + 1; j < len(g.enemies); j++ {
						aliveEnemies = append(aliveEnemies, g.enemies[j])
					}
					g.enemies = aliveEnemies
					return
				}
			}

		case stateShellStill:
			e.timer--
			if e.timer <= 0 {
				e.state = stateWalk
				e.vx = -1.5
				aliveEnemies = append(aliveEnemies, *e)
				continue
			}
			if g.invincible <= 0 && g.stompGrace <= 0 &&
				overlap(marioL+2, marioT+4, float64(blk)-4, MH-8, e.wx, e.wy, float64(blk), float64(blk)) {
				dir := 10.0
				if g.mwx+float64(blk)/2 > e.wx+float64(blk)/2 {
					dir = -10
				}
				e.state = stateShell
				e.vx = dir
				if g.mwx+float64(blk)/2 < e.wx+float64(blk)/2 {
					g.mwx = e.wx - float64(blk) - 2
				} else {
					g.mwx = e.wx + float64(blk) + 2
				}
				g.stompGrace = 15
				g._addScore(100, sx, e.wy-20)
			}

		case stateFuse:
			e.timer--
			if e.timer <= 0 {
				e.state = stateFlat
				e.timer = 20
				blastR := float64(blk) * 3
				for j := range g.enemies {
					other := &g.enemies[j]
					if other == e || other.state == stateFlat {
						continue
					}
					if math.Abs(other.wx-e.wx) < blastR && math.Abs(other.wy-e.wy) < blastR {
						other.state = stateFlat
						other.timer = 18
						g._addScore(200, other.wx-g.cam, other.wy-20)
					}
				}
				if g.invincible <= 0 && g.shrinkTimer <= 0 {
					if math.Abs(g.mwx-e.wx) < blastR && math.Abs(g.my-e.wy) < blastR {
						g._takeHit()
					}
				}
			}

		case stateFlat:
			e.timer--
			if e.timer <= 0 {
				continue
			}

		case stateShell:
			e.wx += e.vx
			for j := range g.enemies {
				other := &g.enemies[j]
				if other == e || other.state == stateFlat || other.state == stateShell {
					continue
				}
				if overlap(e.wx, e.wy, float64(blk), float64(blk), other.wx, other.wy, float64(blk), float64(blk)) {
					if other.state == stateShellStill {
						other.state = stateFlat
						other.timer = 18
					} else {
						other.state = stateFlat
						other.timer = 18
					}
					g._addScore(100, other.wx-g.cam, other.wy-20)
				}
			}
			for _, s := range solids {
				if overlap(e.wx, e.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
					e.vx *= -1
					break
				}
			}
			if g.invincible <= 0 && g.stompGrace <= 0 &&
				overlap(marioL+4, marioT, float64(blk)-8, MH, e.wx, e.wy, float64(blk), float64(blk)) {
				if g.mvy > 0 && marioB < e.wy+float64(blk)*0.5 {
					e.state = stateShellStill
					e.vx = 0
					e.timer = 300
					g.mvy = -10
					g.stompGrace = 25
					g._addScore(100, sx, e.wy-20)
				} else {
					g._takeHit()
					aliveEnemies = append(aliveEnemies, *e)
					for j := i + 1; j < len(g.enemies); j++ {
						aliveEnemies = append(aliveEnemies, g.enemies[j])
					}
					g.enemies = aliveEnemies
					return
				}
			}
			if math.Abs(e.wx-g.mwx) > g.W*2 {
				continue
			}
		}
		aliveEnemies = append(aliveEnemies, *e)
	}
	g.enemies = aliveEnemies

	// ── Score popups ──────────────────────────────────────────────────────────
	alivePopups := g.popups[:0]
	for i := range g.popups {
		p := &g.popups[i]
		p.y -= 3
		p.life--
		if p.life > 0 {
			alivePopups = append(alivePopups, *p)
		}
	}
	g.popups = alivePopups
}

// ── Game draw ─────────────────────────────────────────────────────────────────

func (g *Game) _draw(screen *ebiten.Image) {
	MH := float64(blk)
	if g.isBig {
		MH = 32 * pixelScale
	}
	MT := g.my
	if g.isBig {
		MT = g.my - (MH - blk)
	}
	_ = MT

	// Clouds
	for _, c := range g.clouds {
		drawCloud(screen, c.x, c.y, c.w, c.h)
	}

	// Ground tiles (wrapping, hide over gaps)
	tw := float64(len(g.groundTileXs)) * blk
	for i := range g.groundTileXs {
		gx := float64(i*blk) - math.Mod(g.cam, tw)
		if gx < -blk {
			gx += tw
		}
		gwx := g.cam + gx
		if g._inGap(gwx, blk) {
			continue // gap — don't draw ground tile
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(gx, g.H-48)
		screen.DrawImage(g.spr.groundBlock, op)
	}

	// Bricks
	aliveBricks := g.bricks[:0]
	for _, b := range g.bricks {
		sx := b.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		aliveBricks = append(aliveBricks, b)
		if sx < g.W+blk {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(sx, b.y)
			screen.DrawImage(g.spr.brickBlock, op)
		}
	}
	g.bricks = aliveBricks

	// ?-Blocks
	aliveQBlocks := g.qblocks[:0]
	for _, q := range g.qblocks {
		sx := q.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		aliveQBlocks = append(aliveQBlocks, q)
		if sx < g.W+blk {
			img := g.spr.qBlock
			if q.hit {
				img = g.spr.qBlockUsed
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(sx, q.y)
			screen.DrawImage(img, op)
		}
	}
	g.qblocks = aliveQBlocks

	// Coins
	for _, c := range g.coins {
		sx := c.wx - g.cam
		if sx < -float64(blk)*2 || sx > g.W+blk {
			continue
		}
		img := g.spr.coin1
		if (g.tick/8)%2 == 1 {
			img = g.spr.coin2
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, c.wy)
		screen.DrawImage(img, op)
	}

	// Mushrooms
	for _, m := range g.mushrooms {
		sx := m.wx - g.cam
		if sx < -float64(blk)*2 || sx > g.W+blk {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, m.wy)
		screen.DrawImage(g.spr.mushroomImg, op)
	}

	// Pipes
	pipesAlive := g.pipes[:0]
	for _, p := range g.pipes {
		sx := p.wx - g.cam
		if sx < -p.w*2 {
			continue
		}
		pipesAlive = append(pipesAlive, p)
		if sx <= g.W+p.w {
			drawPipe(screen, sx, p.y, p.w, p.h, p.lipH)
		}
	}
	g.pipes = pipesAlive

	// Fireballs
	fbSz := float64(8 * pixelScale)
	for _, fb := range g.fireballs {
		sx := fb.wx - g.cam
		img := g.spr.fireball1Img
		if (g.tick/3)%2 == 1 {
			img = g.spr.fireball2Img
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx-fbSz/2, fb.wy-fbSz/2)
		screen.DrawImage(img, op)
	}

	// Enemies
	for _, e := range g.enemies {
		sx := e.wx - g.cam
		if sx < -blk || sx > g.W+blk {
			continue
		}
		img := g._enemyImage(&e)
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(sx, e.wy)
			screen.DrawImage(img, op)
		}
	}

	// Mario
	if !g.dead {
		blink := (g.invincible > 0 || g.shrinkTimer > 0) && g.tick%4 < 2
		if !blink {
			msx := g.mwx - g.cam
			marioImg := g._marioImage()
			drawY := g.my
			if g.isBig {
				drawY = g.my - (32*pixelScale - blk)
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(msx, drawY)
			screen.DrawImage(marioImg, op)
		}
	} else {
		// Dead Mario pose
		msx := g.mwx - g.cam
		faceOff := 0
		if !g.facingRight {
			faceOff = 4
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(msx, g.my)
		screen.DrawImage(g.spr.marioSmall[3+faceOff], op)
	}

	// Score popups
	for _, p := range g.popups {
		ebitenutil.DebugPrintAt(screen, p.text, int(p.x), int(p.y))
	}

	// HUD
	hud := fmt.Sprintf("ARROWS:Move  SPACE:Jump(hold=higher)  SHIFT:Run  ESC:Hide  F/Z:Fireball  SCORE:%d", g.score)
	ebitenutil.DebugPrintAt(screen, hud, int(g.W/2)-len(hud)*3, 4)
}

// ── Helper methods ────────────────────────────────────────────────────────────

func (g *Game) _marioImage() *ebiten.Image {
	faceOff := 0
	if !g.facingRight {
		faceOff = 4
	}
	imgs := g.spr.marioSmall[:]
	if g.isBig {
		imgs = g.spr.marioBig[:]
	}
	if !g.onGround {
		return imgs[3+faceOff]
	}
	if math.Abs(g.mvx) > 0.5 {
		frame := 1 + (g.tick/4)%2
		return imgs[frame+faceOff]
	}
	return imgs[0+faceOff]
}

func (g *Game) _enemyImage(e *enemy) *ebiten.Image {
	t := g.tick
	switch e.kind {
	case kindGoomba:
		switch e.state {
		case stateWalk:
			return g.spr.goomba[(t/6)%2]
		case stateFuse:
			return g.spr.bobomb[(t/3)%2]
		default:
			return g.spr.goomba[2]
		}
	case kindKoopa:
		switch e.state {
		case stateWalk:
			return g.spr.koopa[(t/6)%2]
		default:
			return g.spr.koopa[2]
		}
	case kindRedKoopa:
		switch e.state {
		case stateWalk:
			if e.vx > 0 {
				return g.spr.redKoopa[3+(t/6)%2]
			}
			return g.spr.redKoopa[(t/6)%2]
		default:
			return g.spr.redKoopa[2]
		}
	case kindBobomb:
		switch e.state {
		case stateWalk:
			return g.spr.bobomb[(t/6)%2]
		case stateFuse:
			return g.spr.bobomb[(t/3)%2]
		default:
			return g.spr.bobomb[2]
		}
	}
	return nil
}

func (g *Game) _allSolids() [][4]float64 {
	out := make([][4]float64, 0, len(g.bricks)+len(g.qblocks)+len(g.pipes))
	for _, b := range g.bricks {
		out = append(out, [4]float64{b.wx, b.y, float64(blk), float64(blk)})
	}
	for _, q := range g.qblocks {
		out = append(out, [4]float64{q.wx, q.y, float64(blk), float64(blk)})
	}
	for _, p := range g.pipes {
		out = append(out, [4]float64{p.wx, p.y, p.w, p.h})
	}
	return out
}

func (g *Game) _inGap(wx, w float64) bool {
	cx := wx + w/2
	for _, gap := range g.gaps {
		if gap.startX < cx && cx < gap.endX {
			return true
		}
	}
	return false
}

func (g *Game) _die() {
	g.dead = true
	g.deadTimer = 40
	g.mvy = -14
	g.mvx = 0
	g.isBig = false
}

func (g *Game) _respawn() {
	g.dead = false
	g.mwx = math.Max(g.lastSafeWX, g.cam+200)
	for g._inGap(g.mwx, float64(blk)) {
		g.mwx += float64(blk)
	}
	for _, e := range g.enemies {
		if math.Abs(e.wx-g.mwx) < float64(blk)*2 {
			g.mwx = e.wx + float64(blk)*3
		}
	}
	g.my = g.groundY()
	g.mvx = 0
	g.mvy = 0
	g.onGround = true
	g.jumping = false
	g.invincible = 60
	g.stompGrace = 0
	g.isBig = false
	g.shrinkTimer = 0
}

func (g *Game) _takeHit() {
	if g.shrinkTimer > 0 || g.invincible > 0 {
		return
	}
	if g.isBig {
		g.isBig = false
		g.shrinkTimer = 60
		g.my += float64(blk)
	} else {
		g._die()
	}
}

func (g *Game) _addScore(pts int, sx, sy float64) {
	g.score += pts
	g.popups = append(g.popups, scorePopup{
		x: sx, y: sy, text: fmt.Sprintf("+%d", pts), life: 20,
	})
}

// ── Level generator ───────────────────────────────────────────────────────────

func (g *Game) _generate(lo, hi int) {
	x := float64(lo)
	if g.genX > x {
		x = g.genX
	}
	B := float64(blk)
	for x < float64(hi) {
		x += float64(rand.Intn(4)+3) * B
		r := rand.Float64()
		switch {
		case r < 0.11:
			// Floating brick row with one ?-block
			n := rand.Intn(2) + 3
			y := g.groundY() - B*3
			qi := rand.Intn(n)
			for i := 0; i < n; i++ {
				bx := x + float64(i)*B
				if i == qi {
					g.qblocks = append(g.qblocks, qblock{wx: bx, y: y, reward: "coin"})
					g.coins = append(g.coins, coin{wx: bx + B/4, wy: y - B})
				} else {
					g.bricks = append(g.bricks, brick{wx: bx, y: y})
				}
			}
			x += float64(n) * B
		case r < 0.17:
			// Mushroom ?-block
			g.qblocks = append(g.qblocks, qblock{
				wx: x, y: g.groundY() - B*3, reward: "mushroom",
			})
			x += B
		case r < 0.22:
			// Single ?-block (coin)
			g.qblocks = append(g.qblocks, qblock{wx: x, y: g.groundY() - B*3, reward: "coin"})
			x += B
		case r < 0.35:
			g.enemies = append(g.enemies, enemy{
				kind: kindGoomba, wx: x, wy: g.groundY(), vx: -1.5,
			})
		case r < 0.46:
			g.enemies = append(g.enemies, enemy{
				kind: kindKoopa, wx: x, wy: g.groundY(), vx: -1.5,
			})
		case r < 0.54:
			g.enemies = append(g.enemies, enemy{
				kind: kindRedKoopa, wx: x, wy: g.groundY(), vx: -1.5,
			})
		case r < 0.61:
			g.enemies = append(g.enemies, enemy{
				kind: kindBobomb, wx: x, wy: g.groundY(), vx: -1.0,
			})
		case r < 0.67:
			// Staircase
			h := rand.Intn(3) + 2
			for step := 0; step < h; step++ {
				g.bricks = append(g.bricks, brick{
					wx: x + float64(step)*B,
					y:  g.groundY() - float64(step+1)*B,
				})
			}
			x += float64(h) * B
		case r < 0.73:
			// Floating coins
			for i := 0; i < 3; i++ {
				cy := g.groundY() - B*2 - math.Sin(float64(i)/2*math.Pi)*B
				g.coins = append(g.coins, coin{wx: x + float64(i)*B, wy: cy})
			}
			x += 3 * B
		case r < 0.81:
			// Pipe
			ph := float64(rand.Intn(2)+2) * B
			pw := 2 * B
			py := g.groundY() - ph + B
			lipH := B / 3
			g.pipes = append(g.pipes, pipe{wx: x, y: py, w: pw, h: ph, lipH: lipH})
			x += pw
		case r < 0.87:
			// Gap
			gapW := float64(rand.Intn(2)+3) * B
			g.gaps = append(g.gaps, gap{startX: x, endX: x + gapW})
			x += gapW
		}
	}
	g.genX = x
}

// ── Overlap check (AABB) ──────────────────────────────────────────────────────

func overlap(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax+aw > bx && ax < bx+bw && ay+ah > by && ay < by+bh
}

// ── Drawing helpers ───────────────────────────────────────────────────────────

var (
	pipeGreen      = color.RGBA{0x20, 0xA0, 0x10, 0xFF}
	pipeDarkGreen  = color.RGBA{0x00, 0x68, 0x0C, 0xFF}
	pipeLightGreen = color.RGBA{0x80, 0xE0, 0x80, 0xFF}
	cloudWhite     = color.RGBA{0xF0, 0xF0, 0xF0, 0xFF}
)

func drawPipe(screen *ebiten.Image, sx, py, pw, ph, lipH float64) {
	lipExtra := float64(blk) / 4
	// Body
	ebitenutil.DrawRect(screen, sx, py+lipH, pw, ph-lipH, pipeGreen)
	// Lip
	ebitenutil.DrawRect(screen, sx-lipExtra, py, pw+2*lipExtra, lipH, pipeGreen)
	// Highlight strip
	ebitenutil.DrawRect(screen, sx+pw/3, py, 4, ph, pipeLightGreen)
	// Left edge
	ebitenutil.DrawRect(screen, sx-lipExtra, py, 2, lipH, pipeDarkGreen)
	// Right edge
	ebitenutil.DrawRect(screen, sx+pw+lipExtra-2, py, 2, lipH, pipeDarkGreen)
}

func drawCloud(screen *ebiten.Image, x, y, w, h float64) {
	// Three overlapping ellipses approximated as circles + a fill rect
	r1 := h * 0.45
	ebitenutil.DrawRect(screen, x+w*0.12, y+h*0.3, w*0.76, h*0.55, cloudWhite)
	drawFilledCircle(screen, x+w*0.25, y+h*0.55, r1)
	drawFilledCircle(screen, x+w*0.5, y+h*0.3, r1)
	drawFilledCircle(screen, x+w*0.75, y+h*0.55, r1)
}

// drawFilledCircle draws a filled circle using scanlines.
func drawFilledCircle(screen *ebiten.Image, cx, cy, r float64) {
	ir := int(math.Ceil(r))
	for dy := -ir; dy <= ir; dy++ {
		dx := math.Sqrt(r*r - float64(dy*dy))
		ebitenutil.DrawRect(screen, cx-dx, cy+float64(dy), 2*dx, 1, cloudWhite)
	}
}
