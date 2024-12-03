#include <algorithm>
#include <fstream>
#include <iostream>
#include <iterator>
#include <regex>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>

using Passport = std::unordered_map<std::string, std::string>;

std::vector<std::string> split_string(const std::string& str)
{
    auto iss = std::istringstream{str};
    // istream_iterator++ draws out of the input stream, until the next space.
    auto start = std::istream_iterator<std::string>{iss};
    auto end = std::istream_iterator<std::string>{};
    return {start, end};
}

std::vector<Passport> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Passport> data;
    Passport passport;
    std::string line;
    while (std::getline(input, line))
    {
        // A blank line signals the start of a new passport entry.
        if (line == "")
        {
            data.push_back(passport);
            passport.clear();
            continue;
        }

        for (auto&& part : split_string(line))
        {
            auto field_name = part.substr(0, 3);
            auto field_value = part.substr(4);
            passport.emplace(field_name, field_value);
        }
    }
    if (passport.size() > 0)
    {
        data.push_back(passport);
    }

    return data;
}

bool passport_has_required_fields(const Passport& passport)
{
    static const auto required_fields = std::array<std::string, 7>{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"};

    for (auto&& field : required_fields)
    {
        if (passport.find(field) == passport.cend())
        {
            return false;
        }
    }
    return true;
}

bool is_valid_passport(const Passport& passport)
{
    if (!passport_has_required_fields(passport)) return false;

    auto birth_year = std::stoi(passport.at("byr"));
    if (birth_year < 1920 || birth_year > 2002) return false;

    auto issue_year = std::stoi(passport.at("iyr"));
    if (issue_year < 2010 || issue_year > 2020) return false;

    auto expiration_year = std::stoi(passport.at("eyr"));
    if (expiration_year < 2020 || expiration_year > 2030) return false;

    auto height = passport.at("hgt");
    auto units = height.substr(height.size() - 2);
    if (units == "cm")
    {
        auto value = std::stoi(height.substr(0, height.size() - 2));
        if (value < 150 || value > 193) return false;
    }
    else if (units == "in")
    {
        auto value = std::stoi(height.substr(0, height.size() - 2));
        if (value < 59 || value > 76) return false;
    }
    else
    {
        return false;
    }

    static const auto hair_colour_regex = std::regex{"#[0-9a-f]{6}"};
    auto hair_colour = passport.at("hcl");
    std::smatch match;
    if (!std::regex_match(hair_colour, match, hair_colour_regex)) return false;

    static const auto valid_eye_colours = std::array<std::string, 7>{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"};
    auto eye_colour = passport.at("ecl");
    if (std::find(valid_eye_colours.cbegin(), valid_eye_colours.cend(), eye_colour) == valid_eye_colours.cend())
    {
        return false;
    }

    static const auto passport_id_regex = std::regex{"[0-9]{9}"};
    auto passport_id = passport.at("pid");
    if (!std::regex_match(passport_id, match, passport_id_regex)) return false;

    return true;
}

int main()
{
    const auto passports = read_input();

    auto complete_count = std::count_if(passports.cbegin(), passports.cend(), passport_has_required_fields);
    auto valid_count = std::count_if(passports.cbegin(), passports.cend(), is_valid_passport);

    std::cout << "The answer to Part 1 is " << complete_count << "." << std::endl;
    std::cout << "The answer to Part 2 is " << valid_count << "." << std::endl;
}
