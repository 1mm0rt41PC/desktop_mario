import tkinter as tk
import ctypes

# Simple pixel art frames for Mario (Standing and Running modes)
# Colors: _=Transparent, R=Red, B=Brown/Hair/Shoes, S=Skin, O=Blue/Overalls, Y=Yellow
MARIO_COLORS = {
    '_': '',
    'R': '#ff0000',
    'B': '#663300',
    'S': '#ffcc99',
    'O': '#0000ff',
    'Y': '#ffff00'
}

# Frame 1: Standing / Running step 1
FRAME_1 = [
    "____RRRRR___",
    "___RRRRRRRRR",
    "___BBBSSBSS_",
    "__BBBSBSSSSS",
    "__BBBSBSSBSS",
    "__BBBBSSSS__",
    "____RRRRRRR_",
    "___RRROORRR_",
    "__RRROOORRRR",
    "__RROOYOORRR",
    "____OOOOOO__",
    "___OO_OO_OO_",
    "__OO__OO__OO",
    "______RR____",
    "_____BBBB___"
]

# Frame 2: Running step 2
FRAME_2 = [
    "____RRRRR___",
    "___RRRRRRRRR",
    "___BBBSSBSS_",
    "__BBBSBSSSSS",
    "__BBBSBSSBSS",
    "__BBBBSSSS__",
    "____RRRRRRR_",
    "___RRROORRR_",
    "__RRRRROORRR",
    "__RRROOYOORR",
    "__R_OOOOOO__",
    "__RROO_OO___",
    "__ROO__OO___",
    "___OO_RR____",
    "___BBBBBB___"
]

class Mario:
    def __init__(self, canvas, screen_width, screen_height):
        self.canvas = canvas
        self.screen_width = screen_width
        self.screen_height = screen_height
        
        self.pixel_size = 4  # Scale of the pixel art
        self.width = len(FRAME_1[0]) * self.pixel_size
        self.height = len(FRAME_1) * self.pixel_size
        
        # Start at the bottom left (sneaking in)
        self.x = -self.width
        # Assuming taskbar is usually ~40-50 pixels high on Windows 11 at 100% scaling
        self.y = screen_height - 48 - self.height 
        
        self.speed = 3
        self.frame_index = 0
        self.frames = [FRAME_1, FRAME_2]
        self.pixels = [] # Store rect IDs to delete and redraw or move
        self.drawn = False
        
    def draw_frame(self, frame_data):
        # Clear previous frame
        for p in self.pixels:
            self.canvas.delete(p)
        self.pixels.clear()
        
        for row_idx, row in enumerate(frame_data):
            for col_idx, color_code in enumerate(row):
                if color_code != '_':
                    px_x = self.x + (col_idx * self.pixel_size)
                    px_y = self.y + (row_idx * self.pixel_size)
                    rect_id = self.canvas.create_rectangle(
                        px_x, px_y, 
                        px_x + self.pixel_size, px_y + self.pixel_size, 
                        fill=MARIO_COLORS[color_code], outline=MARIO_COLORS[color_code]
                    )
                    self.pixels.append(rect_id)
                    
    def update(self):
        # Move right
        self.x += self.speed
        
        # Wrap around screen
        if self.x > self.screen_width:
            self.x = -self.width
            
        # Swap frames every 5 updates for animation
        self.frame_index += 1
        current_frame = self.frames[(self.frame_index // 5) % len(self.frames)]
        
        self.draw_frame(current_frame)


def main():
    try:
        import ctypes
        ctypes.windll.shcore.SetProcessDpiAwareness(2) # PROCESS_PER_MONITOR_DPI_AWARE
    except Exception:
        pass

    root = tk.Tk()
    
    root.attributes("-fullscreen", True)
    root.overrideredirect(True) 
    root.wm_attributes("-topmost", True)
    root.wm_attributes("-transparentcolor", "black")

    root.update_idletasks()
    screen_width = root.winfo_width()
    screen_height = root.winfo_height()

    canvas = tk.Canvas(root, width=screen_width, height=screen_height, bg='black', highlightthickness=0)
    canvas.pack()

    # Create Mario character
    mario = Mario(canvas, screen_width, screen_height)

    def animate():
        mario.update()
        root.after(30, animate)

    # Exit on escape
    root.bind("<Escape>", lambda e: root.destroy())

    animate()
    root.mainloop()

if __name__ == "__main__":
    main()