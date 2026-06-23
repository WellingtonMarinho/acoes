# Mobile

App Flutter para:

- criar alertas
- listar watchlist
- configurar tokens de push

## Status

O projeto Flutter mínimo já está estruturado neste diretório.

O app já tem:

- home
- sessão provisória com token
- criação de alerta
- registro de device
- persistência local da sessão

Parte da navegação ainda usa dados demo enquanto a integração total com o backend é fechada.

## Como subir no Android Studio

Este é o caminho recomendado para rodar o app Android localmente.

### 1. Abrir o projeto

1. Abra o Android Studio.
2. Clique em `Open`.
3. Selecione a pasta `mobile/` deste repositório.
4. Aguarde o indexamento inicial do projeto.

### 2. Instalar o SDK Android

Se o Android Studio pedir, instale os componentes do SDK.
Caso precise revisar manualmente:

1. Abra `Settings` ou `Preferences`.
2. Entre em `Android SDK`.
3. Instale uma versão recente do Android SDK.
4. Em `SDK Tools`, garanta estes itens:
   - `Android SDK Platform-Tools`
   - `Android SDK Build-Tools`
   - `Android Emulator`

### 3. Criar um emulador

1. Abra o `Device Manager`.
2. Clique em `Create device`.
3. Escolha um aparelho, por exemplo um Pixel.
4. Baixe uma imagem de sistema Android.
5. Finalize a criação e inicie o emulador.

### 4. Validar o ambiente

No terminal, dentro da pasta `mobile/`, rode:

```bash
flutter doctor -v
```

Se o Flutter pedir licenças Android, rode:

```bash
flutter doctor --android-licenses
```

### 5. Executar o app

Com o emulador aberto, você pode:

1. Clicar em `Run` no Android Studio.
2. Ou usar o terminal na raiz do projeto:

```bash
make run-mobile
```

Se o emulador estiver ativo, o Flutter deve detectar o device Android e subir o app.

### 6. Rodar testes

Para executar a suíte do app:

```bash
make test-mobile
```

## Observação

Se o Android Studio abrir sem um device Android disponível, o `flutter run` não vai ter onde iniciar o app.
Nesse caso, volte ao `Device Manager`, inicie um emulador e tente novamente.

## Estrutura

- `lib/app/`: bootstrap da aplicação, tema e navegação
- `lib/core/`: client HTTP, configuração e utilitários compartilhados
- `lib/features/`: funcionalidades por domínio
- `test/`: testes unitários e de widget

## Próximos passos

1. Consolidar a integração do client HTTP com o backend protegido.
2. Melhorar o pós-submit das telas principais.
3. Cobrir os fluxos com testes de widget e de repositório.
