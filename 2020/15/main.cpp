#include <algorithm>
#include <deque>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>

std::vector<std::string> split_on_comma(std::string string)
{
    std::vector<std::string> output;
    std::string part;
    auto ss = std::istringstream{string};
    while (std::getline(ss, part, ','))
    {
        output.push_back(part);
    }
    return output;
}

std::vector<int> read_input()
{
    std::ifstream input("input.txt");

    std::string line;
    std::getline(input, line);

    std::vector<int> numbers;
    auto number_strings = split_on_comma(line);
    std::transform(number_strings.cbegin(), number_strings.cend(), std::back_inserter(numbers),
                   [](const auto& str) { return std::stoi(str); });

    return numbers;
}

int main()
{
    const auto starting_numbers = read_input();

    std::unordered_map<int, std::deque<int>> said_positions_map;

    int i = 1;
    for (auto starting_number : starting_numbers)
    {
        said_positions_map[starting_number].push_back(i);
        ++i;
    }

    int last_number = starting_numbers[starting_numbers.size() - 1];

    while (true)
    {
        int number_to_say;

        const auto& last_said_positions = said_positions_map.find(last_number)->second;
        if (last_said_positions.size() == 1)
        {
            // Last number hadn't been said before so we say 0.
            number_to_say = 0;
        }
        else
        {
            // Last number had been said before so we say the difference between the turn when it was most recently
            // spoken (last turn) and the time before that.
            number_to_say = last_said_positions[1] - last_said_positions[0];
        }

        // Say the number and update things accordingly.
        auto& to_say_positions = said_positions_map[number_to_say];
        to_say_positions.push_back(i);
        last_number = number_to_say;

        // We only need to store the most recent position and the one before that.
        while (to_say_positions.size() > 2)
        {
            to_say_positions.pop_front();
        }

        if (i == 2020)
        {
            std::cout << "The answer to Part 1 is " << number_to_say << "." << std::endl;
        }
        else if (i == 30000000)
        {
            std::cout << "The answer to Part 2 is " << number_to_say << "." << std::endl;
            break;
        }

        ++i;
    }
}
