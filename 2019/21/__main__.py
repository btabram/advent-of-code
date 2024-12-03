from intcode_computer import CompleteIntcodeComputer


with open("21/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = CompleteIntcodeComputer(program, ascii = True)


# Part 1
computer.start()
# The overall logic is:
# (!A || !B || !C) && D
# or in words, jump if there's a hole in one of the three steps in front as
# long as the fourth step is not a hole (because the robot jumps 4 steps).
computer.resume("NOT A J\n")
computer.resume("NOT B T\n")
computer.resume("OR T J\n")
computer.resume("NOT C T\n")
computer.resume("OR T J\n")
computer.resume("AND D J\n")
_, output = computer.resume("WALK\n")
print("The answer to Part 1 is {}.".format(output[-1]))


# Part 2
computer.start()
# The overall logic is:
# (!A || !B || !C) && D && (H || E)
# or in words, jump if there's a hole in one of the three steps in front as
# long as the fourth step is not a hole (because the robot jumps 4 steps) and
# either the eighth step is not a hole (so we can jump again staight after
# landing) or the fifth step is not a hole (so we can take a step before
# needing to jump again).
computer.resume("NOT A J\n")
computer.resume("NOT B T\n")
computer.resume("OR T J\n")
computer.resume("NOT C T\n")
computer.resume("OR T J\n")

computer.resume("AND D J\n")

# Double NOT to get H into T no matter wath the initial value of T is.
computer.resume("NOT H T\n")
computer.resume("NOT T T\n")
computer.resume("OR E T\n")

computer.resume("AND T J\n")

_, output = computer.resume("RUN\n")
print("The answer to Part 2 is {}.".format(output[-1]))
