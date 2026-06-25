import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../data/actions_remote_repository.dart';
import '../data/actions_repository.dart';
import '../domain/action.dart';

final actionsRepositoryProvider = Provider<ActionsRepository>((ref) {
  return ActionsRemoteRepository(ref.watch(apiClientProvider));
});

final actionsControllerProvider =
    AsyncNotifierProvider<ActionsController, List<MarketAction>>(
        ActionsController.new);

class ActionsController extends AsyncNotifier<List<MarketAction>> {
  @override
  Future<List<MarketAction>> build() async {
    return ref.read(actionsRepositoryProvider).listActions();
  }
}
