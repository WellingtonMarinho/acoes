import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../../session/presentation/session_controller.dart';
import '../data/devices_remote_repository.dart';
import '../data/devices_repository.dart';

final devicesRepositoryProvider = Provider<DevicesRepository>((ref) {
  return DevicesRemoteRepository(ref.watch(apiClientProvider));
});

final activeDevicesRepositoryProvider = Provider<DevicesRepository>((ref) {
  final config = ref.watch(appConfigProvider);
  final session = ref.watch(sessionControllerProvider);
  if (config.useDemoData || session == null) {
    return _DemoDevicesRepository();
  }
  return ref.watch(devicesRepositoryProvider);
});

class _DemoDevicesRepository implements DevicesRepository {
  @override
  Future<void> registerDevice({
    required String accessToken,
    required String deviceToken,
    required String platform,
  }) async {
    await Future<void>.delayed(const Duration(milliseconds: 150));
  }
}
