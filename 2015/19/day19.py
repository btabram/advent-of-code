from random import shuffle


def getOneReplacementPossibilities(mol: str, resp: list[tuple[str, str]]) -> set[str]:
    new_distinct_molecules = set()
    for lhs, rhs in resp:
        i = 0
        while i < len(mol):
            match_index = mol.find(lhs, i)
            if match_index == -1:
                break

            new_molecule = mol[:match_index] + rhs + mol[match_index + len(lhs) :]
            new_distinct_molecules.add(new_molecule)

            i = match_index + len(lhs)
    return new_distinct_molecules


# Solve part 2 with a bit of greedy randomness (credit to "What-A-Baller" on the Reddit
# solutions thread). There is a pattern to the data but I didn't spot it. I initially
# tried A* from the medicine molecule back to "e" but the search space is just too big.
def countReplacementsForTargetMolecule(
    start: str, target: str, reps: list[tuple[str, str]]
) -> int:
    replacement_count = 0
    molecule = start
    while True:
        did_replacement = False
        for lhs, rhs in reps:
            if (c := molecule.count(lhs)) > 0:
                molecule = molecule.replace(lhs, rhs)
                replacement_count += c
                did_replacement = True

        if molecule == target:
            return replacement_count

        if not did_replacement:
            # This attempt didn't work, shuffle the list of replacements and try again.
            shuffle(reps)
            return countReplacementsForTargetMolecule(start, target, reps)


with open("input.txt") as f:
    input_lines = [line.strip() for line in f.readlines()]

replacements: list[tuple[str, str]] = []
for line in input_lines:
    if line == "":
        break
    replacements.append(tuple(line.split(" => ")))
reverse_replacements = [(b, a) for (a, b) in replacements]
medicine_molecule = input_lines[-1]

part1 = len(getOneReplacementPossibilities(medicine_molecule, replacements))
part2 = countReplacementsForTargetMolecule(medicine_molecule, "e", reverse_replacements)

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
