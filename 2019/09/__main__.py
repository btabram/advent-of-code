from intcode_computer import CompleteIntcodeComputer

with open("09/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = CompleteIntcodeComputer(program)
print(f"The output for Part 1 is {computer.start(1)}.")
print(f"The output for Part 2 is {computer.start(2)}.")
