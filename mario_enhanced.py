"""
Desktop Mario – A stress-relief mini-game that lives on your desktop.
  - Ctrl+Alt+M  to show/hide the game
  - Arrow keys to move, Space to jump (hold longer = higher!)
  - Shift to run faster
  - ESC to hide
  - Right-click system tray icon to quit
"""
import tkinter as tk
import random
import ctypes
import ctypes.wintypes as wt
import math
import threading
import sys
import os

# ============================================================
#  PALETTE
# ============================================================
PAL = {
    '.': None,
    'R': '#E40818', 'B': '#AC7C00', 'S': '#FCA044', 'O': '#0032EC',
    'Y': '#FFD800', 'G': '#20A010', 'D': '#00680C', 'W': '#FFFFFF',
    'K': '#000000', 'T': '#C84C0C', 'M': '#E89050', 'Q': '#F8B800',
    'C': '#F8D878',
}
PX = 3

# ============================================================
#  SPRITES  (all 16x16 unless noted)
# ============================================================
MARIO_STAND = [
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
]
MARIO_RUN1 = [
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
]
MARIO_RUN2 = [
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
]
MARIO_JUMP = [
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
]
KOOPA_L1 = [
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
]
KOOPA_L2 = [
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
]
SHELL_SPRITE = [
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
]
GOOMBA_1 = [
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
]
GOOMBA_2 = [
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
]
GOOMBA_FLAT = [
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
]
BRICK = [
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
]
QBLOCK = [
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
]
QBLOCK_USED = [
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
]
GROUND_BLOCK = [
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
]
COIN1 = [
    "........", "........", "........",
    "...YY...", "..YYYY..", ".YYCYYY.",
    ".YCCYYY.", ".YYCYYY.", ".YYCYYY.",
    ".YCCYYY.", ".YYCYYY.", "..YYYY..",
    "...YY...", "........", "........", "........",
]
COIN2 = [
    "........", "........", "........",
    "...YY...", "...CY...", "...CY...",
    "...CY...", "...CY...", "...CY...",
    "...CY...", "...CY...", "...CY...",
    "...YY...", "........", "........", "........",
]

# ---- Sprite flipping helper ----
def _flip_rows(sprite):
    """Mirror a sprite horizontally."""
    return [row[::-1] for row in sprite]

MARIO_STAND_L = _flip_rows(MARIO_STAND)
MARIO_RUN1_L  = _flip_rows(MARIO_RUN1)
MARIO_RUN2_L  = _flip_rows(MARIO_RUN2)
MARIO_JUMP_L  = _flip_rows(MARIO_JUMP)

KOOPA_R1 = _flip_rows(KOOPA_L1)
KOOPA_R2 = _flip_rows(KOOPA_L2)

# ---- Big Mario sprites (16x32 - double height) ----
BIG_MARIO_STAND = [
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
]
BIG_MARIO_RUN1 = [
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
]
BIG_MARIO_RUN2 = [
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
]
BIG_MARIO_JUMP = [
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
]

BIG_MARIO_STAND_L = _flip_rows(BIG_MARIO_STAND)
BIG_MARIO_RUN1_L  = _flip_rows(BIG_MARIO_RUN1)
BIG_MARIO_RUN2_L  = _flip_rows(BIG_MARIO_RUN2)
BIG_MARIO_JUMP_L  = _flip_rows(BIG_MARIO_JUMP)

# ---- Mushroom sprite ----
MUSHROOM = [
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
]

# ---- Bob-omb sprite ----
BOBOMB_1 = [
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
]
BOBOMB_2 = [
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
]
BOBOMB_EXPLODE = [
    "................",
    "....YY..YY......",
    "...YOOYY.OY....",
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
]

# ============================================================
#  SPRITE ENGINE  (PhotoImage – 1 canvas item per sprite)
# ============================================================
_img_cache = {}

def _frame_to_photo(frame, px):
    h = len(frame)
    w = max(len(r) for r in frame)
    img = tk.PhotoImage(width=w * px, height=h * px)
    buckets = {}
    for ry, row in enumerate(frame):
        for cx, ch in enumerate(row):
            col = PAL.get(ch)
            if col:
                buckets.setdefault(col, []).append((cx, ry))
    for col, pts in buckets.items():
        row_data = f"{{{col}}} " * px
        for cx, ry in pts:
            x0, y0 = cx * px, ry * px
            for dy in range(px):
                img.put(row_data, to=(x0, y0 + dy, x0 + px, y0 + dy + 1))
    return img

class Sprite:
    __slots__ = ('canvas', 'frames', 'px', '_imgs', 'item', 'cur', 'x', 'y')
    def __init__(self, canvas, frames, px=PX):
        self.canvas = canvas
        self.frames = frames
        self.px = px
        self._imgs = []
        self.item = None
        self.cur = -1
        self.x = 0
        self.y = 0
        for idx, frame in enumerate(frames):
            key = (id(frame[0]), idx, px, len(frames))
            if key in _img_cache:
                self._imgs.append(_img_cache[key])
            else:
                photo = _frame_to_photo(frame, px)
                _img_cache[key] = photo
                self._imgs.append(photo)

    def draw(self, idx):
        idx %= len(self.frames)
        if idx == self.cur:
            return
        self.cur = idx
        if self.item is None:
            self.item = self.canvas.create_image(
                self.x, self.y, image=self._imgs[idx], anchor='nw')
        else:
            self.canvas.itemconfig(self.item, image=self._imgs[idx])

    def move_to(self, x, y):
        if self.item is None:
            self.x, self.y = x, y
            return
        dx, dy = x - self.x, y - self.y
        self.x, self.y = x, y
        if dx or dy:
            self.canvas.move(self.item, dx, dy)

    def destroy(self):
        if self.item is not None:
            self.canvas.delete(self.item)
            self.item = None

