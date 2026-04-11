import tkinter as tk
import random
import ctypes
import time
from ctypes import wintypes

# Setup structure for Windows GetLastInputInfo
class LASTINPUTINFO(ctypes.Structure):
    _fields_ = [
        ('cbSize', wintypes.UINT),
        ('dwTime', wintypes.DWORD),
    ]

def get_idle_time():
    lastInputInfo = LASTINPUTINFO()
    lastInputInfo.cbSize = ctypes.sizeof(lastInputInfo)
    ctypes.windll.user32.GetLastInputInfo(ctypes.byref(lastInputInfo))
    millis = ctypes.windll.kernel32.GetTickCount() - lastInputInfo.dwTime
    return millis / 1000.0

class MystifyLines:
    def __init__(self, canvas, width, height):
        self.canvas = canvas
        self.width = width
        self.height = height
        self.points = []
        self.velocities = []
        self.color = '#%02x%02x%02x' % (random.randint(50, 255), random.randint(50, 255), random.randint(50, 255))
        
        # Create 4 points for the polygon
        for _ in range(4):
            self.points.append([random.randint(0, width), random.randint(0, height)])
            self.velocities.append([random.choice([-4, -3, 3, 4]), random.choice([-4, -3, 3, 4])])
            
        self.line_id = canvas.create_polygon(
            self.points[0][0], self.points[0][1],
            self.points[1][0], self.points[1][1],
            self.points[2][0], self.points[2][1],
            self.points[3][0], self.points[3][1],
            outline=self.color, fill='', width=2
        )
        
    def move(self):
        for i in range(4):
            # Update position
            self.points[i][0] += self.velocities[i][0]
            self.points[i][1] += self.velocities[i][1]
            
            # Bounce off walls
            if self.points[i][0] <= 0 or self.points[i][0] >= self.width:
                self.velocities[i][0] *= -1
            if self.points[i][1] <= 0 or self.points[i][1] >= self.height:
                self.velocities[i][1] *= -1
                
        # Update canvas polygon
        coords = []
        for p in self.points:
            coords.extend(p)
        self.canvas.coords(self.line_id, *coords)
        
        # Color shifting
        if random.random() < 0.05:
            self.color = '#%02x%02x%02x' % (random.randint(50, 255), random.randint(50, 255), random.randint(50, 255))
            self.canvas.itemconfig(self.line_id, outline=self.color)

class FallingHeart:
    def __init__(self, canvas, width, height):
        self.canvas = canvas
        self.width = width
        self.height = height
        
        self.x = random.randint(0, width)
        self.y = random.randint(-1000, -50)
        self.speed_y = random.uniform(2, 6)
        self.speed_x = random.uniform(-0.5, 0.5)
        self.size = random.randint(20, 80)
        
        colors = ['#ff3366', '#ff6699', '#cc0033', '#ff99aa']
        chars = ['❤', '💖', '💗', '💕', '💘']
        
        self.id = canvas.create_text(
            self.x, self.y, 
            text=random.choice(chars), 
            fill=random.choice(colors), 
            font=('Segoe UI Emoji', self.size)
        )
        
    def fall(self):
        self.canvas.move(self.id, self.speed_x, self.speed_y)
        pos = self.canvas.coords(self.id)
        
        if len(pos) >= 2 and pos[1] > self.height + 100:
            self.canvas.coords(self.id, random.randint(0, self.width), random.randint(-200, -50))

class ScreensaverApp:
    def __init__(self):
        self.root = None
        self.IDLE_TIMEOUT = 60 # 60 seconds
        self.running_screensaver = False
        
        # Try DPI awareness
        try:
            ctypes.windll.shcore.SetProcessDpiAwareness(2)
        except Exception:
            pass
            
    def start_screensaver(self):
        if self.running_screensaver or get_idle_time() < 0.5:
            return
            
        self.running_screensaver = True
        self.root = tk.Tk()
        self.root.attributes("-fullscreen", True)
        self.root.overrideredirect(True) 
        self.root.wm_attributes("-topmost", True)
        self.root.wm_attributes("-transparentcolor", "black")
        
        # Bind any key/mouse movement to exit screensaver
        self.root.bind("<Key>", self.stop_screensaver)
        self.root.bind("<Motion>", self.stop_screensaver)
        self.root.bind("<Button>", self.stop_screensaver)
        
        self.root.update_idletasks()
        screen_width = self.root.winfo_screenwidth()
        screen_height = self.root.winfo_screenheight()
        
        self.canvas = tk.Canvas(self.root, width=screen_width, height=screen_height, bg='black', highlightthickness=0)
        self.canvas.pack()
        
        # Elements
        self.mystify_shapes = [MystifyLines(self.canvas, screen_width, screen_height) for _ in range(5)]
        self.hearts = [FallingHeart(self.canvas, screen_width, screen_height) for _ in range(50)]
        
        self.animate()
        
    def stop_screensaver(self, event=None):
        if self.running_screensaver and self.root:
            self.running_screensaver = False
            self.root.destroy()
            self.root = None
            
    def animate(self):
        if not self.running_screensaver:
            return
            
        # If user moved mouse or typed right before/after
        if get_idle_time() < 1:
            self.stop_screensaver()
            return
            
        for shape in self.mystify_shapes:
            shape.move()
            
        for heart in self.hearts:
            heart.fall()
            
        self.root.after(20, self.animate)

    def monitor_idle(self):
        idle = get_idle_time()
        print(f"Idle time: {idle:.1f}s")
        
        if idle >= self.IDLE_TIMEOUT and not self.running_screensaver:
            self.start_screensaver()
            
        # We need another loop to keep checking idle time
        if self.root and self.running_screensaver:
            # tkinter mainloop is handling things 
            pass
        else:
            # We are not in tkinter mainloop, sleep and loop
            time.sleep(1)

def main():
    print("Screensaver monitor running. Waiting for 1 minute of idle time...")
    app = ScreensaverApp()
    while True:
        try:
            app.monitor_idle()
            if app.root and app.running_screensaver:
                app.root.mainloop()
        except KeyboardInterrupt:
            break

if __name__ == "__main__":
    main()