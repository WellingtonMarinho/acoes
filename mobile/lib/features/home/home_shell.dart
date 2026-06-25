import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../actions/domain/action.dart' as stock_action;
import '../actions/presentation/actions_controller.dart';
import '../alerts/domain/alert.dart';
import '../alerts/presentation/alerts_controller.dart';
import '../../core/network/api_client.dart';
import '../watchlist/presentation/watchlist_controller.dart';
import 'pages/alerts_page.dart';
import 'pages/settings_page.dart';
import 'pages/watchlist_page.dart';

class HomeShell extends ConsumerStatefulWidget {
  const HomeShell({super.key});

  @override
  ConsumerState<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends ConsumerState<HomeShell> {
  int _index = 0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _bootstrapSession());
  }

  Future<void> _bootstrapSession() async {
    return;
  }

  Future<void> _refreshCurrent() async {
    switch (_index) {
      case 0:
        await ref.read(watchlistControllerProvider.notifier).refresh();
        return;
      case 1:
        await ref.read(alertsControllerProvider.notifier).refresh();
        return;
      default:
        return;
    }
  }

  Future<void> _openActionPicker() async {
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (context) => _ActionPickerSheet(onCreateAlert: _openCreateAlert),
    );
  }

  void _openCreateAlert(stock_action.MarketAction action) {
    context.push('/alerts/new', extra: action);
  }

  void _openEditAlert(Alert alert) {
    context.push('/alerts/edit', extra: alert);
  }

  @override
  Widget build(BuildContext context) {
    final titles = ['Monitoradas', 'Alertas', 'Ajustes'];
    final isWatchlist = _index == 0;

    return Scaffold(
      appBar: AppBar(
        title: Text(titles[_index]),
        actions: [
          if (_index < 2)
            IconButton(
              onPressed: _refreshCurrent,
              icon: const Icon(Icons.refresh),
              tooltip: 'Atualizar',
            ),
          if (isWatchlist)
            IconButton(
              onPressed: _openActionPicker,
              icon: const Icon(Icons.add_circle_outline),
              tooltip: 'Adicionar ação',
            ),
        ],
      ),
      body: IndexedStack(
        index: _index,
        children: [
          WatchlistPageView(
            onAddAction: _openActionPicker,
            onCreateAlert: _openCreateAlert,
          ),
          AlertsPageView(
            onEditAlert: _openEditAlert,
          ),
          const SettingsPageView(),
        ],
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (value) {
          setState(() {
            _index = value;
          });
        },
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.list_alt_outlined),
            selectedIcon: Icon(Icons.list_alt),
            label: 'Monitoradas',
          ),
          NavigationDestination(
            icon: Icon(Icons.notifications_none),
            selectedIcon: Icon(Icons.notifications),
            label: 'Alertas',
          ),
          NavigationDestination(
            icon: Icon(Icons.tune_outlined),
            selectedIcon: Icon(Icons.tune),
            label: 'Ajustes',
          ),
        ],
      ),
    );
  }
}

class _ActionPickerSheet extends ConsumerStatefulWidget {
  const _ActionPickerSheet({
    required this.onCreateAlert,
  });

  final void Function(stock_action.MarketAction action) onCreateAlert;

  @override
  ConsumerState<_ActionPickerSheet> createState() => _ActionPickerSheetState();
}

