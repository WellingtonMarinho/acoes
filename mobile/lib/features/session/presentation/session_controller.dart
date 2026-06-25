import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
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
    return null;
  }

  Future<void> signIn({required String userId}) async {
    state = Session(userId: userId, accessToken: '');
  }

  Future<void> signOut() async {
    state = null;
  }
}
