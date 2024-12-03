#include <algorithm>
#include <fstream>
#include <iostream>
#include <numeric>
#include <string>
#include <set>
#include <vector>

// Use ordered sets because set_intersection requires sorted ranges.
using Declarations = std::set<char>;
using GroupDeclarations = std::vector<std::set<char>>;

std::vector<GroupDeclarations> read_input()
{
    std::ifstream input("input.txt");

    std::vector<GroupDeclarations> data;
    GroupDeclarations group_declarations;
    std::string line;
    while (std::getline(input, line))
    {
        // A blank line signals the start of a new group.
        if (line == "")
        {
            data.push_back(group_declarations);
            group_declarations.clear();
            continue;
        }
        Declarations declarations;
        std::for_each(line.cbegin(), line.cend(), [&declarations](char c) { declarations.insert(c); });
        group_declarations.push_back(declarations);
    }
    if (group_declarations.size() > 0)
    {
        data.push_back(group_declarations);
    }

    return data;
}

int main()
{
    const auto data = read_input();

    int part1 = 0;
    // For Part 1 we're interested in the number of questions anyone in a given group said yes to.
    for (auto&& group : data)
    {
        auto group_questions = std::reduce(++group.cbegin(), group.cend(), group[0],
            [](const Declarations& a, const Declarations& b) {
                Declarations u;
                std::set_union(a.begin(), a.end(), b.begin(), b.end(), std::inserter(u, u.begin()));
                return u;
            }
        );
        part1 += group_questions.size();
    }
    std::cout << "The answer to Part 1 " << part1 << "." << std::endl;

    int part2 = 0;
    // For Part 2 we're interested in the number of questions everyone in a given group said yes to.
    for (auto&& group : data)
    {
        auto group_questions = std::reduce(++group.cbegin(), group.cend(), group[0],
            [](const Declarations& a, const Declarations& b) {
                Declarations i;
                std::set_intersection(a.begin(), a.end(), b.begin(), b.end(), std::inserter(i, i.begin()));
                return i;
            }
        );
        part2 += group_questions.size();
    }
    std::cout << "The answer to Part 2 " << part2 << "." << std::endl;
}
