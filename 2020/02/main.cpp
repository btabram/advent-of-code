#include <algorithm>
#include <fstream>
#include <iostream>
#include <regex>
#include <string>
#include <vector>

struct PasswordRule {
    PasswordRule(int first, int second, char letter) : first(first), second(second), letter(letter) {}

    bool is_valid_part1(const std::string& password) const
    {
        auto count = std::count_if(password.cbegin(), password.cend(), [this](char c) { return c == letter; });
        // In Part 1 the two numbers are the min and max occurances of |letter|.
        return count >= first && count <= second;
    }

    bool is_valid_part2(const std::string& password) const
    {
        auto letter_at_first = password[first - 1] == letter;
        auto letter_at_second = password[second - 1] == letter;
        // In Part 2 the numbers are indicies (starting at 1) and a valid password has |letter| at exactly one of the
        // two positions.
        return letter_at_first != letter_at_second;
    }

    const char letter;
    const int first;
    const int second;
};

std::vector<std::pair<PasswordRule, std::string>> read_input()
{
    std::ifstream input("input.txt");

    std::vector<std::pair<PasswordRule, std::string>> data;

    static auto re = std::regex{R"(^([0-9]+)-([0-9]+) ([a-z]): ([a-z]+)$)"};
    std::string line;
    while (std::getline(input, line))
    {
        std::smatch match;
        std::regex_match(line, match, re);
        auto rule = PasswordRule(std::stoi(match[1]), std::stoi(match[2]), std::string{match[3]}[0]);
        data.emplace_back(rule, match[4]);
    }

    return data;
}

int main()
{
    const auto data = read_input();

    auto part1 = std::count_if(data.cbegin(), data.cend(),
                               [](const auto& pair) { return pair.first.is_valid_part1(pair.second); });
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;

    auto part2 = std::count_if(data.cbegin(), data.cend(),
                               [](const auto& pair) { return pair.first.is_valid_part2(pair.second); });
    std::cout << "The answer to Part 2 is " << part2 << "." << std::endl;
}
