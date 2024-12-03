import itertools


class Body:
    def __init__(self, name):
        self.name = name
        self.orbitee = None
        self.jumps_to_com = None

    def __repr__(self):
        return f"Body({self.name}, {self.jumps_to_com})"

    def addInfo(self, orbitee, jumps_to_com):
        assert self.orbitee is None
        self.orbitee = orbitee
        self.jumps_to_com = jumps_to_com


with open("06/input.txt") as f:
    input = f.readlines()

# Strip whitespace and validate input.
input = [line.strip() for line in input]
for line in input:
    assert len(line) == 7 and line[3] == ")"

bodies = set(map(Body, [line[4:] for line in input]))
bodies_dict = dict()
for body in bodies:
    bodies_dict[body.name] = body
# Every body should orbit exactly one other.
assert len(bodies) == len(input)

# Dict of bodies -> their orbiters.
# Note that a body can have more than one other body orbiting it.
orbiters_dict = dict()
for line in input:
    [body, orbiter] = line.split(")") 
    if body not in orbiters_dict:
        orbiters_dict[body] = [orbiter]
    else:
        orbiters_dict[body].append(orbiter)

centre_of_mass = Body("COM")
bodies_dict["COM"] = centre_of_mass

# Work through all bodies, moving away from the COM, and assign orbitees.
bodies_at_prev_distance = ["COM"]
for jumps_to_com in itertools.count(1):
    bodies_at_this_distance = []
    for body_name in bodies_at_prev_distance:
        if body_name not in orbiters_dict:
            # Body |body_name| has no orbiters.
            continue
        body = bodies_dict[body_name]
        for orbiter_name in orbiters_dict[body_name]:
            orbiter = bodies_dict[orbiter_name]
            orbiter.addInfo(body, jumps_to_com)
            bodies_at_this_distance.append(orbiter_name)
    bodies_at_prev_distance = bodies_at_this_distance
    
    # We're finished when every body has an orbitee.
    if len([body for body in bodies if body.orbitee is None]) == 0:
        break

orbit_count_checksum = sum([body.jumps_to_com for body in bodies])
print(f"The answer for Part 1 is {orbit_count_checksum}.")

# Every body has only one orbitee so there's only one path (without repetition)
# between any two given bodies. We can easily work out the paths for me and
# santa to get back to COM. At some point, these paths will overlap, giving us
# a path from me to santa.
santa = bodies_dict["SAN"]
me = bodies_dict["YOU"]

def path_back_to_com(start, path):
    prev = start
    while True:
        prev = prev.orbitee
        if prev.name == "COM":
            break
        path.append(prev)

my_path_to_com = []
path_back_to_com(me, my_path_to_com)
santas_path_to_com = []
path_back_to_com(santa, santas_path_to_com)
common_path = set(my_path_to_com).intersection(set(santas_path_to_com))

# The route from me to santa is |me.orbitee| -> the points where our paths to
# COM join -> |santa.orbitee|.
max_common_path_jumps_to_com = max([b.jumps_to_com for b in common_path])
path_length = (me.orbitee.jumps_to_com - max_common_path_jumps_to_com) \
                + (santa.orbitee.jumps_to_com - max_common_path_jumps_to_com)
print(f"The answer for Part 2 is {path_length}.")
