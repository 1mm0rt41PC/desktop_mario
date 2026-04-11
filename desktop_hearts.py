import tkinter as tk
import random
import ctypes

class FallingHeart:
    def __init__(self, canvas, width, height):
        self.canvas = canvas
        self.width = width
        self.height = height
        
        # Randomize attributes
        self.x = random.randint(0, width)
        self.y = random.randint(-1000, -50)
        self.speed_y = random.uniform(2, 6)
        self.speed_x = random.uniform(-0.5, 0.5)
        self.size = random.randint(20, 80)
        
        colors = ['#ff3366', '#ff6699', '#cc0033', '#ff99aa']
        chars = ['❤', '💖', '💗', '💕', '💘']
        
        # Create heart text
        self.id = canvas.create_text(
            self.x, self.y, 
            text=random.choice(chars), 
            fill=random.choice(colors), 
            font=('Segoe UI Emoji', self.size)
        )
        
    def fall(self):
        # Move heart
        self.canvas.move(self.id, self.speed_x, self.speed_y)
        pos = self.canvas.coords(self.id)
        
        # If it goes off screen, reset it back to the top
        if len(pos) >= 2 and pos[1] > self.height + 100:
            self.canvas.coords(self.id, random.randint(0, self.width), random.randint(-200, -50))


def main():
    # Attempt to fix DPI scaling issues on Windows
    try:
        import ctypes
        ctypes.windll.shcore.SetProcessDpiAwareness(2) # PROCESS_PER_MONITOR_DPI_AWARE
    except Exception:
        pass

    root = tk.Tk()
    
    # Needs to be fullscreen and borderless
    root.attributes("-fullscreen", True)
    root.overrideredirect(True) 
    
    # Make window always on top and transparent
    root.wm_attributes("-topmost", True)
    
    # Any black pixels will become fully transparent on Windows
    root.wm_attributes("-transparentcolor", "black")

    # Update to get the actual full screen dimensions
    root.update_idletasks()
    screen_width = root.winfo_width()
    screen_height = root.winfo_height()

    canvas = tk.Canvas(root, width=screen_width, height=screen_height, bg='black', highlightthickness=0)
    canvas.pack()

    # Create falling hearts
    num_hearts = 60
    hearts = [FallingHeart(canvas, screen_width, screen_height) for _ in range(num_hearts)]

    # Animation loop
    def animate():
        for heart in hearts:
            heart.fall()
        root.after(20, animate) # 20 ms -> ~50 FPS

    # Press Escape to close it
    root.bind("<Escape>", lambda e: root.destroy())
    
    # Add a small invisible text label to tell users how to exit
    canvas.create_text(
        screen_width // 2, 50, 
        text="Press ESC to close the hearts overlay", 
        fill="white", 
        font=('Arial', 12, 'bold')
    )

    animate()
    root.mainloop()

if __name__ == "__main__":
    main()
