class ConfigMetadataResponse {
  final ConfigInfo configInfo;
  final ServiceInfo serviceInfo;

  ConfigMetadataResponse({
    required this.configInfo,
    required this.serviceInfo,
  });

  /// Convert a JSON Map to a ConfigResponse instance
  factory ConfigMetadataResponse.fromJson(Map<String, dynamic> json) {
    return ConfigMetadataResponse(
      configInfo: ConfigInfo.fromJson(json['configInfo']),
      serviceInfo: ServiceInfo.fromJson(json['serviceInfo']),
    );
  }

  /// Convert a ConfigResponse instance to a JSON Map
  Map<String, dynamic> toJson() {
    return {
      'configInfo': configInfo.toJson(),
      'serviceInfo': serviceInfo.toJson(),
    };
  }
}

class ConfigInfo {
  final List<String> configKeys;

  ConfigInfo({required this.configKeys});

  factory ConfigInfo.fromJson(Map<String, dynamic> json) {
    return ConfigInfo(
      configKeys: List<String>.from(json['configKeys']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'configKeys': configKeys,
    };
  }
}

class ServiceInfo {
  final String description;
  final String name;

  ServiceInfo({required this.description, required this.name});

  factory ServiceInfo.fromJson(Map<String, dynamic> json) {
    return ServiceInfo(
      description: json['description'],
      name: json['name'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'description': description,
      'name': name,
    };
  }
}
