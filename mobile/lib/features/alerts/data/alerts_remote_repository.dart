import '../../../core/network/api_client.dart';
import '../domain/alert.dart';
import 'alerts_repository.dart';

class AlertsRemoteRepository implements AlertsRepository {
  AlertsRemoteRepository(this._client, {required this.accessToken});

  final ApiClient _client;
  final String accessToken;

  @override
  Future<List<Alert>> listAlerts() async {
    final data = await _client.getJson('/alerts');
    final items = data['items'];
    if (items is List<dynamic>) {
      final out = <Alert>[];
      for (final item in items) {
        out.add(Alert.fromJson(item as Map<String, dynamic>));
      }
      return out;
    }
    final alerts = data['alerts'];
    if (alerts is List<dynamic>) {
      final out = <Alert>[];
      for (final item in alerts) {
        out.add(Alert.fromJson(item as Map<String, dynamic>));
      }
      return out;
    }
    return <Alert>[];
  }

  @override
  Future<Alert> createAlert({
    required String userId,
    required String symbol,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final data = await _client.postJson(
      '/alerts',
      headers: {
        'Authorization': 'Bearer $accessToken',
      },
      body: {
        'user_id': userId,
        'symbol': symbol,
        'target_price': targetPrice,
        'direction': switch (direction) {
          AlertDirection.above => 'above',
          AlertDirection.below => 'below',
        },
        'device_token': '',
      },
    );
    return Alert.fromJson(data);
  }
}
