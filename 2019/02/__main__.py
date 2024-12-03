VALID_OPCODES = [1, 2]


def run_intcode(program):
    instruction_pointer = 0
    while True:
        opcode = program[instruction_pointer]
        if opcode == 99:
            return program
        elif opcode not in VALID_OPCODES:
            raise Exception("Invalid opcode: {}".format(opcode))

        input_1_pos = program[instruction_pointer + 1]
        input_2_pos = program[instruction_pointer + 2]
        output_pos = program[instruction_pointer + 3]

        input_1 = program[input_1_pos]
        input_2 = program[input_2_pos]

        if opcode == 1:
            program[output_pos] = input_1 + input_2
        else: # opcode == 2
            program[output_pos] = input_1 * input_2

        instruction_pointer += 4


with open("02/input.txt") as f:
    input_program = list(map(int, f.readline().split(",")))


def run_gravity_assist_program(noun, verb):
    program = input_program.copy()
    program[1] = noun
    program[2] = verb
    run_intcode(program)
    return program[0]


print("The answer to Part 1 is {}".format(run_gravity_assist_program(12, 2)))


def find_inputs_for_output(desired_output):
    for j in range(100):
        for i in range(100):
            if run_gravity_assist_program(i, j) == desired_output:
                return i, j


noun, verb = find_inputs_for_output(19690720)
print("The answer to Part 2 is {}".format(100 * noun + verb))
