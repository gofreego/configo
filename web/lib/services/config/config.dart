import 'dart:convert';

import 'package:web/models/config/config.dart';
import 'package:web/models/config/metadata.dart';
import 'package:web/models/config/object.dart';
import 'package:web/models/response/response.dart';
import 'package:web/utils/api.dart';
import 'package:http/http.dart' as http;

class ConfigService extends BaseAPIService {
  // Private constructor
  ConfigService._privateConstructor();

  // Singleton instance
  static final ConfigService _instance = ConfigService._privateConstructor();

  // Factory constructor to return the singleton instance
  factory ConfigService() {
    return _instance;
  }

  Future<ApiResponse<ConfigMetadataResponse>> getConfigMetadata() async {
    final String url =
        '${BaseAPIService.BASE_URL}/configo/metadata'; // Removed extra space
    http.Response response;
    response = await http.get(
      Uri.parse(url),
      headers: BaseAPIService.getHeaders(),
    );
    if (response.statusCode == 200) {
      return ApiResponse.fromSuccessJson(
        json.decode(response.body),
        ConfigMetadataResponse.fromJson,
      );
    } else {
      return ApiResponse.fromErrorJson(json.decode(response.body));
    }
  }

  Future<ApiResponse<GetConfigResponse>> getConfig(String id) async {
    final String url = '${BaseAPIService.BASE_URL}/configo/config?id=$id';

    final response = await http.get(
      Uri.parse(url),
      headers: BaseAPIService.getHeaders(),
    );

    if (response.statusCode == 200) {
      return ApiResponse.fromSuccessJson(
        json.decode(response.body),
        GetConfigResponse.fromJson,
      );
    } else {
      return ApiResponse.fromErrorJson(json.decode(response.body));
    }
  }

  Future<ApiResponse> updateConfig(
    String id,
    List<ConfigObject> configs,
  ) async {
    final String url = '${BaseAPIService.BASE_URL}/configo/config';

    final response = await http.patch(
      Uri.parse(url),
      headers: BaseAPIService.getHeaders(),
      body: json.encode(
        UpdateConfigRequest(id: id, configs: configs).toJson(),
      ),
    );

    if (response.statusCode == 200) {
      return ApiResponse.fromMessageJson(json.decode(response.body));
    } else {
      return ApiResponse.fromErrorJson(json.decode(response.body));
    }
  }
}
