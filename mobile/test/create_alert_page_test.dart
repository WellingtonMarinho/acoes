import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:ideacoes_mobile/features/alerts/presentation/create_alert_page.dart';

void main() {
  testWidgets('shows create alert form', (tester) async {
    await tester.pumpWidget(
      const ProviderScope(
        child: MaterialApp(home: CreateAlertPage()),
      ),
    );

    expect(find.text('Novo alerta'), findsOneWidget);
    expect(find.text('Salvar alerta'), findsOneWidget);
    expect(find.text('Ativo'), findsOneWidget);
    expect(find.text('Preço alvo'), findsOneWidget);
  });
}
