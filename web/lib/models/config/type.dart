/// Configuration types
enum ConfigType { string, number, boolean, json, bigText, choice, parent, list }

/// Extension for ConfigType to handle JSON conversion
extension ConfigTypeExtension on ConfigType {
  static ConfigType fromString(String type) {
    switch (type) {
      case 'string':
        return ConfigType.string;
      case 'number':
        return ConfigType.number;
      case 'boolean':
        return ConfigType.boolean;
      case 'json':
        return ConfigType.json;
      case 'big_text':
        return ConfigType.bigText;
      case 'choice':
        return ConfigType.choice;
      case 'parent':
        return ConfigType.parent;
      case 'list':
        return ConfigType.list;
      default:
        throw ArgumentError('Unknown ConfigType: $type');
    }
  }

  String toJson() {
    return toString().split('.').last;
  }
}
