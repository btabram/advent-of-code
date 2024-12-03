from collections import Counter
from copy import deepcopy
from math import ceil


class Reaction:
    def __init__(self, reaction_string):
        # Map of chemical -> amount
        self.inputs = Counter()
        parts = reaction_string.split(" ")
        parts.reverse()
        while True:
            amount = int(parts.pop())
            chemical = parts.pop().replace(",", "")
            self.inputs[chemical] = amount
            if parts[-1] == "=>":
                parts.pop()
                break
        output_amount = int(parts.pop())
        output_chemical = parts.pop()
        # Map of chemical -> amount
        self.outputs = Counter({ output_chemical: output_amount })
        assert len(parts) == 0

    def __repr__(self):
        repr = ", ".join(map(str, self.inputs.items()))
        repr += " => "
        repr += ", ".join(map(str, self.outputs.items()))
        return repr

    @property
    def input_chems(self):
        return list(self.inputs.keys())

    @property
    def output_chems(self):
        return list(self.outputs.keys())

    def substitue_input(self, other_reaction):
        assert len(other_reaction.outputs) == 1
        chem_to_sub = other_reaction.output_chems[0]
        assert chem_to_sub in self.input_chems

        our_amount = self.inputs[chem_to_sub]
        other_amount = other_reaction.outputs[chem_to_sub]

        # Work out how many multiples of |other_reaction| we need to substitute
        # all of |chem_to_sub| in our reaction.
        count_of_other_reaction = ceil(our_amount / other_amount)

        amount_made = count_of_other_reaction * other_amount
        excess = amount_made - our_amount
        assert excess >= 0

        # Substitue |chem_to_sub| using |other_reaction|.
        del self.inputs[chem_to_sub]
        for chem, amount in other_reaction.inputs.items():
            self.inputs[chem] += count_of_other_reaction * amount
        # Any excess |chem_to_sub| adds to our output.
        if excess != 0:
            self.outputs[chem_to_sub] += excess

        # Cancel anything that appears in both input and output.
        inputs_set = set(self.inputs.keys())
        outputs_set = set(self.outputs.keys())
        intersection = inputs_set & outputs_set
        for chem in intersection:
            input_amount = self.inputs[chem]
            output_amount = self.outputs[chem]
            if input_amount == output_amount:
                del self.inputs[chem]
                del self.outputs[chem]
            elif input_amount > output_amount:
                self.inputs[chem] -= output_amount
                del self.outputs[chem]
            else: # input_amount < output_amount
                del self.inputs[chem]
                self.outputs[chem] -= input_amount


def calculate_ore_requirement(reactions_original, fuel_quantity):
    reactions = deepcopy(reactions_original)
    fuel_reactions = [r for r in reactions if "FUEL" in r.output_chems]
    assert len(fuel_reactions) == 1
    fuel_reaction = fuel_reactions[0]
    assert len(fuel_reaction.outputs) == 1
    reactions.remove(fuel_reaction)

    # Multiply up the fuel reaction to make the desired quantity.
    assert fuel_reaction.outputs["FUEL"] == 1
    for c in fuel_reaction.input_chems:
        fuel_reaction.inputs[c] *= fuel_quantity
    for c in fuel_reaction.output_chems:
        fuel_reaction.outputs[c] *= fuel_quantity

    while True:
        chemicals_to_substitute = \
            [c for c in fuel_reaction.input_chems if c != "ORE"]
        if len(chemicals_to_substitute) == 0:
            break

        chem_to_sub = chemicals_to_substitute[0]
        ways_to_generate = [r for r in reactions if chem_to_sub in r.output_chems]
        assert len(ways_to_generate) == 1
        to_generate = ways_to_generate[0]

        fuel_reaction.substitue_input(to_generate)

    return fuel_reaction.inputs["ORE"]


reactions = []
with open("14/input.txt") as f:
    for line in f.readlines():
        reactions.append(Reaction(line.strip()))

ore_for_one_fuel = calculate_ore_requirement(reactions, 1)
print(f"The answer to Part 1 is {ore_for_one_fuel}.")

# For Part 2 we need to calculate how much fuel we can make with
# 1000000000000 (one trillion) ore. We can use our answer for Part 1 to
# calculate a lower bound. The actual answer is bigger because excess output
# from the one fuel equation can be used to make more fuel when the equation is
# multiplied up.
TRILLION = 1000000000000
lower_bound = TRILLION // ore_for_one_fuel
# Choose an upper bound.
upper_bound = 2 * lower_bound
upper = calculate_ore_requirement(reactions, upper_bound)
while upper <= TRILLION:
    upper_bound += lower_bound
    upper = calculate_ore_requirement(reactions, upper_bound)

# Binary search.
while True:
    if upper_bound - lower_bound == 1:
        break
    middle_point = (upper_bound + lower_bound) // 2
    middle = calculate_ore_requirement(reactions, middle_point)
    if middle < TRILLION:
        lower_bound = middle_point
    elif middle > TRILLION:
        upper_bound = middle_point
    else:
        break
print(f"The answer to Part 2 is {lower_bound}.")
