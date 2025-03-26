# Flart - Flutter Artifact Generator

[![Release](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/frankmungnodev/d5877c86cd581fe08db77ebf0623c409/raw/flart_version.json)](https://github.com/frankmungnodev/flart/releases)

A command-line tool to streamline Flutter development by generating models, screens, and managing build runner tasks.

## Features

- 🎯 Generate Models
  - Support for Freezed annotations
  - Automatic test file generation
  - Equatable integration

- 📱 Generate Screens
  - BLoC/Cubit support
  - Freezed state management

- 🔄 Build Runner Management
  - One-time build
  - Watch mode
  - Automatic conflict resolution

## Configuration

Create a `flart_config.json` file in your project root. This is optional.

```json
{
    "projectDir": "~/path/to/your/flutter/project",
    "models": {
        "useFreezed": false
    },
    "screens": {
        "useCubit": false,
        "useFreezed": true
    }
}
```

### Configuration Options

- `projectDir`: Path to your Flutter project (default to current directory)
- `models.useFreezed`: Enable Freezed for model generation (default to false)
- `screens.useCubit`: Use Cubit instead of BLoC (default to false)
- `screens.useFreezed`: Enable Freezed for state classes (default to false)

## Usage

### CLI Mode

Generate a model:
```bash
flart make:model User
```

Generate a screen:
```bash
flart make:screen Login
```

Run build_runner:
```bash
flart build:runner    # One-time build
flart watch:runner    # Watch mode
```

### Interactive Mode

Run without arguments for interactive mode:
```bash
flart
```

## Generated Structure

```
lib/
├── models/
│   └── user.dart
├── screens/
│   └── login/
│       ├── cubit/
│       │   ├── login_cubit.dart
│       │   └── login_cubit_state.dart
│       ├── login_screen.dart
└── test/
    ├── models/
    │   └── user_test.dart
```

## Releases

### Latest Release
v0.1.0 - Initial Release

[View all releases](https://github.com/frankmungnodev/flart/releases)

### Installation from Binary

macOS (Apple Silicon):
```bash
curl -L https://github.com/frankmungnodev/flart/releases/download/v0.1.0/flart_0.1.0_darwin_arm64.tar.gz | tar xz
sudo mv flart_0.1.0_darwin_arm64 /usr/local/bin/flart
```

macOS (Intel):
```bash
curl -L https://github.com/frankmungnodev/flart/releases/download/v0.1.0/flart_0.1.0_darwin_amd64.tar.gz | tar xz
sudo mv flart_0.1.0_darwin_amd64 /usr/local/bin/flart
```

Linux:
```bash
curl -L https://github.com/frankmungnodev/flart/releases/download/v0.1.0/flart_0.1.0_linux_amd64.tar.gz | tar xz
sudo mv flart_0.1.0_linux_amd64 /usr/local/bin/flart
```

Windows:
Download and extract `flart_0.1.0_windows_amd64.zip` and add to your PATH.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.