import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../../../core/storage/session_store.dart';
import '../data/auth_remote_repository.dart';
import '../data/auth_repository.dart';
import '../domain/session.dart';

final sessionControllerProvider =
    NotifierProvider<SessionController, Session?>(SessionController.new);

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRemoteRepository(ref.watch(apiClientProvider));
});

class SessionController extends Notifier<Session?> {
  @override
  Session? build() {
    _restoreSession();
    return null;
  }

  Future<void> _restoreSession() async {
    final store = await ref.read(sessionStoreProvider.future);
    final restored = await store.read();
    state = restored;
  }

  Future<void> signIn({required String userId}) async {
    final token = await ref.read(authRepositoryProvider).issueToken(userId: userId);
    final session = Session(userId: userId, accessToken: token);
    final store = await ref.read(sessionStoreProvider.future);
    await store.write(session);
    state = session;
  }

  Future<void> signOut() async {
    final store = await ref.read(sessionStoreProvider.future);
    await store.clear();
    state = null;
  }
}
