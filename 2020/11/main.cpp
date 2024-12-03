#include <algorithm>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

enum class Type
{
    Floor,
    EmptySeat,
    OccupiedSeat,
};

struct Position
{
    Position(Type type, size_t x, size_t y) : type(type), x(x), y(y) {}

    Type type;
    size_t x;
    size_t y;
};

class WaitingArea
{
public:
    WaitingArea(const std::vector<char>& input, size_t width) : width(width), height(input.size() / width)
    {
        for (size_t y = 0; y < height; ++y)
        {
            for (size_t x = 0; x < width; ++x)
            {
                // The input only contains floor and empty seats.
                auto type = input[x + y*width] == 'L' ? Type::EmptySeat : Type::Floor;
                data.emplace_back(type, x, y);
            }
        }
    }

    void print() const
    {
        for (size_t i = 0; i < data.size(); ++i)
        {
            switch (data[i].type)
            {
                case Type::Floor:
                    std::cout << '.';
                    break;
                case Type::EmptySeat:
                    std::cout << 'L';
                    break;
                case Type::OccupiedSeat:
                    std::cout << '#';
                    break;
            }

            if ((i + 1) % width == 0)
            {
                std::cout << std::endl;
            }
        }
    }

    int part1()
    {
        // Keep advancing until we reach a steady state.
        while (advance(4, true) > 0) {}

        return std::count_if(data.cbegin(), data.cend(),
                             [](const auto& position) { return position.type == Type::OccupiedSeat; });
    }

    int part2()
    {
        // Keep advancing until we reach a steady state.
        while (advance(5, false) > 0) {}

        return std::count_if(data.cbegin(), data.cend(),
                             [](const auto& position) { return position.type == Type::OccupiedSeat; });
    }

private:
    size_t advance(int threshold, bool nearest_neighbours_only)
    {
        std::vector<size_t> to_change;
        for (size_t i = 0; i < data.size(); ++i)
        {
            const auto& position = data[i];
            if (position.type == Type::Floor) {
                continue;
            }

            auto neighbours = count_neighbouring_occupied_seats(position.x, position.y, nearest_neighbours_only);
            switch (position.type)
            {
                case Type::EmptySeat:
                    if (neighbours == 0)
                    {
                        to_change.push_back(i);
                    }
                    break;
                case Type::OccupiedSeat:
                    if (neighbours >= threshold)
                    {
                        to_change.push_back(i);
                    }
                    break;
            }
        }

        for (auto i : to_change)
        {
            auto& position = data[i];
            switch (position.type)
            {
                case Type::EmptySeat:
                    position.type = Type::OccupiedSeat;
                    break;
                case Type::OccupiedSeat:
                    position.type = Type::EmptySeat;
                    break;
            }
        }

        // Return the number of changed positions.
        return to_change.size();
    }

    const Position& get(size_t x, size_t y) const
    {
        return data[x + y*width];
    }

    // Use int here because we might have negative numbers;
    int count_neighbouring_occupied_seats(int x, int y, bool nearest_neighbours_only) const
    {
        static std::vector<std::pair<int, int>> directions = {
            {1, 1}, {1, 0}, {1, -1}, {0, 1}, {0, -1}, {-1, 1}, {-1, 0}, {-1, -1},
        };

        auto count = 0;
        for (auto [d_x, d_y] : directions)
        {
            // If |nearest_neighbours_only| is true then only consider i = 1.
            for (auto i = 1; nearest_neighbours_only ? i < 2 : true; ++i)
            {
                auto n_x = x + i*d_x;
                auto n_y = y + i*d_y;

                if (n_x < 0 || n_x >= width || n_y < 0 || n_y >= height)
                {
                    break;
                }

                // Stop when we hit a seat of some kind;
                auto type = get(n_x, n_y).type;
                if (type != Type::Floor)
                {
                    if (type == Type::OccupiedSeat)
                    {
                        ++count;
                    }
                    break;
                }
            }
        }

        return count;
    }

    std::vector<Position> data;
    const size_t width;
    const size_t height;
};

std::pair<std::vector<char>, size_t> read_input()
{
    std::ifstream input("input.txt");

    std::vector<char> data;
    size_t width = 0;
    std::string line;
    while (std::getline(input, line))
    {
        std::for_each(line.cbegin(), line.cend(), [&data](char c) { data.push_back(c); });
        if (width == 0)
        {
            width = line.size();
        }
    }

    return {data, width};
}

int main()
{
    const auto [input, width] = read_input();
    std::cout << "The answer to Part 1 is " << WaitingArea{input, width}.part1() << "." << std::endl;
    std::cout << "The answer to Part 2 is " << WaitingArea{input, width}.part2() << "." << std::endl;
}
