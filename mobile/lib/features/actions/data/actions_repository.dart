import '../domain/action.dart';

abstract class ActionsRepository {
  Future<List<MarketAction>> listActions({String query = ''});
  Future<MarketAction> createAction({
    required String symbol,
    required String name,
    String exchange = '',
  });
}
