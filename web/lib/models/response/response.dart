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

  /// Convert a JSON Map to an ApiResponse instance for error responses
  factory ApiResponse.fromErrorJson(Map<String, dynamic> json) {
    try {
      return ApiResponse<T>(error: json['error']);
    } catch (e) {
      throw Exception('Somethig went wrong');
    }
  }
}
