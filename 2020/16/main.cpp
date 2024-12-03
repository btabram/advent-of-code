#include <algorithm>
#include <fstream>
#include <iostream>
#include <regex>
#include <sstream>
#include <string>
#include <tuple>
#include <unordered_map>
#include <unordered_set>
#include <vector>

using Ticket = std::vector<int>;

class Field {
public:
    Field(std::string name, int min1, int max1, int min2, int max2)
         : name(name), min1(min1), max1(max1), min2(min2), max2(max2) {}

    const std::string name;

    bool is_valid(int value) const
    {
        return (value >= min1 && value <= max1) || (value >= min2 && value <= max2);
    }

private:
    const int min1;
    const int max1;
    const int min2;
    const int max2;
};

std::vector<int> parse_csv(std::string data)
{
    std::vector<int> output;
    std::string part;
    auto ss = std::istringstream{data};
    while (std::getline(ss, part, ','))
    {
        output.push_back(std::stoi(part));
    }
    return output;
}

std::tuple<std::vector<Field>, Ticket, std::vector<Ticket>> read_input()
{
    std::ifstream input("input.txt");
    std::string line;

    std::vector<Field> fields;
    static auto field_regex = std::regex{R"(^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$)"};

    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        std::smatch match;
        std::regex_match(line, match, field_regex);
        fields.emplace_back(match[1],               // name
                            std::stoi(match[2]),    // min1
                            std::stoi(match[3]),    // max1
                            std::stoi(match[4]),    // min2
                            std::stoi(match[5]));   // max2
    }

    Ticket my_ticket;
    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        if (line != "your ticket:")
        {
            my_ticket = parse_csv(line);
        }
    }

    std::vector<Ticket> nearby_tickets;
    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        if (line != "nearby tickets:")
        {
            nearby_tickets.push_back(parse_csv(line));
        }
    }

    return { fields, my_ticket, nearby_tickets };
}

int main()
{
    const auto [fields, my_ticket, nearby_tickets] = read_input();
    std::vector<Ticket> valid_nearby_tickets;


    // Part 1.
    auto ticket_scanning_error_rate = 0;
    // Check the validity of every value on every nearby ticket.
    for (auto& ticket : nearby_tickets)
    {
        auto valid_ticket = true;
        for (auto value : ticket)
        {
            auto valid_value = false;
            for (auto& field : fields)
            {
               if (field.is_valid(value))
               {
                   valid_value = true;
                   break;
               }
            }
            // If a value is not valid for any field then add it to the ticket scanning error rate.
            if (!valid_value)
            {
                ticket_scanning_error_rate += value;
                valid_ticket = false;
            }
        }

        if (valid_ticket)
        {
            valid_nearby_tickets.push_back(ticket);
        }
    }
    std::cout << "The answer to Part 1 is " << ticket_scanning_error_rate << "." << std::endl;


    // Part 2.
    // Map of the index of a field on the tickets => the valid values we have for that field.
    std::unordered_map<int, std::unordered_set<int>> field_values;
    for (const auto& ticket : valid_nearby_tickets)
    {
        for (auto i = 0; i < ticket.size(); ++i)
        {
            field_values[i].insert(ticket[i]);
        }
    }

    // Map of the index of a field on the tickets => the possible field names that field could correspond to.
    std::unordered_map<int, std::unordered_set<std::string>> field_names;
    for (const auto& [index, values] : field_values)
    {
        for (auto& field : fields)
        {
            if (std::all_of(values.cbegin(), values.cend(), [&field](auto value) { return field.is_valid(value); }))
            {
                field_names[index].insert(field.name);
            }
        }
    }

    std::unordered_set<std::string> solved;
    // Keep looping until we've matched up all the field names.
    while (std::any_of(field_names.cbegin(), field_names.cend(), [](auto& pair) { return pair.second.size() > 1; }))
    {
        // Find an index which has only one possible field and where that field is one we haven't yet solved.
        std::string newly_solved;
        for (const auto& [index, names] : field_names)
        {
            if (names.size() != 1)
            {
                continue;
            }
            auto name = *(names.cbegin());
            if (solved.find(name) == solved.end())
            {
                newly_solved = name;
                break;
            }
        }

        // Remove the newly solved field from the list of possible names for other indicies.
        for (auto& [index, names] : field_names)
        {
            if (names.size() == 1)
            {
                continue; // This index has already been solved
            }
            names.erase(newly_solved);
        }

        solved.insert(newly_solved);
    }

    long my_departure_fields_product = 1; // The answer is too big for an int
    for (const auto& [index, names] : field_names)
    {
        auto name = *(names.cbegin());
        if (name.rfind("departure", 0) == 0) // startswith
        {
            my_departure_fields_product *= my_ticket[index];
        }
    }
    std::cout << "The answer to Part 2 is " << my_departure_fields_product << "." << std::endl;
}
