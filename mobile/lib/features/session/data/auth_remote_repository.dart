import '../../../core/network/api_client.dart';
import 'auth_repository.dart';

class AuthRemoteRepository implements AuthRepository {
  AuthRemoteRepository(this._client);

  final ApiClient _client;

  @override
  Future<String> issueToken({required String userId}) async {
    final data = await _client.postJson(
      '/auth/token',
      body: {'user_id': userId},
    );
    return data['access_token'] as String? ?? '';
  }
}
