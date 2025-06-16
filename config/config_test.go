package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Setenv("TEST_KEY", "test_value")

	value, err := getString("TEST_KEY")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", value)

	_, err = getString("NON_EXISTING_KEY")
	assert.Error(t, err)
}

func TestGetBool(t *testing.T) {
	t.Setenv("TEST_BOOL", "true")

	value, err := getBool("TEST_BOOL")
	assert.NoError(t, err)
	assert.Equal(t, true, value)

	_, err = getBool("NON_EXISTING_BOOL")
	assert.Error(t, err)

	t.Setenv("INVALID_BOOL", "not_a_bool")
	_, err = getBool("INVALID_BOOL")
	assert.Error(t, err)
}

func TestLoad(t *testing.T) {
	setTestEnv(t)

	err := Load("test")
	assert.NoError(t, err)

	t.Run("TestLoadWithInvalidGinMode", func(t *testing.T) {
		setTestEnv(t)
		t.Setenv("GIN_MODE", "invalid_mode")

		err := Load("test")
		assert.NoError(t, err)

		assert.Equal(t, "debug", GinMode)
	})

	t.Run("TestLoadWithMissingNotificationID", func(t *testing.T) {
		setTestEnv(t)
		t.Setenv("NOTIFICATION_ID", "")

		err := Load("test")
		assert.Error(t, err)
	})

	t.Run("TestLoadWithMissingFVersion", func(t *testing.T) {

		setTestEnv(t)
		t.Setenv("F_VERSION", "")

		err := Load("test")
		assert.Error(t, err)
	})
	t.Run("TestLoadWithJsonLogs", func(t *testing.T) {
		setTestEnv(t)
		t.Setenv("JSON_LOGS", "true")

		err := Load("test")
		assert.NoError(t, err)

		assert.Equal(t, true, JsonLogs)
	})

	t.Run("TestLoadWithJsonLogsFalse", func(t *testing.T) {
		setTestEnv(t)
		t.Setenv("JSON_LOGS", "false")

		err := Load("test")
		assert.NoError(t, err)

		assert.Equal(t, false, JsonLogs)
	})

	t.Run("TestLoadWithInvalidJsonLogs", func(t *testing.T) {
		setTestEnv(t)
		t.Setenv("JSON_LOGS", "not_a_bool")

		err := Load("test")
		assert.Error(t, err)
	})
	t.Run("TestLoadWithDefaultValues", func(t *testing.T) {
		setTestEnv(t)

		err := Load("test")
		assert.NoError(t, err)

		assert.Equal(t, "debug", GinMode)
		assert.Equal(t, "test_notification_id", NotificationID)
		assert.Equal(t, "1.0.0", FVersion)
		assert.Equal(t, true, JsonLogs)
	})
	t.Run("TestLoadWithEmptyEnvironment", func(t *testing.T) {
		err := Load("")
		assert.NoError(t, err)
	})
	t.Run("TestLoadWithInvalidEnvironment", func(t *testing.T) {
		err := Load("invalid_env")
		assert.NoError(t, err)
	})
	t.Run("TestLoadWithValidEnvironment", func(t *testing.T) {
		setTestEnv(t)

		err := Load("test")
		assert.NoError(t, err)

		assert.Equal(t, "debug", GinMode)
		assert.Equal(t, "test_notification_id", NotificationID)
		assert.Equal(t, "1.0.0", FVersion)
		assert.Equal(t, true, JsonLogs)
	})

}

func setTestEnv(t *testing.T) {
	t.Setenv("GIN_MODE", "debug")
	t.Setenv("NOTIFICATION_ID", "test_notification_id")
	t.Setenv("F_VERSION", "1.0.0")
	t.Setenv("JSON_LOGS", "true")
}
