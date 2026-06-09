import 'dart:convert';

class CleanupItem {
  final String value;
  final int count;

  const CleanupItem({required this.value, required this.count});

  factory CleanupItem.fromJson(Map<String, dynamic> json) => CleanupItem(
        value: json['value'] as String,
        count: json['count'] as int,
      );

  String get countLabel => count >= 21 ? '20+' : '$count';
}

List<CleanupItem> cleanupItemsFromJson(String str) =>
    List<CleanupItem>.from(json.decode(str).map((x) => CleanupItem.fromJson(x)));

class LegumeReference {
  final String legume;
  final List<String> varietes;

  const LegumeReference({required this.legume, required this.varietes});

  factory LegumeReference.fromJson(Map<String, dynamic> json) => LegumeReference(
        legume: json['legume'] as String,
        varietes: List<String>.from(json['varietes']),
      );
}

List<LegumeReference> legumeReferencesFromJson(String str) =>
    List<LegumeReference>.from(
        json.decode(str).map((x) => LegumeReference.fromJson(x)));
