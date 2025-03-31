
import 'package:flutter/foundation.dart' show kIsWeb;
class BaseAPIService {

  static const String BASE_URL = 'http://localhost:8085/myservice';

  static String? _cachedBaseUrl;
  //get base url
  // if the browser url contains /configo/v1, then the base url will be split and the first part will be returned
   // Get base URL from the current environment
  static String getBaseURL() {
    if (_cachedBaseUrl != null) {
      return _cachedBaseUrl!;
    }

    if (kIsWeb) {
      try {
        // Web-specific code to get current browser URL
        final currentUrl = Uri.base.toString();
        
        if (currentUrl.contains('/configo/v1')) {
          _cachedBaseUrl = currentUrl.split('/configo/v1')[0];
          return _cachedBaseUrl!;
        }
      } catch (e) {
        print('Error getting web URL: $e');
      }
    }
    
    // Fallback to default URL
    _cachedBaseUrl = BASE_URL;
    return BASE_URL;
  }
  //get headers 
  static Map<String, String> getHeaders() {
    return {
      'Content-Type': 'application/json',
    };
  }
}