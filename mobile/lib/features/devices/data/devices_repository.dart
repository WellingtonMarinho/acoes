abstract class DevicesRepository {
  Future<void> registerDevice({
    required String accessToken,
    required String deviceToken,
    required String platform,
  });
}
