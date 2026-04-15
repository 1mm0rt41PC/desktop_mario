package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Color palette (mirrors PAL dict in mario_enhanced.py).
// 'K' uses near-black (8,8,8) to avoid coinciding with the Win32 color-key
// for transparency (0,0,0 = pure black).
var palette = map[byte]color.RGBA{
	'.': {0, 0, 0, 0},
	'R': {0xE4, 0x08, 0x18, 0xFF},
	'B': {0xAC, 0x7C, 0x00, 0xFF},
	'S': {0xFC, 0xA0, 0x44, 0xFF},
	'O': {0x00, 0x32, 0xEC, 0xFF},
	'Y': {0xFF, 0xD8, 0x00, 0xFF},
	'G': {0x20, 0xA0, 0x10, 0xFF},
	'D': {0x00, 0x68, 0x0C, 0xFF},
	'W': {0xFF, 0xFF, 0xFF, 0xFF},
	'K': {0x08, 0x08, 0x08, 0xFF},
	'T': {0xC8, 0x4C, 0x0C, 0xFF},
	'M': {0xE8, 0x90, 0x50, 0xFF},
	'Q': {0xF8, 0xB8, 0x00, 0xFF},
	'C': {0xF8, 0xD8, 0x78, 0xFF},
}

// pixelScale is the rendering scale factor: 1 logical pixel → N×N screen pixels.
const pixelScale = 3

// frameToImage converts a []string pixel-art frame into an *ebiten.Image.
// Each character maps to a color in palette; '.' means transparent.
func frameToImage(frame []string, scale int) *ebiten.Image {
	h := len(frame)
	if h == 0 {
		return ebiten.NewImage(1, 1)
	}
	w := 0
	for _, row := range frame {
		if len(row) > w {
			w = len(row)
		}
	}
	imgW := w * scale
	imgH := h * scale
	pixels := make([]byte, imgW*imgH*4)
	for y, row := range frame {
		for x := 0; x < len(row); x++ {
			c := palette[row[x]]
			if c.A == 0 {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					idx := ((y*scale+dy)*imgW + (x*scale+dx)) * 4
					pixels[idx] = c.R
					pixels[idx+1] = c.G
					pixels[idx+2] = c.B
					pixels[idx+3] = c.A
				}
			}
		}
	}
	img := ebiten.NewImage(imgW, imgH)
	img.WritePixels(pixels)
	return img
}

// flipFrame mirrors a sprite horizontally.
func flipFrame(frame []string) []string {
	out := make([]string, len(frame))
	for i, row := range frame {
		b := []byte(row)
		for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
			b[l], b[r] = b[r], b[l]
		}
		out[i] = string(b)
	}
	return out
}

// scalePNGImage scales an ebiten.Image by the given integer factor using
// nearest-neighbour filtering, returning the enlarged image.
func scalePNGImage(src *ebiten.Image, scale int) *ebiten.Image {
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	dst := ebiten.NewImage(sw*scale, sh*scale)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.Filter = ebiten.FilterNearest
	dst.DrawImage(src, op)
	return dst
}

// loadEmbeddedPNG loads a PNG from the embedded filesystem and returns a scaled
// ebiten.Image (scaled by pixelScale).  Returns nil on any error.
func loadEmbeddedPNG(path string) *ebiten.Image {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		log.Printf("loadEmbeddedPNG: %s: %v", path, err)
		return nil
	}
	raw, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("loadEmbeddedPNG decode: %s: %v", path, err)
		return nil
	}
	img := ebiten.NewImageFromImage(raw)
	return scalePNGImage(img, pixelScale)
}

// loadSheetSprite crops a sub-image from a sprite sheet PNG and scales it.
// cropBox is (x1, y1, x2, y2) in the original (unscaled) sheet.
func loadSheetSprite(sheetPath string, sheet *ebiten.Image, cropBox [4]int) *ebiten.Image {
	if sheet == nil {
		return nil
	}
	x1, y1, x2, y2 := cropBox[0], cropBox[1], cropBox[2], cropBox[3]
	sub := sheet.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image)
	return scalePNGImage(sub, pixelScale)
}

// flipEbitenImage returns a horizontally mirrored copy of an ebiten.Image.
func flipEbitenImage(src *ebiten.Image) *ebiten.Image {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dst := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(w), 0)
	dst.DrawImage(src, op)
	return dst
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

