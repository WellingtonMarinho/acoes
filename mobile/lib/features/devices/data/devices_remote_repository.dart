import '../../../core/network/api_client.dart';
import 'devices_repository.dart';

class DevicesRemoteRepository implements DevicesRepository {
  DevicesRemoteRepository(this._client, {required this.accessToken});

  final ApiClient _client;
  final String accessToken;

  @override
  Future<void> registerDevice({
    required String accessToken,
    required String deviceToken,
    required String platform,
  }) async {
    await _client.postJson(
      '/devices/register',
      headers: {
        'Authorization': 'Bearer $accessToken',
      },
      body: {
        'device_token': deviceToken,
        'platform': platform,
      },
    );
  }
}
