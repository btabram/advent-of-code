#include <fstream>
#include <iostream>
#include <map>
#include <vector>

struct Cups
{
    Cups(int current, std::map<int, int> map) : current(current), map(map) {}

    int current;
    std::map<int, int> map;
};

Cups read_input()
{
    std::ifstream input("input.txt");
    std::string line;
    std::getline(input, line);

    std::vector<int> cups_vector;
    for (auto digit : line)
    {
        cups_vector.push_back(static_cast<int>(digit) - static_cast<int>('0'));
    }

    // Create a map where the keys are cup labels and each value is label of the next cup.
    std::map<int, int> cups_map;
    int current = cups_vector[0], next;
    for (auto i = 1; i < cups_vector.size(); i++)
    {
        next = cups_vector[i];
        cups_map[current] = next;
        current = next;
    }
    cups_map[current] = cups_vector[0];

    return Cups{ cups_vector[0], cups_map };
}

// Storing the cup positions in a map of cup -> next cup makes this function O(log N), where N is the number of cups.
// A more naive initial implementation with std::deque was O(N). Good scaling is essential for Part 2.
void do_move(Cups& cups)
{
    // We remove the three cups immediately after the current cup.
    int removed_1 = cups.map.at(cups.current);
    int removed_2 = cups.map.at(removed_1);
    int removed_3 = cups.map.at(removed_2);

    // After removing the three cups above |next| is now the cup immediately after the current cup.
    int next = cups.map.at(removed_3);
    cups.map[cups.current] = next; // Join the circle back up

    // Find the destination cup label.
    int destination = cups.current - 1;
    while (true)
    {
        if (destination == 0) destination = cups.map.size(); // Wrap around

        if (destination == removed_1 || destination == removed_2 || destination == removed_3)
        {
            destination--;
        }
        else
        {
            break;
        }
    }

    // Find the cup immediately after the destination.
    int after_destination = cups.map.at(destination);

    // Insert the three cups we removed between |destination| and |after_destination| in the circle.
    cups.map[destination] = removed_1;
    // |removed_1| still points to |removed_2| so don't modfy it.
    // |removed_2| still points to |removed_3| so don't modfy it.
    cups.map[removed_3] = after_destination;

    // Advance the current cup.
    cups.current = next;
}

void part1(const Cups& initial_cups)
{
    auto cups = initial_cups;
    for (auto i = 0; i < 100; i++)
    {
        do_move(cups);
    }

    // The answer is the sequence of cups after 100 moves, starting at the cup labelled '1'.
    std::string answer;
    auto pos = 1;
    for (auto i = 0; i < 8; i++)
    {
        pos = cups.map[pos];
        answer += std::to_string(pos);
    }
    std::cout << "The answer to Part 1 is " << answer << "." << std::endl;
}

void part2(const Cups& initial_cups)
{
    auto cups = initial_cups;

    int first_initial_cup = initial_cups.current; // Initially the current cup is the first in the circle

    int last_initial_cup = first_initial_cup;
    for (auto i = 0; i < 8; i++) last_initial_cup = cups.map[last_initial_cup];

    // For Part 2 we add cups to our circle till we have a million.
    cups.map[last_initial_cup] = 10;
    for (auto i = 10; i < 1000000; i++)
    {
        cups.map[i] = i + 1;
    }
    cups.map[1000000] = first_initial_cup;

    // Then we do 10 million moves.
    for (auto i = 0; i < 10000000; i++)
    {
        do_move(cups);
    }

    // The answer is the product of the two cups immediately after the cup labelled '1'.
    long answer = 1;
    auto pos = 1;
    for (auto i = 0; i < 2; i++)
    {
        pos = cups.map[pos];
        answer *= pos;
    }
    std::cout << "The answer to Part 2 is " << answer << "." << std::endl;
}

int main()
{
    const auto initial_cups = read_input();
    part1(initial_cups);
    part2(initial_cups);
}
