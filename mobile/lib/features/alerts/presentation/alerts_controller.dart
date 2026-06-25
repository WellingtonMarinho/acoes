import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../data/alerts_repository.dart';
import '../data/alerts_remote_repository.dart';
import '../domain/alert.dart';

final alertsRepositoryProvider = Provider<AlertsRepository>((ref) {
  return AlertsRemoteRepository(ref.watch(apiClientProvider));
});

final alertsControllerProvider =
    AsyncNotifierProvider<AlertsController, List<Alert>>(AlertsController.new);

class AlertsController extends AsyncNotifier<List<Alert>> {
  @override
  Future<List<Alert>> build() async {
    return ref.watch(activeAlertsRepositoryProvider).listAlerts();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      return ref.read(activeAlertsRepositoryProvider).listAlerts();
    });
  }
}

final activeAlertsRepositoryProvider = Provider<AlertsRepository>((ref) {
  return ref.watch(alertsRepositoryProvider);
});
