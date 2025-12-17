import tkinter as tk
from math import cos, sin, pi, hypot

class HamiltonGraphGame:
    def __init__(self, master):
        self.master = master
        self.master.title("Гамильтонова игра: додекаэдр (реальный граф)")
        self.vertex_radius = 14
        self.padding = 40
        self.canvas_size = 700

        self.canvas = tk.Canvas(master, width=self.canvas_size, height=self.canvas_size, bg="white")
        self.canvas.pack()

        self.info = tk.Label(
            master,
            text="Нажмите вершину в порядке попытки построения гамильтонова цикла на додекаэдре",
            anchor="w"
        )
        self.info.pack(fill="x")

        self.reset_button = tk.Button(master, text="Сброс", command=self.reset_game)
        self.reset_button.pack()

        self.init_graph()
        self.reset_game_state()
        self.draw_graph()

        self.master.bind("<Key>", self.on_key)
        self.bind_vertex_clicks()

    def init_graph(self):
        # Точные координаты вершин додекаэдра (плоская проекция)
        cx = cy = self.canvas_size / 2
        
        # Радиусы двух колец для классической проекции додекаэдра
        r1 = 220  # внешнее кольцо
        r2 = 140  # среднее кольцо  
        r3 = 80   # внутреннее кольцо
        
        self.vertices = [
            # Вершины 1-5: внешнее кольцо (углы каждые 72°)
            (cx + r1 * cos(2*pi*0/5 - pi/2), cy + r1 * sin(2*pi*0/5 - pi/2)),  # 0
            (cx + r1 * cos(2*pi*1/5 - pi/2), cy + r1 * sin(2*pi*1/5 - pi/2)),  # 1
            (cx + r1 * cos(2*pi*2/5 - pi/2), cy + r1 * sin(2*pi*2/5 - pi/2)),  # 2
            (cx + r1 * cos(2*pi*3/5 - pi/2), cy + r1 * sin(2*pi*3/5 - pi/2)),  # 3
            (cx + r1 * cos(2*pi*4/5 - pi/2), cy + r1 * sin(2*pi*4/5 - pi/2)),  # 4
            
            # Вершины 6-10: теперь ВНУТРЕННЕЕ кольцо (используем r3, но оставляем смещение формы)
            (cx + r2 * cos(2*pi*0/5 - pi/2), cy + r2 * sin(2*pi*0/5 - pi/2)),  # 5
            (cx + r2 * cos(2*pi*3/5 + pi/2), cy + r2 * sin(2*pi*3/5 + pi/2)),  # 6
            (cx + r2 * cos(2*pi*1/5 - pi/2), cy + r2 * sin(2*pi*1/5 - pi/2)),  # 7
            (cx + r2 * cos(2*pi*4/5 + pi/2), cy + r2 * sin(2*pi*4/5 + pi/2)),  # 8
            (cx + r2 * cos(2*pi*2/5 - pi/2), cy + r2 * sin(2*pi*2/5 - pi/2)),  # 9
            
            # Вершины 11-15: теперь СРЕДНЕЕ кольцо (используем r2, форма как у прежнего внутреннего)
            (cx + r2 * cos(2*pi*0/5 + pi/2), cy + r2 * sin(2*pi*0/5 + pi/2)),  # 10
            (cx + r2 * cos(2*pi*3/5 - pi/2), cy + r2 * sin(2*pi*3/5 - pi/2)),  # 11
            (cx + r2 * cos(2*pi*1/5 + pi/2), cy + r2 * sin(2*pi*1/5 + pi/2)),  # 12
            (cx + r2 * cos(2*pi*4/5 - pi/2), cy + r2 * sin(2*pi*4/5 - pi/2)),  # 13
            (cx + r2 * cos(2*pi*2/5 + pi/2), cy + r2 * sin(2*pi*2/5 + pi/2)),  # 14
            
            # Вершины 16-20: нижнее кольцо оставляем как было (на r2)
            (cx + r3 * cos(2*pi*0/5 - pi/2 + pi/5), cy + r3 * sin(2*pi*0/5 - pi/2 + pi/5)),  # 15
            (cx + r3 * cos(2*pi*4/5 - pi/2 + pi/5), cy + r3 * sin(2*pi*4/5 - pi/2 + pi/5)),  # 16
            (cx + r3 * cos(2*pi*2/5 - pi/2 + pi/5), cy + r3 * sin(2*pi*2/5 - pi/2 + pi/5)),  # 17
            (cx + r3 * cos(2*pi*1/5 - pi/2 + pi/5), cy + r3 * sin(2*pi*1/5 - pi/2 + pi/5)),  # 18
            (cx + r3 * cos(2*pi*3/5 - pi/2 + pi/5), cy + r3 * sin(2*pi*3/5 - pi/2 + pi/5))   # 19
        ]

        # дальше — без изменений
        edge_list = [
            (1,2),(1,5),(1,6), (2,8),(2,3), (3,10),(3,4), (4,12),(4,5),
            (5,14), (6,7),(6,15), (7,8),(7,17), (8,9), (9,18),(9,10),
            (10,11),(11,19),(11,12), (12,13), (13,20),(13,14), (14,15),
            (15,16),(16,17),(17,18),(18,19),(19,20),(20,16)
        ]
        
        self.edges = set()
        for a, b in edge_list:
            self.edges.add(tuple(sorted((a-1, b-1))))  # 1-based → 0-based
        
        self.vertex_ids = list(range(20))
        self.n = 20
        self.visited_order = []
        self.selected = None


    def draw_graph(self):
        self.canvas.delete("all")
        
        # Рёбра (серые)
        for a, b in self.edges:
            ax, ay = self.vertices[a]
            bx, by = self.vertices[b]
            self.canvas.create_line(ax, ay, bx, by, fill="#ccc", width=2)
        
        # Путь игрока (синим)
        if len(self.visited_order) >= 2:
            for i in range(len(self.visited_order)-1):
                a = self.visited_order[i]
                b = self.visited_order[i+1]
                ax, ay = self.vertices[a]
                bx, by = self.vertices[b]
                self.canvas.create_line(ax, ay, bx, by, fill="blue", width=4)
        
        # Замыкание цикла (зелёное)
        if len(self.visited_order) == self.n:
            first = self.visited_order[0]
            last = self.visited_order[-1]
            edge = (min(first, last), max(first, last))
            if edge in self.edges:
                ax, ay = self.vertices[last]
                bx, by = self.vertices[first]
                self.canvas.create_line(ax, ay, bx, by, fill="green", width=4, dash=(4, 2))
        
        # Вершины с цветовой индикацией
        for idx, (x, y) in enumerate(self.vertices):
            if not self.visited_order:
                color = "white"
            elif idx == self.visited_order[0]:
                color = "#ffe4b5"  # стартовая
            elif idx in self.visited_order:
                color = "lightblue"
            else:
                color = "white"
            
            self.canvas.create_oval(
                x - self.vertex_radius, y - self.vertex_radius,
                x + self.vertex_radius, y + self.vertex_radius,
                fill=color, outline="black", width=2
            )
            self.canvas.create_text(x, y, text=str(idx+1), font=("Arial", 10, "bold"))

    def reset_game_state(self):
        self.visited_order = []
        self.selected = None
        self.update_info()

    def reset_game(self):
        self.reset_game_state()
        self.draw_graph()

    def update_info(self):
        path_text = ' -> '.join(map(str, [v+1 for v in self.visited_order])) if self.visited_order else "(пусто)"
        self.info.config(
            text=f"Путь: {path_text} | посещено: {len(self.visited_order)} из {self.n} "
                 f"(клик мышью по вершинам, R — сброс, цифры 0-9)"
        )

    def on_click_vertex(self, idx):
        if len(self.visited_order) == self.n:
            return

        if not self.visited_order:
            self.visited_order.append(idx)
        else:
            if idx in self.visited_order:
                return
            prev = self.visited_order[-1]
            edge = (min(prev, idx), max(prev, idx))
            if edge in self.edges:
                self.visited_order.append(idx)
            else:
                self.info.config(text="❌ Нет ребра! Выберите соседнюю вершину.", fg="red")
                return

        self.draw_graph()
        self.update_info()
        self.check_completion()

    def check_completion(self):
        if len(self.visited_order) == self.n:
            first = self.visited_order[0]
            last = self.visited_order[-1]
            edge = (min(first, last), max(first, last))
            if edge in self.edges:
                self.info.config(text="🎉 ГАМИЛЬТОНОВ ЦИКЛ НА ДОДЕКАЭДРЕ ПОСТРОЕН! 🎉", fg="green")
            else:
                self.info.config(text="Все вершины посещены, но цикл не замкнут.", fg="orange")
            self.draw_graph()
            return True
        return False

    def on_key(self, event):
        if event.char.lower() == 'r':
            self.reset_game()
        elif event.char.isdigit():
            idx = int(event.char)
            if 0 <= idx < self.n:
                self.on_click_vertex(idx)

    def bind_vertex_clicks(self):
        def on_canvas_click(e):
            min_dist = float('inf')
            chosen = None
            for i, (x, y) in enumerate(self.vertices):
                d = hypot(e.x - x, e.y - y)
                if d < min_dist and d <= self.vertex_radius + 8:
                    min_dist = d
                    chosen = i
            if chosen is not None:
                self.on_click_vertex(chosen)
        self.canvas.bind("<Button-1>", on_canvas_click)

    def run(self):
        self.master.mainloop()

if __name__ == "__main__":
    root = tk.Tk()
    game = HamiltonGraphGame(root)
    game.run()