// buildSprites constructs all sprite images, preferring PNG assets from the
// embedded filesystem and falling back to the string-based pixel art.
func buildSprites() spriteImages {
	var s spriteImages
	px := pixelScale

	// ── Load enemy sprite sheet (used for goomba, koopa, shell) ──────────────
	enemySheetImg := loadRawEmbeddedPNG("assets/enemies.png")

	// Enemy sheet crop boxes (x1,y1,x2,y2) in sheet pixels – same as Python.
	goombaCrops := [3][4]int{
		{0, 19, 18, 36},
		{18, 19, 36, 36},
		{36, 27, 54, 36},
	}
	greenKoopaCrops := [2][4]int{
		{54, 126, 72, 168},
		{72, 126, 90, 168},
	}
	redKoopaCrops := [2][4]int{
		{54, 168, 72, 210},
		{72, 168, 90, 210},
	}
	greenShellCrop := [4]int{126, 126, 144, 168}
	redShellCrop := [4]int{126, 168, 144, 210}

	// ── Small Mario ──────────────────────────────────────────────────────────
	// PNG frames: small_0=stand, small_2=run1, small_3=run2, small_5=jump
	smPNG := [4]*ebiten.Image{
		loadEmbeddedPNG("assets/sprites/small_0.png"),
		loadEmbeddedPNG("assets/sprites/small_2.png"),
		loadEmbeddedPNG("assets/sprites/small_3.png"),
		loadEmbeddedPNG("assets/sprites/small_5.png"),
	}
	smFallback := [4][]string{marioStand, marioRun1, marioRun2, marioJump}
	for i := 0; i < 4; i++ {
		if smPNG[i] != nil {
			s.marioSmall[i] = smPNG[i]
			s.marioSmall[i+4] = flipEbitenImage(smPNG[i])
		} else {
			s.marioSmall[i] = frameToImage(smFallback[i], px)
			s.marioSmall[i+4] = frameToImage(flipFrame(smFallback[i]), px)
		}
	}

	// ── Big Mario ────────────────────────────────────────────────────────────
	// PNG frames: big_0=stand, big_2=run1, big_3=run2, big_5=jump
	bgPNG := [4]*ebiten.Image{
		loadEmbeddedPNG("assets/sprites/big_0.png"),
		loadEmbeddedPNG("assets/sprites/big_2.png"),
		loadEmbeddedPNG("assets/sprites/big_3.png"),
		loadEmbeddedPNG("assets/sprites/big_5.png"),
	}
	bgFallback := [4][]string{bigMarioStand, bigMarioRun1, bigMarioRun2, bigMarioJump}
	for i := 0; i < 4; i++ {
		if bgPNG[i] != nil {
			s.marioBig[i] = bgPNG[i]
			s.marioBig[i+4] = flipEbitenImage(bgPNG[i])
		} else {
			s.marioBig[i] = frameToImage(bgFallback[i], px)
			s.marioBig[i+4] = frameToImage(flipFrame(bgFallback[i]), px)
		}
	}

	// ── Enemies ──────────────────────────────────────────────────────────────
	for i := 0; i < 3; i++ {
		img := loadSheetSprite("assets/enemies.png", enemySheetImg, goombaCrops[i])
		if img != nil {
			s.goomba[i] = img
		}
	}
	if s.goomba[0] == nil {
		s.goomba[0] = frameToImage(goomba1, px)
		s.goomba[1] = frameToImage(goomba2, px)
		s.goomba[2] = frameToImage(goombaFlat, px)
	}

	for i := 0; i < 2; i++ {
		img := loadSheetSprite("assets/enemies.png", enemySheetImg, greenKoopaCrops[i])
		if img != nil {
			s.koopa[i] = img
		}
	}
	shellGreen := loadSheetSprite("assets/enemies.png", enemySheetImg, greenShellCrop)
	if shellGreen != nil {
		s.koopa[2] = shellGreen
	}
	if s.koopa[0] == nil {
		s.koopa[0] = frameToImage(koopaL1, px)
		s.koopa[1] = frameToImage(koopaL2, px)
		s.koopa[2] = frameToImage(shellSprite, px)
	}

	for i := 0; i < 2; i++ {
		img := loadSheetSprite("assets/enemies.png", enemySheetImg, redKoopaCrops[i])
		if img != nil {
			s.redKoopa[i] = img
			s.redKoopa[i+3] = flipEbitenImage(img)
		}
	}
	shellRed := loadSheetSprite("assets/enemies.png", enemySheetImg, redShellCrop)
	if shellRed != nil {
		s.redKoopa[2] = shellRed
	}
	if s.redKoopa[0] == nil {
		s.redKoopa[0] = frameToImage(koopaL1, px)
		s.redKoopa[1] = frameToImage(koopaL2, px)
		s.redKoopa[2] = frameToImage(shellSprite, px)
		s.redKoopa[3] = frameToImage(flipFrame(koopaL1), px)
		s.redKoopa[4] = frameToImage(flipFrame(koopaL2), px)
	}

	// Bobombs always use string sprites (no PNG asset).
	s.bobomb[0] = frameToImage(bobomb1, px)
	s.bobomb[1] = frameToImage(bobomb2, px)
	s.bobomb[2] = frameToImage(bobombExplode, px)

	// ── Blocks / items ────────────────────────────────────────────────────────
	s.groundBlock = frameToImage(groundBlockFrame, px)
	s.brickBlock = frameToImage(brickFrame, px)
	s.qBlock = frameToImage(qblockFrame, px)
	s.qBlockUsed = frameToImage(qblockUsedFrame, px)
	s.coin1 = frameToImage(coin1Frame, px)
	s.coin2 = frameToImage(coin2Frame, px)

	// Mushroom: try PNG (item_0.png), fall back to string sprite.
	mush := loadEmbeddedPNG("assets/sprites/item_0.png")
	if mush != nil {
		s.mushroomImg = mush
	} else {
		s.mushroomImg = frameToImage(mushroomFrame, px)
	}

	s.fireball1Img = frameToImage(fireball1, px)
	s.fireball2Img = frameToImage(fireball2, px)
	return s
}

// loadRawEmbeddedPNG loads a PNG from the embedded filesystem WITHOUT scaling.
// Used when we need to crop sub-images from a sprite sheet.
func loadRawEmbeddedPNG(path string) *ebiten.Image {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		log.Printf("loadRawEmbeddedPNG: %s: %v", path, err)
		return nil
	}
	raw, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("loadRawEmbeddedPNG decode: %s: %v", path, err)
		return nil
	}
	return ebiten.NewImageFromImage(raw)
}
