from hashlib import md5


with open("input.txt") as f:
    key = f.readline().strip()

i = 1
part1 = None
part2 = None
while True:
    input = (key + str(i)).encode()
    hex_hash = md5(input).hexdigest()
    if not part1 and hex_hash.startswith("00000"):
        part1 = i
    if not part2 and hex_hash.startswith("000000"):
        part2 = i
    if part1 and part2:
        break
    i += 1

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
