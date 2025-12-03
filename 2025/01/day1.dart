import 'dart:io';

import '../lib/utils.dart';

enum Direction { left, right }

void main() {
  final lines = ReadInput();

  final rotations = <({Direction direction, int distance})>[];
  for (final line in lines) {
    rotations.add((
      direction: line[0] == 'L' ? Direction.left : Direction.right,
      distance: int.parse(line.substring(1)),
    ));
  }

  var position = 50;
  var part1 = 0;
  var part2 = 0;

  for (final rotation in rotations) {
    // Move one click at a time. Part 2 wants the count of visits to 0.
    final delta = switch (rotation.direction) {
      Direction.left => -1,
      Direction.right => 1,
    };

    for (var i = 0; i < rotation.distance; i++) {
      position += delta;
      position %= 100;

      if (position == 0) {
        part2++;
      }
    }

    // Part 1 only cares about rotations ending at 0.
    if (position == 0) {
      part1++;
    }
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
