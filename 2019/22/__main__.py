from functools import partial
from random import randint


# Multiplcation under modulo is straightfoward but division is more complicated. We can define the modular
# multiplicative inverse x for integer a such that a * x = 1 (mod m) and use multiplication by it as an equivalent to
# division. Fortunately Python 3.8+ makes it very easy to calculate the modular multiplicative inverse.
def modmulinv(value, modulo):
    return pow(value, -1, modulo)


def deal_into_new_stack(position, deck_size):
    return (deck_size - 1) - position


undo_deal_into_new_stack = deal_into_new_stack


def cut_n_cards(position, deck_size, n):
    return ((position - n) + deck_size) % deck_size


def undo_cut_n_cards(position, deck_size, n):
    return ((position + n) + deck_size) % deck_size


def deal_with_increment_n(position, deck_size, n):
    return (position * n) % deck_size


def undo_deal_with_increment_n(position, deck_size, n):
    return (position * modmulinv(n, deck_size)) % deck_size


def parse_shuffles(string, undo=False):
    parts = string.split(" ")
    first = parts[0]
    last = parts[-1]

    if first == "cut":
        return partial(undo_cut_n_cards if undo else cut_n_cards, n = int(last))

    assert first == "deal"

    if last == "stack":
        return undo_deal_into_new_stack if undo else deal_into_new_stack
    else:
        return partial(undo_deal_with_increment_n if undo else deal_with_increment_n, n = int(last))


with open("22/input.txt") as f:
    input = f.read().splitlines()

shuffles = list(map(parse_shuffles, input))
undo_shuffles = list(map(partial(parse_shuffles, undo=True), reversed(input)))


# Part 1 - find the final position of card #2019 in a 10007 card deck after applying all the shuffles.
position = 2019
for shuffle in shuffles:
    position = shuffle(position, 10007)
print("The answer to Part 1 is {}.".format(position))


# Part 2 - find the number (AKA starting position) of the card in position 2020 after a 119315717514047 card deck has
# had all the shuffles applied 101741582076661 times in a row. Thanks to etotheipi1 on Reddit for explaining this one!
deck_size = 119315717514047
times_shuffled = 101741582076661
def undo_one_round_of_shuffles(position):
    for undo_shuffle in undo_shuffles:
        position = undo_shuffle(position, deck_size)
    return position

# Each shuffle is a linear function of position and so the composition of many shuffles will also be a linear function
# of position. Therefore we can express |undo_one_round_of_shuffles| as f(pos) = A*pos + B (mod deck_size). Now to find
# A and B:
f_zero = undo_one_round_of_shuffles(0)
f_one = undo_one_round_of_shuffles(1)
A = (f_one - f_zero) % deck_size
B = f_zero

test = randint(0, deck_size - 1)
assert undo_one_round_of_shuffles(test) == (A*test + B) % deck_size

# We need to undo lots of rounds of shuffling so we need to repeatedly apply our linear function. In general:
#    f(f(f(... n times ...f(x)))) = A^n*x + A^(n-1)*B + A^(n-2)*B + ... + B
#                                 = A^n*x + (A^(n-1) + A^(n-2) + ... + 1)*B
#                                 = A^n*x + ((A^n - 1) / (A - 1))*B
# where we've recognised a geometric series and subsituted its sum in the final step. Remembering everything is modulo
# the deck size in our case we can therefore write:
def undo_all_shuffles(position):
    A_to_n = pow(A, times_shuffled, deck_size)
    return ((A_to_n * position) + ((A_to_n - 1) * modmulinv(A - 1, deck_size) * B)) % deck_size

print("The answer to Part 2 is {}.".format(undo_all_shuffles(2020)))