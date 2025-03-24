import 'package:web/models/config/type.dart';

class ConfigObject {
  final String name;
  final ConfigType type;
  final String description;
  final bool required;
  final List<String>? choices;
   dynamic value;
  final List<ConfigObject>? children;

  ConfigObject({
    required this.name,
    required this.type,
    required this.description,
    required this.required,
    this.choices,
    this.value,
    this.children,
  });

  /// Convert a JSON Map to a ConfigObject instance
  factory ConfigObject.fromJson(Map<String, dynamic> json) {
    return ConfigObject(
      name: json['name'],
      type: ConfigTypeExtension.fromString(json['type']),
      description: json['description'],
      required: json['required'],
      choices:
          (json['choices'] as List<dynamic>?)?.map((e) => e as String).toList(),
      value: json['value'],
      children:
          (json['children'] as List<dynamic>?)
              ?.map((e) => ConfigObject.fromJson(e as Map<String, dynamic>))
              .toList(),
    );
  }

  /// Convert a ConfigObject instance to a JSON Map
  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'type': type.toJson(),
      'description': description,
      'required': required,
      'choices': choices,
      'value': value,
      'children': children?.map((e) => e.toJson()).toList(),
    };
  }
}
