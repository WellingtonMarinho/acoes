import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../actions/domain/action.dart';
import '../domain/alert.dart';

class AlertDraft {
  const AlertDraft({
    this.action,
    this.targetPrice = '',
    this.direction = AlertDirection.above,
  });

  final MarketAction? action;
  final String targetPrice;
  final AlertDirection direction;

  AlertDraft copyWith({
    MarketAction? action,
    String? targetPrice,
    AlertDirection? direction,
  }) {
    return AlertDraft(
      action: action ?? this.action,
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

  void setAction(MarketAction value) {
    state = state.copyWith(action: value);
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
