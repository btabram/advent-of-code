with open("01/input.txt") as f:
    data = f.readlines()

masses = list(map(int, data))

def get_naive_fuel_requirement(mass):
    return (mass // 3) - 2

naive_total_fuel = sum(map(get_naive_fuel_requirement, masses))

print("The fuel requirement for Part 1 is {}".format(naive_total_fuel))

def get_fuel_requirement(mass):
    new_fuel = get_naive_fuel_requirement(mass)
    total_fuel = new_fuel
    while True:
        new_fuel = get_naive_fuel_requirement(new_fuel)
        if new_fuel <= 0:
            break
        total_fuel += new_fuel
    return total_fuel

total_fuel = sum(map(get_fuel_requirement, masses))

print("The fuel requirement for Part 2 is {}".format(total_fuel))
