def get_required_paper(x: int, y: int, z: int) -> int:
    face_areas = [x * y, y * z, z * x]
    surface_area = 2 * sum(face_areas)
    return surface_area + min(face_areas)


def get_required_ribbon(x: int, y: int, z: int) -> int:
    perimeters = map(lambda a: 2 * a, [x + y, y + z, z + x])
    return min(perimeters) + (x * y * z)


with open("input.txt") as f:
    dimensions = [[int(c) for c in line.strip().split("x")] for line in f.readlines()]

required_paper = sum(map(lambda x: get_required_paper(*x), dimensions))
required_ribbon = sum(map(lambda x: get_required_ribbon(*x), dimensions))

print(f"The answer to Part 1 is {required_paper}.")
print(f"The answer to Part 2 is {required_ribbon}.")
