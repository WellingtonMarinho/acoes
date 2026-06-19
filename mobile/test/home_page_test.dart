import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:ideacoes_mobile/app/app.dart';

void main() {
  testWidgets('shows the home shell', (tester) async {
    await tester.pumpWidget(const ProviderScope(child: IdeacoesApp()));
    await tester.pumpAndSettle();

    expect(find.text('Ideacoes'), findsOneWidget);
    expect(find.text('MVP de alertas para ações'), findsOneWidget);
    expect(find.text('Entrar no app'), findsOneWidget);
  });
}
