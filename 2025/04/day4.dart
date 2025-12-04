import '../lib/utils.dart';

typedef Position = ({int x, int y});

void main() {
  final lines = ReadInput();

  final paperPositions = <Position>{};
  for (final (y, line) in lines.indexed) {
    for (final (x, char) in line.split('').indexed) {
      if (char == '@') {
        paperPositions.add((x: x, y: y));
      }
    }
  }

  bool isRemovable(Position p) {
    var count = 0;
    for (var dy = -1; dy <= 1; dy++) {
      for (var dx = -1; dx <= 1; dx++) {
        if (dx == 0 && dy == 0) {
          continue;
        }
        if (paperPositions.contains((x: p.x + dx, y: p.y + dy))) {
          count++;
        }
      }
    }
    return count < 4;
  }

  final part1 = paperPositions.where(isRemovable).length;

  var part2 = 0;
  while (true) {
    final toRemove = paperPositions.where(isRemovable).toList();

    if (toRemove.isEmpty) {
      break;
    }

    for (final p in toRemove) {
      paperPositions.remove(p);
    }

    part2 += toRemove.length;
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
