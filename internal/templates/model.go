package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateModel creates a Dart model class template with Equatable or Freezed implementation
func GenerateModel(name string, useFreezed bool) string {
	if useFreezed {
		return fmt.Sprintf(`
import 'package:freezed_annotation/freezed_annotation.dart';

part '%[1]s.freezed.dart';
part '%[1]s.g.dart';

@freezed
abstract class %[2]s with _$%[2]s {
    const factory %[2]s({
        required String id,
    }) = _%[2]s;

    factory %[2]s.fromJson(Map<String, dynamic> json) => 
        _$%[2]sFromJson(json);
}`, strings.ToLower(name), name)
	}

	return fmt.Sprintf(`
import 'package:equatable/equatable.dart';

class %s extends Equatable {
    final String id;

    const %s({
        required this.id,
    });

    @override
    List<Object?> get props => [id];
}`, name, name)
}

func GenerateModelTest(modelName string, useFreezed bool, projectDir string) string {
	modelFileName := strings.ToLower(modelName)
	packageName, err := GetPackageName(projectDir)
	if err != nil {
		packageName = "flutter_app"
	}

	imports := []string{
		"package:flutter_test/flutter_test.dart",
		fmt.Sprintf("package:%s/models/%s.dart", packageName, modelFileName),
	}

	// Common test cases for both Equatable and Freezed
	testCases := []struct {
		name     string
		testCase string
	}{
		{"should create instance correctly", createInstanceTest(modelName)},
		{"should support value comparison", valueEqualityTest(modelName)},
		{"should have correct string representation", toStringTest(modelName)},
	}

	// Add props test only for Equatable models
	if !useFreezed {
		testCases = append(testCases, struct {
			name     string
			testCase string
		}{"should have correct props", propsTest(modelName)})
	}

	var tests []string
	for _, tc := range testCases {
		tests = append(tests, fmt.Sprintf(`test('%s', () {
            %s
        });`, tc.name, tc.testCase))
	}

	// Add Freezed-specific tests
	if useFreezed {
		freezedTests := []struct {
			name     string
			testCase string
		}{
			{"should convert to and from JSON", jsonTest(modelName)},
			{"should support copyWith", copyWithTest(modelName)},
		}

		for _, tc := range freezedTests {
			tests = append(tests, fmt.Sprintf(`test('%s', () {
                %s
            });`, tc.name, tc.testCase))
		}
	}

	return fmt.Sprintf(`import '%s';

void main() {
    group('%s', () {
        %s
    });
}`, strings.Join(imports, "';\nimport '"), modelName, strings.Join(tests, "\n\n        "))
}

func createInstanceTest(modelName string) string {
	return fmt.Sprintf(`final model = %s(
        id: '1',
    );
    
    expect(model.id, equals('1'));`, modelName)
}

func valueEqualityTest(modelName string) string {
	return fmt.Sprintf(`final model1 = %s(
        id: '1',
    );
    final model2 = %s(
        id: '1',
    );
    
    expect(model1, equals(model2));
    expect(model1.hashCode, equals(model2.hashCode));`, modelName, modelName)
}

func propsTest(modelName string) string {
	return fmt.Sprintf(`final model = %s(
        id: '1',
    );
    
    expect(model.props, equals([model.id]));`, modelName)
}

func toStringTest(modelName string) string {
	return fmt.Sprintf(`final model = %s(
        id: '1',
    );
    
    expect(model.toString(), contains('%s'));
    expect(model.toString(), contains('1'));`, modelName, modelName)
}

func jsonTest(modelName string) string {
	return fmt.Sprintf(`final model = %s(
        id: '1',
    );
    final json = model.toJson();
    final fromJson = %s.fromJson(json);
    
    expect(fromJson, equals(model));
    expect(json['id'], equals('1'));`, modelName, modelName)
}

func copyWithTest(modelName string) string {
	return fmt.Sprintf(`final model = %s(
        id: '1',
    );
    final copy = model.copyWith(id: '2');
    
    expect(copy.id, equals('2'));
    expect(model.id, equals('1'));`, modelName)
}

func GetPackageName(projectDir string) (string, error) {
	pubspecPath := filepath.Join(projectDir, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return "", fmt.Errorf("failed to read pubspec.yaml: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "name:") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "name:")), nil
		}
	}

	return "", fmt.Errorf("package name not found in pubspec.yaml")
}
