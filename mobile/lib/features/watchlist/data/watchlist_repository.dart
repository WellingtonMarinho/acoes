import '../domain/watchlist_item.dart';

abstract class WatchlistRepository {
  Future<List<WatchlistItem>> listWatchlist();
  Future<WatchlistItem> addWatchlist(String actionId);
  Future<void> removeWatchlist(String actionId);
}
