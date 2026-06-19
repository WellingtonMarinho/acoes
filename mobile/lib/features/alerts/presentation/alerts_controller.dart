import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../data/alerts_repository.dart';
import '../data/demo_alerts_repository.dart';
import '../data/alerts_remote_repository.dart';
import '../domain/alert.dart';
import '../../session/presentation/session_controller.dart';

final alertsRepositoryProvider = Provider<AlertsRepository>((ref) {
  final session = ref.watch(sessionControllerProvider);
  return AlertsRemoteRepository(
    ref.watch(apiClientProvider),
    accessToken: session?.accessToken ?? '',
  );
});

final demoAlertsRepositoryProvider = Provider<DemoAlertsRepository>((ref) {
  return DemoAlertsRepository();
});

final alertsControllerProvider =
    AsyncNotifierProvider<AlertsController, List<Alert>>(AlertsController.new);

class AlertsController extends AsyncNotifier<List<Alert>> {
  @override
  Future<List<Alert>> build() async {
    return ref.read(activeAlertsRepositoryProvider).listAlerts();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      return ref.read(activeAlertsRepositoryProvider).listAlerts();
    });
  }
}

final activeAlertsRepositoryProvider = Provider<AlertsRepository>((ref) {
  final config = ref.watch(appConfigProvider);
  final session = ref.watch(sessionControllerProvider);
  if (config.useDemoData || session == null) {
    return ref.watch(demoAlertsRepositoryProvider);
  }
  return ref.watch(alertsRepositoryProvider);
});
