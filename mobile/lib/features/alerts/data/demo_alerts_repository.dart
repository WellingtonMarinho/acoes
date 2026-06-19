import '../domain/alert.dart';
import 'alerts_repository.dart';

class DemoAlertsRepository implements AlertsRepository {
  DemoAlertsRepository()
      : _alerts = [
          const Alert(
            id: 'alert-1',
            userId: 'user-demo',
            symbol: 'PETR4',
            targetPrice: 40.50,
            direction: AlertDirection.above,
            status: AlertStatus.open,
          ),
          const Alert(
            id: 'alert-2',
            userId: 'user-demo',
            symbol: 'VALE3',
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
    required String symbol,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    final alert = Alert(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      userId: userId,
      symbol: symbol,
      targetPrice: targetPrice,
      direction: direction,
      status: AlertStatus.open,
    );
    _alerts.insert(0, alert);
    return alert;
  }
}
