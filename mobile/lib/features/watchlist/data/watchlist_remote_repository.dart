import '../../../core/network/api_client.dart';
import '../domain/watchlist_item.dart';
import 'watchlist_repository.dart';

class WatchlistRemoteRepository implements WatchlistRepository {
  WatchlistRemoteRepository(this._client);

  final ApiClient _client;

  @override
  Future<WatchlistItem> addWatchlist(String actionId) async {
    final data = await _client.postJson(
      '/watchlist',
      body: {'action_id': actionId},
    );
    return WatchlistItem.fromJson(_asMap(data));
  }

  @override
  Future<List<WatchlistItem>> listWatchlist() async {
    final data = await _client.getJsonWithHeaders(
      '/watchlist',
    );
    final items = _itemsFromResponse(data);
    return [
      for (final item in items) WatchlistItem.fromJson(item),
    ];
  }

  @override
  Future<void> removeWatchlist(String actionId) async {
    await _client.deleteJson(
      '/watchlist/$actionId',
    );
  }

  Map<String, dynamic> _asMap(dynamic data) {
    if (data is Map<String, dynamic>) {
      return data;
    }
    return <String, dynamic>{};
  }

  List<Map<String, dynamic>> _itemsFromResponse(dynamic data) {
    if (data is Map<String, dynamic>) {
      final items = data['items'];
      if (items is List<dynamic>) {
        return [
          for (final item in items)
            if (item is Map<String, dynamic>) item,
        ];
      }
    }
    return <Map<String, dynamic>>[];
  }
}
