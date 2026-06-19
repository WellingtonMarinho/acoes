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
      symbol: 'B3SA3',
      targetPrice: 13.2,
      direction: AlertDirection.below,
    );

    expect(created.symbol, 'B3SA3');

    final alerts = await repository.listAlerts();
    expect(alerts.first.symbol, 'B3SA3');
  });
}
