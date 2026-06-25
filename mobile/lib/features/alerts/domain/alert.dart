enum AlertDirection { above, below }

enum AlertStatus { open, triggered }

class Alert {
  const Alert({
    required this.id,
    required this.userId,
    required this.actionId,
    required this.symbol,
    required this.actionName,
    required this.targetPrice,
    required this.direction,
    required this.status,
    this.updatedAt,
    this.triggeredAt,
  });

  final String id;
  final String userId;
  final String actionId;
  final String symbol;
  final String actionName;
  final double targetPrice;
  final AlertDirection direction;
  final AlertStatus status;
  final DateTime? updatedAt;
  final DateTime? triggeredAt;

  factory Alert.fromJson(Map<String, dynamic> json) {
    return Alert(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      actionId: json['action_id'] as String? ?? '',
      symbol: json['symbol'] as String? ?? '',
      actionName: json['action_name'] as String? ?? '',
      targetPrice: (json['target_price'] as num?)?.toDouble() ?? 0,
      direction: _directionFromJson(json['direction'] as String?),
      status: _statusFromJson(json['status'] as String?),
      updatedAt: DateTime.tryParse(json['updated_at'] as String? ?? ''),
      triggeredAt: DateTime.tryParse(json['triggered_at'] as String? ?? ''),
    );
  }

  Alert copyWith({
    String? id,
    String? userId,
    String? actionId,
    String? symbol,
    String? actionName,
    double? targetPrice,
    AlertDirection? direction,
    AlertStatus? status,
    DateTime? updatedAt,
    DateTime? triggeredAt,
  }) {
    return Alert(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      actionId: actionId ?? this.actionId,
      symbol: symbol ?? this.symbol,
      actionName: actionName ?? this.actionName,
      targetPrice: targetPrice ?? this.targetPrice,
      direction: direction ?? this.direction,
      status: status ?? this.status,
      updatedAt: updatedAt ?? this.updatedAt,
      triggeredAt: triggeredAt ?? this.triggeredAt,
    );
  }
}

AlertDirection _directionFromJson(String? value) {
  return switch (value) {
    'above' => AlertDirection.above,
    'below' => AlertDirection.below,
    _ => AlertDirection.above,
  };
}

AlertStatus _statusFromJson(String? value) {
  return switch (value) {
    'triggered' => AlertStatus.triggered,
    _ => AlertStatus.open,
  };
}
