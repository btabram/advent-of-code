from intcode_computer import CompleteIntcodeComputer, Status
import os
import time


class ArcadeCabinet:
    def __init__(self, program):
        self.original_program = program.copy()
        # Dict of screen coordinates -> tile id
        self.screen_data = dict()
        self.width = None
        self.height = None
        self.score = None

    def get_tile_positions(self, tile_id):
        if self.screen_data is None:
            return None
        return [p for (p, id) in self.screen_data.items() if id == tile_id]

    @property
    def ball_position(self):
        positions = self.get_tile_positions(4)
        assert len(positions) == 1
        return positions[0]

    @property
    def paddle_position(self):
        positions = self.get_tile_positions(3)
        assert len(positions) == 1
        return positions[0]

    def play(self, quarters = None):
        self.screen_data = dict()
        self.width = None
        self.height = None
        self.score = None

        program = self.original_program.copy()
        if quarters is not None:
            program[0] = quarters
        computer = CompleteIntcodeComputer(program)

        status, output = computer.start()
        self.update_screen_data(output)
        self.draw_screen()

        # An initial paddle movement to the left sets things up correctly so
        # that we now just need to move the paddle in whatever direction the
        # ball is moving.
        paddle_command = -1
        prev_ball_position = self.ball_position
        while status == Status.INPUT_REQUIRED:
            status, output = computer.resume(paddle_command)
            self.update_screen_data(output)
            self.draw_screen()

            b_x, b_y = self.ball_position
            p_x, p_y = self.paddle_position

            # Don't move the paddle if the ball is directly above.
            if b_x == p_x and b_y == p_y - 1:
                paddle_command = 0
            # Otherwise move the paddle in the same direction as the ball.
            elif b_x > prev_ball_position[0]:
                paddle_command = 1
            else:
                paddle_command = -1

            prev_ball_position = (b_x, b_y)
            time.sleep(0.01)

    def update_screen_data(self, new_data):
        assert len(new_data) % 3 == 0
        for i in range(0, len(new_data), 3):
            x = new_data[i]
            y = new_data[i + 1]
            if x == -1 and y == 0:
                self.score = new_data[i + 2]
                continue
            tile_id = new_data[i + 2]
            self.screen_data[(x, y)] = tile_id

        if self.width == None:
            max_x = max(self.screen_data.keys(), key = lambda pos: pos[0])[0]
            min_x = min(self.screen_data.keys(), key = lambda pos: pos[0])[0]
            max_y = max(self.screen_data.keys(), key = lambda pos: pos[1])[1]
            min_y = min(self.screen_data.keys(), key = lambda pos: pos[1])[1]
            assert min_x == 0
            assert min_y == 0
            self.width = max_x + 1
            self.height = max_y + 1

    def draw_screen(self):
        # First clear the screen.
        if os.name == 'posix': # Mac and Linux
            os.system('clear')
        else: # Windows
            os.system('cls')
        print()
        score_line = "##########     SCORE = "
        score_line += str(self.score)
        padding = self.width - len(score_line) - 10
        score_line += " " * padding
        score_line += "##########"
        print(score_line)
        for y in range(self.height):
            line = ""
            for x in range(self.width):
                tile_id = self.screen_data[(x, y)]
                assert tile_id in [0, 1, 2, 3, 4]
                if tile_id == 0: # Empty
                    line += " "
                elif tile_id == 1: # Wall
                    line += "|"
                elif tile_id == 2: # Block
                    line += "="
                elif tile_id == 3: # Horizontal paddle
                    line += "‾" # U+203E OVERLINE
                else: # Ball
                    line += "o"
            print(line)


with open("13/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

# Part 1.
cabinet = ArcadeCabinet(program)
cabinet.play()
number_of_blocks = len(cabinet.get_tile_positions(2))

# Part 2.
cabinet.play(2)

print()
print(f"The answer to Part 1 is {number_of_blocks}")
print(f"The answer to Part 2 is {cabinet.score}")
