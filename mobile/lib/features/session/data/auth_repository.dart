abstract class AuthRepository {
  Future<String> issueToken({required String userId});
}
