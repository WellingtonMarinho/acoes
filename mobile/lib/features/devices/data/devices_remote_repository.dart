import '../../../core/network/api_client.dart';
import 'devices_repository.dart';

class DevicesRemoteRepository implements DevicesRepository {
  DevicesRemoteRepository(this._client);

  final ApiClient _client;

  @override
  Future<void> registerDevice({
    required String deviceToken,
    required String platform,
  }) async {
    await _client.postJson(
      '/devices/register',
      body: {
        'device_token': deviceToken,
        'platform': platform,
      },
    );
  }
}
