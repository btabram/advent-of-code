from collections import namedtuple
import itertools


# Satutus codes
SUCCESS = 0
WAIT_FOR_INPUT = 1


Param = namedtuple('Param', ['mode', 'value'])


class IntcodeComputer:
    VALID_OPCODES = [1, 2, 3, 4, 5, 6, 7, 8]

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
    }

    def __init__(self, program):
        self._program = program
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

    def start(self, input):
        self._instruction_pointer = 0
        self._jumped = False
        self._memory = self._program.copy()
        self._input = input.copy()
        self._output = []
        status_code = self.run()
        return status_code, self._output

    def resume(self, additional_input):
        self._input.extend(additional_input)
        status_code = self.run()
        return status_code, self._output

    def run(self):
        while True:
            instruction = self.read()
            if instruction == 99: # Halt
                return SUCCESS
            opcode, parameter_modes = self.parse_instruction(instruction)

            parameter_count = self.PARAMETER_COUNTS[opcode]

            parameters = self.read_params(parameter_count, parameter_modes)

            opcode_fn = getattr(self, "op{}".format(opcode))
            status_code = opcode_fn(parameters)
            if status_code != SUCCESS:
                # An opcode function which doesn't reutrn SUCCESS must not
                # modify any state so that the computer can be resumed.
                return status_code

            self.advance_to_next_instruction(parameter_count)

    # Add
    def op1(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) + self.read_param(p2)
        self.write_to_param(p3, output)
        return SUCCESS

    # Mulitply
    def op2(self, parameters):
        [p1, p2, p3] = parameters
        output = self.read_param(p1) * self.read_param(p2)
        self.write_to_param(p3, output)
        return SUCCESS

    # Input
    def op3(self, parameters):
        if len(self._input) == 0:
            return WAIT_FOR_INPUT
        [p1] = parameters
        self.write_to_param(p1, self._input.pop(0))
        return SUCCESS

    # Output
    def op4(self, parameters):
        [p1] = parameters
        self._output.append(self.read_param(p1))
        return SUCCESS

    # Jump if true
    def op5(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) != 0:
            self.jump(self.read_param(p2))
        return SUCCESS

    # Jump if false
    def op6(self, parameters):
        [p1, p2] = parameters
        if self.read_param(p1) == 0:
            self.jump(self.read_param(p2))
        return SUCCESS

    # Less than
    def op7(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) < self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)
        return SUCCESS

    # Equals
    def op8(self, parameters):
        [p1, p2, p3] = parameters
        if self.read_param(p1) == self.read_param(p2):
            self.write_to_param(p3, 1)
        else:
            self.write_to_param(p3, 0)
        return SUCCESS


with open("07/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

amp_a = IntcodeComputer(program)
amp_b = IntcodeComputer(program)
amp_c = IntcodeComputer(program)
amp_d = IntcodeComputer(program)
amp_e = IntcodeComputer(program)
initial_input = 0

def calculate_thruster_signal(phase_a, phase_b, phase_c, phase_d, phase_e):
    (_,      [output_a]) = amp_a.start([phase_a, initial_input])
    (_,      [output_b]) = amp_b.start([phase_b, output_a])
    (_,      [output_c]) = amp_c.start([phase_c, output_b])
    (_,      [output_d]) = amp_d.start([phase_d, output_c])
    (status, [output_e]) = amp_e.start([phase_e, output_d])
    # We have halted after one run through and should just return.
    if status == SUCCESS:
        return output_e
    # Feedback loop mode:
    prev_output = output_e
    while True:
        for (i, amp) in enumerate([amp_a, amp_b, amp_c, amp_d, amp_e]):
            (status, output) = amp.resume([prev_output])
            prev_output = output[-1]
            # Return when the final amp halts.
            if i == 4 and status == SUCCESS:
                return prev_output

def maximise_thruster_signal(possible_phases):
    max_signal = 0
    best_phase_sequence = None
    for phases in itertools.permutations(possible_phases):
        signal = calculate_thruster_signal(*phases)
        if signal > max_signal:
            max_signal = signal
            best_phase_sequence = phases
    return max_signal, best_phase_sequence

print("The max signal and associated phase sequence for Part 1 is {}.".format(
    maximise_thruster_signal([0, 1, 2, 3, 4])))
print("The max signal and associated phase sequence for Part 2 is {}.".format(
    maximise_thruster_signal([5, 6, 7, 8, 9])))
