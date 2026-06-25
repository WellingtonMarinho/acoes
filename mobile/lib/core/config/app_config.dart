import 'dart:io';

class AppConfig {
  const AppConfig({
    required this.apiBaseUrl,
  });

  final String apiBaseUrl;

  factory AppConfig.fromEnvironment() {
    return AppConfig(
      apiBaseUrl: String.fromEnvironment(
        'API_BASE_URL',
        defaultValue: _platformDefaultApiBaseUrl(),
      ),
    );
  }
}

String _platformDefaultApiBaseUrl() {
  return Platform.isAndroid ? 'http://10.0.2.2:8080' : 'http://localhost:8080';
}
