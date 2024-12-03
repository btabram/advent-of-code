#include <algorithm>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

struct Seat {
    Seat(int row, int column) : row(row), column(column) {};

    int id() const
    {
        return 8*row + column;
    }

    int row;
    int column;
};

Seat parse_position(std::string position_str)
{
    std::string row_binary_str, column_binary_str;
    for (auto&& val : position_str)
    {
        // Front in row can be mapped to binary 1.
        if (val ==  'B') row_binary_str += "1";
        // Back in row can be mapped to binary 0.
        else if (val ==  'F') row_binary_str += "0";
        // Right in column can be mapped to binary 1.
        else if (val ==  'R') column_binary_str += "1";
        // Left in column can be mapped to binary 0.
        else if (val ==  'L') column_binary_str += "0";
    }

    auto row = std::stoi(row_binary_str, nullptr, 2);
    auto column = std::stoi(column_binary_str, nullptr, 2);

    return {row, column};
}

std::vector<Seat> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Seat> data;
    std::string line;
    while (std::getline(input, line))
    {
        data.push_back(parse_position(line));
    }

    return data;
}

int main()
{
    const auto seats = read_input();

    std::vector<int> seat_ids;
    std::transform(seats.cbegin(), seats.cend(), std::back_inserter(seat_ids), [](Seat s) { return s.id(); });
    std::sort(seat_ids.begin(), seat_ids.end());

    std::cout << "The answer to Part 1 is " << seat_ids[seat_ids.size() - 1] << "." << std::endl;

    int last_id = 0;
    for (auto&& id : seat_ids)
    {
        // Look for a gap in the seats which has neighbours.
        if (last_id + 2 == id)
        {
            std::cout << "The answer to Part 1 is " << id - 1 << "." << std::endl;
            break;
        }
        last_id = id;
    }
}
