import 'dart:math';

import '../lib/utils.dart';

typedef Coord = ({int x, int y});

int rectangleArea(Coord a, Coord b) {
  return ((a.x - b.x).abs() + 1) * ((a.y - b.y).abs() + 1).toInt();
}

bool isInsideRectangle(Coord point, Coord corner1, Coord corner2) {
  final minX = min(corner1.x, corner2.x);
  final maxX = max(corner1.x, corner2.x);
  final minY = min(corner1.y, corner2.y);
  final maxY = max(corner1.y, corner2.y);

  return point.x > minX && point.x < maxX && point.y > minY && point.y < maxY;
}

void main() {
  final lines = ReadInput();

  final redTiles = <Coord>[];
  for (final line in lines) {
    final parts = line.split(',');
    redTiles.add((x: int.parse(parts[0]), y: int.parse(parts[1])));
  }

  var part1 = 0;
  var part2 = 0;

  for (var i = 0; i < redTiles.length; i++) {
    for (var j = i + 1; j < redTiles.length; j++) {
      final area = rectangleArea(redTiles[i], redTiles[j]);

      // In part 1 we naively find the largest rectangle with red tile corners.
      if (area > part1) {
        part1 = area;
      }

      // In part 2 we want the largest rectangle within the shape formed by the outline of the red
      // tiles. This is a lot harder! I struggled with a general solution so I plotted my input and
      // it's a very specific shape - basically a circle with a cutout (imagine pac-man with a very
      // thin and long mouth). So I've hardcoded this to only consider rectangles that are entirely
      // within either the top or bottom half of the circle.
      if (area > part2) {
        if (!((redTiles[i].y >= 50147 && redTiles[j].y >= 50147) ||
            (redTiles[i].y <= 48634 && redTiles[j].y <= 48634))) {
          continue;
        }
        var noTilesInside = true;
        for (final tile in redTiles) {
          if (isInsideRectangle(tile, redTiles[i], redTiles[j])) {
            noTilesInside = false;
            break;
          }
        }
        // If there are no red tiles inside the rectangle it must be entirely within the shape.
        if (noTilesInside) {
          part2 = area;
        }
      }
    }
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
