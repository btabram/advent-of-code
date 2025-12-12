import 'dart:math';

import '../lib/utils.dart';

typedef Coord = ({int x, int y, int z});

num distance(Coord a, Coord b) {
  return sqrt(pow((a.x - b.x), 2) + pow((a.y - b.y), 2) + pow((a.z - b.z), 2));
}

void main() {
  final lines = ReadInput();

  final junctionBoxes = <Coord>[];
  for (final line in lines) {
    final parts = line.split(',');
    junctionBoxes.add((
      x: int.parse(parts[0]),
      y: int.parse(parts[1]),
      z: int.parse(parts[2]),
    ));
  }

  final distancesMap = <(int, int), num>{};

  for (var i = 0; i < junctionBoxes.length; i++) {
    for (var j = i + 1; j < junctionBoxes.length; j++) {
      // Standarise key order, we don't care about directionality.
      final key = (min(i, j), max(i, j));
      distancesMap[key] = distance(junctionBoxes[i], junctionBoxes[j]);
    }
  }

  final distances = distancesMap.entries.toList();
  distances.sort((a, b) => a.value.compareTo(b.value));

  var part1 = 0;
  var part2 = 0;

  var circuits = <Set<int>>[];
  var connectionsMade = 0;

  for (final MapEntry(key: (a, b)) in distances) {
    if (connectionsMade == 1000) {
      final lengths = circuits.map((c) => c.length).toList()..sort();
      part1 = lengths.reversed.toList().sublist(0, 3).reduce(product);
    }
    connectionsMade++;

    if (circuits.any((c) => c.contains(a) && c.contains(b))) {
      continue; // Already connected to each other, no point carrying on
    }

    var newCircuits = <Set<int>>[];
    var toMerge = <Set<int>>[];
    for (final circuit in circuits) {
      if (circuit.contains(a) || circuit.contains(b)) {
        circuit.addAll([a, b]);
        toMerge.add(circuit);
      } else {
        newCircuits.add(circuit);
      }
    }

    if (toMerge.isEmpty) {
      newCircuits.add({a, b});
    } else {
      newCircuits.add(toMerge.expand((c) => c).toSet());
    }

    circuits = newCircuits;

    if (newCircuits.length == 1 &&
        newCircuits[0].length == junctionBoxes.length) {
      // All connected into one circuit!
      part2 = junctionBoxes[a].x * junctionBoxes[b].x;
      break;
    }
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
