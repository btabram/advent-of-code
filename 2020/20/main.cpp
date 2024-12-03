#include <exception>
#include <fstream>
#include <optional>
#include <iostream>
#include <regex>
#include <set>
#include <string>
#include <map>
#include <vector>

template<typename T>
T reverse(const T& container)
{
    return T{ container.rbegin(), container.rend() };
}

class Image
{
public:
    Image(const std::vector<std::string>& data) : data(data) {}

    virtual std::vector<std::string> image() const
    {
        return data;
    }

    virtual void rotate90()
    {
        auto size = data.size(); // Assume square
        auto rotated_data = data;
        for (auto i = 0; i < size; ++i)
        {
            for (auto j = 0; j < size; ++j)
            {
                rotated_data[i][j] = data[size - 1 - j][i];
            }
        }
        data = rotated_data;
    }

    virtual void reflect()
    {
        data = reverse(data);
    }

protected:
    std::vector<std::string> data;
};

struct Border
{
    Border(const std::string& str) : str(str), matching_id(std::nullopt) {};

    std::string str;
    std::optional<long> matching_id;
};

class Tile : public Image
{
public:
    Tile(const std::vector<std::string>& dataWithBorder) : Image(dataWithBorder)
    {
        borders = generate_border_strings();
    }

    std::vector<std::string> image() const override
    {
        // Trim the borders.
        std::vector<std::string> image{ data.cbegin() + 1, data.cend() - 1 };
        std::transform(image.begin(), image.end(), image.begin(),
                       [](const auto& s) { return std::string{ s.cbegin() + 1, s.cend() - 1 }; });
        return image;
    }

    std::vector<Border> borders;

    const Border& top() const { return borders[0]; }
    const Border& left() const { return borders[1]; }
    const Border& bottom() const { return borders[2]; }
    const Border& right() const { return borders[3]; }

    const std::set<std::string>& get_possible_borders() const
    {
        if (!possible_borders)
        {
            possible_borders = std::set<std::string>{};
            for (const auto& border : borders)
            {
                (*possible_borders).insert(border.str);
                (*possible_borders).insert(reverse(border.str));
            }
        }
        return *possible_borders;
    }

    void rotate90() override
    {
        Image::rotate90();

        auto new_borders = generate_border_strings();
        // Rotate everything round.
        new_borders[0].matching_id = borders[1].matching_id;
        new_borders[1].matching_id = borders[2].matching_id;
        new_borders[2].matching_id = borders[3].matching_id;
        new_borders[3].matching_id = borders[0].matching_id;
        borders = new_borders;
    }

    void reflect() override
    {
        Image::reflect();

        auto new_borders = generate_border_strings();
        // Swap top and bottom.
        new_borders[0].matching_id = borders[2].matching_id;
        new_borders[1].matching_id = borders[1].matching_id;
        new_borders[2].matching_id = borders[0].matching_id;
        new_borders[3].matching_id = borders[3].matching_id;
        borders = new_borders;
    }

private:
    std::vector<Border> generate_border_strings() const
    {
        std::string top = data.front();
        std::string bottom = data.back();
        std::string left, right;
        for (const auto& line : data) {
            left.push_back(line.front());
            right.push_back(line.back());
        }
        return { top, left, bottom, right };
    }

    mutable std::optional<std::set<std::string>> possible_borders;
};

std::map<long, Tile> read_input()
{
    std::ifstream input("input.txt");
    std::map<long, Tile> tiles;
    std::string line;

    static auto tile_id_regex = std::regex{"^Tile ([0-9]+):$"};

    while (std::getline(input, line))
    {
        std::smatch match;
        std::regex_match(line, match, tile_id_regex);
        auto id = std::stol(match[1]);

        std::vector<std::string> image;
        while (std::getline(input, line) && line.size() != 0) // Read until a blank line
        {
            image.push_back(line);
        }

        tiles.emplace(id, Tile{ image });
    }

    return tiles;
}

bool check_border_match(const std::string& requirement, const Border& border)
{
    // Check for the edges of the image which should have unique borders.
    if (requirement.empty())
    {
        return !border.matching_id.has_value();
    }
    return requirement == border.str;
}

Tile find_matching_orientation(const std::string& required_top, const std::string& required_left, const Tile& ref_tile)
{
    auto check_tile_match = [&required_top, &required_left](const Tile& t)
    {
        return check_border_match(required_top, t.top()) && check_border_match(required_left, t.left());
    };

    // Consider all possible orientations.
    auto tile = ref_tile;
    for (auto i = 0; i < 4; ++i)
    {
        if (check_tile_match(tile)) return tile;

        tile.reflect();
        if (check_tile_match(tile)) return tile;
        tile.reflect();

        tile.rotate90();
    }

    throw std::runtime_error("Tile cannot be matched");
}

