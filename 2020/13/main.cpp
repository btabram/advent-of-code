#include <algorithm>
#include <fstream>
#include <iostream>
#include <optional>
#include <sstream>
#include <string>
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

std::pair<int, std::vector<int>> read_input()
{
    std::ifstream input("input.txt");
    std::string line;

    std::getline(input, line);
    auto timestamp = std::stoi(line);

    std::getline(input, line);
    auto raw_buses = split_on_comma(line);
    std::vector<int> buses;
    std::transform(raw_buses.cbegin(), raw_buses.cend(), std::back_inserter(buses),
                                [](const auto& bus_str) { return bus_str == "x" ? 0 : std::stoi(bus_str); });

    return { timestamp, buses };
}

int main()
{
    const auto [timestamp, buses] = read_input();

    // Part 1.
    for (int t = timestamp; ; ++t)
    {
        std::optional<int> departing_bus;
        for (auto b : buses)
        {
            if (b != 0 && (t % b) == 0)
            {
                departing_bus = b;
                break;
            }
        }
        if (departing_bus.has_value())
        {
            std::cout << "The answer to Part 1 is " << (t - timestamp) * *departing_bus << "." << std::endl;
            break;
        }
    }

    // Part 2.
    // Vector of pairs of bus ID and the time offset of that bus.
    std::vector<std::pair<int, int>> buses_offsets;
    for (int i = 0; i < buses.size(); ++i)
    {
        if (buses[i] != 0)
        {
            buses_offsets.emplace_back(buses[i], i);
        }
    }

    // To make the problem manageable we find the time till the first alignment for a subset of the buses and work out
    // the repeat period of that subset of buses. Then to add in the next bus we can quickly iterate forwards in time,
    // jumping from one time where the already solved subset of buses all align to the next.
    long partial_period = 1;
    long first_alignment = 1;
    std::vector<std::pair<int, int>> partial_buses_offsets;
    for (auto [bus, offset] : buses_offsets)
    {
        partial_buses_offsets.emplace_back(bus, offset);
        for (long t = first_alignment; ; t += partial_period)
        {
            if (std::all_of(partial_buses_offsets.cbegin(), partial_buses_offsets.cend(),
                            [&t](const auto& pair) { return ((t + pair.second) % pair.first) == 0; }))
            {
                // We've solved the problem for the new subset and can now update the outer variables.
                first_alignment = t;
                partial_period *= bus; // The numbers are mutually coprime so multiplying them together gives the LCM
                break;
            }
        }
    }
    std::cout << "The answer to Part 2 is " << first_alignment << "." << std::endl;
}
