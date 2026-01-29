package main

import (
	"os"
	"testing"
)

func TestInitMasterKey(t *testing.T) {
	testKeyPath := "test_master.key"
	defer os.Remove(testKeyPath)

	t.Run("Генерация нового ключа", func(t *testing.T) {
		os.Remove(testKeyPath)

		if _, err := os.Stat(testKeyPath); os.IsNotExist(err) {
			t.Error("Файл ключа не был создан")
		}
	})

	t.Run("Загрузка существующего ключа", func(t *testing.T) {
		testKey := []byte("test-32-byte-key-for-testing!123")
		if err := os.WriteFile(testKeyPath, testKey, 0600); err != nil {
			t.Fatalf("Не удалось создать тестовый ключ: %v", err)
		}

		// Читаем ключ
		data, err := os.ReadFile(testKeyPath)
		if err != nil {
			t.Fatalf("Ошибка чтения ключа: %v", err)
		}

		if len(data) != 32 {
			t.Errorf("Загруженный ключ имеет неверную длину: %d", len(data))
		}

		if string(data) != string(testKey) {
			t.Error("Загруженный ключ не совпадает с оригиналом")
		}
	})
}

// обработка переменных окружения
func TestGetEnv(t *testing.T) {
	t.Run("Возврат значения переменной окружения", func(t *testing.T) {
		testKey := "TEST_VAR_KEY"
		testValue := "test_value"

		os.Setenv(testKey, testValue)
		defer os.Unsetenv(testKey)

		result := getEnv(testKey, "default")

		if result != testValue {
			t.Errorf("Ожидалось %s, получено %s", testValue, result)
		}
	})

	t.Run("Возврат значения по умолчанию", func(t *testing.T) {
		testKey := "NON_EXISTENT_VAR"
		defaultValue := "default_value"

		os.Unsetenv(testKey)

		result := getEnv(testKey, defaultValue)

		if result != defaultValue {
			t.Errorf("Ожидалось %s, получено %s", defaultValue, result)
		}
	})

	t.Run("Пустая переменная возвращает default", func(t *testing.T) {
		testKey := "EMPTY_VAR"
		defaultValue := "default"

		os.Setenv(testKey, "")
		defer os.Unsetenv(testKey)

		result := getEnv(testKey, defaultValue)

		if result != defaultValue {
			t.Errorf("Пустая переменная должна вернуть default: %s, получено %s", defaultValue, result)
		}
	})
}

// проверяет формирование строки подключения
func TestGetDBConnectionString(t *testing.T) {
	// Сохраняем текущие переменные окружения
	oldVars := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
	}

	// после теста
	defer func() {
		for key, value := range oldVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("Значения по умолчанию", func(t *testing.T) {
		// Очищаем все переменные
		for key := range oldVars {
			os.Unsetenv(key)
		}

		connStr := getDBConnectionString()

		expected := "host=localhost port=5432 user=postgres password=postgres dbname=password sslmode=disable"
		if connStr != expected {
			t.Errorf("Неверная строка подключения.\nОжидалось: %s\nПолучено: %s", expected, connStr)
		}
	})

	t.Run("Кастомные значения", func(t *testing.T) {
		os.Setenv("DB_HOST", "192.168.1.100")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_USER", "myuser")
		os.Setenv("DB_PASSWORD", "mypass")
		os.Setenv("DB_NAME", "mydb")
		os.Setenv("DB_SSLMODE", "require")

		connStr := getDBConnectionString()

		expected := "host=192.168.1.100 port=5433 user=myuser password=mypass dbname=mydb sslmode=require"
		if connStr != expected {
			t.Errorf("Неверная строка подключения.\nОжидалось: %s\nПолучено: %s", expected, connStr)
		}
	})
}
