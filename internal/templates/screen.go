package templates

import (
	"flart/internal/utils"
	"fmt"
	"strings"
)

// GenerateScreen creates a Flutter screen template with BLoC or Cubit integration
func GenerateScreen(screenName string, useCubit bool) string {
	if useCubit {
		return fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'cubit/%[1]s_cubit.dart';
import 'cubit/%[1]s_state.dart';

class %[2]sScreen extends StatelessWidget {
  const %[2]sScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => %[2]sCubit(),
      child: const %[2]sView(),
    );
  }
}

class %[2]sView extends StatefulWidget {
  const %[2]sView({super.key});

  @override
  State<%[2]sView> createState() => _%[2]sViewState();
}

class _%[2]sViewState extends State<%[2]sView> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('%[2]s')),
      body: BlocBuilder<%[2]sCubit, %[2]sState>(
        builder: (context, state) {
          return const Center(child: Text('%[2]s Screen'));
        },
      ),
    );
  }
}`, utils.ToSnakeCase(screenName), screenName)
	}
	return fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'bloc/%[1]s_bloc.dart';
import 'bloc/%[1]s_event.dart';
import 'bloc/%[1]s_state.dart';

class %[2]sScreen extends StatelessWidget {
  const %[2]sScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => %[2]sBloc()..add(const %[2]sInitialEvent()),
      child: const %[2]sView(),
    );
  }
}

class %[2]sView extends StatefulWidget {
  const %[2]sView({super.key});

  @override
  State<%[2]sView> createState() => _%[2]sViewState();
}

class _%[2]sViewState extends State<%[2]sView> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('%[2]s')),
      body: BlocBuilder<%[2]sBloc, %[2]sState>(
        builder: (context, state) {
          return const Center(child: Text('%[2]s Screen'));
        },
      ),
    );
  }
}`, utils.ToSnakeCase(screenName), screenName)
}

// GenerateBloc creates a BLoC template with initial setup
func GenerateBloc(screenName string) string {
	return fmt.Sprintf(`import 'package:flutter_bloc/flutter_bloc.dart';

import '%[1]s_event.dart';
import '%[1]s_state.dart';

class %[2]sBloc extends Bloc<%[2]sEvent, %[2]sState> {
  %[2]sBloc() : super(const %[2]sState()) {
    on<%[2]sInitialEvent>(_onInitial);
  }

  Future<void> _onInitial(
    %[2]sInitialEvent event,
    Emitter<%[2]sState> emit,
  ) async {
    // TODO: Add your logic here
  }
}`, utils.ToSnakeCase(screenName), screenName)
}

// GenerateCubit creates a Cubit template with initial setup
func GenerateCubit(screenName string) string {
	return fmt.Sprintf(`import 'package:flutter_bloc/flutter_bloc.dart';

import '%[1]s_state.dart';

class %[2]sCubit extends Cubit<%[2]sState> {
  %[2]sCubit() : super(const %[2]sState());

  Future<void> init() async {
    // TODO: Add your logic here
  }
}`, utils.ToSnakeCase(screenName), screenName)
}

// GenerateEvent creates event classes for the BLoC
func GenerateEvent(screenName string) string {
	return fmt.Sprintf(`
import 'package:equatable/equatable.dart';

abstract class %[2]sEvent extends Equatable {
  const %[2]sEvent();

  @override
  List<Object> get props => [];
}

class %[2]sInitialEvent extends %[2]sEvent {
  const %[2]sInitialEvent();
}

class %[2]sRefreshEvent extends %[2]sEvent {
  const %[2]sRefreshEvent();
}`, strings.ToLower(screenName), screenName)
}

// GenerateState creates state classes for the BLoC or Cubit
func GenerateState(screenName string, useCubit bool, useFreezed bool) string {
	parentFile := "bloc"
	if useCubit {
		parentFile = "cubit"
	}

	if useFreezed {
		return fmt.Sprintf(`
import 'package:freezed_annotation/freezed_annotation.dart';

part '%[1]s_state.freezed.dart';

@freezed
abstract class %[2]sState with _$%[2]sState {
  const factory %[2]sState({
    @Default(false) bool isLoading,
  }) = _%[2]sState;
}`, strings.ToLower(screenName), screenName, parentFile)
	}

	return fmt.Sprintf(`
import 'package:equatable/equatable.dart';

class %[2]sState extends Equatable {
  final bool isLoading;

  const %[2]sState({
    this.isLoading = false,
  });

  @override
  List<Object?> get props => [isLoading];

  %[2]sState copyWith({
    bool? isLoading,
  }) {
    return %[2]sState(
      isLoading: isLoading ?? this.isLoading,
    );
  }
}`, utils.ToSnakeCase(screenName), screenName, parentFile)
}
