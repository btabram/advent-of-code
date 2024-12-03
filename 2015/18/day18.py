class LightsGrid:
    def __init__(
        self,
        initial_grid: list[str],
        always_on_lights: None | list[tuple[int, int]] = None,
    ) -> None:
        assert len(initial_grid) == len(initial_grid[0])  # Check square
        self._size = len(initial_grid)

        self._grid: list[list[bool]] = []
        for initial_row in initial_grid:
            row = []
            for initial_value in initial_row:
                row.append(True if initial_value == "#" else False)
            self._grid.append(row)

        self._always_on_lights = (
            always_on_lights if always_on_lights is not None else []
        )
        for (x, y) in self._always_on_lights:
            self._grid[x][y] = True

    def __repr__(self) -> str:
        lines = []
        for row in self._grid:
            line = ""
            for is_on in row:
                line += "#" if is_on else "."
            lines.append(line)
        return "\n".join(lines) + "\n"

    def _get_value(self, x: int, y: int) -> bool:
        if x < 0 or x >= self._size:
            return False
        if y < 0 or y >= self._size:
            return False
        return self._grid[x][y]

    def _count_on_neighbours(self, x: int, y: int) -> int:
        neighbours = []
        for x_i in range(x - 1, x + 2):
            for y_i in range(y - 1, y + 2):
                if x_i == x and y_i == y:
                    continue
                neighbours.append((x_i, y_i))
        return [self._get_value(*n) for n in neighbours].count(True)

    def advance(self) -> None:
        lights_to_change = []
        for x, row in enumerate(self._grid):
            for y, is_on in enumerate(row):
                on_neighbours = self._count_on_neighbours(x, y)
                if is_on:
                    if on_neighbours not in [2, 3]:
                        lights_to_change.append((x, y))  # Turn off
                else:
                    if on_neighbours == 3:
                        lights_to_change.append((x, y))  # Turn on

        for (x, y) in lights_to_change:
            self._grid[x][y] = not self._grid[x][y]

        for (x, y) in self._always_on_lights:
            self._grid[x][y] = True

    def count_on_lights(self) -> int:
        return sum([row.count(True) for row in self._grid])


with open("input.txt") as f:
    input = [line.strip() for line in f.readlines()]

part1_grid = LightsGrid(input)

corners = [
    (0, 0),
    (0, len(input) - 1),
    (len(input) - 1, 0),
    (len(input) - 1, len(input) - 1),
]
part2_grid = LightsGrid(input, always_on_lights=corners)

for _ in range(100):
    part1_grid.advance()
    part2_grid.advance()

print(f"The answer to Part 1 is {part1_grid.count_on_lights()}.")
print(f"The answer to Part 2 is {part2_grid.count_on_lights()}.")
