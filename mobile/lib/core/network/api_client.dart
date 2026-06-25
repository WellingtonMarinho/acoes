import 'dart:convert';

import 'package:http/http.dart' as http;

class ApiException implements Exception {
  ApiException({
    required this.statusCode,
    required this.code,
    required this.message,
  });

  final int statusCode;
  final String code;
  final String message;

  @override
  String toString() => message;
}

class ApiClient {
  ApiClient({
    required this.baseUrl,
    http.Client? client,
  }) : _client = client ?? http.Client();

  final String baseUrl;
  final http.Client _client;

  Uri _uri(String path) => Uri.parse('$baseUrl$path');

  Future<dynamic> getJson(String path) {
    return _request(() => _client.get(_uri(path)));
  }

  Future<dynamic> getJsonWithHeaders(
    String path, {
    Map<String, String>? headers,
  }) {
    return _request(() => _client.get(_uri(path), headers: headers));
  }

  Future<dynamic> postJson(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) {
    return _request(
      () => _client.post(
        _uri(path),
        headers: {
          'Content-Type': 'application/json',
          if (headers != null) ...headers,
        },
        body: jsonEncode(body ?? const <String, dynamic>{}),
      ),
    );
  }

  Future<dynamic> patchJson(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) {
    return _request(
      () => _client.patch(
        _uri(path),
        headers: {
          'Content-Type': 'application/json',
          if (headers != null) ...headers,
        },
        body: jsonEncode(body ?? const <String, dynamic>{}),
      ),
    );
  }

  Future<dynamic> deleteJson(
    String path, {
    Map<String, String>? headers,
  }) {
    return _request(() => _client.delete(_uri(path), headers: headers));
  }

  Future<dynamic> _request(Future<http.Response> Function() send) async {
    late final http.Response response;
    try {
      response = await send();
    } catch (error) {
      throw ApiException(
        statusCode: 0,
        code: 'network_error',
        message: 'Falha de conexão. Tente novamente.',
      );
    }

    late final dynamic decoded;
    try {
      decoded = _decode(response.body);
    } on FormatException {
      throw ApiException(
        statusCode: response.statusCode,
        code: 'invalid_response',
        message: 'Resposta inválida do servidor.',
      );
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw ApiException(
        statusCode: response.statusCode,
        code: _extractErrorCode(decoded),
        message: _extractErrorMessage(decoded),
      );
    }
    return decoded;
  }

  dynamic _decode(String body) {
    if (body.trim().isEmpty) {
      return <String, dynamic>{};
    }
    return jsonDecode(body);
  }

  String _extractErrorCode(dynamic decoded) {
    if (decoded is Map<String, dynamic>) {
      final error = decoded['error'];
      if (error is Map<String, dynamic>) {
        return error['code']?.toString() ?? 'http_error';
      }
      return decoded['code']?.toString() ?? 'http_error';
    }
    return 'http_error';
  }

  String _extractErrorMessage(dynamic decoded) {
    if (decoded is Map<String, dynamic>) {
      final error = decoded['error'];
      if (error is Map<String, dynamic>) {
        return error['message']?.toString() ?? 'Falha na requisição.';
      }
      return decoded['message']?.toString() ?? 'Falha na requisição.';
    }
    return 'Falha na requisição.';
  }

  static String friendlyMessageFor(ApiException error) {
    switch (error.code) {
      case 'unauthorized':
        return 'Sessão expirada. Entre novamente.';
      case 'network_error':
      case 'invalid_response':
      case 'invalid_json':
      case 'invalid_alert':
      case 'invalid_watchlist_item':
      case 'invalid_price_snapshot':
      case 'not_found':
      case 'action_not_found':
      case 'watchlist_item_not_found':
      case 'alert_not_found':
      case 'conflict':
      case 'alert_not_editable':
        return error.message;
      default:
        return error.message;
    }
  }
}
