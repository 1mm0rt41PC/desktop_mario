package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

// ── Sprite selection helpers ──────────────────────────────────────────────────

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

// ── Drawing primitives ────────────────────────────────────────────────────────

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
