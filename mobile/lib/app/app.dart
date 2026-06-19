import 'package:flutter/material.dart';

import 'theme.dart';
import 'router.dart';

class IdeacoesApp extends StatelessWidget {
  const IdeacoesApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      title: 'Ideacoes',
      theme: buildAppTheme(),
      routerConfig: appRouter,
    );
  }
}
