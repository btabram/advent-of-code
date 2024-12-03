#include <fstream>
#include <iostream>
#include <string>
#include <vector>

using Map = std::vector<std::vector<char>>;

Map read_input()
{
    std::ifstream input("input.txt");

    Map grid;
    std::string line;
    while (std::getline(input, line))
    {
        grid.emplace_back(line.cbegin(), line.cend());
    }

    return grid;
}

int count_trees(const Map& terrain, int right_step, int down_step)
{
    const auto height = terrain.size();
    const auto width = terrain[0].size();

    auto count = 0;
    for (int y = 0; y < height; y += down_step)
    {
        auto x = (y * right_step / down_step) % width; // The terrain repeats horizontally
        if (terrain[y][x] == '#')
        {
            ++count;
        }
    }

    return count;
}

int main()
{
    const auto terrain = read_input();

    auto part1 = count_trees(terrain, 3, 1);
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;

    long part2 = part1; // The answer is greater than 2^31 - 1 so an int would overflow
    part2 *= count_trees(terrain, 1, 1);
    part2 *= count_trees(terrain, 5, 1);
    part2 *= count_trees(terrain, 7, 1);
    part2 *= count_trees(terrain, 1, 2);
    std::cout << "The answer to Part 2 is " << part2 << "." << std::endl;
}
