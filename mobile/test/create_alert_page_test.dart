import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:ideacoes_mobile/features/actions/data/actions_repository.dart';
import 'package:ideacoes_mobile/features/actions/domain/action.dart'
    as stock_action;
import 'package:ideacoes_mobile/features/actions/presentation/actions_controller.dart';
import 'package:ideacoes_mobile/features/alerts/data/alerts_repository.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';
import 'package:ideacoes_mobile/features/alerts/presentation/alerts_controller.dart';
import 'package:ideacoes_mobile/features/alerts/presentation/create_alert_page.dart';
import 'package:ideacoes_mobile/features/session/data/auth_repository.dart';
import 'package:ideacoes_mobile/features/session/presentation/session_controller.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('shows create alert form', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          authRepositoryProvider.overrideWithValue(_FakeAuthRepository()),
          actionsRepositoryProvider.overrideWithValue(_FakeActionsRepository()),
          activeAlertsRepositoryProvider
              .overrideWithValue(_FakeAlertsRepository()),
        ],
        child: const MaterialApp(home: CreateAlertPage()),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.text('Novo alerta'), findsOneWidget);
    expect(find.text('Salvar alerta'), findsOneWidget);
    expect(find.text('Ação'), findsOneWidget);
    expect(find.text('Preço alvo'), findsOneWidget);
  });
}

class _FakeAuthRepository implements AuthRepository {
  @override
  Future<String> issueToken({required String userId}) async {
    return 'token-for-$userId';
  }
}

class _FakeActionsRepository implements ActionsRepository {
  @override
  Future<List<stock_action.MarketAction>> listActions(
      {String query = ''}) async {
    return const [
      stock_action.MarketAction(
          id: 'action-petr4',
          symbol: 'PETR4',
          name: 'Petrobras PN',
          exchange: 'B3'),
    ];
  }

  @override
  Future<stock_action.MarketAction> createAction({
    required String symbol,
    required String name,
    String exchange = '',
  }) async {
    return stock_action.MarketAction(
      id: 'action-${symbol.toLowerCase()}',
      symbol: symbol,
      name: name,
      exchange: exchange,
    );
  }
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
      id: 'alert-1',
      userId: userId,
      actionId: actionId,
      symbol: 'PETR4',
      actionName: 'Petrobras PN',
      targetPrice: targetPrice,
      direction: direction,
      status: AlertStatus.open,
    );
  }

  @override
  Future<List<Alert>> listAlerts() async {
    return const [];
  }

  @override
  Future<Alert> updateAlert({
    required String alertId,
    required double targetPrice,
    required AlertDirection direction,
  }) async {
    return createAlert(
      userId: 'user-1',
      actionId: alertId,
      targetPrice: targetPrice,
      direction: direction,
    );
  }

  @override
  Future<void> deleteAlert({required String alertId}) async {}
}
