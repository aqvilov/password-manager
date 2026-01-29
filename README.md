# 🔐 Password manager 'ppocket'

- [Russian Version README](#Ru-Guide)
- [English Version README](#En-Guide)

---


### Ru-Guide
Безопасный консольный менеджер паролей с шифрованием AES-256 и безопасным хранением данных.

## Быстрый старт

### Вариант 1: С Docker

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/your-username/password-manager.git
   cd password-manager
   ```

2. **Запустите DOcker Desktop:**
   ```bash
   docker-compose up -d
   ```

3. **Установите зависимости:**
   ```bash
   go mod download
   ```

4. **Запустите приложение:**
   ```bash
   docker attach password_manager_app
   ```

При первом запуске автоматически создастся:
- База данных `password`
- Мастер-ключ шифрования `master.key`
- Таблица для хранения паролей

### Вариант 2: С локальной базой данных

1. **Установите PostgreSQL** (если еще не установлен)
   ``` https://www.postgresql.org/download/ ```

2. **Отредактируйте `.env`** (при необходимости):
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=name
   DB_SSLMODE=disable
   ```

3. **Запустите приложение:**
   ```bash
   go build password-manager

   ./password-manager
   ```

## 📦 Зависимости

```bash
go get github.com/lib/pq
go get github.com/joho/godotenv
```

## 🔧 Конфигурация

### Мастер-ключ шифрования

Мастер-ключ автоматически генерируется при первом запуске и сохраняется в файл `master.key`.

⚠️ **ВАЖНО:**
- **Храните `master.key` в безопасном месте!**
- Без этого файла вы не сможете расшифровать ваши пароли
- Не загружайте `master.key` в git (уже в `.gitignore`)

## 📖 Использование

После запуска программы вы увидите главное меню:

```
══════════════════════════════════════════════════
МЕНЕДЖЕР ПАРОЛЕЙ
══════════════════════════════════════════════════

 ГЛАВНОЕ МЕНЮ:
1.  Показать все пароли
2.  Добавить новый пароль
3.  Изменить существующий пароль
4.  Удалить пароль
5.  Поиск паролей
6.  Очистить экран
0.  Выход
```

### Примеры использования

**Добавить новый пароль:**
```
Выберите действие: 2
Введите название сервиса: GitHub
Введите логин: myusername
Введите пароль: mySecurePassword123!
Введите описание: Мой основной аккаунт
```

**Поиск паролей:**
```
Выберите действие: 5
Введите ключевые слова для поиска: git
```

Программа найдет все записи, содержащие "git" в названии сервиса, логине или описании.


## 🔒 Безопасность

1. **Шифрование:** Все пароли шифруются с использованием AES-256 в режиме GCM
2. **Мастер-ключ:** 32-байтный ключ, генерируемый с помощью `crypto/rand`
3. **Права доступа:** Файл `master.key` создается с правами `0600` (права владельца)
4. **База данных:** Пароли хранятся в зашифрованном виде
5. **Переменные окружения:** Конфиденциальные данные не хардкодятся в коде
6. 

## 🛠️ Разработка

### Сборка проекта

```bash
go build -o password-manager
```


### Остановка БД (Docker)

```bash
docker-compose down
```

### Полная очистка (включая данные)

```bash
docker-compose down -v
```

## ⚠️ Troubleshooting

### PostgreSQL не запущен

```
PostgreSQL не запущен: connection refused

Попробуйте запустить:
  docker-compose up -d
или установите PostgreSQL
```

**Решение:** Запустите `docker-compose up -d` или установите PostgreSQL локально.

### Ошибка подключения к БД

Проверьте настройки в файле `.env`:
```bash
cat .env
```

### DANGER ZONE

Если вы потеряли файл `master.key`, расшифровать существующие пароли **невозможно**.

**Решение:** 
- Или удалите старую БД и создайте новую (потеряете все пароли)


### En-Guide

Secure console password manager with AES-256 encryption and secure data storage.

## Quick Start

### Option 1: With Docker

1. **Clone the repository:**
```bash
   git clone https://github.com/your-username/password-manager.git
   cd password-manager
   ```

2. **Run DOcker Desktop:**
```bash
   docker-compose up -d
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Run the application:**
   ```bash
   docker attach password_manager_app
   ```

The following will be created automatically on first launch:
- `password` database
- `master.key` encryption master key
- Table for storing passwords

### Option 2: With a local database

1. **Install PostgreSQL** (if not already installed)
```https://www.postgresql.org/download/```

2. **Edit `.env`** (if necessary):
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=name
   DB_SSLMODE=disable
   ```

3. **Run the application:**
```bash
   go build password-manager

   ./password-manager
   ```

## 📦 Dependencies

```bash
go get github.com/lib/pq
go get github.com/joho/godotenv
```

## 🔧 Configuration

### Master encryption key

The master key is automatically generated on first launch and stored in the `master.key` file.

⚠️ **IMPORTANT:**
- **Keep `master.key` in a safe place!**
- Without this file, you will not be able to decrypt your passwords.
- Do not upload `master.key` to git (already in `.gitignore`).

## 📖 Usage

After launching the program, you will see the main menu:

```
═════════════════════════════════ ═════════════════
PASSWORD MANAGER
═══════════ ═════════════════════════════════ Main Menu:
1.  Show all passwords
2.  Add new password
3.  Change existing password
4.  Delete password
5.  Search passwords
6.  Clear screen
0.  Exit

### Examples of use

**Add new password:**
```
Select action: 2
Enter service name: GitHub
Enter login: myusername
Enter password: mySecurePassword123!
Enter description: My main account
```

**Search for passwords:**
```
Select action: 5
Enter keywords to search for: git
```

The program will find all entries containing “git” in the service name, login, or description.


## 🔒 Security

1. **Encryption:** All passwords are encrypted using AES-256 in GCM mode
2. **Master key:** 32-byte key generated using `crypto/rand`
3. **Access rights:** The `master.key` file is created with `0600` permissions (owner permissions)
4. **Database:** Passwords are stored in encrypted form
5. **Environment variables:** Confidential data is not hardcoded in the code
6.

## 🛠️ Development

### Building the project

```bash
go build -o password-manager
```


### Stopping the database (Docker)

```bash
docker-compose down
```

### Complete cleanup (including data)

```bash
docker-compose down -v
```

## ⚠️ Troubleshooting

### PostgreSQL is not running

```
PostgreSQL is not running: connection refused

Try running:
  docker-compose up -d
or install PostgreSQL
```

**Solution:** Run `docker-compose up -d` or install PostgreSQL locally.

### Database connection error

Check the settings in the `.env` file:
```bash
cat .env
```

### DANGER ZONE

If you have lost the `master.key` file, it is **impossible** to decrypt existing passwords.

**Solution:**
- Or delete the old database and create a new one (you will lose all passwords)

Translated with DeepL.com (free version)




## aqvilov.
