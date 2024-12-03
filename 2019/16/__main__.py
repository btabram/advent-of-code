from math import ceil
from collections import deque


def get_pattern(length, element_number):
    pattern = [0] * element_number
    pattern += [1] * element_number
    pattern += [0] * element_number
    pattern += [-1] * element_number
    if len(pattern) <= length + 1:
        pattern *= ceil((length + 1) / len(pattern))
    # Discard the very first value.
    return pattern[1:length + 1]


def apply_pattern(signal, pattern):
    assert len(signal) == len(pattern)
    value = sum([signal[i] * pattern[i] for i in range(len(signal))])
    # Keep only the ones digit.
    return abs(value) % 10


def do_phase(signal):
    length = len(signal)
    output = []
    for element_number in range(1, length + 1):
        pattern = get_pattern(length, element_number)
        output.append(apply_pattern(signal, pattern))
    return output


def fft(signal, phases):
    for _ in range(phases):
        signal = do_phase(signal)
    return signal


"""
The trick is that the pattern for the second half of the signal is always of the following form:
    00001111 (example for signal length 8)
    00000111
    00000011
    00000001
So the fft results for the second half only depend on the numbers in the second half. The fft for
the final n digits (assuming n is greater than half the length) only depend on the final n digits
of the input signal. Given the form of the pattern we can also simplify our fft calculation by
reverse summing the input digits.
"""
def fft_with_offset(signal, phases, offset):
    if offset < len(signal)//2:
        raise Exception("The logic requires an offset over half the length")

    signal_region_of_interest = signal[offset:]
    current = deque(signal_region_of_interest)
    for i in range(phases):
        sum = 0
        next = deque()
        for val in reversed(current):
            sum = (sum + val) % 10
            next.appendleft(sum)
        current = next
    return list(current)


with open("16/input.txt") as f:
    input_string = f.readline().strip()
input_signal = [int(c) for c in input_string]

part1_output = fft(input_signal, 100)
print(f"The answer to Part 1 is {''.join(map(str, part1_output[:8]))}.")

offset = int(input_string[:7])
real_input_signal = input_signal * 10000
part2_output = fft_with_offset(real_input_signal, 100, offset)
print(f"The answer to Part 2 is {''.join(map(str, part2_output[:8]))}.")