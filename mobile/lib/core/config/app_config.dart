class AppConfig {
  const AppConfig({
    required this.apiBaseUrl,
    required this.useDemoData,
  });

  final String apiBaseUrl;
  final bool useDemoData;

  factory AppConfig.fromEnvironment() {
    return const AppConfig(
      apiBaseUrl: String.fromEnvironment(
        'API_BASE_URL',
        defaultValue: 'http://localhost:8080',
      ),
      useDemoData: bool.fromEnvironment(
        'USE_DEMO_DATA',
        defaultValue: false,
      ),
    );
  }
}
