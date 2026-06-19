import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../domain/alert.dart';

class AlertDraft {
  const AlertDraft({
    this.symbol = '',
    this.targetPrice = '',
    this.direction = AlertDirection.above,
  });

  final String symbol;
  final String targetPrice;
  final AlertDirection direction;

  AlertDraft copyWith({
    String? symbol,
    String? targetPrice,
    AlertDirection? direction,
  }) {
    return AlertDraft(
      symbol: symbol ?? this.symbol,
      targetPrice: targetPrice ?? this.targetPrice,
      direction: direction ?? this.direction,
    );
  }
}

final alertFormControllerProvider =
    NotifierProvider<AlertFormController, AlertDraft>(AlertFormController.new);

class AlertFormController extends Notifier<AlertDraft> {
  @override
  AlertDraft build() => const AlertDraft();

  void setSymbol(String value) {
    state = state.copyWith(symbol: value);
  }

  void setTargetPrice(String value) {
    state = state.copyWith(targetPrice: value);
  }

  void setDirection(AlertDirection value) {
    state = state.copyWith(direction: value);
  }

  void reset() {
    state = const AlertDraft();
  }
}
