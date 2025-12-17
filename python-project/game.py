import tkinter as tk
from math import cos, sin, pi, hypot

class HamiltonDodecahedronGame:
    def __init__(self, master):
        self.master = master
        self.master.title("Гамильтонова игра: додекаэдр")
        self.vertex_radius = 14
        self.canvas_size = 700

        self.canvas = tk.Canvas(master, width=self.canvas_size, height=self.canvas_size, bg="white")
        self.canvas.pack()

        self.info = tk.Label(master, text="Кликайте по вершинам додекаэдра, чтобы построить гамильтонов цикл", anchor="w")
        self.info.pack(fill="x")

        self.reset_button = tk.Button(master, text="Сброс", command=self.reset_game)
        self.reset_button.pack()

        self.init_graph()
        self.reset_game_state()
        self.draw_graph()

        self.master.bind("<Key>", self.on_key)
        self.bind_vertex_clicks()

    def init_graph(self):
        # 20 вершин додекаэдра зададим "ручными" координатами в виде двух колец + верх/низ.
        # Это не геометрически точный 3D, а удобная плоская проекция.
        cx = cy = self.canvas_size / 2
        r_outer = 260
        r_inner = 140

        # 10 вершин внешнего десятиугольника
        outer = []
        for i in range(10):
            angle = 2 * pi * i / 10 - pi / 2
            x = cx + r_outer * cos(angle)
            y = cy + r_outer * sin(angle)
            outer.append((x, y))

        # 10 вершин внутреннего десятиугольника
        inner = []
        for i in range(10):
            angle = 2 * pi * (i + 0.5) / 10 - pi / 2
            x = cx + r_inner * cos(angle)
            y = cy + r_inner * sin(angle)
            inner.append((x, y))

        # Список вершин: сначала внешние 0–9, затем внутренние 10–19
        self.vertices = outer + inner

        # Рёбра додекаэдра (20 вершин, 3 ребра у каждой).
        # Конкретная нумерация подобрана так, чтобы получился «додекаэдроподобный» 3-regular граф:
        edges = set()

        # Рёбра внешнего кольца (десятиугольник)
        for i in range(10):
            edges.add((i, (i + 1) % 10))

        # Рёбра внутреннего кольца (десятиугольник)
        for i in range(10, 20):
            edges.add((i, 10 + (i - 9) % 10))

        # Рёбра между внешними и внутренними вершинами (имитируем грани додекаэдра)
        # Каждая внешняя соединена с двумя соседними внутренними
        for i in range(10):
            inner1 = 10 + i
            inner2 = 10 + ((i - 1) % 10)
            edges.add((i, inner1))
            edges.add((i, inner2))

        self.edges = set()
        for a, b in edges:
            if a != b:
                self.edges.add(tuple(sorted((a, b))))

        self.vertex_ids = list(range(len(self.vertices)))
        self.n = len(self.vertices)  # 20

    def draw_graph(self):
        self.canvas.delete("all")

        # Рёбра
        for a, b in self.edges:
            ax, ay = self.vertices[a]
            bx, by = self.vertices[b]
            self.canvas.create_line(ax, ay, bx, by, fill="#ccc", width=2)

        # Текущий путь (синим)
        if len(self.visited_order) >= 2:
            for i in range(len(self.visited_order) - 1):
                a = self.visited_order[i]
                b = self.visited_order[i + 1]
                ax, ay = self.vertices[a]
                bx, by = self.vertices[b]
                self.canvas.create_line(ax, ay, bx, by, fill="blue", width=4)
            # При замыкании цикла выделим последнюю дугу
            if len(self.visited_order) == self.n:
                first = self.visited_order[0]
                last = self.visited_order[-1]
                if (min(first, last), max(first, last)) in self.edges:
                    ax, ay = self.vertices[last]
                    bx, by = self.vertices[first]
                    self.canvas.create_line(ax, ay, bx, by, fill="green", width=4, dash=(4, 2))

        # Вершины
        for idx, (x, y) in enumerate(self.vertices):
            if not self.visited_order:
                color = "white"
            elif idx == self.visited_order[0]:
                color = "#ffe4b5"  # стартовая вершина
            elif idx in self.visited_order:
                color = "lightblue"
            else:
                color = "white"

            self.canvas.create_oval(
                x - self.vertex_radius, y - self.vertex_radius,
                x + self.vertex_radius, y + self.vertex_radius,
                fill=color, outline="black", width=2
            )
            self.canvas.create_text(x, y, text=str(idx), font=("Arial", 10, "bold"))

    def reset_game_state(self):
        self.visited_order = []
        self.update_info()

    def reset_game(self):
        self.reset_game_state()
        self.draw_graph()

    def update_info(self):
        path_text = " -> ".join(map(str, self.visited_order)) if self.visited_order else "(пусто)"
        self.info.config(
            text=f"Путь: {path_text} | посещено: {len(self.visited_order)}/{self.n} "
                 f"(R — сброс, клики по вершинам)"
        )

    def on_click_vertex(self, idx):
        # уже есть полный путь — игнорируем клики
        if len(self.visited_order) == self.n:
            return

        if not self.visited_order:
            # первая вершина — старт цикла
            self.visited_order.append(idx)
        else:
            if idx in self.visited_order:
                # нельзя заходить повторно, кроме как замыкать цикл в самом конце
                if len(self.visited_order) == self.n and idx == self.visited_order[0]:
                    self.check_completion()
                return
            prev = self.visited_order[-1]
            edge = (min(prev, idx), max(prev, idx))
            if edge in self.edges:
                self.visited_order.append(idx)
            else:
                self.info.config(text="Нет ребра между выбранными вершинами. Попробуйте другую вершину.")
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
                self.info.config(text="Гамильтонов цикл на графе додекаэдра построен! Поздравляю!", fg="green")
            else:
                self.info.config(text="Все вершины посещены, но цикл не замкнут на стартовую вершину.", fg="orange")
            return True
        return False

    def on_key(self, event):
        ch = event.char.lower()
        if ch == 'r':
            self.reset_game()
        elif ch.isdigit():
            idx = int(ch)
            if 0 <= idx < self.n:
                self.on_click_vertex(idx)

    def bind_vertex_clicks(self):
        def on_canvas_click(e):
            min_dist = float("inf")
            chosen = None
            for i, (x, y) in enumerate(self.vertices):
                d = hypot(e.x - x, e.y - y)
                if d < min_dist and d <= self.vertex_radius + 8:
                    min_dist = d
                    chosen = i
            if chosen is not None:
                self.on_click_vertex(chosen)
        self.canvas.bind("<Button-1>", on_canvas_click)

if __name__ == "__main__":
    root = tk.Tk()
    game = HamiltonDodecahedronGame(root)
    root.mainloop()
