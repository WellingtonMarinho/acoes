class WatchlistItem {
  const WatchlistItem({
    required this.actionId,
    required this.symbol,
    required this.name,
    required this.exchange,
    required this.openAlertsCount,
    this.currentPrice,
    this.lastPriceAt,
  });

  final String actionId;
  final String symbol;
  final String name;
  final String exchange;
  final int openAlertsCount;
  final double? currentPrice;
  final DateTime? lastPriceAt;

  factory WatchlistItem.fromJson(Map<String, dynamic> json) {
    return WatchlistItem(
      actionId: json['action_id'] as String? ?? '',
      symbol: json['symbol'] as String? ?? '',
      name: json['name'] as String? ?? '',
      exchange: json['exchange'] as String? ?? '',
      openAlertsCount: (json['open_alerts_count'] as num?)?.toInt() ?? 0,
      currentPrice: (json['current_price'] as num?)?.toDouble(),
      lastPriceAt: DateTime.tryParse(json['last_price_at'] as String? ?? ''),
    );
  }
}
