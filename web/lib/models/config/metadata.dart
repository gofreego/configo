class ConfigMetadataResponse {
  final String name;
  final String description;
  final List<String> keys;

  ConfigMetadataResponse({
    required this.name,
    required this.description,
    required this.keys,
  });

  /// Convert a JSON Map to a ConfigMetadataResponse instance
  factory ConfigMetadataResponse.fromJson(Map<String, dynamic> json) {
    return ConfigMetadataResponse(
      name: json['name'],
      description: json['description'],
      keys: List<String>.from(json['keys']),
    );
  }

  /// Convert a ConfigMetadataResponse instance to a JSON Map
  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'description': description,
      'keys': keys,
    };
  }
}
