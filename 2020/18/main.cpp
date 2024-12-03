#include <algorithm>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

enum class Type {
    number,
    plus,
    times,
    left_bracket,
    right_bracket,
};

struct Token {
    Token(Type type, long value = 0) : type(type), value(value) {}

    Type type;
    long value;
};

std::ostream& operator<< (std::ostream& os, const Token& token)
{
    switch (token.type)
    {
        case Type::number:
            os << token.value;
            break;
        case Type::plus:
            os << '+';
            break;
        case Type::times:
            os << '*';
            break;
        case Type::left_bracket:
            os << '(';
            break;
        case Type::right_bracket:
            os << ')';
            break;
    }
    return os;
}

using Equation = std::vector<Token>;

std::vector<Equation> read_input()
{
    std::ifstream input("input.txt");

    std::vector<Equation> equations;

    std::string line;
    while (std::getline(input, line))
    {
        Equation equation;
        for (auto c : line)
        {
            switch (c)
            {
                case ' ':
                    break;
                case '+':
                    equation.emplace_back(Type::plus);
                    break;
                case '*':
                    equation.emplace_back(Type::times);
                    break;
                case '(':
                    equation.emplace_back(Type::left_bracket);
                    break;
                case ')':
                    equation.emplace_back(Type::right_bracket);
                    break;
                default:
                    // The input equations only contain one digit numbers.
                    equation.emplace_back(Type::number, std::stol(std::string{c}));
                    break;
            }
        }
        equations.push_back(equation);
    }

    return equations;
}

size_t find_matching_bracket(const Equation& equation, size_t opening_bracket_pos)
{
    auto right_brackets_to_find = 1;
    for (auto i = opening_bracket_pos + 1; i < equation.size(); ++i)
    {
        switch (equation[i].type)
        {
            case Type::left_bracket:
                right_brackets_to_find += 1;
                break;
            case Type::right_bracket:
                right_brackets_to_find -= 1;
                if (right_brackets_to_find == 0)
                {
                    return i;
                }
                break;
        }
    }
    std::cout << "Failed to find matching bracket" << std::endl;
    exit(1);
}


long solve(const Equation& initial_equation, bool plus_is_higher_priority)
{
    auto equation = initial_equation;

    // Evaluate and replace bracketed sub-expressions.
    while (std::any_of(equation.cbegin(), equation.cend(), [](auto& t) { return t.type == Type::left_bracket; }))
    {
        size_t i = 0;
        for (; i < equation.size(); ++i)
        {
            if (equation[i].type == Type::left_bracket)
            {
                break;
            }
        }
        auto matching_bracket = find_matching_bracket(equation, i);
        auto sub_expression = Equation({equation.cbegin() + i + 1, equation.cbegin() + matching_bracket});
        // Replace the sub-expression with with its value.
        equation[i] = Token{Type::number, solve(sub_expression, plus_is_higher_priority)};
        equation.erase(equation.cbegin() + i + 1, equation.cbegin() + matching_bracket + 1);
    }

    while (true)
    {
        size_t operation_to_resolve = 1;
        if (plus_is_higher_priority)
        {
            // If the order matters than resolve the higher priority opertions first.
            for (size_t i = 1; i < equation.size(); ++i)
            {
                if (equation[i].type == Type::plus)
                {
                    operation_to_resolve = i;
                    break;
                }
            }
        }
        auto a_pos = operation_to_resolve - 1;
        auto b_pos = operation_to_resolve + 1;

        auto val_a = equation[a_pos].value;
        auto op = equation[operation_to_resolve];
        auto val_b = equation[b_pos].value;

        long value;
        switch (op.type)
        {
            case Type::plus:
                value = val_a + val_b;
                break;
            case Type::times:
                value = val_a * val_b;
                break;
            default:
                std::cout << "Expected operator, got: " << op << std::endl;
                exit(1);
        }

        // Replace the binary operation with its value unless this was the final operation.
        if (equation.size() == 3)
        {
            return value;
        }
        equation[a_pos] = Token{Type::number, value};
        equation.erase(equation.cbegin() + a_pos + 1, equation.cbegin() + b_pos + 1);
    }
}

int main()
{
    const auto equations = read_input();
    long part1 = 0;
    long part2 = 0;
    for (auto& equation : equations)
    {
        part1 += solve(equation, false);
        part2 += solve(equation, true);
    }
    std::cout << "The answer to Part 1 is " << part1 << "." << std::endl;
    std::cout << "The answer to Part 2 is " << part2 << "." << std::endl;
}
