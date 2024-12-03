#include <algorithm>
#include <fstream>
#include <iostream>
#include <map>
#include <regex>
#include <set>
#include <string>
#include <vector>

struct BagRule {
    BagRule(std::string containing_bag, std::map<std::string, int> contained_bags)
        : containing_bag{containing_bag}, contained_bags{contained_bags} {}

    bool contains(std::string bag_name) const
    {
        for (auto&& [name, number] : contained_bags)
        {
            if (name == bag_name)
            {
                return true;
            }
        }
        return false;
    }

    std::string containing_bag;
    std::map<std::string, int> contained_bags;
};

// Make a regex which can accomodate up to |max_containing_bags|.
std::regex make_regex(int max_containing_bags)
{
    std::string regex_string = R"(^(\w+ \w+) bags contain (?:no other bags|(\d+) (\w+ \w+) bags?)";
    for (int i = 0; i < max_containing_bags; ++i)
    {
        regex_string += R"((?:, (\d+) (\w+ \w+) bags?)?)";
    }
    regex_string += R"().$)";
    return std::regex{regex_string};
}

std::vector<BagRule> read_input()
{
    std::ifstream input("input.txt");

    auto most_commas = 0;
    std::vector<std::string> lines;
    std::string line;
    while (std::getline(input, line))
    {
        auto comma_count = std::count_if(line.cbegin(), line.cend(), [](char c) { return c == ','; });
        if (comma_count > most_commas) most_commas = comma_count;
        lines.push_back(line);
    }

    std::vector<BagRule> rules;
    // The maximum number of containing bags is equal to the maximum numbers of commas in a line.
    auto re = make_regex(most_commas);
    std::smatch match;
    for (auto&& line : lines)
    {
        if (!std::regex_match(line, match, re))
        {
            std::cout << "Regex failed for: " << line << std::endl;
            exit(1);
        }

        auto matches = match.cbegin();
        // We don't care about the first, full string match.
        ++matches;

        auto containing_bag = *matches;
        ++matches;

        std::map<std::string, int> contained_bags;
        while (matches != match.cend())
        {
            // Some input strings are don't match all the groups.
            if (*matches == "")
            {
                break;
            }

            auto number = std::stoi(*matches);
            ++matches;
            auto name = *matches;
            ++matches;
            contained_bags.emplace(name, number);
        }

        rules.emplace_back(containing_bag, contained_bags);
    }

    return rules;
}

void part1(const std::vector<BagRule>& rules)
{
    // Shiny Gold Containing Bags AKA sgcb.
    std::set<std::string> sgcb{"shiny gold"};
    // Loop over all rules, adding to our set of bags which can ultimately contain a shiny gold bag and then checking
    // for any new rules which match now that we've added to our set and so on until there's no more changes.
    while (true)
    {
        std::set<std::string> to_add;
        for (auto&& rule : rules)
        {
            for (auto&& bag : sgcb)
            {
                auto rule_name = rule.containing_bag;
                if (rule.contains(bag) && sgcb.find(rule_name) == sgcb.end())
                {
                    to_add.insert(rule_name);
                }
            }
        }
        if (to_add.size() == 0)
        {
            break;
        }
        std::for_each(to_add.cbegin(), to_add.cend(), [&sgcb](std::string name) { sgcb.insert(name); });
        to_add.clear();
    }

    // -1 because we don't count the shiny gold bag itself.
    std::cout << "The answer to Part 1 is " << sgcb.size() - 1 << "." << std::endl;
}

void part2(const std::vector<BagRule>& rules)
{
    // Map of names of bags contained in the shiny gold bag => the number of them in the shiny gold bag.
    std::map<std::string, int> bags;
    // Set of bags which don't contain any other bags.
    std::set<std::string> leaf_bags;
    const BagRule* shiny_gold;
    for (auto&& rule : rules)
    {
        if (rule.containing_bag == "shiny gold")
        {
            shiny_gold = &rule;
        }
        if (rule.contained_bags.size() == 0)
        {
            leaf_bags.insert(rule.containing_bag);
        }
    }
    for (auto&& bag : (*shiny_gold).contained_bags)
    {
        bags.insert(bag);
    }

    auto bag_count = 0;
    // Every iteration we substitute a bag type in |bags| for the bags it contains until we can't substitute any more.
    while (true)
    {
        auto bag_it = bags.cbegin();
        // Ignore bags which don't contain any other bags for now.
        while (leaf_bags.find((*bag_it).first) != leaf_bags.end())
        {
            ++bag_it;
        }
        auto [name, number] = *bag_it;

        for (auto&& rule : rules)
        {
            if (rule.containing_bag == name)
            {
                // Substitute |name| for the bags it contains.
                for (auto&& [sub_name, sub_number] : rule.contained_bags)
                {
                    // If |sub_name| isn't already in |bags| then the [] operator will insert it and default construct
                    // the mapped value, which will be 0 in this case.
                    bags[sub_name] += number * sub_number;
                }
                bags.erase(name);
                // Add the substituted bags to our total count.
                bag_count += number;
            }
        }

        // Stop when we've substituted enough such that only bags which don't contain any other bags remain.
        if (std::all_of(bags.cbegin(), bags.cend(),
                        [&leaf_bags](auto pair) { return leaf_bags.find(pair.first) != leaf_bags.end(); }))
        {
            break;
        }
    }

    // Finally add the leaf bags to our total count;
    std::for_each(bags.cbegin(), bags.cend(), [&bag_count](auto pair) { bag_count += pair.second; });

    std::cout << "The answer to Part 2 is " << bag_count << "." << std::endl;
}

int main()
{
    const auto rules = read_input();
    part1(rules);
    part2(rules);
}
