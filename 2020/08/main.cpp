#include <fstream>
#include <iostream>
#include <string>
#include <unordered_set>
#include <vector>

using Instruction = std::pair<std::string, int>;

class HandheldGameConsole {
public:
    std::pair<bool, int> run(const std::vector<Instruction>& instructions)
    {
        accumulator = 0;
        position = 0;
        std::unordered_set<int> visited;

        while (position != instructions.size())
        {
            auto [_, success] = visited.insert(position);
            // If the insertion wasn't a success then it's becuse |position| was already in the set.
            if (!success)
            {
                // Infinite loop detected - bail out!
                return { false, accumulator };
            }

            auto& [command, argument] = instructions[position];
            if (command == "acc")
            {
                acc(argument);
            }
            else if (command == "jmp")
            {
                jmp(argument);
            }
            else
            {
                nop(argument);
            }
        }

        return { true, accumulator };
    }

private:
    void acc(int argument)
    {
        accumulator += argument;
        position += 1;
    }

    void jmp(int argument)
    {
        position += argument;
    }

    void nop(int argument)
    {
        position += 1;
    }

    int accumulator = 0;
    int position = 0;
};

std::vector<Instruction> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Instruction> instructions;
    std::string line;
    while (std::getline(input, line))
    {
        auto operation = line.substr(0, 3);
        auto argument = std::stoi(line.substr(3));
        instructions.emplace_back(operation, argument);
    }

    return instructions;
}

int part1(const std::vector<Instruction>& instructions)
{
    HandheldGameConsole console;
    auto [_, accumulator] = console.run(instructions);
    return accumulator;
}

int part2(const std::vector<Instruction>& instructions)
{
    HandheldGameConsole console;

    std::vector<int> jmps, nops;
    for (int i = 0; i < instructions.size(); ++i)
    {
        if (instructions[i].first == "jmp")
        {
            jmps.push_back(i);
        }
        else if (instructions[i].first == "nop")
        {
            nops.push_back(i);
        }
    }

    // Try swapping every possible jmp for a nop and vice versa until we can run without an infinite loop.
    for (auto jmp_pos : jmps)
    {
        std::vector<Instruction> instructions_copy = instructions;
        instructions_copy[jmp_pos].first = "nop";
        auto [success, accumulator] = console.run(instructions_copy);
        if (success)
        {
            return accumulator;
        }
    }
    for (auto nop_pos : nops)
    {
        std::vector<Instruction> instructions_copy = instructions;
        instructions_copy[nop_pos].first = "jmp";
        auto [success, accumulator] = console.run(instructions_copy);
        if (success)
        {
            return accumulator;
        }
    }

    std::cout << "Part 2 failed" << std::endl;
    exit(1);
}

int main()
{
    const auto instructions = read_input();
    std::cout << "The answer to Part 1 is " << part1(instructions) << "." << std::endl;
    std::cout << "The answer to Part 2 is " << part2(instructions) << "." << std::endl;
}
