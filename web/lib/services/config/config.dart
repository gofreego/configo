import 'dart:convert';

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
   
   Future<ConfigMetadataResponse> getConfigMetadata() async {
    final String url = '${BaseAPIService.BASE_URL}/configo/metadata'; // Removed extra space

    final response = await http.get(
      Uri.parse(url),
      headers: BaseAPIService.getHeaders(),
    );

    if (response.statusCode == 200) {
      return  ApiResponse.fromSuccessJson(json.decode(response.body),  ConfigMetadataResponse.fromJson).data!;
    } else {
      throw Exception('Failed to fetch config metadata');
    }
  }

  Future<ConfigObject> getConfig(String id) async {
    final String url = '${BaseAPIService.BASE_URL}/configo/config/$id'; 

    final response = await http.get(
      Uri.parse(url),
      headers: BaseAPIService.getHeaders(),
    );

    if (response.statusCode == 200) {
      return  ApiResponse.fromSuccessJson(json.decode(response.body),  ConfigObject.fromJson).data!;
    } else {
      throw Exception('Failed to fetch config object');
    }
  }
}