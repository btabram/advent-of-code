IMAGE_WIDTH = 25
IMAGE_HEIGHT = 6

with open("08/input.txt") as f:
    encoded_image = list(map(int, [c for c in f.readline().strip()]))

pixels_per_frame = IMAGE_WIDTH * IMAGE_HEIGHT
assert len(encoded_image) % pixels_per_frame == 0

frames = []
for (i, p) in enumerate(encoded_image):
    frame_number = i // pixels_per_frame;
    if len(frames) == frame_number:
        frames.append([])
    frames[frame_number].append(p)

min_zeroes_frame = \
    min(frames, key = lambda frame: sum([1 for pixel in frame if pixel == 0]))

ones = sum([1 for pixel in min_zeroes_frame if pixel == 1])
twos = sum([1 for pixel in min_zeroes_frame if pixel == 2])
print(f"The answer for Part 1 is {ones * twos}.")

decoded_image = []
for i in range(pixels_per_frame):
    # Iterate through frames from front to back, the overall value is the
    # first non-transparent value.
    for frame in frames:
        value = frame[i]
        if value == 2: # transparent
            continue
        else:
            decoded_image.append(value)
            break

print("The decoded image for Part 2 is:")
for i in range(IMAGE_HEIGHT):
    line = ""
    for j in range(IMAGE_WIDTH):
        pixel = decoded_image[(i * IMAGE_WIDTH) + j]
        line += "#" if pixel == 1 else " "
    print(line)
