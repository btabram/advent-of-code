#include <algorithm>
#include <fstream>
#include <iostream>
#include <memory>
#include <regex>
#include <string>
#include <vector>

class IRule {
public:
    IRule(size_t id): id(id) {}
    virtual ~IRule() = default;

    bool match(const std::string& str) const
    {
        size_t i = 0;
        return match(str, i) && i == str.size();
    }

    virtual bool match(const std::string& str, size_t& pos) const = 0;

    const size_t id;
};

using RuleDirectory = std::vector<std::unique_ptr<IRule>>;
RuleDirectory rule_directory;

class NormalRule : public IRule {
public:
    NormalRule(size_t id, std::vector<size_t> first, std::vector<size_t> second)
        : IRule(id), first_subrules_group(first), second_subrules_group(second) {}

    bool match(const std::string& str, size_t& pos) const override
    {
        auto first_match = true;
        auto pos_copy = pos;
        for (auto& rule_id : first_subrules_group)
        {
            if (!rule_directory[rule_id]->match(str, pos_copy))
            {
                first_match = false;
                break;
            }
        }

        if (first_match)
        {
            pos = pos_copy;
            return true;
        }
        else if (second_subrules_group.size() == 0)
        {
            return false;
        }

        for (auto& rule_id : second_subrules_group)
        {
            if (!rule_directory[rule_id]->match(str, pos))
            {
                return false;
            }
        }

        return true;
    }

private:
    const std::vector<size_t> first_subrules_group;
    const std::vector<size_t> second_subrules_group;
};

class LeafRule : public IRule {
public:
    LeafRule(size_t id, char letter) : IRule(id), letter(letter) {}

    bool match(const std::string& str, size_t& pos) const override
    {
        if (str[pos] == letter)
        {
            ++pos;
            return true;
        }
        return false;
    }

private:
    const char letter;
};

class RepeatRule : public IRule {
public:
    RepeatRule(size_t id, size_t to_repeat, int repeats) : IRule(id), rule_to_repeat(to_repeat), repeats(repeats) {}

    // Match |repeats| repetitions of |rule_to_repeat|.
    bool match(const std::string& str, size_t& pos) const override
    {
        auto pos_copy = pos;
        for (int i = 0; i < repeats; ++i)
        {
            if (!rule_directory[rule_to_repeat]->match(str, pos_copy))
            {
                return false;
            }
        }
        pos = pos_copy;
        return true;
    }

private:
    const size_t rule_to_repeat;
    const int repeats;
};

class SymmetricRepeatRule : public IRule {
public:
    SymmetricRepeatRule(size_t id, size_t first, size_t second, int repeats)
                        : IRule(id), first_rule(first), second_rule(second), repeats(repeats) {}

    // Match |repeats| repetitions of |first_rule| followed by the same number of repetitions of |second_rule|.
    bool match(const std::string& str, size_t& pos) const override
    {
        auto pos_copy = pos;
        for (int i = 0; i < repeats; ++i)
        {
            if (!rule_directory[first_rule]->match(str, pos_copy))
            {
                return false;
            }
        }

        for (int i = 0; i < repeats; ++i)
        {
            if (!rule_directory[second_rule]->match(str, pos_copy))
            {
                return false;
            }
        }

        pos = pos_copy;
        return true;
    }

private:
    const size_t first_rule;
    const size_t second_rule;
    const int repeats;
};

std::vector<std::string> read_input()
{
    std::ifstream input("input.txt");
    std::string line;

    static auto leaf_rule_regex = std::regex{R"raw(^([0-9]+): "([a-z])"$)raw"};
    static auto rule_regex = std::regex{R"(^([0-9]+): ([0-9]+)(?: ([0-9]+))?(?: ([0-9]+))?(?: \| ([0-9]+)(?: ([0-9]+))?)?$)"};

    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        std::smatch match;
        if (std::regex_match(line, match, rule_regex))
        {
            size_t i = 1;
            auto id = std::stoul(match[i]);

            std::vector<size_t> first_subrules_group; // Rule IDs before "|"
            for (i = 2; i < 5; ++i)
            {
                auto value = match[i];
                if (value != "")
                {
                    first_subrules_group.push_back(std::stoul(value));
                }
            }

            std::vector<size_t> second_subrules_group; // Rule IDs after "|"
            for (; i < 7; ++i)
            {
                auto value = match[i];
                if (value != "")
                {
                    second_subrules_group.push_back(std::stoul(value));
                }
            }

            rule_directory.push_back(std::make_unique<NormalRule>(id, first_subrules_group, second_subrules_group));
        }
        else if (std::regex_match(line, match, leaf_rule_regex))
        {
            auto id = std::stoul(match[1]);
            auto letter = std::string{match[2]}[0];
            rule_directory.push_back(std::make_unique<LeafRule>(id, letter));
        }
        else
        {
            std::cout << "Error parsing input: " << line << std::endl;
            exit(1);
        }
    }

    std::sort(rule_directory.begin(), rule_directory.end(), [](auto& a, auto& b) { return a->id < b->id; });

    std::vector<std::string> messages;
    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        messages.push_back(line);
    }

    return messages;
}

int main()
{
    const auto messages = read_input();

    auto part1 = 0;
    for (auto& m : messages)
    {
        if (rule_directory[0]->match(m))
        {
            ++part1;
        }
    }
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;

    // For Part 2 we add in a couple of special rules which involve loops and therefore can match a number of repeats
    // of a pattern. The rules aren't greedy so it's not easy to fit these rules into my solver for Part 1 in a genearl
    // way. Instead we can repeatedly run the Part 1 solver, looking for a set number of repeats of each of the two
    // looping rules. The upper limit of 99 repeats is enough to solve Part 2 for my input.
    auto part2 = 0;
    for (auto& m : messages)
    {
        for (int i = 1; i < 100; ++i)
        {
            auto should_break = false;
            for (int j = 1; j < 100; ++j)
            {
                rule_directory[8] = std::make_unique<RepeatRule>(8, 42, i);
                rule_directory[11] = std::make_unique<SymmetricRepeatRule>(11, 42, 31, j);
                if (rule_directory[0]->match(m))
                {
                    ++part2;
                    should_break = true;
                    break;
                }
            }
            if (should_break) break;
        }
    }
    std::cout << "The answer to Part 2 is " << part2 << "." << std::endl;
}
