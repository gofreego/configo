class BaseAPIService {
  static const String BASE_URL = 'http://localhost:8085/myservice';

  //get headers 
  static Map<String, String> getHeaders() {
    return {
      'Content-Type': 'application/json',
    };
  }
}