import 'package:shared_preferences/shared_preferences.dart';

import '../../features/session/domain/session.dart';

class SessionStore {
  SessionStore(this._prefs);

  final SharedPreferences _prefs;

  static const _userIdKey = 'session.user_id';
  static const _accessTokenKey = 'session.access_token';

  Future<Session?> read() async {
    final userId = _prefs.getString(_userIdKey);
    final accessToken = _prefs.getString(_accessTokenKey);

    if (userId == null || accessToken == null || userId.isEmpty || accessToken.isEmpty) {
      return null;
    }

    return Session(userId: userId, accessToken: accessToken);
  }

  Future<void> write(Session session) async {
    await _prefs.setString(_userIdKey, session.userId);
    await _prefs.setString(_accessTokenKey, session.accessToken);
  }

  Future<void> clear() async {
    await _prefs.remove(_userIdKey);
    await _prefs.remove(_accessTokenKey);
  }
}
