import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:ideacoes_mobile/features/alerts/presentation/alert_form_controller.dart';
import 'package:ideacoes_mobile/features/alerts/domain/alert.dart';

void main() {
  test('alert form controller updates state', () {
    final container = ProviderContainer();
    addTearDown(container.dispose);

    final controller = container.read(alertFormControllerProvider.notifier);
    controller.setSymbol('PETR4');
    controller.setTargetPrice('41.10');
    controller.setDirection(AlertDirection.below);

    final draft = container.read(alertFormControllerProvider);
    expect(draft.symbol, 'PETR4');
    expect(draft.targetPrice, '41.10');
    expect(draft.direction, AlertDirection.below);
  });
}
