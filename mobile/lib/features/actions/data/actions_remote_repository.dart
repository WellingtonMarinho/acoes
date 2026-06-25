import '../../../core/network/api_client.dart';
import '../domain/action.dart';
import 'actions_repository.dart';

class ActionsRemoteRepository implements ActionsRepository {
  ActionsRemoteRepository(this._client);

  final ApiClient _client;

  @override
  Future<List<MarketAction>> listActions({String query = ''}) async {
    final path = query.trim().isEmpty
        ? '/actions'
        : '/actions?query=${Uri.encodeComponent(query.trim())}';
    final data = await _client.getJson(path);
    final items = data['items'];
    final source = items is List<dynamic> ? items : <dynamic>[];
    return [
      for (final item in source)
      MarketAction.fromJson(item as Map<String, dynamic>),
    ];
  }

  @override
  Future<MarketAction> createAction({
    required String symbol,
    required String name,
    String exchange = '',
  }) async {
    final data = await _client.postJson(
      '/actions',
      body: {
        'symbol': symbol,
        'name': name,
        'exchange': exchange,
      },
    );
    return MarketAction.fromJson(data as Map<String, dynamic>);
  }
}
