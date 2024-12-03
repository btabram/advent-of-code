#include <algorithm>
#include <fstream>
#include <iostream>
#include <iterator>
#include <set>
#include <sstream>
#include <string>
#include <map>
#include <vector>

// Use ordered sets because set_intersection requires sorted ranges.
using StringSet = std::set<std::string>;

struct Food
{
    Food(const StringSet& ingredients, const StringSet& alergens) : ingredients(ingredients), alergens(alergens) {}

    const StringSet ingredients;
    const StringSet alergens;
};

std::vector<std::string> split_string(const std::string& str)
{
    auto iss = std::istringstream{str};
    // istream_iterator++ draws out of the input stream, until the next space.
    auto start = std::istream_iterator<std::string>{iss};
    auto end = std::istream_iterator<std::string>{};
    return {start, end};
}

std::vector<Food> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Food> data;
    std::string line;
    while (std::getline(input, line))
    {
        StringSet ingredients;
        StringSet alergens;
        auto finished_ingedients = false;
        for (auto&& part : split_string(line))
        {
            if (!finished_ingedients)
            {
                if (part == "(contains")
                {
                    finished_ingedients = true;
                    continue;
                }
                ingredients.insert(part);
            }
            else
            {
                // Trim the last char, it's either "," or ")".
                alergens.insert(part.substr(0, part.size() - 1));
            }
        }
        data.emplace_back(ingredients, alergens);
    }

    return data;
}

int main()
{
    const auto food = read_input();

    StringSet all_ingredients;
    std::map<std::string, int> ingredient_occurances;
    for (auto& f : food)
    {
        for (auto& ingredient : f.ingredients)
        {
            all_ingredients.insert(ingredient);
            ingredient_occurances[ingredient] += 1;
        }
    }

    // Map of alergen -> set of ingredients which could contain that alergen.
    std::map<std::string, StringSet> possible_alergens;
    for (auto& f : food)
    {
        for (auto& alergen : f.alergens)
        {
            auto new_possibilities = StringSet{f.ingredients.begin(), f.ingredients.end()};
            if (possible_alergens.find(alergen) == possible_alergens.end())
            {
                // If we haven't encountered this alergen before then all ingredients are possibilities.
                possible_alergens[alergen] = new_possibilities;
            }
            else
            {
                // We already have some possibilities, use the new information to narrow it down.
                const StringSet& existing_possibilities = possible_alergens[alergen];
                StringSet intersection;
                std::set_intersection(existing_possibilities.begin(), existing_possibilities.end(),
                                      new_possibilities.begin(), new_possibilities.end(),
                                      std::inserter(intersection, intersection.begin()));
                possible_alergens[alergen] = intersection;
            }
        }
    }

    // The set of ingredients which may contain alergens.
    StringSet suspicious_ingredients;
    for (const auto& [alergen, possibilities] : possible_alergens)
    {
        suspicious_ingredients.insert(possibilities.begin(), possibilities.end());
    }

    // The set of ingredients which definitely don't contain alergens.
    StringSet good_ingredients;
    std::set_difference(all_ingredients.begin(), all_ingredients.end(),
                        suspicious_ingredients.begin(), suspicious_ingredients.end(),
                        std::inserter(good_ingredients, good_ingredients.begin()));

    auto part1 = 0;
    for (const auto& good_ingredient : good_ingredients)
    {
        part1 += ingredient_occurances[good_ingredient];
    }
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;

    // Loop until we've determined the single ingredient which corresponds to each alergen.
    StringSet solved_alergens;
    while (std::any_of(possible_alergens.begin(), possible_alergens.end(),
                       [](auto& pa) { return pa.second.size() > 1; }))
    {
        // Identify an unsolved alergen which we can solve because there's only one possibility.
        std::string alergen_to_solve;
        std::string solved_ingredient;
        for (const auto& [alergen, possibilities] : possible_alergens)
        {
            if (possibilities.size() == 1 && solved_alergens.find(alergen) == solved_alergens.end())
            {
                alergen_to_solve = alergen;
                solved_ingredient = *possibilities.begin();
            }
        }

        // Remove the solved alergen ingredient as a possibility for other alergens.
        for (auto& [alergen, possibilities] : possible_alergens)
        {
            if (possibilities.size() > 1)
            {
                possibilities.erase(solved_ingredient);
            }
        }

        solved_alergens.insert(alergen_to_solve);
    }

    // The answer for Part 2 is the list of alergen ingredients, sorted alphabetically by alergen. Fortunately
    // |possible_alergens| is already sorted alphabetically by alergen.
    std::string part2 = "";
    for (const auto& [alergen, possibilities] : possible_alergens)
    {
        part2 += *possibilities.begin() + ",";
    }
    part2 = part2.substr(0, part2.size() - 1); // Remove the trailing comma
    std::cout << "The answer to Part 2 is:\n\t" << part2 << std::endl;
}
