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
