#include <algorithm>
#include <deque>
#include <fstream>
#include <iostream>
#include <numeric>
#include <string>
#include <vector>

std::vector<long> read_input()
{
    std::ifstream input("input.txt");

    std::vector<long> data;
    std::string line;
    while (std::getline(input, line))
    {
        // Use long because some of the input data is too big for int.
        data.push_back(std::stol(line));
    }

    return data;
}

long find_first_invalid(const std::vector<long>& data)
{
    std::deque<long> prev_values;

    // Read in the preamble.
    int i = 0;
    for (; i < 25; ++i)
    {
        prev_values.push_back(data[i]);
    }

    for (; i < data.size(); ++i)
    {
        auto value = data[i];
        auto valid = false;
        // Try to find two previous values which sum to |value|.
        for (auto&& prev_value : prev_values)
        {
            auto complement = value - prev_value;
            if (std::find(prev_values.cbegin(), prev_values.cend(), complement) != prev_values.cend())
            {
                valid = true;
                break;
            }
        }

        if (!valid)
        {
            return value;
        }

        prev_values.push_back(value);
        prev_values.pop_front();
    }

    std::cout << "Failed to find an invalid value!" << std::endl;
    exit(1);
}

std::deque<long> find_contiguous_sum_block(const std::vector<long>& data, long target)
{
    std::deque<long> current_contiguous_block;

    for (auto&& value : data)
    {
        current_contiguous_block.push_back(value);
        auto sum = std::reduce(current_contiguous_block.cbegin(), current_contiguous_block.cend(), 0);

        // Shrink our contiguous block if we're overshooting.
        while (sum > target)
        {
            current_contiguous_block.pop_front();
            sum = std::reduce(current_contiguous_block.cbegin(), current_contiguous_block.cend(), 0);
        }

        if (sum == target)
        {
            return current_contiguous_block;
        }
    }

    std::cout << "Failed to find a contiguous sum!" << std::endl;
    exit(1);
}

int main()
{
    const auto data = read_input();
    auto invalid = find_first_invalid(data);
    std::cout << "The answer to Part 1 is " << invalid << "." << std::endl;

    auto block = find_contiguous_sum_block(data, invalid);
    auto block_min = *(std::min_element(block.cbegin(), block.cend()));
    auto block_max = *(std::max_element(block.cbegin(), block.cend()));
    std::cout << "The answer to Part 2 is " << (block_min + block_max) << "." << std::endl;
}
