import '../../../core/network/api_client.dart';
import '../domain/alert.dart';
import 'alerts_repository.dart';

class AlertsRemoteRepository implements AlertsRepository {
  AlertsRemoteRepository(this._client);

  final ApiClient _client;

  @override
  Future<List<Alert>> listAlerts() async {
    final data = await _client.getJsonWithHeaders(
      '/alerts',
    );
    final source = switch (data) {
      final List<dynamic> list => list,
      final Map<String, dynamic> map when map['items'] is List<dynamic> =>
        map['items'] as List<dynamic>,
      final Map<String, dynamic> map when map['alerts'] is List<dynamic> =>
        map['alerts'] as List<dynamic>,
      _ => const <dynamic>[],
    };
    return [
      for (final item in source) Alert.fromJson(item as Map<String, dynamic>),
    ];
  }

  @override
  Future<Alert> createAlert({
    required String userId,
    required String actionId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final data = await _client.postJson(
      '/alerts',
      body: {
        'user_id': userId,
        'action_id': actionId,
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

  @override
  Future<Alert> updateAlert({
    required String alertId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final data = await _client.patchJson(
      '/alerts/$alertId',
      body: {
        'target_price': targetPrice,
        'direction': switch (direction) {
          AlertDirection.above => 'above',
          AlertDirection.below => 'below',
        },
      },
    );
    return Alert.fromJson(data as Map<String, dynamic>);
  }

  @override
  Future<void> deleteAlert({
    required String alertId,
  }) async {
    await _client.deleteJson(
      '/alerts/$alertId',
    );
  }
}
