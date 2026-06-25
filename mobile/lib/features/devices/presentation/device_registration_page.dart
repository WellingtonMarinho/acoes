import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'devices_controller.dart';

class DeviceRegistrationPage extends ConsumerStatefulWidget {
  const DeviceRegistrationPage({super.key});

  @override
  ConsumerState<DeviceRegistrationPage> createState() =>
      _DeviceRegistrationPageState();
}

class _DeviceRegistrationPageState
    extends ConsumerState<DeviceRegistrationPage> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _tokenController;
  late final TextEditingController _platformController;
  bool _submitting = false;

  @override
  void initState() {
    super.initState();
    _tokenController = TextEditingController();
    _platformController = TextEditingController(text: 'android');
  }

  @override
  void dispose() {
    _tokenController.dispose();
    _platformController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    setState(() {
      _submitting = true;
    });

    try {
      await ref.read(activeDevicesRepositoryProvider).registerDevice(
            deviceToken: _tokenController.text.trim(),
            platform: _platformController.text.trim(),
          );

      if (!mounted) {
        return;
      }
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Device registrado.')),
      );
      Navigator.of(context).pop();
    } catch (error) {
      if (!mounted) {
        return;
      }
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Falha ao registrar device: $error')),
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
    return Scaffold(
      appBar: AppBar(title: const Text('Registrar device')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Text(
              'Registre o device para receber notificações.',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _tokenController,
              decoration: const InputDecoration(
                labelText: 'Device token',
                hintText: 'fcm-token',
              ),
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return 'Informe o token.';
                }
                return null;
              },
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _platformController,
              decoration: const InputDecoration(
                labelText: 'Plataforma',
                hintText: 'android',
              ),
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return 'Informe a plataforma.';
                }
                return null;
              },
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: _submitting ? null : _submit,
              child: Text(_submitting ? 'Salvando...' : 'Salvar device'),
            ),
          ],
        ),
      ),
    );
  }
}
