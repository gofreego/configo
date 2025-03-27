import 'dart:convert';

import 'package:configo_ui/models/config/config.dart';
import 'package:configo_ui/models/config/metadata.dart';
import 'package:configo_ui/models/config/object.dart';
import 'package:configo_ui/models/response/response.dart';
import 'package:configo_ui/utils/api.dart';
import 'package:dio/dio.dart';

class ConfigService extends BaseAPIService {
  // Private constructor
  ConfigService._privateConstructor() {
    _dio = Dio();
    _dio.options.headers = BaseAPIService.getHeaders();
  }

  // Singleton instance
  static final ConfigService _instance = ConfigService._privateConstructor();
  late final Dio _dio;

  // Factory constructor to return the singleton instance
  factory ConfigService() {
    return _instance;
  }

  Future<ApiResponse<ConfigMetadataResponse>> getConfigMetadata() async {
    try {
      final String url = '${BaseAPIService.BASE_URL}/configo/metadata';
      final response = await _dio.get(url);
      return ApiResponse.fromSuccessJson(
        json.decode(response.data),
        ConfigMetadataResponse.fromJson,
      );
    } catch (e) {
      return ApiResponse.fromDioException(e);
    }
  }

  Future<ApiResponse<GetConfigResponse>> getConfig(String id) async {
    try {
      final String url = '${BaseAPIService.BASE_URL}/configo/config?id=$id';
      final response = await _dio.get(url);

      return ApiResponse.fromSuccessJson(
        json.decode(response.data),
        GetConfigResponse.fromJson,
      );
    } catch (e) {
      return ApiResponse.fromDioException(e);
    }
  }

  Future<ApiResponse> updateConfig(
    String id,
    List<ConfigObject> configs,
  ) async {
    try {
      final String url = '${BaseAPIService.BASE_URL}/configo/config';
      final data = UpdateConfigRequest(id: id, configs: configs).toJson();

      final response = await _dio.post(url, data: json.encode(data));

      return ApiResponse.fromMessageJson(json.decode(response.data));
    } catch (e) {
      return ApiResponse.fromDioException(e);
    }
  }
}
