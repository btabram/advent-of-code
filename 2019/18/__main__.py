from itertools import count
import sys


class VaultSolver:
    def __init__(self, grid_list):
        # Make a grid dict of position tuples -> values
        self.grid = dict()
        self.width = len(grid_list[0])
        self.height = len(grid_list)
        for x in range(self.width):
            for y in range(self.height):
                self.grid[(x, y)] = grid_list[y][x]

        self.number_of_keys = 0
        self.start_positions = []
        for pos, val in self.grid.items():
            if val == "@":
                self.start_positions.append(pos)
                self.grid[pos] = "."
            elif 97 <= ord(val) <= 122: # Keys are lower case letters
                self.number_of_keys += 1

        self.keys_cache = dict()
        self.steps_cache = dict()

        self.done = 0
        self.answer = sys.maxsize

    def get_valid_neighbours(self, position):
        x, y = position
        neighbours = [(x + 1, y), (x, y + 1), (x - 1, y),  (x, y - 1)]
        def is_valid(neighbour):
            try:
                return self.grid[neighbour] != "#" # Can't walk through walls
            except KeyError:
                return False
        return set(filter(is_valid, neighbours))

    def find_accessible_keys(self, current_position, keys_collected):
        keys = []

        # Go through open passages, already collected keys, and unlocked doors.
        passable_values = set(keys_collected)
        [passable_values.add(key.upper()) for key in keys_collected]
        passable_values.add(".")

        visited = set([current_position])
        boundary = set([current_position])
        for i in count(1):
            new_boundary = set()
            for pos in boundary:
                new_boundary |= self.get_valid_neighbours(pos)
            new_boundary -= visited

            to_remove = set()
            # Look for interesting things in the new boundary.
            for pos in new_boundary:
                val = self.grid[pos]
                if val not in passable_values:
                    to_remove.add(pos)
                    if 97 <= ord(val) <= 122: # Keys are lower case letters
                        keys.append((val, i, pos))
            # Remove things we don't want to or can't move past (like doors).
            new_boundary -= to_remove

            if len(new_boundary) == 0:
                break

            visited |= new_boundary
            boundary = new_boundary
        return keys

    def solve(self):
        self.solve_impl(self.start_positions)

    def solve_impl(self, positions, keys_collected = "", steps_so_far = 0,
                   first_subgrid_choice = None, first_branch_choice = None):

        bad = False
        while True:
            keys = []

            cache_key = (keys_collected, *positions)
            try:
                keys = self.keys_cache[cache_key]
            except KeyError:
                keys = []
                for pos in positions:
                    keys.append(self.find_accessible_keys(pos, keys_collected))
                self.keys_cache[cache_key] = keys

            multiple_subgrids = len(keys) > 1

            # We need to branch, trying moves in all subgrids which have
            # accessible keys, to ensure we find the optimal solution.
            if multiple_subgrids and first_subgrid_choice is None:
                for i in range(1, len(keys)):
                    if len(keys[i]) > 0:
                        self.solve_impl(positions.copy(), keys_collected,
                                        steps_so_far, i)
                if len(keys[0]) == 0:
                    return

            subgrid_id = 0
            # We have been spawned to follow a specific branch.
            if multiple_subgrids and first_subgrid_choice is not None:
                subgrid_id = first_subgrid_choice
                first_subgrid_choice = None

            subgrid_keys = keys[subgrid_id]

            branching = len(subgrid_keys) > 1

            # We need to branch to ensure we find the optimal solution.
            if branching and first_branch_choice is None:
                for i in range(1, len(subgrid_keys)):
                    subgrid_keys[i]
                    self.solve_impl(positions.copy(), keys_collected,
                                    steps_so_far, subgrid_id, i)

            choice = 0
            # We have been spawned to follow a specific branch.
            if branching and first_branch_choice is not None:
                choice = first_branch_choice
                first_branch_choice = None

            val, distance, key_pos = subgrid_keys[choice]
            keys_collected = "".join(sorted(keys_collected + val))
            steps_so_far += distance
            positions[subgrid_id] = key_pos

            if len(keys_collected) == self.number_of_keys:
                break

            cache_key = (keys_collected, *positions)
            try:
                cached_steps = self.steps_cache[cache_key]
                if steps_so_far >= cached_steps:
                    return
            except KeyError:
                pass
            self.steps_cache[cache_key] = steps_so_far

        if steps_so_far < self.answer:
            self.answer = steps_so_far


grid_list = []
with open("18/input.txt") as f:
    for line in f.readlines():
        grid_list.append([char for char in line.strip()])

vaultSolver = VaultSolver(grid_list)
vaultSolver.solve()
print("The answer to Part 1 is {}.".format(vaultSolver.answer))

# We need to modify the maze for Part 2.
entrance = None
for x in range(len(grid_list[0])):
    for y in range(len(grid_list)):
        if grid_list[y][x] == "@":
            entrance = (x, y)
assert entrance is not None

x, y = entrance
grid_list[y][x] = "#"
grid_list[y + 1][x] = "#"
grid_list[y - 1][x] = "#"
grid_list[y][x + 1] = "#"
grid_list[y][x - 1] = "#"
grid_list[y + 1][x + 1] = "@"
grid_list[y + 1][x - 1] = "@"
grid_list[y - 1][x + 1] = "@"
grid_list[y - 1][x - 1] = "@"

vaultSolver = VaultSolver(grid_list)
vaultSolver.solve()
print("The answer to Part 2 is {}.".format(vaultSolver.answer))
