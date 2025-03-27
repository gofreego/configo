import 'package:configo_ui/models/config/object.dart';

class GetConfigResponse {
  final List<ConfigObject>? configs;

  GetConfigResponse({this.configs});

  factory GetConfigResponse.fromJson(Map<String, dynamic> json) {
    return GetConfigResponse(
      configs:
          json['configs'] != null
              ? (json['configs'] as List)
                  .map((e) => ConfigObject.fromJson(e))
                  .toList()
              : [],
    );
  }
}

class UpdateConfigRequest {
  final String key;
  final List<ConfigObject> configs;

  UpdateConfigRequest({required this.key, required this.configs});

  Map<String, dynamic> toJson() {
    return {"key": key, 'configs': configs.map((e) => e.toJson()).toList()};
  }
}
