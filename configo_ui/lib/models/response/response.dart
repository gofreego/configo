import 'package:dio/dio.dart';

class ApiResponse<T> {
  final String? message;
  final String? error;
  final T? data;

  ApiResponse({this.message, this.error, this.data});

  /// Convert a JSON Map to an ApiResponse instance for success responses
  factory ApiResponse.fromSuccessJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) fromJson,
  ) {
    return ApiResponse<T>(
      message: json['message'],
      data: json['data'] != null ? fromJson(json['data']) : null,
    );
  }

  // only message
  factory ApiResponse.fromMessageJson(Map<String, dynamic> json) {
    return ApiResponse<T>(message: json['message']);
  }

  /// Convert a JSON Map to an ApiResponse instance for error responses
  factory ApiResponse.fromErrorJson(Map<String, dynamic> json) {
    try {
      return ApiResponse(error: json['error']);
    } catch (e) {
      throw Exception('Somethig went wrong');
    }
  }

  // from dio exception
  factory ApiResponse.fromDioException(Object e) {
    try {
      if (e is DioError && e.response != null) {
        return ApiResponse.fromErrorJson(e.response!.data);
      }
      return ApiResponse(error: "An error occurred. Please try again later.");
    } catch (e) {
      throw Exception('Somethig went wrong');
    }
  }
}
