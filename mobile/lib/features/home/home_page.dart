import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../alerts/presentation/alerts_controller.dart';
import '../alerts/presentation/widgets/alert_list_item.dart';
import '../session/presentation/session_card.dart';
import 'widgets/action_tile.dart';
import 'widgets/section_card.dart';

class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final alerts = ref.watch(alertsControllerProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Ideacoes'),
        actions: [
          IconButton(
            onPressed: () => ref.read(alertsControllerProvider.notifier).refresh(),
            icon: const Icon(Icons.refresh),
            tooltip: 'Atualizar alertas',
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          const _HeroCard(),
          const SizedBox(height: 20),
          const SessionCard(),
          const SizedBox(height: 20),
          SectionCard(
            title: 'Alertas',
            child: alerts.when(
              data: (items) {
                if (items.isEmpty) {
                  return const Text('Nenhum alerta encontrado por enquanto.');
                }

                return Column(
                  children: [
                    for (var i = 0; i < items.length; i++) ...[
                      AlertListItem(alert: items[i]),
                      if (i < items.length - 1) const SizedBox(height: 12),
                    ],
                  ],
                );
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (error, stackTrace) => Text('Falha ao carregar alertas: $error'),
            ),
          ),
          const SizedBox(height: 20),
          SectionCard(
            title: 'Próximas ações',
            child: Column(
              children: [
                ActionTile(
                  title: 'Registrar device',
                  subtitle: 'Preparar o token de push do usuário',
                  onTap: () => context.push('/devices/register'),
                ),
                const SizedBox(height: 12),
                ActionTile(
                  title: 'Criar alerta',
                  subtitle: 'Definir ativo, direção e preço alvo',
                  onTap: () => context.push('/alerts/new'),
                ),
                const SizedBox(height: 12),
                ActionTile(
                  title: 'Conectar API',
                  subtitle: 'Buscar alertas e estados no backend',
                  onTap: () {},
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _HeroCard extends StatelessWidget {
  const _HeroCard();

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;

    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [colors.primary, colors.primaryContainer],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(28),
      ),
      child: const Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'MVP de alertas para ações',
            style: TextStyle(
              color: Colors.white,
              fontSize: 26,
              fontWeight: FontWeight.w700,
            ),
          ),
          SizedBox(height: 8),
          Text(
            'Base pronta para criar alertas, registrar device e integrar o backend.',
            style: TextStyle(color: Colors.white70, height: 1.4),
          ),
        ],
      ),
    );
  }
}
