package main

import (
	"image/color"

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

// ── Sprite frame definitions ─────────────────────────────────────────────────

var marioStand = []string{
	"................",
	"......RRRRR.....",
	".....RRRRRRRRR..",
	".....BBBSSBSK...",
	"....BSBSSSKSSK..",
	"....BSBSSKSSSSK.",
	"....BBKSSSSK....",
	"......SSSSSSSS..",
	"....RRRORRRRR...",
	"...RRRRORRRORR..",
	"...RRRROOYOORR..",
	"...RRROOYOYORRR.",
	".....OOOYOOO....",
	".....OOO.OOO....",
	"....BBB...BBB...",
	"...BBBB...BBBB..",
}
var marioRun1 = []string{
	"................",
	"......RRRRR.....",
	".....RRRRRRRRR..",
	".....BBBSSBSK...",
	"....BSBSSSKSSK..",
	"....BSBSSKSSSSK.",
	"....BBKSSSSK....",
	"......SSSSSSSS..",
	"....RRRORRRRR...",
	"...RRRRORRRORR..",
	"...RRRROOYOORR..",
	"...RRROOYOYORRR.",
	".....OOOYOOO....",
	"......OOO.OO....",
	".......BBB.B....",
	"......BBB.BB....",
}
var marioRun2 = []string{
	"................",
	"......RRRRR.....",
	".....RRRRRRRRR..",
	".....BBBSSBSK...",
	"....BSBSSSKSSK..",
	"....BSBSSKSSSSK.",
	"....BBKSSSSK....",
	".....RSSSSSSSS..",
	"....RRRORRRRR...",
	"...RRRRORRRORR..",
	"...RRRROOYOORR..",
	"...RRROOYOYOO...",
	".....OOOYOOO....",
	"....OOO..OOO....",
	"...BBB....BBB...",
	"..BBBB.....BB...",
}
var marioJump = []string{
	"................",
	"......RRRRR.....",
	".....RRRRRRRRR..",
	".....BBBSSBSK...",
	"....BSBSSSKSSK..",
	"....BSBSSKSSSSK.",
	"....BBKSSSSK....",
	"......SSSSSSSS..",
	"...RRRRORRRRRR..",
	"..RRRRRORRRORR..",
	"..RRRRROOYOORR..",
	"..RRRROOYOYOO...",
	".....OOOYOOO.R..",
	"....OOO.OOOORR..",
	"...BBB...RRRR...",
	"..BBBB..........",
}

var bigMarioStand = []string{
	"................",
	".....RRRRR......",
	"....RRRRRRRRR...",
	"....BBBSSBSK....",
	"...BSBSSSKSSK...",
	"...BSBSSKSSSSK..",
	"...BBKSSSSK.....",
	".....SSSSSSSS...",
	"....RRSRRRS.....",
	"...RRRSSRRRSSS..",
	"...RRRSSSSRRSSS.",
	"...RRSSSSSSS....",
	".....SSSSSSS....",
	"....RRRORRR.....",
	"...RRRRORRRRR...",
	"...RRRROORRR....",
	".....OOOOOO.....",
	"....OOOOOOO.....",
	"...OOOOOOOO.....",
	"...OO.OOOO.OO...",
	"..OOO.OOOO.OOO..",
	"..OOO......OOO..",
	"......OOOO......",
	".....OOOOOO.....",
	"....BBBBBBBB....",
	"...BBBBBBBBBB...",
	"...BBBB..BBBB...",
	"...BBB....BBB...",
	"..BBBB....BBBB..",
	"..BBBB....BBBB..",
	"................",
	"................",
}
var bigMarioRun1 = []string{
	"................",
	".....RRRRR......",
	"....RRRRRRRRR...",
	"....BBBSSBSK....",
	"...BSBSSSKSSK...",
	"...BSBSSKSSSSK..",
	"...BBKSSSSK.....",
	".....SSSSSSSS...",
	"....RRSRRRS.....",
	"...RRRSSRRRSSS..",
	"...RRRSSSSRRSSS.",
	"...RRSSSSSSS....",
	".....SSSSSSS....",
	"....RRRORRR.....",
	"...RRRRORRRRR...",
	"...RRRROORRR....",
	".....OOOOOO.....",
	"....OOOOOOO.....",
	"...OOOOOOOO.....",
	"...OO.OOOO.OO...",
	"..OOO.OOOO.OOO..",
	"..OOO......OOO..",
	"......OOOO......",
	".....OOOOO......",
	"....BBBBB.......",
	"...BBBBBBBB.....",
	"...BBBB..BBBB...",
	"....BBB...BBB...",
	".....BBB..BBBB..",
	"......BB...BBB..",
	"................",
	"................",
}
var bigMarioRun2 = []string{
	"................",
	".....RRRRR......",
	"....RRRRRRRRR...",
	"....BBBSSBSK....",
	"...BSBSSSKSSK...",
	"...BSBSSKSSSSK..",
	"...BBKSSSSK.....",
	".....SSSSSSSS...",
	"....RRSRRRS.....",
	"...RRRSSRRRSSS..",
	"...RRRSSSSRRSSS.",
	"...RRSSSSSSS....",
	".....SSSSSSS....",
	"....RRRORRR.....",
	"...RRRRORRRRR...",
	"...RRRROORRR....",
	".....OOOOOO.....",
	"....OOOOOOO.....",
	"...OOOOOOOO.....",
	"...OO.OOOO.OO...",
	"..OOO.OOOO.OOO..",
	"..OOO......OOO..",
	".......OOOO.....",
	"......OOOOO.....",
	".......BBBBB....",
	".....BBBBBBBB...",
	"...BBBB..BBBB...",
	"...BBB...BBB....",
	"..BBBB..BBB.....",
	"..BBB...BB......",
	"................",
	"................",
}
var bigMarioJump = []string{
	"................",
	".....RRRRR......",
	"....RRRRRRRRR...",
	"....BBBSSBSK....",
	"...BSBSSSKSSK...",
	"...BSBSSKSSSSK..",
	"...BBKSSSSK.....",
	".....SSSSSSSS...",
	"....RRSRRRS.....",
	"...RRRSSRRRSSS..",
	"...RRRSSSSRRSSS.",
	"...RRSSSSSSS....",
	".....SSSSSSS....",
	"..RRRRORRRRRR...",
	".RRRRRORRRORRR..",
	".RRRRROOROORRR..",
	".RRRROOYOYOO....",
	"....OOOYOOO..R..",
	"...OOO.OOOOORRR.",
	"..OOO..OOOOORR..",
	"..OOO.....RRR...",
	"......OOO.......",
	".......OOOOO....",
	"....BBBBB..BBB..",
	"...BBBBBB...BB..",
	"...BBBBB........",
	"..BBBB..........",
	"..BBB...........",
	"................",
	"................",
	"................",
	"................",
}

var koopaL1 = []string{
	"................",
	"................",
	".GG.............",
	"GGGG............",
	"GGGGGG.SS.......",
	"GGGDGDG.SS......",
	"GDWDWDGSKKS.....",
	"GDWDWDGSKSS.....",
	"GDWDWDGSKS......",
	"GGDGDGDGSS......",
	"..G..G..........",
	".GG..GG.........",
	".BB..BB.........",
	"................",
	"................",
	"................",
}
var koopaL2 = []string{
	"................",
	"................",
	".GG.............",
	"GGGG............",
	"GGGGGG.SS.......",
	"GGGDGDG.SS......",
	"GDWDWDGSKKS.....",
	"GDWDWDGSKSS.....",
	"GDWDWDGSKS......",
	"GGDGDGDGSS......",
	"...G..G.........",
	"..GG..GG........",
	"..BB..BB........",
	"................",
	"................",
	"................",
}
var shellSprite = []string{
	"................",
	"................",
	"................",
	"................",
	"................",
	"....GGGG........",
	"...GGGDGDG......",
	"..GGDWDWDWG.....",
	"..GGDWDWDWG.....",
	"..GGDWDWDWG.....",
	"...GGGDGDG......",
	"....GGGG........",
	"................",
	"................",
	"................",
	"................",
}
var goomba1 = []string{
	"................",
	"................",
	"................",
	"......BBBB......",
	".....BBBBBB.....",
	"....BKBBBBKB....",
	"....BKWBBWKB....",
	"...BBKKBBKKBB...",
	"...BBBSBBSBBB...",
	"...BBBBBBBBB....",
	"....BBBBBBB.....",
	".....BBBBB......",
	"....SSSSSSSS....",
	"...SSSSSSSSSS...",
	"...BB......BB...",
	"..BBB......BBB..",
}
var goomba2 = []string{
	"................",
	"................",
	"................",
	"......BBBB......",
	".....BBBBBB.....",
	"....BKBBBBKB....",
	"....BKWBBWKB....",
	"...BBKKBBKKBB...",
	"...BBBSBBSBBB...",
	"...BBBBBBBBB....",
	"....BBBBBBB.....",
	".....BBBBB......",
	"....SSSSSSSS....",
	"...SSSSSSSSSS...",
	"....BB....BB....",
	"...BBB....BBB...",
}
var goombaFlat = []string{
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"................",
	"...BBBBBBBBBB...",
	"..BKWBBBBWKBB..",
	"..BBKKBBKKBBBB..",
	"...BBBBBBBBB....",
}
var brickFrame = []string{
	"TTTTTTTTTTTTTTTT",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"MMMMMMMMMMMMMMMM",
	"TTTMTTTTTTTMTTTT",
	"TTTMTTTTTTTMTTTT",
	"MMMMMMMMMMMMMMMM",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"MMMMMMMMMMMMMMMM",
	"TTTMTTTTTTTMTTTT",
	"TTTMTTTTTTTMTTTT",
	"MMMMMMMMMMMMMMMM",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"TTTTTTTTTTTTTTTT",
}
var qblockFrame = []string{
	"KKKKKKKKKKKKKKKK",
	"KQQQQQQQQQQQQQQK",
	"KQQQQQQQQQQQQQQK",
	"KQQQQKKKKKKQQQQK",
	"KQQQKKQQQQKKQQQK",
	"KQQQQQQQQQKKQQQK",
	"KQQQQQQQQKKQQQQK",
	"KQQQQQQQKKQQQQQK",
	"KQQQQQQKKQQQQQQK",
	"KQQQQQQKKQQQQQQK",
	"KQQQQQQQQQQQQQQK",
	"KQQQQQQKKQQQQQQK",
	"KQQQQQQKKQQQQQQK",
	"KQQQQQQQQQQQQQQK",
	"KQQQQQQQQQQQQQQK",
	"KKKKKKKKKKKKKKKK",
}
var qblockUsedFrame = []string{
	"KKKKKKKKKKKKKKKK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KMMMMMMMMMMMMMMK",
	"KKKKKKKKKKKKKKKK",
}
var groundBlockFrame = []string{
	"TTTTTTTTTTTTTTTT",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"MMMMMMMMMMMMMMMM",
	"TTTMTTTTTTTMTTTT",
	"TTTMTTTTTTTMTTTT",
	"MMMMMMMMMMMMMMMM",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"MMMMMMMMMMMMMMMM",
	"TTTMTTTTTTTMTTTT",
	"TTTMTTTTTTTMTTTT",
	"MMMMMMMMMMMMMMMM",
	"TMTTTTTMTTTTTTMT",
	"TMTTTTTMTTTTTTMT",
	"TTTTTTTTTTTTTTTT",
}
var coin1Frame = []string{
	"........",
	"........",
	"........",
	"...YY...",
	"..YYYY..",
	".YYCYYY.",
	".YCCYYY.",
	".YYCYYY.",
	".YYCYYY.",
	".YCCYYY.",
	".YYCYYY.",
	"..YYYY..",
	"...YY...",
	"........",
	"........",
	"........",
}
var coin2Frame = []string{
	"........",
	"........",
	"........",
	"...YY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...CY...",
	"...YY...",
	"........",
	"........",
	"........",
}
var mushroomFrame = []string{
	"................",
	"......RRRR......",
	"....RRRRRRRR....",
	"...RRWWRRWWRR...",
	"..RRRWWRRWWRRR..",
	"..RRRWWRRWWRRR..",
	".RRRRRRRRRRRRR..",
	".RRRRRRRRRRRRR..",
	"..SSSSMMMMSSS...",
	"..SSSMMMMMMSSS..",
	"..SSSMMMMMMSSS..",
	"..SSSMMMMMMSSS..",
	"...SSSMMMMSSS...",
	"....SSSSSSSS....",
	"................",
	"................",
}
var bobomb1 = []string{
	"................",
	"........YY......",
	".......YOY......",
	"........KK......",
	"......KKKK......",
	".....KKKKKK.....",
	"....KWKKKWKK....",
	"....KWWKKWWK....",
	"...KKKKKKKKKK...",
	"...KKKKKKKKKK...",
	"...KKKKKKKKKK...",
	"....KKKKKKKK....",
	".....YKKKY......",
	"....YY...YY.....",
	"................",
	"................",
}
var bobomb2 = []string{
	"................",
	".......YYY......",
	"......YOOY......",
	"........KK......",
	"......KKKK......",
	".....KKKKKK.....",
	"....KWKKKWKK....",
	"....KWWKKWWK....",
	"...KKKKKKKKKK...",
	"...KKKKKKKKKK...",
	"...KKKKKKKKKK...",
	"....KKKKKKKK....",
	"....YKKKKKY.....",
	"...YY.....YY....",
	"................",
	"................",
}
var bobombExplode = []string{
	"................",
	"....YY..YY......",
	"...YOOYY.OY.....",
	"..YOOOOOOOY.....",
	"..YOOOOOOOOOY...",
	".YOOOOOOOOOOY...",
	".YROOROOORORY...",
	".YRRRROOORORY...",
	".YRRRRROORRRY...",
	"..YRRRRRRROY....",
	"..YRRRRRRRRY....",
	"...YRRRRRYY.....",
	"....YYRRYY......",
	".....YYY........",
	"................",
	"................",
}
var fireball1 = []string{
	"........",
	"........",
	"...YY...",
	"..YRYY..",
	"..YRRY..",
	"..YYYY..",
	"...YY...",
	"........",
}
var fireball2 = []string{
	"........",
	"........",
	"..YYY...",
	"..YRRY..",
	".YRRRY..",
	"..YRRY..",
	"..YYY...",
	"........",
}
