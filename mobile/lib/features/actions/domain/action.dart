class MarketAction {
  const MarketAction({
    required this.id,
    required this.symbol,
    required this.name,
    required this.exchange,
  });

  final String id;
  final String symbol;
  final String name;
  final String exchange;

  factory MarketAction.fromJson(Map<String, dynamic> json) {
    return MarketAction(
      id: json['id'] as String? ?? '',
      symbol: json['symbol'] as String? ?? '',
      name: json['name'] as String? ?? '',
      exchange: json['exchange'] as String? ?? '',
    );
  }
}
