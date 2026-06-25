import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../data/devices_remote_repository.dart';
import '../data/devices_repository.dart';

final devicesRepositoryProvider = Provider<DevicesRepository>((ref) {
  return DevicesRemoteRepository(ref.watch(apiClientProvider));
});

final activeDevicesRepositoryProvider = Provider<DevicesRepository>((ref) {
  return ref.watch(devicesRepositoryProvider);
});
