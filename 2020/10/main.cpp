#include <algorithm>
#include <fstream>
#include <iostream>
#include <map>
#include <string>
#include <vector>

std::vector<int> read_input()
{
    std::ifstream input("input.txt");

    std::vector<int> data;
    std::string line;
    while (std::getline(input, line))
    {
        data.push_back(std::stoi(line));
    }

    // Sort and add the wall and builtin adapter joltage.
    data.push_back(0);
    std::sort(data.begin(), data.end());
    data.push_back(data[data.size() - 1] + 3);

    return data;
}

int count_arrangements(const std::vector<int>& adapters)
{
    // The first and last adapters are fixed.
    if (adapters.size() < 3)
    {
        return 1;
    }
    const auto adapters_to_vary = adapters.size() - 2;
    // Every variable adapter can be present or absent in an arrangement. Note that we can't reorder the adapters.
    const auto number_of_arrangments = 1 << adapters_to_vary; // 2^adapters_to_vary

    const auto first = adapters[0];
    const auto last = adapters[adapters.size() - 1];

    auto valid_arrangements = 0;
    // Use a bit mask to iterate through all possible combinations of omitting / including different adapters.
    for (size_t mask = 0; mask < number_of_arrangments; ++mask)
    {
        std::vector<int> arrangement_to_try;
        arrangement_to_try.push_back(first);
        for (size_t i = 0; i < adapters_to_vary; ++i)
        {
            if ((1 << i) & mask) {
                arrangement_to_try.push_back(adapters[i + 1]);
            }
        }
        arrangement_to_try.push_back(last);

        // Work out if the arrangement is valid.
        auto valid = true;
        auto prev = arrangement_to_try[0];
        for (size_t i = 1; i < arrangement_to_try.size(); ++i)
        {
            auto value = arrangement_to_try[i];
            if ((value - prev) > 3)
            {
                valid = false;
                break;
            }
            prev = value;
        }

        if (valid)
        {
            ++valid_arrangements;
        }
    }

    return valid_arrangements;
}

int main()
{
    const auto adapters = read_input();

    // For Part 1.
    std::map<int, int> diff_counts;

    // For Part 2. Use long because the answer would cause an int to overflow.
    long total_arrangements = 1;

    int prev = adapters[0];
    size_t prev_3_diff_pos = 0;
    for (size_t i = 1; i < adapters.size(); ++i)
    {
        auto adapter = adapters[i];
        auto diff = adapter - prev;

        // Part 1. If a key isn't in the map then [] will default construct a value (which is 0 for an int).
        diff_counts[diff] += 1;

        // Part 2. We can only have different adapter arrangements if the gap bettween adjacent adapters is < 3 so we
        // can split the problem up into the number of arrangements in between joltage jumps of 3.
        if (diff ==3)
        {
            auto slice = std::vector<int>{adapters.begin() + prev_3_diff_pos, adapters.begin() + i};
            total_arrangements *= count_arrangements(slice);
            prev_3_diff_pos = i;
        }

        prev = adapter;
    }

    std::cout << "The answer to Part 1 is " << (diff_counts[1] * diff_counts[3]) << "." << std::endl;
    std::cout << "The answer to Part 2 is " << total_arrangements << "." << std::endl;
}
