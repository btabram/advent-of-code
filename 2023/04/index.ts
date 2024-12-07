import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const cards = lines.map((line, i) => {
  const numbers = line.split(":")[1]!;
  // Lines are like "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53".
  const [winningNumbers, numbersYouHave] = numbers
    .split("|")
    .map((numbersStr) =>
      numbersStr
        .split(" ")
        .filter(Boolean)
        .map((number) => Number.parseInt(number, 10)),
    ) as [number[], number[]];
  const id = i + 1; // Assume inputs are ordered
  return { id, winningNumbers, numbersYouHave };
});

const countWinners = ({
  winningNumbers,
  numbersYouHave,
}: (typeof cards)[number]) =>
  numbersYouHave.filter((n) => winningNumbers.includes(n)).length;

const scoreCard = (card: (typeof cards)[number]) => {
  const winnersYouHave = countWinners(card);
  return Math.max(1 << (winnersYouHave - 1), 0);
};

const part1 = cards.map(scoreCard).reduce((acc, score) => acc + score, 0);

const cardsWithCounts = cards.map((card) => ({ ...card, count: 1 }));

cardsWithCounts.forEach((cwc, i) => {
  const numberOfWinners = countWinners(cwc);
  const amountOfThisCard = cwc.count;

  // Copy subsequent cards when a card (and all its copies) win.
  let remainingToCopyCount = numberOfWinners;
  let indexToCopy = i + 1;
  while (remainingToCopyCount && indexToCopy < cardsWithCounts.length) {
    cardsWithCounts[indexToCopy]!.count += amountOfThisCard;
    remainingToCopyCount -= 1;
    indexToCopy += 1;
  }
});

const part2 = cardsWithCounts.reduce((acc, { count }) => acc + count, 0);

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