int count_sea_monsters(const Image& image)
{
    static const auto sea_monster_length = 20;
    static const auto sea_monster_top       = std::regex{"..................#."};
    static const auto sea_monster_middle    = std::regex{"#....##....##....###"};
    static const auto sea_monster_bottom    = std::regex{".#..#..#..#..#..#..."};
    // Use a lookahead so that the matching part of the regex is a single character because this works nicely with the
    // regex iterator to check every single position in the string for a possible match.
    static const auto sea_monster_middle_lookahead  = std::regex{"#(?=....##....##....###)"};

    auto count = 0;
    auto image_data = image.image();
    for (auto i = 1; i < image_data.size() - 1; ++i)
    {
        const auto& line = image_data[i];
        for (auto it = std::sregex_iterator(line.begin(), line.end(), sea_monster_middle_lookahead);
                it != std::sregex_iterator{}; ++it)
        {
            auto middle_match_position = (*it).position(0);
            std::smatch _;
            auto above = image_data[i - 1].substr(middle_match_position, sea_monster_length);
            auto below = image_data[i + 1].substr(middle_match_position, sea_monster_length);
            if (std::regex_match(above, _, sea_monster_top) && std::regex_match(below, _, sea_monster_bottom))
            {
                count++;
            }
        }
    }
    return count;
}

int find_most_sea_monsters(const Image& ref_image)
{
    auto max_count = 0;
    // Consider all possible orientations.
    auto image = ref_image;
    for (auto i = 0; i < 4; ++i)
    {
        auto count = count_sea_monsters(image);
        if (count > max_count) max_count = count;

        image.reflect();
        count = count_sea_monsters(image);
        if (count > max_count) max_count = count;
        image.reflect();

        image.rotate90();
    }
    return max_count;
}

int main()
{
    auto tiles = read_input();


    long part1 = 1;
    std::optional<Tile> top_left_corner;
    for (auto& [id, tile] : tiles)
    {
        // Go through all tiles, working out which ones fit together.
        for (const auto& [other_id, other_tile] : tiles)
        {
            if (id == other_id) continue;

            for (auto& border : tile.borders)
            {
                auto possible_borders = other_tile.get_possible_borders();
                if (possible_borders.find(border.str) != possible_borders.end())
                {
                    border.matching_id = other_id;
                    break;
                }
            }
        }

        // Count the number of borders which matched up with another tile.
        auto border_matches = 0;
        for (const auto& b : tile.borders)
        {
            if (b.matching_id) border_matches++;
        }
        // If there's two unmatched borders then this must be a corner tile.
        if (border_matches == 2)
        {
            part1 *= id;
            // Choose an arbitrary corner tile to become our top left corner tile for Part 2..
            if (!top_left_corner) top_left_corner = tile;
        }
    }
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;


    std::vector<std::vector<Tile>> full_image_mosaic{1, {{}}};
    // For Part 2 we need the full image. We build up the image as a mosaic of tiles. Starting with our abitrary choice
    // of a top left corner we keep adding matching tiles to the right / bottom until we have the full image.
    auto i = 0, j = 0 ;
    std::string required_top = "", required_left = "";
    Tile next_tile = *top_left_corner;
    while (true)
    {
        auto oriented_tile = find_matching_orientation(required_top, required_left, next_tile);
        full_image_mosaic[j].push_back(oriented_tile);

        long next_tile_id;
        if (auto ot_r = oriented_tile.right(); ot_r.matching_id.has_value())
        {
            // Not the end of a row.
            i++;
            required_left = ot_r.str;
            next_tile_id = *(ot_r.matching_id);
        }
        else
        {
            // End of the image.
            if (!oriented_tile.bottom().matching_id.has_value()) break;

            // End of a row.
            i = 0;
            j++;
            full_image_mosaic.push_back({});
            required_left = "";
        }

        // Don't need to consider j==0 because |required_top| starts at the correct value for this case.
        if (j != 0)
        {
            auto bottom_of_tile_above = full_image_mosaic[j - 1][i].bottom();
            required_top = bottom_of_tile_above.str;
            // We may have evaluated this earlier, but we must do it here for the case where we're starting a new row.
            next_tile_id = *(bottom_of_tile_above.matching_id);
        }

        next_tile = (*(tiles.find(next_tile_id))).second;
    }

    // Now we have our mosaic of tiles we can stitch together the full image.
    const auto tile_image_size = full_image_mosaic[0][0].image()[0].size(); // Assume square
    std::vector<std::string> full_image_data;
    auto row_num = 0;
    for (const auto& row : full_image_mosaic)
    {
        // Create blank rows...
        for (auto i = 0; i < tile_image_size; ++i) full_image_data.push_back("");
        // ...and then append the data from the tiles on that row.
        for (const auto& tile : row)
        {
            auto tile_image = tile.image();
            for (auto i = 0; i < tile_image_size; ++i)
            {
                full_image_data[(row_num * tile_image_size) + i] += tile_image[i];
            }
        }
        row_num++;
    }
    auto full_image = Image{ full_image_data };

    // Count the total number of rough water pixels.
    auto total_rough_waters = 0;
    for (auto line : full_image.image())
    {
        for (auto c : line)
        {
            if (c == '#') total_rough_waters++;
        }
    }
    // Work out the number of rough water pixels due to sea monsters.
    auto sea_monsters = find_most_sea_monsters(full_image);
    auto rough_sea_monster_tiles = sea_monsters * 15; // Assume that the sea monsters do not overlap
    std::cout << "The answer to Part 2 is " << total_rough_waters - rough_sea_monster_tiles << "." << std::endl;
}
