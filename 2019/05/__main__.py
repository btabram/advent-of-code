from collections import namedtuple


Param = namedtuple('Param', ['mode', 'value'])


class IntcodeComputer:
    VALID_OPCODES = [1, 2, 3, 4, 5, 6, 7, 8]

    # Map of opcode -> number of parameters
    PARAMETER_COUNTS = {
        1: 3, # Add
        2: 3, # Mulitply
        3: 1, # Input
        4: 1, # Output
        5: 2, # Jump if true
        6: 2, # Jump if false
        7: 3, # Less than
        8: 3, # Equals
    }

    def __init__(self):
        self._instruction_pointer = 0
        self._jumped = False
        self._memory = []
        self._input = []
        self._output = []

    def advance_to_next_instruction(self, parameter_count):
        if not self._jumped:
            self._instruction_pointer += 1 + parameter_count
        # Reset per-instruction state.
        self._jumped = False

    def jump(self, position):
        self._instruction_pointer = position
        self._jumped = True

    def read(self, position = None, offset = 0):
        """
        Return the value of the program at |position| + |offset|
        |position| defaults to |self._instruction_pointer|
        """
        if position is None:
            position = self._instruction_pointer
        return self._memory[position + offset]

    def write(self, position, value):
        self._memory[position] = value

    def parse_instruction(self, instruction):
        """
        An instruction is up to 5 digits ABCDE where DE is the opcode and ABC
        are the parameter modes which may be omitted if they are zero.
        """
        opcode = instruction % 100
        if opcode not in self.VALID_OPCODES:
            raise Exception("Invalid opcode: {}".format(opcode))
        parameter_modes = instruction // 100
        return opcode, parameter_modes

    def read_params(self, parameter_count, parameter_modes):
        """
        |parameter_modes| is up to 3 digits ABC where C is the mode of the 1st
        parameter, B is the mode of the 2nd parameter etc. Leading zeroes may
        be ommitted.

        Parameters are always immediately after their instruction, so the 1st
        parameter is at offset = 1 and so on.
        """
        parameters = []
        if parameter_count >= 1:
            first_param_mode = parameter_modes % 10
            parameters.append(Param(first_param_mode, self.read(offset = 1)))
        if parameter_count >= 2:
            second_param_mode = (parameter_modes // 10) % 10
            parameters.append(Param(second_param_mode, self.read(offset = 2)))
        if parameter_count >= 3:
            third_param_mode = (parameter_modes // 100) % 10
            parameters.append(Param(third_param_mode, self.read(offset = 3)))
        return parameters

    def read_param(self, parameter):
        if parameter.mode == 0: # Position mode
            return self.read(parameter.value)
        elif parameter.mode == 1: # Immediate mode
            return parameter.value
        else:
            raise Exception("Invalid parameter mode: {}".format(
                parameter.mode))

    def write_to_param(self, parameter, write_value):
        if parameter.mode == 0: # Position mode
            return self.write(parameter.value, write_value)
        else:
            # Note that immediate mode does not support writing.
            raise Exception("Invalid parameter mode for writing: {}".format(
                parameter.mode))

    def run(self, program, input):
        self._instruction_pointer = 0
        self._memory = program.copy()
        self._input = input.copy()
        self._output = []

        while True:
            instruction = self.read()
            if instruction == 99: # Halt
                break
            opcode, parameter_modes = self.parse_instruction(instruction)

            parameter_count = self.PARAMETER_COUNTS[opcode]

            parameters = self.read_params(parameter_count, parameter_modes)

            opcode_fn = getattr(self, "op{}".format(opcode))
            opcode_fn(parameters)

            self.advance_to_next_instruction(parameter_count)

        return self._output

    # Add
    def op1(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) + self.read_param(p2)
        self.write_to_param(p3, output)

    # Mulitply
    def op2(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) * self.read_param(p2)
        self.write_to_param(p3, output)

    # Input
    def op3(self, parameters):
        [p1] = parameters
        self.write_to_param(p1, self._input.pop(0))

    # Output
    def op4(self, parameters):
        [p1] = parameters
        self._output.append(self.read_param(p1))

    # Jump if true
    def op5(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) != 0:
            self.jump(self.read_param(p2))

    # Jump if false
    def op6(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) == 0:
            self.jump(self.read_param(p2))

    # Less than
    def op7(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) < self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)

    # Equals
    def op8(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) == self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)


with open("05/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = IntcodeComputer()

print("The output for Part 1 is {}.".format(computer.run(program, [1])))
print("The output for Part 2 is {}.".format(computer.run(program, [5])))
