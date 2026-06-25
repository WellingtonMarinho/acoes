import '../domain/alert.dart';

abstract class AlertsRepository {
  Future<List<Alert>> listAlerts();
  Future<Alert> createAlert({
    required String userId,
    required String actionId,
    required double targetPrice,
    required AlertDirection direction,
  });
  Future<Alert> updateAlert({
    required String alertId,
    required double targetPrice,
    required AlertDirection direction,
  });
  Future<void> deleteAlert({
    required String alertId,
  });
}
