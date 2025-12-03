import 'dart:io';

/// Reads lines from today's input file (or a test file if [testing] is true).
List<String> ReadInput({bool testing = false}) {
  final filename = testing ? 'test.txt' : 'input.txt';
  final scriptDir = File(Platform.script.toFilePath()).parent;
  return File('${scriptDir.path}/$filename').readAsLinesSync();
}

int sum(int a, int b) {
  return a + b;
}
