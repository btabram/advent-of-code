import '../lib/utils.dart';

typedef Range = ({int start, int end});

void main() {
  final lines = ReadInput();

  final ranges = <Range>[];
  final ingredients = <int>[];

  var seenEmptyLine = false;
  for (final line in lines) {
    if (line.isEmpty) {
      seenEmptyLine = true;
      continue;
    }

    if (!seenEmptyLine) {
      final parts = line.split('-');
      ranges.add((start: int.parse(parts[0]), end: int.parse(parts[1])));
    } else {
      ingredients.add(int.parse(line));
    }
  }

  bool isFresh(int ingredientId) {
    for (final range in ranges) {
      if (ingredientId >= range.start && ingredientId <= range.end) {
        return true;
      }
    }
    return false;
  }

  final part1 = ingredients.where(isFresh).length;

  Range? mergeRanges(Range a, Range b) {
    // +1 is important because ranges are inclusive at both ends.
    if (a.end + 1 < b.start || b.end + 1 < a.start) {
      return null;
    }

    final merged = (
      start: a.start < b.start ? a.start : b.start,
      end: a.end > b.end ? a.end : b.end,
    );

    return merged;
  }

  var mergedRanges = List<Range>.from(ranges);

  var nextIndexToTry = 0;
  var tryingToMerge = mergedRanges.removeAt(nextIndexToTry);

  while (true) {
    var didMerge = false;

    for (var i = 0; i < mergedRanges.length; i++) {
      final merged = mergeRanges(tryingToMerge, mergedRanges[i]);
      if (merged != null) {
        mergedRanges.removeAt(i);
        tryingToMerge = merged; // Immediately try to merge the new range

        didMerge = true;
        nextIndexToTry = -1; // Something's changed, start over to be safe

        break;
      }
    }

    if (!didMerge) {
      mergedRanges.add(tryingToMerge); // No merges, put it back
      nextIndexToTry++;
      if (nextIndexToTry >= mergedRanges.length - 1) {
        break; // Done, we've tried merging all combinations
      }
      tryingToMerge = mergedRanges.removeAt(nextIndexToTry);
    }
  }

  var part2 = 0;
  for (final range in mergedRanges) {
    part2 += (range.end - range.start + 1);
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
