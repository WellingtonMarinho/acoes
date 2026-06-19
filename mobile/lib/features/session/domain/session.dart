class Session {
  const Session({
    required this.userId,
    required this.accessToken,
  });

  final String userId;
  final String accessToken;

  Session copyWith({
    String? userId,
    String? accessToken,
  }) {
    return Session(
      userId: userId ?? this.userId,
      accessToken: accessToken ?? this.accessToken,
    );
  }
}
