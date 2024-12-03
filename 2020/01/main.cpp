#include <algorithm>
#include <fstream>
#include <iostream>
#include <iterator>
#include <sstream>
#include <string>
#include <unordered_map>
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

    return data;
}

void part1(const std::vector<int>& numbers)
{
    for (int a : numbers) 
    {
        for (int b : numbers) 
        {
            if (a + b == 2020)
            {
                std::cout << "The answer to Part 1 is " << a * b << "." << std::endl;
                return;
            }
        }
    }
}

void part2(const std::vector<int>& numbers)
{
    for (int a : numbers) 
    {
        for (int b : numbers) 
        {
            for (int c : numbers) 
            {
                if (a + b + c == 2020)
                {
                    std::cout << "The answer to Part 2 is " << a * b * c << "." << std::endl;
                    return;
                }
            }
        }
    }
}

int main()
{
    const auto numbers = read_input();
    part1(numbers);
    part2(numbers);
}
