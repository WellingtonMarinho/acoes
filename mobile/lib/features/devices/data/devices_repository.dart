abstract class DevicesRepository {
  Future<void> registerDevice({
    required String deviceToken,
    required String platform,
  });
}
