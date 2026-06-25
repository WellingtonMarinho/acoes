import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/network/api_client.dart';
import '../../actions/domain/action.dart' as stock_action;
import '../../actions/presentation/actions_controller.dart';
import '../../session/presentation/session_controller.dart';
import '../../watchlist/presentation/watchlist_controller.dart';
import '../domain/alert.dart';
import 'alerts_controller.dart';

class CreateAlertPage extends StatelessWidget {
  const CreateAlertPage({super.key, this.initialAction});

  final stock_action.MarketAction? initialAction;

  @override
  Widget build(BuildContext context) {
    return AlertFormPage(
      title: 'Novo alerta',
      submitLabel: 'Salvar alerta',
      initialAction: initialAction,
    );
  }
}

class EditAlertPage extends StatelessWidget {
  const EditAlertPage({
    super.key,
    required this.alert,
  });

  final Alert alert;

  @override
  Widget build(BuildContext context) {
    return AlertFormPage(
      title: 'Editar alerta',
      submitLabel: 'Salvar alterações',
      alert: alert,
      allowActionSelection: false,
    );
  }
}

class AlertFormPage extends ConsumerStatefulWidget {
  const AlertFormPage({
    super.key,
    required this.title,
    required this.submitLabel,
    this.initialAction,
    this.alert,
    this.allowActionSelection = true,
  });

  final String title;
  final String submitLabel;
  final stock_action.MarketAction? initialAction;
  final Alert? alert;
  final bool allowActionSelection;

  @override
  ConsumerState<AlertFormPage> createState() => _AlertFormPageState();
}

class _AlertFormPageState extends ConsumerState<AlertFormPage> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _priceController;
  stock_action.MarketAction? _selectedAction;
  AlertDirection _direction = AlertDirection.above;
  bool _submitting = false;

  bool get _editing => widget.alert != null;

  @override
  void initState() {
    super.initState();
    _selectedAction = widget.initialAction ??
        (widget.alert == null
            ? null
            : stock_action.MarketAction(
                id: widget.alert!.actionId,
                symbol: widget.alert!.symbol,
                name: widget.alert!.actionName.isEmpty
                    ? widget.alert!.symbol
                    : widget.alert!.actionName,
                exchange: '',
              ));
    _direction = widget.alert?.direction ?? AlertDirection.above;
    _priceController = TextEditingController(
      text: widget.alert?.targetPrice.toStringAsFixed(2) ?? '',
    );
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _bootstrapAction();
    });
  }

  Future<void> _bootstrapAction() async {
    if (_selectedAction != null || !_editing) {
      return;
    }
    final actions = await ref.read(actionsRepositoryProvider).listActions();
    if (!mounted) {
      return;
    }
    final matched =
        actions.where((action) => action.id == widget.alert!.actionId).toList();
    if (matched.isNotEmpty) {
      setState(() {
        _selectedAction = matched.first;
      });
    }
  }

  @override
  void dispose() {
    _priceController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (!_editing && _selectedAction == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Selecione uma ação.')),
      );
      return;
    }

    final session = ref.read(sessionControllerProvider);
    if (session == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Entre no app para salvar o alerta.')),
      );
      return;
    }

    final targetPrice =
        double.parse(_priceController.text.replaceAll(',', '.'));
    final repo = ref.read(activeAlertsRepositoryProvider);

    setState(() {
      _submitting = true;
    });

    try {
      if (_editing) {
        await repo.updateAlert(
          alertId: widget.alert!.id,
          targetPrice: targetPrice,
          direction: _direction,
        );
      } else {
        await repo.createAlert(
          userId: session.userId,
          actionId: _selectedAction!.id,
          targetPrice: targetPrice,
          direction: _direction,
        );
        await ref
            .read(activeWatchlistRepositoryProvider)
            .addWatchlist(_selectedAction!.id);
      }

      ref.invalidate(alertsControllerProvider);
      ref.invalidate(watchlistControllerProvider);

      if (!mounted) {
        return;
      }
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
            content: Text(
                _editing ? 'Alerta atualizado.' : 'Alerta salvo com sucesso.')),
      );
      context.pop();
    } on ApiException catch (error) {
      if (!mounted) {
        return;
      }
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(ApiClient.friendlyMessageFor(error))),
      );
    } catch (error) {
      if (!mounted) {
        return;
      }
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
    final actions = ref.watch(actionsControllerProvider);
    final colors = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            if (_editing)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Ação monitorada',
                        style: Theme.of(context).textTheme.labelLarge?.copyWith(
                              color: colors.primary,
                              fontWeight: FontWeight.w700,
                            ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        '${widget.alert!.symbol} · ${widget.alert!.actionName.isEmpty ? widget.alert!.symbol : widget.alert!.actionName}',
                        style:
                            Theme.of(context).textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.w700,
                                ),
                      ),
                    ],
                  ),
                ),
              )
            else if (widget.allowActionSelection)
              actions.when(
                data: (items) => _ActionField(
                  actions: items,
                  selected: _selectedAction,
                  onChanged: (value) => setState(() {
                    _selectedAction = value;
                  }),
                ),
                loading: () => const LinearProgressIndicator(),
                error: (error, stackTrace) =>
                    Text('Falha ao carregar ações: $error'),
              ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _priceController,
              decoration: const InputDecoration(
                labelText: 'Preço alvo',
                hintText: '40.50',
              ),
              keyboardType:
                  const TextInputType.numberWithOptions(decimal: true),
              validator: (value) {
                final parsed =
                    double.tryParse((value ?? '').replaceAll(',', '.'));
                if (parsed == null || parsed <= 0) {
                  return 'Informe um preço válido.';
                }
                return null;
              },
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
              selected: {_direction},
              onSelectionChanged: (value) {
                setState(() {
                  _direction = value.first;
                });
              },
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: _submitting ? null : _submit,
              child: Text(_submitting ? 'Salvando...' : widget.submitLabel),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActionField extends StatelessWidget {
  const _ActionField({
    required this.actions,
    required this.selected,
    required this.onChanged,
  });

  final List<stock_action.MarketAction> actions;
  final stock_action.MarketAction? selected;
  final ValueChanged<stock_action.MarketAction> onChanged;

  @override
  Widget build(BuildContext context) {
    return DropdownButtonFormField<String>(
      initialValue: selected?.id,
      decoration: const InputDecoration(labelText: 'Ação'),
      items: [
        for (final action in actions)
          DropdownMenuItem(
            value: action.id,
            child: Text('${action.symbol} · ${action.name}'),
          ),
      ],
      validator: (value) => value == null ? 'Selecione uma ação.' : null,
      onChanged: (value) {
        if (value == null) {
          return;
        }
        for (final action in actions) {
          if (action.id == value) {
            onChanged(action);
            return;
          }
        }
      },
    );
  }
}
