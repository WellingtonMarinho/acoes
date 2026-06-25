import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:ideacoes_mobile/app/app.dart';
import 'package:ideacoes_mobile/features/alerts/data/alerts_repository.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';
import 'package:ideacoes_mobile/features/alerts/presentation/alerts_controller.dart';
import 'package:ideacoes_mobile/features/watchlist/data/watchlist_repository.dart';
import 'package:ideacoes_mobile/features/watchlist/domain/watchlist_item.dart';
import 'package:ideacoes_mobile/features/watchlist/presentation/watchlist_controller.dart';
import 'package:ideacoes_mobile/features/session/data/auth_repository.dart';
import 'package:ideacoes_mobile/features/session/presentation/session_controller.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('shows the app shell', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          authRepositoryProvider.overrideWithValue(_FakeAuthRepository()),
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
    expect(find.text('Alertas'), findsAtLeastNWidgets(1));
    expect(find.text('Ajustes'), findsAtLeastNWidgets(1));
    expect(find.byTooltip('Adicionar ação'), findsOneWidget);
    expect(find.text('Nenhuma ação monitorada ainda.'), findsOneWidget);
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
