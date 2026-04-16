package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// _update runs one game-logic tick.
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
		if g.deadTimer > 50 {
			g.my -= 4
		} else {
			g.my += 5
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

	g._updatePhysics(MT, MH)
	if g.dead {
		return
	}

	solids := g._allSolids()
	g._updateCoins(solids, MT, MH)
	g._updateMushrooms(solids, MT, MH)
	g._updateFireballs(solids)

	marioB := g.my + float64(blk)
	if g._updateEnemies(solids, g.mwx, MT, marioB, MH) {
		return
	}

	g._updatePopups()
}

// _updatePhysics handles player input, movement, collision resolution, and camera.
// Physics constants are tuned for 60 TPS.
func (g *Game) _updatePhysics(MT, MH float64) {
	// Input → velocity
	running := g.keys[ebiten.KeyShiftLeft] || g.keys[ebiten.KeyShiftRight]
	accel := 0.4
	maxSpeed := 2.5
	friction := 0.3
	if running {
		accel = 0.6
		maxSpeed = 4.0
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
		g.mvy = -8
		g.onGround = false
		g.jumping = true
	}
	if g.jumping && g.keys[ebiten.KeySpace] && g.mvy < 0 {
		g.mvy += 0.275
	} else {
		g.mvy += 0.6
		if g.mvy >= 0 {
			g.jumping = false
		}
	}
	if g.mvy > 9 {
		g.mvy = 9
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
			vx: 3.5 * dir, vy: -2.0,
		})
		g.fireCooldown = 24
	}

	solids := g._allSolids()

	// Move X → check solids
	g.mwx += g.mvx
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
}
