from collections import namedtuple
from enum import Enum


class Status(Enum):
    SUCCESS = 0
    INPUT_REQUIRED = 1


Param = namedtuple('Param', ['mode', 'value'])


class CompleteIntcodeComputer:
    VALID_OPCODES = [1, 2, 3, 4, 5, 6, 7, 8, 9]

    # Dict of opcode -> number of parameters
    PARAMETER_COUNTS = {
        1: 3, # Add
        2: 3, # Mulitply
        3: 1, # Input
        4: 1, # Output
        5: 2, # Jump if true
        6: 2, # Jump if false
        7: 3, # Less than
        8: 3, # Equals
        9: 1, # Adjust relative base
    }

    def __init__(self, program, ascii = False):
        self._program = program
        self._instruction_pointer = 0
        self._jumped = False
        self._relative_base = 0
        self._memory = []
        self._input = []
        self._output = []
        self._ascii = ascii

    def advance_to_next_instruction(self, parameter_count):
        if not self._jumped:
            self._instruction_pointer += 1 + parameter_count
        # Reset per-instruction state.
        self._jumped = False

    def jump(self, position):
        self._instruction_pointer = position
        self._jumped = True

    def extend_memory_if_required(self, address):
        """
        Running a program may need much more memory than the program itself
        takes up so we need to be able to expand |self._memory| as required.
        """
        if address >= len(self._memory):
            # Memory is always initialized to zero.
            memory_to_add = (1 + address - len(self._memory)) * [0]
            self._memory.extend(memory_to_add)

    def read(self, position = None, offset = 0):
        """
        Return the value of the program at |position| + |offset|
        |position| defaults to |self._instruction_pointer|
        """
        if position is None:
            position = self._instruction_pointer
        address = position + offset
        self.extend_memory_if_required(address)
        return self._memory[address]

    def write(self, address, value):
        self.extend_memory_if_required(address)
        self._memory[address] = value

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
            return self.read(position = parameter.value)
        elif parameter.mode == 1: # Immediate mode
            return parameter.value
        elif parameter.mode == 2: # Relative mode
            return self.read(position = parameter.value + self._relative_base)
        else:
            raise Exception("Invalid parameter mode: {}".format(
                parameter.mode))

    def write_to_param(self, parameter, value):
        if parameter.mode == 0: # Position mode
            return self.write(parameter.value, value)
        elif parameter.mode == 2: # Relative mode
            return self.write(parameter.value + self._relative_base, value)
        else:
            # Note that immediate mode does not support writing.
            raise Exception("Invalid parameter mode for writing: {}".format(
                parameter.mode))

    def start(self, input = None):
        """
        |input| can be an int, iterable or string (if ascii is set).
        """
        self._instruction_pointer = 0
        self._jumped = False
        self._relative_base = 0
        self._memory = self._program.copy()
        if input is None:
            self._input = []
        elif isinstance(input, int):
            self._input = [input]
        elif self._ascii and isinstance(input, str):
            self._input = [ord(c) for c in input]
        else:
            self._input = input.copy()
        self._output = []
        status_code = self.run()
        # Output a string when configured, if the values are in the ascii range.
        if self._ascii and all([0 < val < 128 for val in self._output]):
            self._output = "".join(map(chr, self._output))
        return status_code, self._output

    def resume(self, additional_input):
        """
        |additional_input| can be an int, iterable or string (if ascii is set).
        """
        if isinstance(additional_input, int):
            self._input.append(additional_input)
        elif self._ascii and isinstance(additional_input, str):
            self._input.extend([ord(c) for c in additional_input])
        else:
            self._input.extend(additional_input)
        # Only care about new output.
        self._output = []
        status_code = self.run()
        # Output a string when configured, if the values are in the ascii range.
        if self._ascii and all([0 < val < 128 for val in self._output]):
            self._output = "".join(map(chr, self._output))
        return status_code, self._output

    def run(self):
        while True:
            instruction = self.read()
            if instruction == 99: # Halt
                return Status.SUCCESS
            opcode, parameter_modes = self.parse_instruction(instruction)

            parameter_count = self.PARAMETER_COUNTS[opcode]

            parameters = self.read_params(parameter_count, parameter_modes)

            opcode_fn = getattr(self, "op{}".format(opcode))
            status_code = opcode_fn(parameters)
            if status_code != Status.SUCCESS:
                # An opcode function which doesn't return SUCCESS must not
                # modify any state so that the computer can be resumed.
                return status_code

            self.advance_to_next_instruction(parameter_count)

    # Add
    def op1(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) + self.read_param(p2)
        self.write_to_param(p3, output)
        return Status.SUCCESS

    # Mulitply
    def op2(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) * self.read_param(p2)
        self.write_to_param(p3, output)
        return Status.SUCCESS

    # Input
    def op3(self, parameters):
        if len(self._input) == 0:
            return Status.INPUT_REQUIRED
        [p1] = parameters
        self.write_to_param(p1, self._input.pop(0))
        return Status.SUCCESS

    # Output
    def op4(self, parameters):
        [p1] = parameters
        self._output.append(self.read_param(p1))
        return Status.SUCCESS

    # Jump if true
    def op5(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) != 0:
            self.jump(self.read_param(p2))
        return Status.SUCCESS

    # Jump if false
    def op6(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) == 0:
            self.jump(self.read_param(p2))
        return Status.SUCCESS

    # Less than
    def op7(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) < self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)
        return Status.SUCCESS

    # Equals
    def op8(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) == self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)
        return Status.SUCCESS

    # Adjust relative base
    def op9(self, parameters):
        [p1] = parameters
        self._relative_base += self.read_param(p1)
        return Status.SUCCESS
