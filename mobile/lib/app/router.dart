import 'package:go_router/go_router.dart';

import '../features/actions/domain/action.dart' as stock_action;
import '../features/alerts/domain/alert.dart';
import '../features/alerts/presentation/create_alert_page.dart';
import '../features/devices/presentation/device_registration_page.dart';
import '../features/home/home_shell.dart';

final appRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) => const HomeShell(),
    ),
    GoRoute(
      path: '/alerts/new',
      builder: (context, state) => CreateAlertPage(
        initialAction: state.extra is stock_action.MarketAction
            ? state.extra as stock_action.MarketAction
            : null,
      ),
    ),
    GoRoute(
      path: '/alerts/edit',
      builder: (context, state) => EditAlertPage(
        alert: state.extra as Alert,
      ),
    ),
    GoRoute(
      path: '/devices/register',
      builder: (context, state) => const DeviceRegistrationPage(),
    ),
  ],
);
