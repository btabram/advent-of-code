#include <algorithm>
#include <fstream>
#include <iostream>
#include <map>
#include <string>
#include <vector>

// Global state.
std::pair<int, int> g_position = { 0, 0 }; // (x, y)
std::pair<int, int> g_waypoint = { 10, 1 }; // (x, y)
char g_direction = 'E'; // North, East, South, West

void reset()
{
    g_position = { 0, 0 };
    g_waypoint = { 10, 1 };
    g_direction = 'E';
}

void north(int value, bool use_waypoint = false)
{
    use_waypoint ? g_waypoint.second += value : g_position.second += value;
}

void east(int value, bool use_waypoint = false)
{
    use_waypoint ? g_waypoint.first += value : g_position.first += value;
}

void south(int value, bool use_waypoint = false)
{
    use_waypoint ? g_waypoint.second -= value : g_position.second -= value;
}

void west(int value, bool use_waypoint = false)
{
    use_waypoint ? g_waypoint.first -= value : g_position.first -= value;
}

const std::array<char, 4> CARDINALS = { 'N', 'E', 'S', 'W' };

// Rotate (x, y) about the origin anticlockwise for theta in degrees.
std::pair<int, int> rotate(std::pair<int, int> position, int theta)
{
    auto [x, y] = position;
    switch (theta % 360)
    {
        case 0:
            return { x, y };
        case 90:
        case -270:
            return { -y, x };
        case 180:
        case -180:
            return { -x, -y };
        case 270:
        case -90:
            return { y, -x };
    }

    std::cout << "Invalid angle " << theta << "." << std::endl;
    exit(1);
}

void left(int value, bool use_waypoint = false)
{
    if (use_waypoint)
    {
        g_waypoint = rotate(g_waypoint, value);
    }
    else
    {
        auto cardinal_index =
            std::distance(CARDINALS.cbegin(), std::find(CARDINALS.cbegin(), CARDINALS.cend(), g_direction));
        cardinal_index -= (value % 360) / 90;
        if (cardinal_index < 0)
        {
            cardinal_index += 4;
        }
        g_direction = CARDINALS[cardinal_index];
    }
}

void right(int value, bool use_waypoint = false)
{
    if (use_waypoint)
    {
        g_waypoint = rotate(g_waypoint, -value);
    }
    else
    {
        auto cardinal_index =
            std::distance(CARDINALS.cbegin(), std::find(CARDINALS.cbegin(), CARDINALS.cend(), g_direction));
        cardinal_index += (value % 360) / 90;
        cardinal_index %= 4;
        g_direction = CARDINALS[cardinal_index];
    }
}

void forwards(int value, bool use_waypoint = false)
{
    if (use_waypoint)
    {
        auto [d_x, d_y] = g_waypoint;
        g_position.first += (d_x * value);
        g_position.second += (d_y * value);
    }
    else
    {
        switch (g_direction)
        {
            case 'N':
                north(value, use_waypoint);
                break;
            case 'E':
                east(value, use_waypoint);
                break;
            case 'S':
                south(value, use_waypoint);
                break;
            case 'W':
                west(value, use_waypoint);
                break;
        }
    }
}

std::vector<std::pair<char, int>> read_input()
{
    std::ifstream input("input.txt");

    std::vector<std::pair<char, int>> data;
    std::string line;
    while (std::getline(input, line))
    {
        data.emplace_back(line[0], std::stoi(line.substr(1)));
    }

    return data;
}

int main()
{
    const auto instructions = read_input();

    const std::map<char, void (*)(int, bool)> instruction_fns = {
        { 'N', north }, { 'E', east }, { 'S', south }, { 'W', west }, { 'L', left }, { 'R', right }, { 'F', forwards }
    };

    for (auto [instuction, value] : instructions)
    {
        instruction_fns.at(instuction)(value, false);
    }
    auto manhattan_distance = std::abs(g_position.first) + std::abs(g_position.second);
    std::cout << "The answer to Part 1 is " << manhattan_distance << "." << std::endl;

    reset();

    for (auto [instuction, value] : instructions)
    {
        instruction_fns.at(instuction)(value, true);
    }
    manhattan_distance = std::abs(g_position.first) + std::abs(g_position.second);
    std::cout << "The answer to Part 2 is " << manhattan_distance << "." << std::endl;
}
