import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../data/watchlist_remote_repository.dart';
import '../data/watchlist_repository.dart';
import '../domain/watchlist_item.dart';

final watchlistRepositoryProvider = Provider<WatchlistRepository>((ref) {
  return WatchlistRemoteRepository(ref.watch(apiClientProvider));
});

final watchlistControllerProvider =
    AsyncNotifierProvider<WatchlistController, List<WatchlistItem>>(
        WatchlistController.new);

class WatchlistController extends AsyncNotifier<List<WatchlistItem>> {
  @override
  Future<List<WatchlistItem>> build() async {
    return ref.watch(activeWatchlistRepositoryProvider).listWatchlist();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      return ref.read(activeWatchlistRepositoryProvider).listWatchlist();
    });
  }
}

final activeWatchlistRepositoryProvider = Provider<WatchlistRepository>((ref) {
  return ref.watch(watchlistRepositoryProvider);
});
