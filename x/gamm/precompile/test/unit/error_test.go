package gamm_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

// TestErrorCodes tests that all error codes are defined
func TestErrorCodes(t *testing.T) {
	require.Equal(t, gamm.ErrorCode(1), gamm.ErrorCodeInvalidInput)
	require.Equal(t, gamm.ErrorCode(2), gamm.ErrorCodePoolNotFound)
	require.Equal(t, gamm.ErrorCode(3), gamm.ErrorCodeSwapFailed)
	require.Equal(t, gamm.ErrorCode(4), gamm.ErrorCodeQueryFailed)
	require.Equal(t, gamm.ErrorCode(5), gamm.ErrorCodeInternalError)
	require.Equal(t, gamm.ErrorCode(6), gamm.ErrorCodeUnauthorized)
	require.Equal(t, gamm.ErrorCode(7), gamm.ErrorCodeJoinPoolFailed)
	require.Equal(t, gamm.ErrorCode(8), gamm.ErrorCodeExitPoolFailed)
	require.Equal(t, gamm.ErrorCode(9), gamm.ErrorCodeIBCTransferFailed)
}

// TestNewPrecompileError tests error creation
func TestNewPrecompileError(t *testing.T) {
	err := gamm.NewPrecompileError(gamm.ErrorCodeInvalidInput, "test message", "test details")
	require.NotNil(t, err)
	require.Equal(t, gamm.ErrorCodeInvalidInput, err.Code)
	require.Equal(t, "test message", err.Message)
	require.Equal(t, "test details", err.Details)
	require.Contains(t, err.Error(), "precompile error")
	require.Contains(t, err.Error(), "test message")
	require.Contains(t, err.Error(), "test details")
}

// TestMapCosmosErrorToPrecompileError tests error mapping
func TestMapCosmosErrorToPrecompileError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode gamm.ErrorCode
		shouldMap    bool
	}{
		{
			name:         "ErrPoolNotFound",
			err:          gammtypes.ErrPoolNotFound,
			expectedCode: gamm.ErrorCodePoolNotFound,
			shouldMap:    true,
		},
		// Note: poolmanagertypes.ErrInvalidRoute doesn't exist, using gammtypes errors instead
		{
			name:         "unmapped_error",
			err:          fmt.Errorf("some random error"),
			expectedCode: 0,
			shouldMap:    false,
		},
		{
			name:         "nil_error",
			err:          nil,
			expectedCode: 0,
			shouldMap:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, found := gamm.MapCosmosErrorToPrecompileError(tt.err)
			require.Equal(t, tt.shouldMap, found, "mapping should match expected")
			if tt.shouldMap {
				require.Equal(t, tt.expectedCode, code, "error code should match")
			}
		})
	}
}

// TestWrapError tests error wrapping
func TestWrapError(t *testing.T) {
	// Test that mapped errors use the mapped code
	err := gammtypes.ErrPoolNotFound
	wrapped := gamm.WrapError(err, gamm.ErrorCodeQueryFailed, "test message")
	require.Equal(t, gamm.ErrorCodePoolNotFound, wrapped.Code, "should use mapped code, not default")

	// Test that unmapped errors use the default code
	unmappedErr := fmt.Errorf("some random error")
	wrapped2 := gamm.WrapError(unmappedErr, gamm.ErrorCodeInternalError, "test message")
	require.Equal(t, gamm.ErrorCodeInternalError, wrapped2.Code, "should use default code for unmapped errors")

	// Test that nil errors use the default code
	wrapped3 := gamm.WrapError(nil, gamm.ErrorCodeQueryFailed, "test message")
	require.Equal(t, gamm.ErrorCodeQueryFailed, wrapped3.Code, "should use default code for nil errors")
}

// TestErrorConstructors tests all error constructors
func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		createFn func() *gamm.PrecompileError
		code     gamm.ErrorCode
	}{
		{"ErrInvalidInput", func() *gamm.PrecompileError { return gamm.ErrInvalidInput("test") }, gamm.ErrorCodeInvalidInput},
		{"ErrPoolNotFound", func() *gamm.PrecompileError { return gamm.ErrPoolNotFound(1) }, gamm.ErrorCodePoolNotFound},
		{"ErrSwapFailed", func() *gamm.PrecompileError { return gamm.ErrSwapFailed("test") }, gamm.ErrorCodeSwapFailed},
		{"ErrJoinPoolFailed", func() *gamm.PrecompileError { return gamm.ErrJoinPoolFailed("test") }, gamm.ErrorCodeJoinPoolFailed},
		{"ErrExitPoolFailed", func() *gamm.PrecompileError { return gamm.ErrExitPoolFailed("test") }, gamm.ErrorCodeExitPoolFailed},
		{"ErrQueryFailed", func() *gamm.PrecompileError { return gamm.ErrQueryFailed("test") }, gamm.ErrorCodeQueryFailed},
		{"ErrInternalError", func() *gamm.PrecompileError { return gamm.ErrInternalError("test") }, gamm.ErrorCodeInternalError},
		{"ErrUnauthorized", func() *gamm.PrecompileError { return gamm.ErrUnauthorized("test") }, gamm.ErrorCodeUnauthorized},
		{"ErrIBCTransferFailed", func() *gamm.PrecompileError { return gamm.ErrIBCTransferFailed("test") }, gamm.ErrorCodeIBCTransferFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createFn()
			require.NotNil(t, err)
			require.Equal(t, tt.code, err.Code)
		})
	}
}

// TestErrorSanitization tests that error sanitization removes sensitive information
func TestErrorSanitization(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string // What should NOT appear in output
		mustContain string // What should appear in output
	}{
		{
			name:        "removes_unix_paths",
			input:       "error at /home/user/file.go:123",
			expected:    "/home/",
			mustContain: "[path]/",
		},
		{
			name:        "removes_windows_paths",
			input:       "error at C:\\Users\\file.go:123",
			expected:    "C:\\",
			mustContain: "[path]\\",
		},
		{
			name:        "removes_go_file_references",
			input:       "error in handlers.go:456",
			expected:    ".go:",
			mustContain: "[file]:",
		},
		{
			name:        "truncates_long_messages",
			input:       string(make([]byte, 1000)), // 1000 bytes
			expected:    "",                          // Should be truncated
			mustContain: "[truncated]",
		},
		{
			name:        "preserves_safe_messages",
			input:       "pool not found: 123",
			expected:    "", // Should not be modified
			mustContain: "pool not found: 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an error with the test input
			err := gamm.NewPrecompileError(gamm.ErrorCodeInternalError, "test", tt.input)

			// Check that sanitized details don't contain sensitive info
			if tt.expected != "" {
				require.NotContains(t, err.Details, tt.expected, "sensitive information should be removed")
			}

			// Check that sanitized details contain safe replacement
			if tt.mustContain != "" {
				require.Contains(t, err.Details, tt.mustContain, "should contain sanitized replacement")
			}
		})
	}
}

