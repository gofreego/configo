import 'package:web/models/config/object.dart';

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
  final String id;
  final List<ConfigObject> configs;

  UpdateConfigRequest({required this.id, required this.configs});

  Map<String, dynamic> toJson() {
    return {"id": id, 'configs': configs.map((e) => e.toJson()).toList()};
  }
}
