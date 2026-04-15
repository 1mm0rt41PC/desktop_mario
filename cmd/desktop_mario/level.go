package main

import (
	"fmt"
	"math"
	"math/rand"
)

// overlap returns true when two axis-aligned rectangles intersect.
func overlap(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax+aw > bx && ax < bx+bw && ay+ah > by && ay < by+bh
}

// _allSolids returns bounding boxes ([x,y,w,h]) for all solid world objects.
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

// _inGap returns true if the world-x centre of the given span falls inside any gap.
func (g *Game) _inGap(wx, w float64) bool {
	cx := wx + w/2
	for _, gap := range g.gaps {
		if gap.startX < cx && cx < gap.endX {
			return true
		}
	}
	return false
}

// _die starts Mario's death animation.
func (g *Game) _die() {
	g.dead = true
	g.deadTimer = 40
	g.mvy = -14
	g.mvx = 0
	g.isBig = false
}

// _respawn resets Mario after a death.
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

// _takeHit handles Mario being hit by an enemy.
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

// _addScore awards points and spawns a floating score popup.
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
