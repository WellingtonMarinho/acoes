import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:ideacoes_mobile/app/app.dart';
import 'package:ideacoes_mobile/features/actions/data/actions_repository.dart';
import 'package:ideacoes_mobile/features/actions/domain/action.dart' as stock_action;
import 'package:ideacoes_mobile/features/actions/presentation/actions_controller.dart';
import 'package:ideacoes_mobile/features/alerts/data/alerts_repository.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';
import 'package:ideacoes_mobile/features/alerts/presentation/alerts_controller.dart';
import 'package:ideacoes_mobile/core/network/api_client.dart';
import 'package:ideacoes_mobile/features/session/data/auth_repository.dart';
import 'package:ideacoes_mobile/features/session/presentation/session_controller.dart';
import 'package:ideacoes_mobile/features/watchlist/data/watchlist_repository.dart';
import 'package:ideacoes_mobile/features/watchlist/domain/watchlist_item.dart';
import 'package:ideacoes_mobile/features/watchlist/presentation/watchlist_controller.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('shows the home shell', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          authRepositoryProvider.overrideWithValue(_FakeAuthRepository()),
          actionsRepositoryProvider.overrideWithValue(_FailingActionsRepository()),
          activeAlertsRepositoryProvider
              .overrideWithValue(_FakeAlertsRepository()),
          activeWatchlistRepositoryProvider
              .overrideWithValue(_FakeWatchlistRepository()),
        ],
        child: const IdeacoesApp(),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.byType(NavigationBar), findsOneWidget);
    expect(find.text('Monitoradas'), findsAtLeastNWidgets(1));
    expect(find.text('Adicionar ação'), findsAtLeastNWidgets(1));
    expect(find.text('Nenhuma ação monitorada ainda.'), findsOneWidget);

    await tester.tap(find.text('Ajustes'));
    await tester.pumpAndSettle();

    expect(find.text('Tema'), findsOneWidget);
    expect(find.text('Registrar device'), findsOneWidget);
  });

  testWidgets('shows inline error when creating action fails', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          authRepositoryProvider.overrideWithValue(_FakeAuthRepository()),
          actionsRepositoryProvider.overrideWithValue(_FailingActionsRepository()),
          activeAlertsRepositoryProvider
              .overrideWithValue(_FakeAlertsRepository()),
          activeWatchlistRepositoryProvider
              .overrideWithValue(_FakeWatchlistRepository()),
        ],
        child: const IdeacoesApp(),
      ),
    );

    await tester.pumpAndSettle();
    await tester.tap(find.text('Adicionar ação').first);
    await tester.pumpAndSettle();

    await tester.enterText(find.byType(TextFormField).at(0), 'ABCD3');
    await tester.enterText(find.byType(TextFormField).at(1), 'Acao Teste');
    await tester.enterText(find.byType(TextFormField).at(2), 'Acao Teste SA');
    await tester.tap(find.text('Cadastrar e adicionar'));
    await tester.pumpAndSettle();

    expect(find.textContaining('Sessão inválida'), findsAtLeastNWidgets(1));
  });
}

class _FakeAuthRepository implements AuthRepository {
  @override
  Future<String> issueToken({required String userId}) async {
    return 'token-for-$userId';
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

class _FakeWatchlistRepository implements WatchlistRepository {
  @override
  Future<WatchlistItem> addWatchlist(String actionId) async {
    return WatchlistItem(
      actionId: actionId,
      symbol: 'PETR4',
      name: 'Petrobras PN',
      exchange: 'B3',
      openAlertsCount: 0,
    );
  }

  @override
  Future<List<WatchlistItem>> listWatchlist() async {
    return const [];
  }

  @override
  Future<void> removeWatchlist(String actionId) async {}
}

class _FailingActionsRepository implements ActionsRepository {
  @override
  Future<List<stock_action.MarketAction>> listActions({String query = ''}) async {
    return const [
      stock_action.MarketAction(
        id: 'action-petr4',
        symbol: 'PETR4',
        name: 'Petrobras PN',
        exchange: 'B3',
      ),
    ];
  }

  @override
  Future<stock_action.MarketAction> createAction({
    required String symbol,
    required String name,
    String exchange = '',
  }) async {
    throw ApiException(
      statusCode: 401,
      code: 'unauthorized',
      message: 'Sessão inválida. Feche e entre novamente para continuar.',
    );
  }
}
