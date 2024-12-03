#include <fstream>
#include <iostream>

std::pair<long, long> read_input()
{
    std::ifstream input("input.txt");
    std::string first;
    std::getline(input, first);
    std::string second;
    std::getline(input, second);
    return { std::stol(first), std::stol(second) };
}

long transform_subject_number(long subject_number, long loop_size)
{
    long value = 1;
    for (auto i = 0; i < loop_size; ++i)
    {
        value *= subject_number;
        value %= 20201227;
    }
    return value;
}

long find_loop_size(long pub_key, long initial_subject_number)
{
    long loop_size = 1;
    long value = 1;
    while (true)
    {
        value *= initial_subject_number;
        value %= 20201227;
        if (value == pub_key)
        {
            return loop_size;
        }
        ++loop_size;
    }
}

int main()
{
    const auto [ card_pub_key, door_pub_key ] = read_input();
    const long initial_subject_number = 7;

    const auto card_loop_size = find_loop_size(card_pub_key, initial_subject_number);
    const auto door_loop_size = find_loop_size(door_pub_key, initial_subject_number);

    const auto encryption_key = transform_subject_number(card_pub_key, door_loop_size);
    if (transform_subject_number(door_pub_key, card_loop_size) != encryption_key)
    {
        std::cout << "Failed to find a consistent encryption key..." << std::endl;
        exit(1);
    }
    std::cout << "The answer to Part 1 is " << encryption_key << "." << std::endl;
}
