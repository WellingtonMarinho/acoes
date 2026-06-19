import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers.dart';
import '../domain/alert.dart';
import 'alert_form_controller.dart';
import 'alerts_controller.dart';
import '../../session/presentation/session_controller.dart';

class CreateAlertPage extends ConsumerStatefulWidget {
  const CreateAlertPage({super.key});

  @override
  ConsumerState<CreateAlertPage> createState() => _CreateAlertPageState();
}

class _CreateAlertPageState extends ConsumerState<CreateAlertPage> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _symbolController;
  late final TextEditingController _priceController;
  bool _submitting = false;

  @override
  void initState() {
    super.initState();
    final draft = ref.read(alertFormControllerProvider);
    _symbolController = TextEditingController(text: draft.symbol);
    _priceController = TextEditingController(text: draft.targetPrice);
  }

  @override
  void dispose() {
    _symbolController.dispose();
    _priceController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    setState(() {
      _submitting = true;
    });

    final draft = ref.read(alertFormControllerProvider);
    final symbol = draft.symbol.trim().toUpperCase();
    final targetPrice = double.parse(draft.targetPrice.replaceAll(',', '.'));
    final session = ref.read(sessionControllerProvider);

    if (session == null && !ref.read(appConfigProvider).useDemoData) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Entre no app para salvar no backend.')),
      );
      setState(() {
        _submitting = false;
      });
      return;
    }

    try {
      await ref.read(activeAlertsRepositoryProvider).createAlert(
            userId: session?.userId ?? 'user-demo',
            symbol: symbol,
            targetPrice: targetPrice,
            direction: draft.direction,
          );
      ref.invalidate(alertsControllerProvider);

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Alerta de $symbol salvo com sucesso.'),
        ),
      );
      ref.read(alertFormControllerProvider.notifier).reset();
      _symbolController.clear();
      _priceController.clear();
      if (mounted) {
        Navigator.of(context).pop();
      }
    } catch (error) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Falha ao salvar alerta: $error')),
      );
    } finally {
      if (mounted) {
        setState(() {
          _submitting = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final draft = ref.watch(alertFormControllerProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Novo alerta')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            TextFormField(
              controller: _symbolController,
              decoration: const InputDecoration(
                labelText: 'Ativo',
                hintText: 'PETR4',
              ),
              textCapitalization: TextCapitalization.characters,
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return 'Informe o ativo.';
                }
                return null;
              },
              onChanged: ref.read(alertFormControllerProvider.notifier).setSymbol,
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _priceController,
              decoration: const InputDecoration(
                labelText: 'Preço alvo',
                hintText: '40.50',
              ),
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              validator: (value) {
                final parsed = double.tryParse((value ?? '').replaceAll(',', '.'));
                if (parsed == null || parsed <= 0) {
                  return 'Informe um preço válido.';
                }
                return null;
              },
              onChanged: ref.read(alertFormControllerProvider.notifier).setTargetPrice,
            ),
            const SizedBox(height: 16),
            SegmentedButton<AlertDirection>(
              segments: const [
                ButtonSegment(
                  value: AlertDirection.above,
                  label: Text('Acima'),
                  icon: Icon(Icons.trending_up),
                ),
                ButtonSegment(
                  value: AlertDirection.below,
                  label: Text('Abaixo'),
                  icon: Icon(Icons.trending_down),
                ),
              ],
              selected: {draft.direction},
              onSelectionChanged: (value) {
                ref.read(alertFormControllerProvider.notifier).setDirection(value.first);
              },
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: _submitting ? null : _submit,
              child: Text(_submitting ? 'Salvando...' : 'Salvar alerta'),
            ),
          ],
        ),
      ),
    );
  }
}
