from collections import deque
from dataclasses import dataclass
from typing import Iterator, List, Tuple

from intcode_computer import CompleteIntcodeComputer, Status


@dataclass(frozen=True)
class Packet(object):
    x: int
    y: int


@dataclass()
class NAT(object):
    last_delivered_y_value: int
    packet: Packet


def get_input(packet_queue: deque[Packet]) -> Iterator[int]:
    while packet_queue:
        packet = packet_queue.popleft()
        yield packet.x
        yield packet.y
    yield -1
    return


with open("23/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computers: List[Tuple[CompleteIntcodeComputer, deque[Packet]]] = []
for i in range(50):
    computer = CompleteIntcodeComputer(program)
    status, output = computer.start(i)
    # Check that it's fine not to consider these initial outputs.
    assert status == Status.INPUT_REQUIRED and output == []
    computers.append((computer, deque()))

nat = NAT(None, None)

done_part_1 = False
while True:
    had_output = True
    # Loop through all computers in the network, routing packets as appropriate.
    for computer, packet_queue in computers:
        for next_input in get_input(packet_queue):
            status, output = computer.resume(next_input)
            assert status == Status.INPUT_REQUIRED
            while output:
                had_output = False
                y = output.pop()
                x = output.pop()
                address = output.pop()
                if address == 255:
                    nat.packet = Packet(x, y)
                    if not done_part_1:
                        print("The answer to Part 1 is {}.".format(y))
                        done_part_1 = True
                else:
                    computers[address][1].append(Packet(x, y))
    # The network is idle if there was no output and all packet queues are empty.
    if not had_output and all(map(lambda item: len(item[1]) == 0, computers)):
        # When idle the NAT delivers a packet computer 0.
        computers[0][1].append(nat.packet)
        last_nat_y = nat.last_delivered_y_value
        if last_nat_y and last_nat_y == nat.packet.y:
            print("The answer to Part 2 is {}.".format(last_nat_y))
            exit(0)
        else:
            nat.last_delivered_y_value = nat.packet.y