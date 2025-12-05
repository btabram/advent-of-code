import '../lib/utils.dart';

void main() {
  final [line] = ReadInput();

  final ranges = line.split(',').map((pairStr) {
    final parts = pairStr.split('-');
    return (start: int.parse(parts[0]), end: int.parse(parts[1]));
  });

  bool isInvalidId(int id, {int? maxRepeatCoount = null}) {
    final idStr = id.toString();

    outerLoop:
    for (
      var repeatCount = 2;
      repeatCount <= (maxRepeatCoount ?? idStr.length);
      repeatCount++
    ) {
      if (idStr.length % repeatCount != 0) {
        continue;
      }
      final repeatLength = idStr.length ~/ repeatCount;

      final firstRepeat = idStr.substring(0, repeatLength);
      for (var i = 1; i < repeatCount; i++) {
        final nextRepeat = idStr.substring(
          i * repeatLength,
          (i + 1) * repeatLength,
        );
        if (nextRepeat != firstRepeat) {
          continue outerLoop;
        }
      }

      return true;
    }

    return false;
  }

  var part1 = 0;
  for (final range in ranges) {
    for (var id = range.start; id <= range.end; id++) {
      if (isInvalidId(id, maxRepeatCoount: 2)) {
        part1 += id;
      }
    }
  }

  var part2 = 0;
  for (final range in ranges) {
    for (var id = range.start; id <= range.end; id++) {
      if (isInvalidId(id)) {
        part2 += id;
      }
    }
  }

  print("The answer to part 1 is $part1");
  print("The answer to part 2 is $part2");
}
