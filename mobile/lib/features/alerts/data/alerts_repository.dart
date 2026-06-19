import '../domain/alert.dart';

abstract class AlertsRepository {
  Future<List<Alert>> listAlerts();
  Future<Alert> createAlert({
    required String userId,
    required String symbol,
    required double targetPrice,
    required AlertDirection direction,
  });
}
