import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

import 'package:ideacoes_mobile/core/network/api_client.dart';

void main() {
  test('extracts structured api errors', () async {
    final client = ApiClient(
      baseUrl: 'http://localhost',
      client: MockClient(
        (_) async => http.Response(
          '{"error":{"code":"alert_not_found","message":"Alerta não encontrado."}}',
          404,
        ),
      ),
    );

    final call = client.getJson('/alerts/missing');

    await expectLater(
      call,
      throwsA(
        isA<ApiException>()
            .having((error) => error.statusCode, 'statusCode', 404)
            .having((error) => error.code, 'code', 'alert_not_found')
            .having(
                (error) => error.message, 'message', 'Alerta não encontrado.'),
      ),
    );
  });

  test('reports malformed json as invalid response', () async {
    final client = ApiClient(
      baseUrl: 'http://localhost',
      client: MockClient((_) async => http.Response('not-json', 200)),
    );

    final call = client.getJson('/actions');

    await expectLater(
      call,
      throwsA(
        isA<ApiException>()
            .having((error) => error.statusCode, 'statusCode', 200)
            .having((error) => error.code, 'code', 'invalid_response'),
      ),
    );
  });
}
