import 'package:flutter_test/flutter_test.dart';

import 'package:ideacoes_mobile/features/alerts/data/demo_alerts_repository.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';

void main() {
  test('demo alerts repository returns seeded alerts', () async {
    final repository = DemoAlertsRepository();

    final alerts = await repository.listAlerts();

    expect(alerts, hasLength(2));
    expect(alerts.first.symbol, 'PETR4');
    expect(alerts.first.direction, AlertDirection.above);
  });

  test('demo alerts repository accepts new alerts', () async {
    final repository = DemoAlertsRepository();

    final created = await repository.createAlert(
      userId: 'user-1',
      actionId: 'action-b3sa3',
      targetPrice: 13.2,
      direction: AlertDirection.below,
    );

    expect(created.actionId, 'action-b3sa3');

    final alerts = await repository.listAlerts();
    expect(alerts.first.symbol, 'action-b3sa3');
  });

  test('demo alerts repository updates and deletes alerts', () async {
    final repository = DemoAlertsRepository();

    final created = await repository.createAlert(
      userId: 'user-1',
      actionId: 'action-petr4',
      targetPrice: 40.5,
      direction: AlertDirection.above,
    );

    final updated = await repository.updateAlert(
      alertId: created.id,
      targetPrice: 44.0,
      direction: AlertDirection.below,
    );

    expect(updated.targetPrice, 44.0);
    expect(updated.direction, AlertDirection.below);

    await repository.deleteAlert(alertId: created.id);

    final alerts = await repository.listAlerts();
    expect(alerts.where((alert) => alert.id == created.id), isEmpty);
  });
}
