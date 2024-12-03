#include <algorithm>
#include <bitset>
#include <fstream>
#include <iostream>
#include <iterator>
#include <regex>
#include <string>
#include <unordered_map>
#include <vector>

class MaskV1 {
public:
    MaskV1(std::string mask_str)
    {
        if (mask_str.size() != 36)
        {
            std::cout << "Invalid input: " << mask_str << std::endl;
            exit(1);
        }
        to_or = 0;
        to_and = (1UL << 37) - 1; // We want 36 1s
        // For every 0 in the mask we need to make that bit 0 when applying the mask. Do this via a bitwise AND with 0.
        // For every 1 in the mask we need to make that bit 1 when applying the mask. Do this via a bitwise OR with 1.
        for (auto i = 0UL; i < 36; ++i)
        {
            switch (mask_str[35 - i]) // Iterate backwards, starting with the least significant bit
            {
                case '0':
                    to_and = to_and & ~(1UL << i); // set the ith bit to 0
                    break;
                case '1':
                    to_or = to_or | (1UL << i); // set the ith bit to 1
                    break;
                case 'X':
                    // Nothing to do.
                    break;
            }
        }
    }

    ulong apply(ulong value)
    {
        value |= to_or;
        value &= to_and;
        return value;
    }

private:
    ulong to_or;
    ulong to_and;
};

class MaskV2 {
public:
    MaskV2(std::string mask_str)
    {
        if (mask_str.size() != 36)
        {
            std::cout << "Invalid input: " << mask_str << std::endl;
            exit(1);
        }
        to_or = 0;
        // 0 and 1 in the mask now have no effect and set the corresponding bit to 1, respectively. We can accomplish
        // this using binary OR with 0 and 1, respectively.
        // x in the mask is now a floating bit which can be either 0 or 1.
        for (auto i = 0UL; i < 36; ++i)
        {
            switch (mask_str[35 - i]) // Iterate backwards, starting with the least significant bit
            {
                case '0':
                    // Nothing to do, all bits in |to_or| are initially 0.
                    break;
                case '1':
                    to_or = to_or | (1UL << i); // set the ith bit to 1
                    break;
                case 'X':
                    floating_bit_positions.push_back(i);
                    break;
            }
        }
    }

    std::vector<ulong> apply(ulong value)
    {
        // Apply the 0s and 1s from the mask.
        value |= to_or;

        const auto initial_bitset = std::bitset<64>{value}; // bitset gives easy access to individual bits

        // Every floating bit can be either 1 or 0 and we need to consider all possible values. There's 2^n
        // combinations where n is the number of floating bits.
        auto number_of_masked_values = 1UL << floating_bit_positions.size();

        // Iterate through all possible floating bit combinations.
        std::vector<ulong> masked_values;
        for (auto i = 0UL; i < number_of_masked_values; ++i)
        {
            auto bitset = initial_bitset;
            for (auto key = 0UL; key < floating_bit_positions.size(); ++key)
            {
                auto floating_bit_position = floating_bit_positions[key];
                auto floating_bit_value = (i & (1UL << key)) == 0 ? 0 : 1;
                bitset[floating_bit_position] = floating_bit_value;
            }
            masked_values.push_back(bitset.to_ulong());
        }
        return masked_values;
    }

private:
    ulong to_or;
    std::vector<ulong> floating_bit_positions;
};

// Each line of the input either sets the mask or writes to memory.
struct InputData
{
    InputData(std::string mask) : mask_str{mask}, is_mask{true} {}
    InputData(ulong address, ulong value) : mem_command{address, value}, is_mask{false} {}

    const bool is_mask;
    const std::string mask_str;
    // The problem involves 36 bit unsigned integers.
    const std::pair<ulong, ulong> mem_command;
};

std::vector<InputData> read_input()
{
    std::ifstream input("input.txt");

    std::vector<InputData> data;

    std::string line;
    while (std::getline(input, line))
    {
        if (line.rfind("mask", 0) == 0) // startswith
        {
           data.emplace_back(line.substr(7));
        }
        else
        {
            static auto mem_regex = std::regex{R"(^mem\[(\d+)\] = (\d+)$)"};
            std::smatch match;
            std::regex_match(line, match, mem_regex);
            auto matches = std::next(match.cbegin(), 1); // Skip the full string match
            auto memory_address = std::stoul(*matches);
            ++matches;
            auto value_to_set = std::stoul(*matches);
            data.emplace_back(memory_address, value_to_set);
        }
    }

    return data;
}

void part1(const std::vector<InputData>& data)
{
    std::unordered_map<ulong, ulong> memory;

    MaskV1 current_mask{"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}; // Start with an empty mask
    for (auto d : data)
    {
        if (d.is_mask)
        {
            current_mask = MaskV1(d.mask_str);
        }
        else
        {
            auto [address, value] = d.mem_command;
            memory[address] = current_mask.apply(value);
        }
    }

    auto sum = 0UL;
    std::for_each(memory.cbegin(), memory.cend(), [&sum](const auto& pair) { sum += pair.second; });
    std::cout << "The answer to Part 1 is " << sum << "." << std::endl;
}

void part2(const std::vector<InputData>& data)
{
    std::unordered_map<ulong, ulong> memory;

    MaskV2 current_mask{"000000000000000000000000000000000000"}; // Start with an empty mask
    for (auto d : data)
    {
        if (d.is_mask)
        {
            current_mask = MaskV2(d.mask_str);
        }
        else
        {
            auto [address, value] = d.mem_command;
            auto addresses = current_mask.apply(address);
            for (auto a : addresses)
            {
                memory[a] = value;
            }
        }
    }

    auto sum = 0UL;
    std::for_each(memory.cbegin(), memory.cend(), [&sum](const auto& pair) { sum += pair.second; });
    std::cout << "The answer to Part 2 is " << sum << "." << std::endl;
}

int main()
{
    const auto data = read_input();
    part1(data);
    part2(data);
}
