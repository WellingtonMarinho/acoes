import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:ideacoes_mobile/features/alerts/data/alerts_repository.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';
import 'package:ideacoes_mobile/features/alerts/presentation/alerts_controller.dart';
import 'package:ideacoes_mobile/features/home/pages/alerts_page.dart';

void main() {
  testWidgets('shows active alerts and history sections', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          activeAlertsRepositoryProvider
              .overrideWithValue(_FakeAlertsRepository()),
        ],
        child: MaterialApp(
          home: Scaffold(
            body: AlertsPageView(
              onEditAlert: (_) {},
            ),
          ),
        ),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.text('Alertas ativos'), findsOneWidget);
    expect(find.textContaining('PETR4'), findsOneWidget);
    expect(find.textContaining('VALE3'), findsNothing);
    expect(find.text('ativo'), findsOneWidget);

    await tester.tap(find.text('Histórico'));
    await tester.pumpAndSettle();

    expect(find.text('Histórico de alertas'), findsOneWidget);
    expect(find.textContaining('VALE3'), findsOneWidget);
    expect(find.textContaining('PETR4'), findsNothing);
    expect(find.text('disparado'), findsOneWidget);
    expect(find.text('Editar'), findsNothing);
  });
}

class _FakeAlertsRepository implements AlertsRepository {
  @override
  Future<Alert> createAlert({
    required String userId,
    required String actionId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    return Alert(
      id: 'alert-new',
      userId: userId,
      actionId: actionId,
      symbol: 'TEST',
      actionName: 'Test',
      targetPrice: targetPrice,
      direction: direction,
      status: AlertStatus.open,
    );
  }

  @override
  Future<void> deleteAlert({required String alertId}) async {}

  @override
  Future<List<Alert>> listAlerts() async {
    return [
      const Alert(
        id: 'alert-open',
        userId: 'user-1',
        actionId: 'action-petr4',
        symbol: 'PETR4',
        actionName: 'Petrobras PN',
        targetPrice: 40.5,
        direction: AlertDirection.above,
        status: AlertStatus.open,
      ),
      Alert(
        id: 'alert-history',
        userId: 'user-1',
        actionId: 'action-vale3',
        symbol: 'VALE3',
        actionName: 'Vale ON',
        targetPrice: 60.0,
        direction: AlertDirection.below,
        status: AlertStatus.triggered,
        triggeredAt: DateTime.utc(2026, 6, 23, 12, 0),
      ),
    ];
  }

  @override
  Future<Alert> updateAlert({
    required String alertId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    return Alert(
      id: alertId,
      userId: 'user-1',
      actionId: 'action-petr4',
      symbol: 'PETR4',
      actionName: 'Petrobras PN',
      targetPrice: targetPrice,
      direction: direction,
      status: AlertStatus.open,
    );
  }
}
