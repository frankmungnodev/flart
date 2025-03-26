
```
# Flart - Flutter Artifact Generator

A command-line tool to streamline Flutter development by generating models, screens, and managing build runner tasks.

## Features

- 🎯 Generate Models
  - Support for Freezed annotations
  - Automatic test file generation
  - Equatable integration

- 📱 Generate Screens
  - BLoC/Cubit support
  - Freezed state management
  - Test file scaffolding

- 🔄 Build Runner Management
  - One-time build
  - Watch mode
  - Automatic conflict resolution

## Installation

```bash
go install github.com/yourusername/flart@latest
```

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
│       ├── bloc/
│       │   ├── login_cubit.dart
│       │   └── login_cubit_state.dart
│       ├── login_screen.dart
└── test/
    ├── models/
    │   └── user_test.dart
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.