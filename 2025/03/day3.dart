import '../lib/utils.dart';

int getMaxJoltage(List<int> batteryBank, int batteriesToTurnOn) {
  final joltage = <int>[];

  // The index of the last battery we turned on, start subsequent searches from here.
  var lastIndex = -1;

  for (var i = 1; i <= batteriesToTurnOn; i++) {
    var bestValue = 0;
    for (final (j, battery) in batteryBank.indexed) {
      if (j <= lastIndex) {
        continue;
      }
      if (j >= batteryBank.length - batteriesToTurnOn + i) {
        continue; // Too near the end, wouldn't be able to turn on enough batteries
      }

      if (battery > bestValue) {
        bestValue = battery;
        lastIndex = j;
      }
    }

    joltage.add(bestValue);
  }

  return joltage.reduce((a, b) => (10 * a) + b);
}

void main() {
  final lines = ReadInput();

  final batteryBanks = <List<int>>[];
  for (final line in lines) {
    final bank = <int>[];
    for (final char in line.split('')) {
      bank.add(int.parse(char));
    }
    batteryBanks.add(bank);
  }

  final part1 = batteryBanks.map((bb) => getMaxJoltage(bb, 2)).reduce(sum);
  final part2 = batteryBanks.map((bb) => getMaxJoltage(bb, 12)).reduce(sum);

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
