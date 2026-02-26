package db

import (
	"testing"
)

func TestRun_NoSubcommand_ReturnsHelp(t *testing.T) {
	getArgs = func() []string {
		return []string{"sbx", "db"}
	}

	result, err := Run()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != subcommandHelp {
		t.Errorf("Expected help text, got %v", result)
	}
}

func TestRun_HelpSubcommand_ReturnsHelp(t *testing.T) {
	getArgs = func() []string {
		return []string{"sbx", "db", "help"}
	}

	result, err := Run()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != subcommandHelp {
		t.Errorf("Expected help text, got %v", result)
	}
}

func TestRun_UnknownSubcommand_ReturnsError(t *testing.T) {
	getArgs = func() []string {
		return []string{"sbx", "db", "unknown"}
	}

	_, err := Run()

	if err == nil {
		t.Error("Expected error for unknown subcommand, got nil")
	}
}
