#include <algorithm>
#include <fstream>
#include <iostream>
#include <set>
#include <vector>

// Using axial hex coordinates here, see https://www.redblobgames.com/grids/hexagons/#coordinates-axial
using Coord = std::pair<int, int>;

enum class Direction { E, SE, SW, W, NW, NE };

std::vector<std::vector<Direction>> read_input()
{
    std::ifstream input("input.txt");
    std::vector<std::vector<Direction>> directions_list;
    std::string line;
    while (std::getline(input, line))
    {
        std::vector<Direction> directions;
        for (auto i = 0; i < line.size(); ++i)
        {
            Direction d;
            auto first_char = line[i];
            if (first_char == 'e')
            {
                d = Direction::E;
            }
            else if (first_char == 'w')
            {
                d = Direction::W;
            }
            else
            {
                ++i;
                auto east = line[i] == 'e';
                if (first_char == 's')
                {
                    d = east ? Direction::SE : Direction::SW;
                }
                else // first_char == 'n'
                {
                    d = east ? Direction::NE : Direction::NW;
                }
            }
            directions.push_back(d);
        }
        directions_list.push_back(std::move(directions));
    }
    return directions_list;
}

std::vector<Coord> get_neighbours(const Coord& position)
{
    auto [x, y] = position;
    return { { x + 1, y }, { x, y + 1 }, { x - 1, y + 1 }, { x - 1, y }, { x, y - 1 }, { x + 1, y - 1} };
}


int main()
{
    const auto directions_list = read_input();

    std::set<Coord> black_tiles;
    auto is_black = [&black_tiles](const Coord& tile) { return black_tiles.find(tile) != black_tiles.end(); };
    auto flip_tile = [&black_tiles](const Coord& tile)
    {
        auto [it, inserted] = black_tiles.insert(tile);
        if (!inserted)
        {
            // The tile was already black so we flip it back to white.
            black_tiles.erase(it);
        }
    };

    // For Part 1 we set up the tiled floor according to the input instructions.
    for (const auto& directions : directions_list)
    {
        Coord position;
        // Go through the directions to work out which tile should be flipped.
        for (auto direction : directions)
        {
            switch (direction)
            {
                case Direction::E:
                    position.first += 1;
                    break;
                case Direction::SE:
                    position.second += 1;
                    break;
                case Direction::SW:
                    position.first -= 1;
                    position.second += 1;
                    break;
                case Direction::W:
                    position.first -= 1;
                    break;
                case Direction::NW:
                    position.second -= 1;
                    break;
                case Direction::NE:
                    position.first += 1;
                    position.second -= 1;
                    break;
            }
        }
        flip_tile(position);
    }
    std::cout << "The answer to Part 1 is " << black_tiles.size() << "." << std::endl;

    // For Part 2 we let the floor evolve over 100 days.
    for (auto i = 1; i <= 100; ++i)
    {
        // Every black tile could flip over to white and every white neighbour of black tiles could flip to black.
        auto black_tiles_and_neighbours = black_tiles;
        for (const auto& tile : black_tiles)
        {
            for (const auto& neighbour : get_neighbours(tile))
            {
                black_tiles_and_neighbours.insert(neighbour);
            }
        }

        std::vector<Coord> tiles_to_flip;
        for (const auto& tile : black_tiles_and_neighbours)
        {
            auto neighbours = get_neighbours(tile);
            auto black_neighbour_count = std::count_if(neighbours.begin(), neighbours.end(), is_black);

            if (is_black(tile))
            {
                if (black_neighbour_count == 0 || black_neighbour_count > 2)
                {
                    tiles_to_flip.push_back(tile);
                }
            }
            else // Tile is white
            {
                if (black_neighbour_count == 2)
                {
                    tiles_to_flip.push_back(tile);
                }
            }
        }

        for (const auto& tile : tiles_to_flip)
        {
            flip_tile(tile);
        }
    }
    std::cout << "The answer to Part 2 is " << black_tiles.size() << "." << std::endl;
}
