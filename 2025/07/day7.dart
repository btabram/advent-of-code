import '../lib/utils.dart';

typedef Coord = ({int x, int y});

void main() {
  final lines = ReadInput();

  final splitters = <Coord>{};
  Coord? start;

  for (var (y, line) in lines.indexed) {
    for (var (x, char) in line.split('').indexed) {
      if (char == 'S') {
        start = (x: x, y: y);
      } else if (char == '^') {
        splitters.add((x: x, y: y));
      }
    }
  }

  // Number of splitters hit.
  var part1 = 0;

  // Number of unique paths arriving at this beam.
  var beamCounts = {start!.x: 1};

  for (var y = start.y + 1; y < lines.length; y++) {
    final newBeamCounts = <int, int>{};

    for (final MapEntry(key: beamX, value: count) in beamCounts.entries) {
      if (splitters.contains((x: beamX, y: y))) {
        part1++;

        newBeamCounts[beamX - 1] = (newBeamCounts[beamX - 1] ?? 0) + count;
        newBeamCounts[beamX + 1] = (newBeamCounts[beamX + 1] ?? 0) + count;
      } else {
        newBeamCounts[beamX] = (newBeamCounts[beamX] ?? 0) + count;
      }
    }

    beamCounts = newBeamCounts;
  }

  // Total number of unique beam paths.
  final part2 = beamCounts.values.reduce(sum);

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
