enum AlertDirection { above, below }

enum AlertStatus { open, triggered }

class Alert {
  const Alert({
    required this.id,
    required this.userId,
    required this.symbol,
    required this.targetPrice,
    required this.direction,
    required this.status,
  });

  final String id;
  final String userId;
  final String symbol;
  final double targetPrice;
  final AlertDirection direction;
  final AlertStatus status;

  factory Alert.fromJson(Map<String, dynamic> json) {
    return Alert(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      symbol: json['symbol'] as String? ?? '',
      targetPrice: (json['target_price'] as num?)?.toDouble() ?? 0,
      direction: _directionFromJson(json['direction'] as String?),
      status: _statusFromJson(json['status'] as String?),
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
