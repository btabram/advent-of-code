#include <algorithm>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <tuple>
#include <set>
#include <vector>

using Point3D = std::tuple<int, int, int>;
using Point4D = std::tuple<int, int, int, int>;

template <typename PointT>
class PocketDimension
{
public:
    PocketDimension(const std::vector<PointT>& initial_active_cubes)
    {
        for (auto& initial_cube : initial_active_cubes)
        {
            active_cubes.insert(initial_cube);
        }
    }

    // We're interested in the number of active cubes after 6 cycles.
    size_t solve()
    {
        for (auto i = 0; i < 6; ++i)
        {
            advance();
        }
        return active_cubes.size();
    }


private:
    void advance()
    {
        // Currently inactive cubes which have active neighbours and may become active.
        std::set<PointT> to_possibly_activate;
        // Currently inactive cubes which become active during this cycle.
        std::vector<PointT> to_activate;
        // Active cubes which will deactivate during this cycle.
        std::vector<PointT> to_deactivate;

        for (const auto& cube : active_cubes)
        {
            auto active_neighbours = 0;
            for (auto& neighbour : get_neighbours(cube))
            {
                if (is_active(neighbour))
                {
                    ++active_neighbours;
                }
                else
                {
                    // The pocket dimension is infinite and any cube which borders active cubes could become active.
                    to_possibly_activate.insert(neighbour);
                }
            }
            // Cubes stay active as long as they have exactly 2 or 3 active neighbours.
            if (!(active_neighbours == 2 || active_neighbours == 3))
            {
                to_deactivate.push_back(cube);
            }
        }

        for (const auto& cube : to_possibly_activate)
        {
            auto neighbours = get_neighbours(cube);
            // Inactive cubes become active if they have exactly 3 active neighbours.
            if (std::count_if(neighbours.cbegin(), neighbours.cend(), [this](auto& n) { return is_active(n); }) == 3)
            {
                to_activate.push_back(cube);
            }
        }

        std::for_each(to_activate.cbegin(), to_activate.cend(), [this](auto& c) { active_cubes.insert(c); });
        std::for_each(to_deactivate.cbegin(), to_deactivate.cend(), [this](auto& c) { active_cubes.erase(c); });
    }

    bool is_active(const PointT& cube) const
    {
        return active_cubes.find(cube) != active_cubes.end();
    }

    std::array<Point3D, 26> get_neighbours(const Point3D& cube) const
    {
        auto [x, y, z] = cube;
        return {{
            {x + 1, y + 1, z + 1}, {x + 1, y, z + 1}, {x + 1, y - 1, z + 1},
            {x,     y + 1, z + 1}, {x,     y, z + 1}, {x,     y - 1, z + 1},
            {x - 1, y + 1, z + 1}, {x - 1, y, z + 1}, {x - 1, y - 1, z + 1},
            {x + 1, y + 1, z    }, {x + 1, y, z    }, {x + 1, y - 1, z    },
            {x,     y + 1, z    },                    {x,     y - 1, z    },
            {x - 1, y + 1, z    }, {x - 1, y, z    }, {x - 1, y - 1, z    },
            {x + 1, y + 1, z - 1}, {x + 1, y, z - 1}, {x + 1, y - 1, z - 1},
            {x,     y + 1, z - 1}, {x,     y, z - 1}, {x,     y - 1, z - 1},
            {x - 1, y + 1, z - 1}, {x - 1, y, z - 1}, {x - 1, y - 1, z - 1},
        }};
    }

    std::vector<Point4D> get_neighbours(const Point4D& cube) const
    {
        auto [x_ref, y_ref, z_ref, w_ref] = cube;
        std::vector<Point4D> neighbours;
        for (auto x = x_ref - 1; x <= x_ref + 1; ++x)
        {
            for (auto y = y_ref - 1; y <= y_ref + 1; ++y)
            {
                for (auto z = z_ref - 1; z <= z_ref + 1; ++z)
                {
                    for (auto w = w_ref - 1; w <= w_ref + 1; ++w)
                    {
                        if (x == x_ref && y == y_ref && z == z_ref && w == w_ref)
                        {
                            continue;
                        }
                        neighbours.emplace_back(x, y, z, w);
                    }
                }
            }
        }
        return neighbours;
    }

    std::set<PointT> active_cubes;
};

std::vector<Point3D> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Point3D> active_cubes;

    auto x = 0, y = 0, z = 0;
    std::string line;
    while (std::getline(input, line))
    {
        x = 0;
        for (auto c : line)
        {
            if (c == '#')
            {
                active_cubes.emplace_back(x, y, z);
            }
            ++x;
        }
        ++y;
    }

    return active_cubes;
}

int main()
{
    const auto initial_active_cubes_in_3D = read_input();

    // Part 2 we need to work in 4D.
    std::vector<Point4D> initial_active_cubes_in_4D;
    for (auto [x, y, z] : initial_active_cubes_in_3D)
    {
        initial_active_cubes_in_4D.emplace_back(x, y, z, 0);
    }

    auto solver_3d = PocketDimension(initial_active_cubes_in_3D);
    std::cout << "The answer to Part 1 is " << solver_3d.solve() << "." << std::endl;

    auto solver_4d = PocketDimension(initial_active_cubes_in_4D);
    std::cout << "The answer to Part 2 is " << solver_4d.solve() << "." << std::endl;
}