# ============================================================
#  SCORE POPUP
# ============================================================
class ScorePopup:
    def __init__(self, canvas, x, y, text="+100"):
        self.canvas = canvas
        self.id = canvas.create_text(x, y, text=text, fill='#FFFFFF',
                                     font=('Courier', 14, 'bold'))
        self.life = 20
    def update(self):
        self.canvas.move(self.id, 0, -3)
        self.life -= 1
        if self.life <= 0:
            self.canvas.delete(self.id)
            return False
        return True

# ============================================================
#  GAME SCENE
# ============================================================
class Game:
    BLK = 16 * PX
    SPR = 16 * PX      # sprite size

    def __init__(self, canvas, W, H):
        self.canvas = canvas
        self.W = W
        self.H = H
        self.tick = 0
        self.score = 0
        self.popups = []

        self.ground_y = H - 48 - self.BLK   # top of ground row

        # ---- INPUT STATE ----
        self.keys = set()

        # ---- MARIO (frames 0-3 right, 4-7 left) ----
        self.mario = Sprite(canvas, [
            MARIO_STAND, MARIO_RUN1, MARIO_RUN2, MARIO_JUMP,
            MARIO_STAND_L, MARIO_RUN1_L, MARIO_RUN2_L, MARIO_JUMP_L,
        ])
        self.big_mario = Sprite(canvas, [
            BIG_MARIO_STAND, BIG_MARIO_RUN1, BIG_MARIO_RUN2, BIG_MARIO_JUMP,
            BIG_MARIO_STAND_L, BIG_MARIO_RUN1_L, BIG_MARIO_RUN2_L, BIG_MARIO_JUMP_L,
        ])
        self.big_mario.move_to(-200, -200)  # hide offscreen until needed
        self.is_big = False        # mushroom power-up state
        self.shrink_timer = 0      # invincibility after getting hit while big
        self.mwx = 200.0          # world-x
        self.my  = self.ground_y  # screen-y (top of sprite)
        self.mvx = 0.0
        self.mvy = 0.0
        self.on_ground = True
        self.facing_right = True
        self.jumping = False       # variable-height jump tracking
        self.stomp_grace = 0       # invincibility frames after stomp/kick
        self.invincible = 60       # spawn invincibility (like original SMB)

        # ---- DEATH / RESPAWN ----
        self.dead = False
        self.dead_timer = 0
        self.last_safe_wx = 200.0  # last world-x where Mario stood on ground

        # ---- CAMERA ----
        self.cam = 0.0

        # ---- CLOUDS ----
        self.clouds = []
        for _ in range(5):
            self.clouds.append(self._mk_cloud(
                random.randint(0, W), random.randint(30, max(60, H // 3))))

        # ---- GROUND (looping) ----
        self.ground_tiles = []
        n = (W // self.BLK) + 3
        for i in range(n):
            g = Sprite(canvas, [GROUND_BLOCK])
            g.draw(0)
            g.move_to(i * self.BLK, H - 48)
            self.ground_tiles.append(g)

        # ---- WORLD LISTS ----
        self.bricks  = []
        self.qblocks = []
        self.coins   = []
        self.enemies = []
        self.pipes   = []          # {'wx', 'y', 'w', 'h', 'lip_h', 'ids'}
        self.gaps    = []          # (start_wx, end_wx) – no ground here
        self.mushrooms = []        # {'s', 'wx', 'wy', 'vx', 'vy', 'active'}
        self.gen_x   = 0

        # ---- HUD ----
        self.hud = canvas.create_text(
            W // 2, 18,
            text="ARROWS: Move  SPACE: Jump (hold=higher)  SHIFT: Run  ESC: Quit  SCORE: 0",
            fill='#FFFFFF', font=('Courier', 12, 'bold'), anchor='n')

        # Start generation well past Mario's spawn (wx=200) — safe runway
        self._generate(600, W + 600)

    # ---- cloud helpers ----
    def _mk_cloud(self, x, y):
        tag = f"cl{id(self)}_{random.randint(0,99999)}"
        w = random.randint(70, 120)
        h = random.randint(26, 40)
        spd = random.uniform(0.2, 0.4)
        col = random.choice(['#f8f8f8', '#f0f0f0', '#e8eef8'])
        self.canvas.create_oval(0, 0, w * .45, h * .9, fill=col, outline='', tags=tag)
        self.canvas.create_oval(w * .2, -h * .25, w * .8, h * .75, fill=col, outline='', tags=tag)
        self.canvas.create_oval(w * .5, 0, w, h * .85, fill=col, outline='', tags=tag)
        self.canvas.create_rectangle(w * .15, h * .3, w * .85, h * .78, fill=col, outline='', tags=tag)
        self.canvas.move(tag, x, y)
        return {'tag': tag, 'x': x, 'y': y, 'speed': spd, 'w': w}

    def _reset_cloud(self, c):
        nx = self.W + random.randint(50, 200)
        ny = random.randint(30, max(60, self.H // 3))
        self.canvas.move(c['tag'], nx - c['x'], ny - c['y'])
        c['x'], c['y'] = nx, ny

    # ---- gap detection ----
    def _is_in_gap(self, wx, w=None):
        """Check if world-x position is over a pit."""
        if w is None:
            w = self.SPR
        cx = wx + w / 2
        for gs, ge in self.gaps:
            if gs < cx < ge:
                return True
        return False

    # ---- LEVEL GENERATOR ----
    def _generate(self, lo, hi):
        x = max(lo, self.gen_x)
        B = self.BLK
        while x < hi:
            # 3-6 blocks spacing between features (tighter, more interesting)
            x += random.randint(3, 6) * B
            r = random.random()
            if r < 0.11:
                # floating brick row (3-4) with one ?-block (coin)
                n = random.randint(3, 4)
                y = self.ground_y - B * 3
                qi = random.randint(0, n - 1)
                for i in range(n):
                    bx = x + i * B
                    if i == qi:
                        s = Sprite(self.canvas, [QBLOCK, QBLOCK_USED])
                        s.draw(0)
                        self.qblocks.append({'s': s, 'wx': bx, 'y': y, 'hit': False, 'reward': 'coin'})
                    else:
                        s = Sprite(self.canvas, [BRICK])
                        s.draw(0)
                        self.bricks.append({'s': s, 'wx': bx, 'y': y})
                co = Sprite(self.canvas, [COIN1, COIN2]); co.draw(0)
                self.coins.append({'s': co, 'wx': x + qi * B + B // 4, 'wy': y - B, 'got': False, 'ft': 0})
                x += n * B
            elif r < 0.17:
                # mushroom ?-block (spawns mushroom!)
                s = Sprite(self.canvas, [QBLOCK, QBLOCK_USED]); s.draw(0)
                self.qblocks.append({'s': s, 'wx': x, 'y': self.ground_y - B * 3, 'hit': False, 'reward': 'mushroom'})
                x += B
            elif r < 0.22:
                # single ?-block (coin)
                s = Sprite(self.canvas, [QBLOCK, QBLOCK_USED]); s.draw(0)
                self.qblocks.append({'s': s, 'wx': x, 'y': self.ground_y - B * 3, 'hit': False, 'reward': 'coin'})
                x += B
            elif r < 0.32:
                # goomba
                s = Sprite(self.canvas, [GOOMBA_1, GOOMBA_2, GOOMBA_FLAT]); s.draw(0)
                self.enemies.append({'s': s, 'wx': float(x), 'wy': self.ground_y,
                                     'kind': 'goomba', 'vx': -1.5, 'state': 'walk', 'timer': 0})
            elif r < 0.40:
                # green koopa
                s = Sprite(self.canvas, [KOOPA_L1, KOOPA_L2, SHELL_SPRITE]); s.draw(0)
                self.enemies.append({'s': s, 'wx': float(x), 'wy': self.ground_y,
                                     'kind': 'koopa', 'vx': -1.5, 'state': 'walk', 'timer': 0})
            elif r < 0.46:
                # red koopa (turns at edges)
                s = Sprite(self.canvas, [KOOPA_L1, KOOPA_L2, SHELL_SPRITE, KOOPA_R1, KOOPA_R2]); s.draw(0)
                self.enemies.append({'s': s, 'wx': float(x), 'wy': self.ground_y,
                                     'kind': 'red_koopa', 'vx': -1.5, 'state': 'walk', 'timer': 0})
            elif r < 0.52:
                # bob-omb
                s = Sprite(self.canvas, [BOBOMB_1, BOBOMB_2, BOBOMB_EXPLODE]); s.draw(0)
                self.enemies.append({'s': s, 'wx': float(x), 'wy': self.ground_y,
                                     'kind': 'bobomb', 'vx': -1.0, 'state': 'walk', 'timer': 0})
            elif r < 0.58:
                # small staircase
                h = random.randint(2, 4)
                for step in range(h):
                    s = Sprite(self.canvas, [BRICK]); s.draw(0)
                    self.bricks.append({'s': s, 'wx': x + step * B, 'y': self.ground_y - (step + 1) * B})
                x += h * B
            elif r < 0.64:
                # floating coins
                for i in range(3):
                    co = Sprite(self.canvas, [COIN1, COIN2]); co.draw(0)
                    cy = self.ground_y - B * 2 - int(math.sin(i / 2 * math.pi) * B)
                    self.coins.append({'s': co, 'wx': x + i * B, 'wy': cy, 'got': False, 'ft': 0})
                x += 3 * B
            elif r < 0.72:
                # pipe
                ph = random.choice([2, 3]) * B
                pw = 2 * B
                py = self.ground_y - ph + B
                lip_h = B // 3
                self.pipes.append({
                    'wx': x, 'y': py, 'w': pw, 'h': ph,
                    'lip_h': lip_h, 'ids': [],
                })
                x += pw
            elif r < 0.78:
                # ground gap (pit)
                gap_w = random.randint(3, 4) * B
                self.gaps.append((x, x + gap_w))
                x += gap_w
            else:
                # empty ground – just the pre-spacing
                pass
        self.gen_x = max(self.gen_x, x)

    # ---- helpers ----
    def _all_solids(self):
        """Returns list of (wx, y, w, h) for all solid blocks + pipes."""
        B = self.BLK
        out = []
        for b in self.bricks:
            out.append((b['wx'], b['y'], B, B))
        for q in self.qblocks:
            out.append((q['wx'], q['y'], B, B))
        for p in self.pipes:
            out.append((p['wx'], p['y'], p['w'], p['h']))
        return out

    def _popup(self, sx, sy, text="+100"):
        self.popups.append(ScorePopup(self.canvas, sx, sy, text))

    def _add_score(self, pts, sx, sy):
        self.score += pts
        self.canvas.itemconfig(self.hud,
            text=f"ARROWS: Move  SPACE: Jump (hold=higher)  SHIFT: Run  ESC: Quit  SCORE: {self.score}")
        self._popup(sx, sy, f"+{pts}")

    # ---- AABB overlap ----
    @staticmethod
    def _overlap(ax, ay, aw, ah, bx, by, bw, bh):
        return ax + aw > bx and ax < bx + bw and ay + ah > by and ay < by + bh

    # ---- DEATH / RESPAWN ----
    def _take_hit(self):
        """Called when Mario touches an enemy. If big, shrink. If small, die."""
        if self.shrink_timer > 0 or self.invincible > 0:
            return  # already invincible
        if self.is_big:
            self.is_big = False
            self.shrink_timer = 60  # blink for ~2 seconds
            # Adjust position (big sprite is taller)
            self.my += self.SPR  # drop down since small sprite is shorter
            # Hide the big sprite
            self.big_mario.move_to(-200, -200)
        else:
            self._die()

    def _die(self):
        self.dead = True
        self.dead_timer = 40
        self.mvy = -14  # pop up
        self.mvx = 0
        self.is_big = False
        self.big_mario.move_to(-200, -200)

    def _respawn(self):
        self.dead = False
        # Respawn ahead of camera
        self.mwx = max(self.last_safe_wx, self.cam + 200)
        # Make sure we don't respawn in a gap
        while self._is_in_gap(self.mwx):
            self.mwx += self.BLK
        # Push away from any nearby enemies so we don't insta-die
        S = self.SPR
        for e in self.enemies:
            if abs(e['wx'] - self.mwx) < S * 2:
                self.mwx = e['wx'] + S * 3
        self.my = self.ground_y
        self.mvx = 0
        self.mvy = 0
        self.on_ground = True
        self.jumping = False
        self.invincible = 60  # brief invincibility after respawn
        self.stomp_grace = 0
        self.is_big = False
        self.shrink_timer = 0
        self.big_mario.move_to(-200, -200)

    # ---- MAIN UPDATE ----
    def update(self):
        # ---- DEATH ANIMATION ----
        if self.dead:
            self.dead_timer -= 1
            # Pop up then fall (original SMB death animation)
            if self.dead_timer > 25:
                self.my -= 8
            else:
                self.my += 10
            face_off = 0 if self.facing_right else 4
            self.mario.draw(3 + face_off)  # jump frame as death pose
            self.mario.move_to(self.mwx - self.cam, self.my)
            if self.dead_timer <= 0:
                self._respawn()
            return

        self.tick += 1
        B = self.BLK
        S = self.SPR

        # ---- invincibility / grace countdowns ----
        if self.stomp_grace > 0:
            self.stomp_grace -= 1
        if self.invincible > 0:
            self.invincible -= 1
        if self.shrink_timer > 0:
            self.shrink_timer -= 1

        # ---- clouds ----
        for c in self.clouds:
            c['x'] -= c['speed']
            self.canvas.move(c['tag'], -c['speed'], 0)
            if c['x'] + c['w'] < -20:
                self._reset_cloud(c)

        # ---- INPUT -> velocity ----
        # Run button (Shift) – from original SMB B-button run mechanic
        running = 'Shift_L' in self.keys or 'Shift_R' in self.keys
        accel = 1.2 if running else 0.8
        max_speed = 8.0 if running else 5.0
        friction = 0.6

        if 'Right' in self.keys:
            self.mvx = min(self.mvx + accel, max_speed)
            self.facing_right = True
        elif 'Left' in self.keys:
            self.mvx = max(self.mvx - accel, -max_speed)
            self.facing_right = False
        else:
            # friction
            if abs(self.mvx) < friction:
                self.mvx = 0
            elif self.mvx > 0:
                self.mvx -= friction
            else:
                self.mvx += friction

        # Variable-height jump (from original SMB JumpSwimTimer mechanic:
        # holding jump applies reduced gravity while ascending)
        if 'space' in self.keys and self.on_ground:
            self.mvy = -16
            self.on_ground = False
            self.jumping = True

        # Gravity – reduced while holding jump and still ascending
        if self.jumping and 'space' in self.keys and self.mvy < 0:
            self.mvy += 0.55   # light gravity – hold for higher jump
        else:
            self.mvy += 1.2
            if self.mvy >= 0:
                self.jumping = False
        if self.mvy > 18:
            self.mvy = 18

        # ---- move X, then check solids ----
        self.mwx += self.mvx
        solids = self._all_solids()
        for sx, sy, sw, sh in solids:
            if self._overlap(self.mwx, self.my, S, S, sx, sy, sw, sh):
                if self.mvx > 0:
                    self.mwx = sx - S
                elif self.mvx < 0:
                    self.mwx = sx + sw
                self.mvx = 0
                break

        # ---- move Y, then check solids ----
        self.my += self.mvy
        self.on_ground = False
        for sx, sy, sw, sh in solids:
            if self._overlap(self.mwx + 2, self.my, S - 4, S, sx, sy, sw, sh):
                if self.mvy > 0:
                    self.my = sy - S
                    self.mvy = 0
                    self.on_ground = True
                elif self.mvy < 0:
                    self.my = sy + sh
                    self.mvy = 1
                    # bump ?-blocks
                    for q in self.qblocks:
                        if not q['hit'] and q['wx'] == sx and q['y'] == sy:
                            q['hit'] = True
                            q['s'].draw(1)
                            scr_x = q['wx'] - self.cam
                            reward = q.get('reward', 'coin')
                            if reward == 'mushroom' and not self.is_big:
                                # Spawn a mushroom above the block
                                ms = Sprite(self.canvas, [MUSHROOM])
                                ms.draw(0)
                                self.mushrooms.append({
                                    's': ms, 'wx': q['wx'], 'wy': q['y'] - B,
                                    'vx': 2.0, 'vy': 0, 'active': True
                                })
                            else:
                                self._add_score(100, scr_x + B // 2, q['y'] - 20)
                break

        # ground floor (with gap check)
        if self.my >= self.ground_y:
            if self._is_in_gap(self.mwx):
                pass  # falling through gap!
            else:
                self.my = self.ground_y
                self.mvy = 0
                self.on_ground = True

        # Track last safe position
        if self.on_ground and not self._is_in_gap(self.mwx):
            self.last_safe_wx = self.mwx

        # Fell off screen = death
        if self.my > self.H + 100:
            self._die()
            return

        # don't go left past camera
        if self.mwx < self.cam:
            self.mwx = self.cam
            self.mvx = 0

        # ---- camera ----
        self.cam = self.mwx - self.W * 0.3

        # generate ahead
        edge = self.cam + self.W + 400
        if edge > self.gen_x:
            self._generate(self.gen_x, edge)

        # ---- draw Mario (small or big, direction-aware, blink) ----
        msx = self.mwx - self.cam
        face_off = 0 if self.facing_right else 4
        if self.is_big:
            # Big Mario (32px tall sprite)
            big_h = 32 * PX
            draw_y = self.my - (big_h - S)  # align feet
            if not self.on_ground:
                self.big_mario.draw(3 + face_off)
            elif abs(self.mvx) > 0.5:
                self.big_mario.draw(1 + (self.tick // 4 % 2) + face_off)
            else:
                self.big_mario.draw(0 + face_off)
            self.big_mario.move_to(msx, draw_y)
            # Hide small sprite
            self.mario.move_to(-200, -200)
            # Blink
            blink = (self.invincible > 0 or self.shrink_timer > 0) and self.tick % 4 < 2
            if blink:
                self.canvas.itemconfigure(self.big_mario.item, state='hidden')
            elif self.big_mario.item:
                self.canvas.itemconfigure(self.big_mario.item, state='normal')
        else:
            # Small Mario
            if not self.on_ground:
                self.mario.draw(3 + face_off)
            elif abs(self.mvx) > 0.5:
                self.mario.draw(1 + (self.tick // 4 % 2) + face_off)
            else:
                self.mario.draw(0 + face_off)
            self.mario.move_to(msx, self.my)
            # Hide big sprite
            self.big_mario.move_to(-200, -200)
            # Blink
            blink = (self.invincible > 0 or self.shrink_timer > 0) and self.tick % 4 < 2
            if blink:
                self.canvas.itemconfigure(self.mario.item, state='hidden')
            elif self.mario.item:
                self.canvas.itemconfigure(self.mario.item, state='normal')

        # ---- ground tiles (hide over gaps) ----
        tw = len(self.ground_tiles) * B
        for i, g in enumerate(self.ground_tiles):
            gx = (i * B) - (int(self.cam) % tw)
            if gx < -B: gx += tw
            gwx = self.cam + gx  # approximate world-x
            if self._is_in_gap(gwx, B):
                g.move_to(gx, self.H + 200)  # hide below screen
            else:
                g.move_to(gx, self.H - 48)

        # ---- bricks ----
        for b in self.bricks[:]:
            sx = b['wx'] - self.cam
            if sx < -B * 2:
                b['s'].destroy(); self.bricks.remove(b)
            elif sx < self.W + B:
                b['s'].move_to(sx, b['y'])

        # ---- ?-blocks ----
        for q in self.qblocks[:]:
            sx = q['wx'] - self.cam
            if sx < -B * 2:
                q['s'].destroy(); self.qblocks.remove(q)
            elif sx < self.W + B:
                q['s'].move_to(sx, q['y'])

        # ---- coins ----
        for c in self.coins[:]:
            sx = c['wx'] - self.cam
            if sx < -B * 2:
                c['s'].destroy(); self.coins.remove(c); continue
            if not c['got']:
                c['s'].draw((self.tick // 8) % 2)
                c['s'].move_to(sx, c['wy'])
                if self._overlap(self.mwx, self.my, S, S, c['wx'], c['wy'], 8 * PX, S):
                    c['got'] = True; c['ft'] = 15
                    self._add_score(100, sx, c['wy'])
            else:
                c['ft'] -= 1; c['wy'] -= 5
                c['s'].move_to(sx, c['wy'])
                if c['ft'] <= 0:
                    c['s'].destroy(); self.coins.remove(c)

        # ---- mushrooms (physics + collection) ----
        for m in self.mushrooms[:]:
            sx = m['wx'] - self.cam
            if sx < -B * 2:
                m['s'].destroy(); self.mushrooms.remove(m); continue
            if m['active']:
                # Gravity
                m['vy'] += 1.0
                if m['vy'] > 10:
                    m['vy'] = 10
                m['wx'] += m['vx']
                m['wy'] += m['vy']
                # Ground collision
                if m['wy'] >= self.ground_y:
                    m['wy'] = self.ground_y
                    m['vy'] = 0
                # Bounce off solids
                for bwx, by, bw, bh in solids:
                    if self._overlap(m['wx'], m['wy'], S, S, bwx, by, bw, bh):
                        if m['vx'] > 0:
                            m['wx'] = bwx - S
                        else:
                            m['wx'] = bwx + bw
                        m['vx'] *= -1
                        break
                m['s'].move_to(sx, m['wy'])
                # Mario collects mushroom
                if self._overlap(self.mwx, self.my, S, S, m['wx'], m['wy'], S, S):
                    m['active'] = False
                    m['s'].destroy()
                    self.mushrooms.remove(m)
                    if not self.is_big:
                        self.is_big = True
                        self._add_score(1000, sx, m['wy'] - 20)
                    continue

        # ---- pipes (canvas rectangles) ----
        for p in self.pipes[:]:
            sx = p['wx'] - self.cam
            if sx < -p['w'] * 2:
                for cid in p['ids']:
                    self.canvas.delete(cid)
                self.pipes.remove(p)
                continue
            if sx > self.W + p['w']:
                for cid in p['ids']:
                    self.canvas.coords(cid, -2000, -2000, -1999, -1999)
                continue
            lip_extra = B // 4
            if not p['ids']:
                # Create pipe canvas items: lip + body + highlight
                lip = self.canvas.create_rectangle(
                    sx - lip_extra, p['y'],
                    sx + p['w'] + lip_extra, p['y'] + p['lip_h'],
                    fill='#20A010', outline='#00680C', width=2)
                body = self.canvas.create_rectangle(
                    sx, p['y'] + p['lip_h'],
                    sx + p['w'], p['y'] + p['h'],
                    fill='#20A010', outline='#00680C', width=2)
                hl = self.canvas.create_rectangle(
                    sx + B // 3, p['y'],
                    sx + B // 3 + 4, p['y'] + p['h'],
                    fill='#80E080', outline='')
                p['ids'] = [lip, body, hl]
            else:
                lip, body, hl = p['ids']
                self.canvas.coords(lip,
                    sx - lip_extra, p['y'],
                    sx + p['w'] + lip_extra, p['y'] + p['lip_h'])
                self.canvas.coords(body,
                    sx, p['y'] + p['lip_h'],
                    sx + p['w'], p['y'] + p['h'])
                self.canvas.coords(hl,
                    sx + B // 3, p['y'],
                    sx + B // 3 + 4, p['y'] + p['h'])

        # ---- enemies (with authentic SMB shell mechanics) ----
        mario_l, mario_t = self.mwx, self.my
        mario_r, mario_b = self.mwx + S, self.my + S
        for e in self.enemies[:]:
            sx = e['wx'] - self.cam
            if sx < -B * 3 or sx > self.W + B * 3:
                e['s'].destroy(); self.enemies.remove(e); continue

            if e['state'] == 'walk':
                e['wx'] += e['vx']
                # Red koopas turn at edges (from original RedKoopa AI)
                if e['kind'] == 'red_koopa':
                    ahead_x = e['wx'] + (S if e['vx'] > 0 else -4)
                    on_solid = False
                    if not self._is_in_gap(ahead_x, 4):
                        if e['wy'] >= self.ground_y:
                            on_solid = True
                        else:
                            for bwx, by, bw, bh in solids:
                                if self._overlap(ahead_x, e['wy'] + S, 4, 4, bwx, by, bw, bh):
                                    on_solid = True
                                    break
                    if not on_solid:
                        e['vx'] *= -1
                # Animation (red koopa uses right-facing frames when going right)
                if e['kind'] == 'red_koopa' and e['vx'] > 0:
                    e['s'].draw(3 + (self.tick // 6) % 2)
                else:
                    e['s'].draw((self.tick // 6) % 2)
                e['s'].move_to(sx, e['wy'])

                # Stomp / contact check (skip during invincibility or stomp grace)
                if self.invincible <= 0 and self.stomp_grace <= 0 and self.shrink_timer <= 0 and self._overlap(mario_l + 4, mario_t, S - 8, S, e['wx'], e['wy'], S, S):
                    if self.mvy > 0 and mario_b < e['wy'] + S * 0.6:
                        # Stomped from above
                        if e['kind'] == 'goomba':
                            e['state'] = 'flat'; e['s'].draw(2); e['timer'] = 18
                        elif e['kind'] == 'bobomb':
                            # Stomp bob-omb: starts fuse countdown
                            e['state'] = 'fuse'; e['timer'] = 90; e['vx'] = 0
                        else:
                            # Koopa → still shell
                            e['state'] = 'shell_still'
                            e['s'].draw(2)
                            e['vx'] = 0
                            e['timer'] = 300
                        self.mvy = -10
                        self.stomp_grace = 25
                        self._add_score(200, sx, e['wy'] - 20)
                    else:
                        # Side/bottom contact = take hit
                        self._take_hit()
                        return

            elif e['state'] == 'shell_still':
                e['s'].move_to(sx, e['wy'])
                e['timer'] -= 1
                # Koopa revives if shell left alone (original mechanic)
                if e['timer'] <= 0:
                    e['state'] = 'walk'
                    e['vx'] = -1.5
                    e['timer'] = 0
                    continue
                # Kick the still shell on contact (skip during grace/invincibility)
                if self.invincible <= 0 and self.stomp_grace <= 0 and self._overlap(mario_l + 2, mario_t + 4, S - 4, S - 8, e['wx'], e['wy'], S, S):
                    # Kick direction based on Mario's side
                    kick_dir = 10 if self.mwx + S / 2 < e['wx'] + S / 2 else -10
                    e['state'] = 'shell'
                    e['vx'] = kick_dir
                    # NO bounce – original SMB doesn't bounce on shell kick
                    # Just nudge Mario back so he doesn't ride the shell
                    if self.mwx + S / 2 < e['wx'] + S / 2:
                        self.mwx = e['wx'] - S - 2
                    else:
                        self.mwx = e['wx'] + S + 2
                    self.stomp_grace = 15
                    self._add_score(100, sx, e['wy'] - 20)

            elif e['state'] == 'fuse':
                # Bob-omb fuse lit – blinks then explodes
                e['timer'] -= 1
                e['s'].draw((self.tick // 3) % 2)  # fast blink
                e['s'].move_to(sx, e['wy'])
                if e['timer'] <= 0:
                    # EXPLODE – kill nearby enemies
                    e['s'].draw(2)  # explosion frame
                    e['state'] = 'flat'; e['timer'] = 20
                    blast_r = S * 3
                    for other in self.enemies[:]:
                        if other is e or other['state'] == 'flat':
                            continue
                        if abs(other['wx'] - e['wx']) < blast_r and abs(other['wy'] - e['wy']) < blast_r:
                            other['state'] = 'flat'
                            other['s'].draw(min(2, len(other['s'].frames) - 1))
                            other['timer'] = 18
                            self._add_score(200, other['wx'] - self.cam, other['wy'] - 20)
                    # Hurt Mario if in blast radius
                    if self.invincible <= 0 and self.shrink_timer <= 0:
                        if abs(self.mwx - e['wx']) < blast_r and abs(self.my - e['wy']) < blast_r:
                            self._take_hit()
                            return

            elif e['state'] == 'flat':
                e['s'].move_to(sx, e['wy'])
                e['timer'] -= 1
                if e['timer'] <= 0:
                    e['s'].destroy(); self.enemies.remove(e); continue

            elif e['state'] == 'shell':
                e['wx'] += e['vx']
                e['s'].move_to(sx, e['wy'])
                # Shell kills other enemies (chain kills from original ShellOrBlockDefeat)
                for other in self.enemies[:]:
                    if other is e or other['state'] in ('flat', 'dead'):
                        continue
                    if self._overlap(e['wx'], e['wy'], S, S, other['wx'], other['wy'], S, S):
                        if other['state'] == 'shell_still':
                            other['state'] = 'flat'; other['timer'] = 18
                        else:
                            other['state'] = 'flat'
                            other['s'].draw(2 if other['kind'] == 'goomba' else 2)
                            other['timer'] = 18
                        self._add_score(100, other['wx'] - self.cam, other['wy'] - 20)
                # Shell bounces off walls
                for bwx, by, bw, bh in solids:
                    if self._overlap(e['wx'], e['wy'], S, S, bwx, by, bw, bh):
                        e['vx'] *= -1
                        break
                # Moving shell can hurt Mario too
                if self.invincible <= 0 and self.stomp_grace <= 0 and self._overlap(mario_l + 4, mario_t, S - 8, S, e['wx'], e['wy'], S, S):
                    if self.mvy > 0 and mario_b < e['wy'] + S * 0.5:
                        # Stomp moving shell → stops it
                        e['state'] = 'shell_still'
                        e['vx'] = 0
                        e['timer'] = 300
                        self.mvy = -10
                        self.stomp_grace = 25
                        self._add_score(100, sx, e['wy'] - 20)
                    else:
                        self._take_hit()
                        return
                # shell disappears after going far off screen
                if abs(e['wx'] - self.mwx) > self.W * 2:
                    e['s'].destroy(); self.enemies.remove(e); continue

        # ---- popups ----
        for p in self.popups[:]:
            if not p.update():
                self.popups.remove(p)


# ============================================================
#  GLOBAL HOTKEY (Ctrl+Alt+M) using Windows RegisterHotKey
# ============================================================
HOTKEY_ID = 1
MOD_CTRL_ALT = 0x0001 | 0x0002  # MOD_ALT | MOD_CONTROL

def _hotkey_listener(callback):
    """Run in a thread – blocks on GetMessage, fires callback on hotkey."""
    user32 = ctypes.windll.user32
    # Unregister any stale hotkey from a previous crashed instance
    user32.UnregisterHotKey(None, HOTKEY_ID)
    if not user32.RegisterHotKey(None, HOTKEY_ID, MOD_CTRL_ALT, 0x4D):  # 0x4D = 'M'
        print("Ctrl+Alt+M taken, trying Ctrl+Alt+Shift+M...")
        MOD_FALLBACK = 0x0001 | 0x0002 | 0x0004  # ALT | CONTROL | SHIFT
        if not user32.RegisterHotKey(None, HOTKEY_ID, MOD_FALLBACK, 0x4D):
            print("Could not register hotkey. Kill other instances first.")
            return
        print("Registered: Ctrl+Alt+Shift+M")
    msg = wt.MSG()
    while user32.GetMessageW(ctypes.byref(msg), None, 0, 0) != 0:
        if msg.message == 0x0312:  # WM_HOTKEY
            callback()
    user32.UnregisterHotKey(None, HOTKEY_ID)


# ============================================================
#  SYSTEM TRAY ICON
# ============================================================
def _create_tray_icon(on_show, on_quit):
    """Create a system tray icon. Runs in its own thread."""
    try:
        import pystray
        from PIL import Image, ImageDraw
    except ImportError:
        return  # silently skip if deps not bundled

    # Draw a tiny Mario-hat icon (16x16)
    img = Image.new('RGBA', (64, 64), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    # Red hat
    d.rectangle([16, 4, 48, 20], fill='#E40818')
    d.rectangle([12, 12, 52, 24], fill='#E40818')
    # Face
    d.rectangle([18, 24, 46, 44], fill='#FCA044')
    # Eyes
    d.rectangle([22, 28, 28, 34], fill='#000000')
    d.rectangle([36, 28, 42, 34], fill='#000000')
    # Body
    d.rectangle([16, 44, 48, 58], fill='#E40818')
    # Feet
    d.rectangle([12, 54, 24, 62], fill='#AC7C00')
    d.rectangle([40, 54, 52, 62], fill='#AC7C00')

    menu = pystray.Menu(
        pystray.MenuItem("Show/Hide (Ctrl+Alt+M)", lambda: on_show()),
        pystray.MenuItem("☕ Donate (PayPal)", lambda: _open_donate()),
        pystray.Menu.SEPARATOR,
        pystray.MenuItem("Quit", lambda: on_quit()),
    )
    icon = pystray.Icon("DesktopMario", img, "Desktop Mario", menu)
    icon.run()


def _open_donate():
    """Open PayPal donate link in default browser."""
    import webbrowser
    webbrowser.open("https://www.paypal.com/paypalme/bhanu1001")


# ============================================================
#  APP
# ============================================================
class App:
    def __init__(self):
        self.root = None
        self.game = None
        self.visible = False
        self._quitting = False

    def toggle(self):
        """Called from hotkey thread – schedule on tkinter thread."""
        if self.root and not self._quitting:
            self.root.after_idle(self._toggle_impl)

    def _toggle_impl(self):
        if self.visible:
            self.root.withdraw()
            self.visible = False
        else:
            self.root.deiconify()
            self.root.focus_force()
            self.visible = True

    def quit_app(self):
        """Called from tray icon – safely shutdown."""
        self._quitting = True
        if self.root:
            self.root.after_idle(self.root.destroy)

    def run(self):
        try:
            ctypes.windll.shcore.SetProcessDpiAwareness(2)
        except Exception:
            pass

        self.root = tk.Tk()
        root = self.root
        root.attributes("-fullscreen", True)
        root.overrideredirect(True)
        root.wm_attributes("-topmost", True)
        root.wm_attributes("-transparentcolor", "black")
        root.update_idletasks()

        W = root.winfo_screenwidth()
        H = root.winfo_screenheight()

        canvas = tk.Canvas(root, width=W, height=H, bg='black', highlightthickness=0)
        canvas.pack()

        # prevent PhotoImage garbage collection
        root._img_refs = _img_cache

        self.game = Game(canvas, W, H)

        # start hidden – press Ctrl+Alt+M to show
        root.withdraw()
        self.visible = False

        # key bindings – ESC hides instead of quitting
        root.bind("<KeyPress>",   lambda e: self.game.keys.add(e.keysym))
        root.bind("<KeyRelease>", lambda e: self.game.keys.discard(e.keysym))
        root.bind("<Escape>",     lambda e: self._toggle_impl())

        # start global hotkey listener
        t = threading.Thread(target=_hotkey_listener, args=(self.toggle,), daemon=True)
        t.start()

        # start system tray icon
        tray_t = threading.Thread(
            target=_create_tray_icon,
            args=(self.toggle, self.quit_app),
            daemon=True
        )
        tray_t.start()

        def animate():
            if self._quitting:
                return
            self.game.update()
            root.after(33, animate)

        animate()
        root.mainloop()


if __name__ == "__main__":
    App().run()