class _ActionPickerSheetState extends ConsumerState<_ActionPickerSheet> {
  late Future<List<stock_action.MarketAction>> _future;
  final _formKey = GlobalKey<FormState>();
  final _queryController = TextEditingController();
  final _symbolController = TextEditingController();
  final _nameController = TextEditingController();
  final _exchangeController = TextEditingController(text: 'B3');
  bool _isSubmitting = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _future = ref.read(actionsRepositoryProvider).listActions();
    _queryController.addListener(() {
      setState(() {});
    });
  }

  @override
  void dispose() {
    _queryController.dispose();
    _symbolController.dispose();
    _nameController.dispose();
    _exchangeController.dispose();
    super.dispose();
  }

  Future<void> _addAction(stock_action.MarketAction action) async {
    try {
      await ref.read(activeWatchlistRepositoryProvider).addWatchlist(action.id);
      if (!mounted) {
        return;
      }
      ref.read(watchlistControllerProvider.notifier).refresh();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${action.symbol} adicionada à watchlist.')),
      );
      if (Navigator.of(context).canPop()) {
        Navigator.of(context).pop();
      }
    } catch (error) {
      if (!mounted) {
        return;
      }
      _showError('Falha ao adicionar ação: $error');
    }
  }

  Future<void> _createActionAndAdd() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    if (_isSubmitting) {
      return;
    }

    final symbol = _symbolController.text.trim();
    final name = _nameController.text.trim();
    final exchange = _exchangeController.text.trim();

    try {
      setState(() {
        _isSubmitting = true;
        _errorMessage = null;
      });
      final repo = ref.read(actionsRepositoryProvider);
      final action = await repo.createAction(
        symbol: symbol,
        name: name,
        exchange: exchange.isEmpty ? 'B3' : exchange,
      );
      await ref.read(activeWatchlistRepositoryProvider).addWatchlist(action.id);
      if (!mounted) {
        return;
      }
      ref.read(watchlistControllerProvider.notifier).refresh();
      setState(() {
        _future = ref.read(actionsRepositoryProvider).listActions();
        _isSubmitting = false;
        _errorMessage = null;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${action.symbol} cadastrada e adicionada.')),
      );
      if (Navigator.of(context).canPop()) {
        Navigator.of(context).pop();
      }
    } on ApiException catch (error) {
      if (!mounted) {
        return;
      }
      setState(() {
        _isSubmitting = false;
      });
      _showError(_friendlyActionError(error));
    } catch (error) {
      if (!mounted) {
        return;
      }
      setState(() {
        _isSubmitting = false;
      });
      _showError('Falha ao criar ação: $error');
    }
  }

  String _friendlyActionError(ApiException error) {
    if (error.code == 'unauthorized') {
      return 'Sessão inválida. Feche e entre novamente para continuar.';
    }
    return ApiClient.friendlyMessageFor(error);
  }

  void _showError(String message) {
    setState(() {
      _errorMessage = message;
    });
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        duration: const Duration(seconds: 4),
        behavior: SnackBarBehavior.floating,
      ),
    );
  }

  void _createAlert(stock_action.MarketAction action) {
    Navigator.of(context).pop();
    widget.onCreateAlert(action);
  }

  @override
  Widget build(BuildContext context) {
    final query = _queryController.text.trim().toLowerCase();

    return Padding(
      padding: EdgeInsets.only(
        left: 20,
        right: 20,
        top: 8,
        bottom: MediaQuery.of(context).viewInsets.bottom + 20,
      ),
      child: FutureBuilder<List<stock_action.MarketAction>>(
        future: _future,
        builder: (context, snapshot) {
          final actions = snapshot.data ?? const <stock_action.MarketAction>[];
          final filtered = query.isEmpty
              ? actions
              : actions
                  .where(
                    (action) =>
                        action.symbol.toLowerCase().contains(query) ||
                        action.name.toLowerCase().contains(query),
                  )
                  .toList();

          return SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
              if (_errorMessage != null) ...[
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Theme.of(context)
                          .colorScheme
                          .errorContainer
                          .withOpacity(0.35),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(
                        color: Theme.of(context).colorScheme.error,
                      ),
                    ),
                    child: Row(
                      children: [
                        const Icon(Icons.error_outline),
                        const SizedBox(width: 12),
                        Expanded(child: Text(_errorMessage!)),
                        TextButton(
                          onPressed: () {
                            setState(() {
                              _errorMessage = null;
                            });
                          },
                          child: const Text('Fechar'),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 12),
                ],
                Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                    Text(
                      'Cadastrar nova ação',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.w800,
                          ),
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          flex: 2,
                          child: TextFormField(
                            controller: _symbolController,
                            textCapitalization: TextCapitalization.characters,
                            decoration: const InputDecoration(
                              labelText: 'Símbolo',
                              hintText: 'PETR4',
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Informe o símbolo';
                              }
                              return null;
                            },
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          flex: 3,
                          child: TextFormField(
                            controller: _exchangeController,
                            decoration: const InputDecoration(
                              labelText: 'Bolsa',
                              hintText: 'B3',
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Informe a bolsa';
                              }
                              return null;
                            },
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _nameController,
                      decoration: const InputDecoration(
                        labelText: 'Nome da ação',
                        hintText: 'Petrobras PN',
                      ),
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) {
                          return 'Informe o nome da ação';
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 12),
                    Align(
                      alignment: Alignment.centerRight,
                      child: FilledButton.icon(
                        onPressed: _isSubmitting ? null : _createActionAndAdd,
                        icon: _isSubmitting
                            ? const SizedBox(
                                height: 16,
                                width: 16,
                                child: CircularProgressIndicator(strokeWidth: 2),
                              )
                            : const Icon(Icons.add),
                        label: const Text('Cadastrar e adicionar'),
                      ),
                    ),
                  ],
                ),
                ),
                const SizedBox(height: 20),
                const Divider(),
                const SizedBox(height: 12),
                TextField(
                  controller: _queryController,
                  decoration: const InputDecoration(
                    labelText: 'Buscar ação',
                    hintText: 'PETR4 ou Petrobras PN',
                    prefixIcon: Icon(Icons.search),
                  ),
                ),
                const SizedBox(height: 16),
                Text(
                  'Catálogo existente',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                ),
                const SizedBox(height: 12),
                if (snapshot.connectionState == ConnectionState.waiting)
                  const Padding(
                    padding: EdgeInsets.symmetric(vertical: 20),
                    child: Center(child: CircularProgressIndicator()),
                  )
                else if (filtered.isEmpty)
                  const Padding(
                    padding: EdgeInsets.symmetric(vertical: 24),
                    child: Text('Nenhuma ação encontrada.'),
                  )
                else
                  SizedBox(
                    height: 420,
                    child: ListView.separated(
                      itemCount: filtered.length,
                      separatorBuilder: (_, __) => const SizedBox(height: 12),
                      itemBuilder: (context, index) {
                        final action = filtered[index];
                        return Card(
                          child: ListTile(
                            title: Text('${action.symbol} · ${action.name}'),
                            subtitle: Text(action.exchange),
                            trailing: Wrap(
                              spacing: 8,
                              children: [
                                TextButton(
                                  onPressed: () => _addAction(action),
                                  child: const Text('Adicionar'),
                                ),
                                FilledButton.tonal(
                                  onPressed: () => _createAlert(action),
                                  child: const Text('Criar alerta'),
                                ),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
                  ),
              ],
            ),
          );
        },
      ),
    );
  }
}
