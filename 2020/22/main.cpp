#include <assert.h>
#include <deque>
#include <fstream>
#include <iostream>
#include <set>
#include <string>

using Hand = std::deque<int>;

std::pair<Hand, Hand> read_input()
{
    std::ifstream input("input.txt");
    Hand hand_1, hand_2;
    std::string line;

    std::getline(input, line);
    assert (line == "Player 1:");
    while (std::getline(input, line) && line.size() != 0) // Read until a blank line
    {
        hand_1.push_back(std::stoi(line));
    }

    std::getline(input, line);
    assert (line == "Player 2:");
    while (std::getline(input, line)) // Can simply read until EOF here
    {
        hand_2.push_back(std::stoi(line));
    }

    return { hand_1, hand_2 };
}

void combat(Hand& hand_1, Hand& hand_2)
{
    while (hand_1.size() > 0 && hand_2.size() > 0)
    {
        auto card_1 = hand_1.front();
        hand_1.pop_front();
        auto card_2 = hand_2.front();
        hand_2.pop_front();

        if (card_1 > card_2)
        {
            hand_1.push_back(card_1);
            hand_1.push_back(card_2);
        }
        else // card_2 > card_1 (all cards are unique so equality isn't possible)
        {
            hand_2.push_back(card_2);
            hand_2.push_back(card_1);
        }
    }
}

Hand freeze_state(const Hand& hand_1, const Hand& hand_2)
{
    auto state = hand_1;
    state.push_back(-1);
    state.insert(state.end(), hand_2.begin(), hand_2.end());
    return state;
}

bool recursive_combat(Hand& hand_1, Hand& hand_2)
{
    std::set<Hand> previous_states;

    while (hand_1.size() > 0 && hand_2.size() > 0)
    {
        // Check for infinite loops.
        auto current_state = freeze_state(hand_1, hand_2);
        auto [ _, inserted ] = previous_states.insert(current_state);
        if (!inserted)
        {
            // Current state has been seen before, player 1 wins.
            return true;
        }

        // Draw cards.
        auto card_1 = hand_1.front();
        hand_1.pop_front();
        auto card_2 = hand_2.front();
        hand_2.pop_front();

        // Work out the winner for the round.
        bool player_1_won_round;
        if (hand_1.size() >= card_1 && hand_2.size() >= card_2)
        {
            // Time to start a sub-game!
            auto sub_hand_1 = Hand{ hand_1.begin(), hand_1.begin() + card_1 };
            auto sub_hand_2 = Hand{ hand_2.begin(), hand_2.begin() + card_2 };
            player_1_won_round = recursive_combat(sub_hand_1, sub_hand_2);
        }
        else
        {
            // Resolve without a sub-game.
            player_1_won_round = card_1 > card_2;
        }

        // Give the winner their cards.
        if (player_1_won_round)
        {
            hand_1.push_back(card_1);
            hand_1.push_back(card_2);
        }
        else
        {
            hand_2.push_back(card_2);
            hand_2.push_back(card_1);
        }
    }

    // Return whether player 1 won.
    return hand_1.size() > 0 ? true : false;
}

int score_game(const Hand& hand_1, const Hand& hand_2)
{
    const auto& winning_hand = hand_1.size() > 0 ? hand_1 : hand_2;
    const auto size = winning_hand.size();
    auto score = 0;
    for (auto i = 1; i <= size; i++)
    {
        score += i * winning_hand[size - i];
    }
    return score;
}

int main()
{
    const auto [ initial_hand_1, initial_hand_2 ] = read_input();

    auto hand_1 = initial_hand_1;
    auto hand_2 = initial_hand_2;
    combat(hand_1, hand_2);
    std::cout << "The answer to Part 1 is " << score_game(hand_1, hand_2) << "." << std::endl;

    hand_1 = initial_hand_1;
    hand_2 = initial_hand_2;
    recursive_combat(hand_1, hand_2);
    std::cout << "The answer to Part 2 is " << score_game(hand_1, hand_2) << "." << std::endl;
}
