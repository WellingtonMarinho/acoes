import '../domain/alert.dart';
import 'alerts_repository.dart';

class DemoAlertsRepository implements AlertsRepository {
  DemoAlertsRepository()
      : _alerts = [
          const Alert(
            id: 'alert-1',
            userId: 'user-demo',
            actionId: 'action-petr4',
            symbol: 'PETR4',
            actionName: 'Petrobras PN',
            targetPrice: 40.50,
            direction: AlertDirection.above,
            status: AlertStatus.open,
          ),
          const Alert(
            id: 'alert-2',
            userId: 'user-demo',
            actionId: 'action-vale3',
            symbol: 'VALE3',
            actionName: 'Vale ON',
            targetPrice: 58.90,
            direction: AlertDirection.below,
            status: AlertStatus.triggered,
          ),
        ];

  final List<Alert> _alerts;

  @override
  Future<List<Alert>> listAlerts() async {
    await Future<void>.delayed(const Duration(milliseconds: 250));
    return List<Alert>.unmodifiable(_alerts);
  }

  @override
  Future<Alert> createAlert({
    required String userId,
    required String actionId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final alert = Alert(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      userId: userId,
      actionId: actionId,
      symbol: actionId,
      actionName: actionId,
      targetPrice: targetPrice,
      direction: direction,
      status: AlertStatus.open,
    );
    _alerts.insert(0, alert);
    return alert;
  }

  @override
  Future<Alert> updateAlert({
    required String alertId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final index = _alerts.indexWhere((alert) => alert.id == alertId);
    if (index < 0) {
      throw StateError('Alert not found');
    }
    final updated = _alerts[index].copyWith(
      targetPrice: targetPrice,
      direction: direction,
      updatedAt: DateTime.now(),
    );
    _alerts[index] = updated;
    return updated;
  }

  @override
  Future<void> deleteAlert({
    required String alertId,
  }) async {
    _alerts.removeWhere((alert) => alert.id == alertId);
  }
}
