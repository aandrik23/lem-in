import json
import argparse
import matplotlib.pyplot as plt
import networkx as nx
from matplotlib.widgets import Button, Slider
from collections import defaultdict

class SimulationViewer:
    def __init__(self, dump_file):
        # Load data
        with open(dump_file) as f:
            data = json.load(f)
        self.rooms      = data['rooms']
        self.moves      = data['moves']
        self.start_room = data.get('start')

        # Build graph and static layout
        self.G = nx.Graph()
        self.pos = {}
        for r in self.rooms:
            self.G.add_node(r['name'])
            self.pos[r['name']] = (r['x'], r['y'])

        # Infer edges (spawn ⇒ start_room)
        edges = set()
        for m in self.moves:
            frm = m['from'] or self.start_room
            to  = m['to']
            edges.add((frm, to))
        self.G.add_edges_from(edges)

        # Organize moves by turn
        self.moves_by_turn = defaultdict(list)
        for m in self.moves:
            self.moves_by_turn[m['turn']].append(m)
        self.max_turn = max(self.moves_by_turn.keys())

        # Playback state
        self.current_turn = 0
        self.playing = False

        # Precompute ant positions per turn
        self.ant_positions_by_turn = {}
        positions = {}  # ant → room
        for t in range(self.max_turn+1):
            for m in self.moves_by_turn.get(t, []):
                positions[m['ant']] = m['to']
            # deep copy for this turn
            self.ant_positions_by_turn[t] = dict(positions)

    def draw(self):
        self.ax.clear()
        # draw nodes & edges *without* labels
        nx.draw(self.G, pos=self.pos, ax=self.ax,
                with_labels=False, node_size=800,
                node_color="#FFD700", edge_color="#555555")
        # draw room names above each node
        # compute a small y-offset based on the overall y-range
        ys = [y for (_, y) in self.pos.values()]
        y_offset = (max(ys) - min(ys)) * 0.05
        for room, (x, y) in self.pos.items():
            self.ax.text(
                x, y + y_offset, room,
                ha='center', va='bottom',
                fontsize=10, fontweight='bold', color='black'
            )
        # overlay ants
        for ant, room in self.ant_positions_by_turn[self.current_turn].items():
            x, y = self.pos[room]
            self.ax.text(x, y, str(ant), fontsize=12,
                         fontweight="bold", color="crimson",
                         ha='center', va='center')
        self.ax.set_title(f"Turn {self.current_turn}", pad=20)
        # redraw the figure
        self.fig.canvas.draw_idle()
        # update slider without triggering its callback
        self.slider.eventson = False
        self.slider.set_val(self.current_turn)
        self.slider.eventson = True

    def on_next(self, event):
        if self.current_turn < self.max_turn:
            self.current_turn += 1
            self.draw()

    def on_prev(self, event):
        if self.current_turn > 0:
            self.current_turn -= 1
            self.draw()

    def on_play(self, event):
        self.playing = not self.playing
        self.btn_play.label.set_text("❚❚" if self.playing else "►")
        if self.playing:
            self._animate()

    def _animate(self):
        if not self.playing:
            return
        if self.current_turn < self.max_turn:
            self.current_turn += 1
            self.draw()
            # call again after 500 ms
            self.fig.canvas.new_timer(
                interval=500, callbacks=[(self._animate, (), {})]
            ).start()
        else:
            # stop at end
            self.playing = False
            self.btn_play.label.set_text("►")

    def on_slider(self, val):
        self.current_turn = int(val)
        self.draw()

    def show(self):
        # set up figure
        plt.style.use('ggplot')
        self.fig, self.ax = plt.subplots(figsize=(8,6))
        plt.subplots_adjust(bottom=0.25)

        # buttons
        ax_prev = plt.axes([0.1, 0.05, 0.1, 0.075])
        ax_play = plt.axes([0.225, 0.05, 0.1, 0.075])
        ax_next = plt.axes([0.35, 0.05, 0.1, 0.075])
        self.btn_prev = Button(ax_prev, "◀ Prev")
        self.btn_play = Button(ax_play, "►")
        self.btn_next = Button(ax_next, "Next ▶")

        self.btn_prev.on_clicked(self.on_prev)
        self.btn_play.on_clicked(self.on_play)
        self.btn_next.on_clicked(self.on_next)

        # slider
        ax_slider = plt.axes([0.55, 0.05, 0.35, 0.04])
        self.slider = Slider(ax_slider, "Turn", 0, self.max_turn,
                             valinit=0, valstep=1)
        self.slider.on_changed(self.on_slider)

        # initial draw
        self.draw()
        plt.show()


def main():
    parser = argparse.ArgumentParser(
        description="Visualize lem-in simulation with controls"
    )
    parser.add_argument('--input', '-i', required=True,
                        help='Simulation JSON file')
    args = parser.parse_args()

    viewer = SimulationViewer(args.input)
    viewer.show()

if __name__ == '__main__':
    main()
