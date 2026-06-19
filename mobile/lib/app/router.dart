import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../features/alerts/presentation/create_alert_page.dart';
import '../features/devices/presentation/device_registration_page.dart';
import '../features/home/home_page.dart';

final appRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) => const HomePage(),
    ),
    GoRoute(
      path: '/alerts/new',
      builder: (context, state) => const CreateAlertPage(),
    ),
    GoRoute(
      path: '/devices/register',
      builder: (context, state) => const DeviceRegistrationPage(),
    ),
  ],
);
