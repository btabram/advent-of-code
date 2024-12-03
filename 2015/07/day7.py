class LogicGate:
    def __init__(self, s: str) -> None:
        [input, output] = s.split(" -> ")
        self.args: list[int | str] = []
        self.op: str | None = None
        self.output = output
        for item in input.split():
            if item[0].isupper():
                self.op = item
            else:
                try:
                    self.args.append(int(item))
                except ValueError:
                    self.args.append(item)

    def __repr__(self) -> str:
        return f"{self.op} {self.args} -> {self.output}"


def solve_for_a(logic_gates: list[LogicGate], initial_wires: dict[str, int]) -> int:
    known_wires = dict(initial_wires)
    while True:
        for lg in logic_gates:
            # Has this logic gate already been solved?
            if lg.output in known_wires:
                continue

            # Do we have enough information to solve this logic gate?
            if all([isinstance(arg, int) or arg in known_wires for arg in lg.args]):
                args = [
                    arg if isinstance(arg, int) else known_wires[arg] for arg in lg.args
                ]
                match lg.op:
                    case None:
                        result = args[0]
                    case "NOT":
                        result = ~args[0]
                    case "AND":
                        result = args[0] & args[1]
                    case "OR":
                        result = args[0] | args[1]
                    case "LSHIFT":
                        result = args[0] << args[1]
                    case "RSHIFT":
                        result = args[0] >> args[1]
                    case _:
                        raise Exception(f"Invalid operation: {lg.op}")
                result &= 65535  # Wire signals are 16-bit so we discard any other bits
                known_wires[lg.output] = result

            if "a" in known_wires:
                return known_wires["a"]


with open("input.txt") as f:
    logic_gates = [LogicGate(line.strip()) for line in f.readlines()]

part1 = solve_for_a(logic_gates, {})
part2 = solve_for_a(logic_gates, {"b": part1})

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
