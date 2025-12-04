package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"password/modules"
	"strconv"
	"strings"
)

func createDB() error {
	// подключаемся в Mysql
	rootDB, err := sql.Open("mysql", "root:SQLpassforCon5@tcp(127.0.0.1:3306)/")
	if err != nil {
		return fmt.Errorf("ошибка подключения к mysql: %v", err)
	}
	defer rootDB.Close()

	if err := rootDB.Ping(); err != nil {
		return fmt.Errorf("mysql не запущен: %v", err)
	}

	_, err = rootDB.Exec("CREATE DATABASE IF NOT EXISTS password")
	if err != nil {
		return fmt.Errorf("не могу создать БД: %v", err)
	}

	return nil
}

func addPassword(pm *modules.PasswordManager) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Добавляем что-то новое...")

	fmt.Println("Введите название сервиса: ")
	service, _ := reader.ReadString('\n')
	service = service[:len(service)-1] // обрезаем переход на новую строку

	fmt.Println("А теперь введите логин ( если не хотите, нажмите Enter )")
	username, _ := reader.ReadString('\n')
	if username == "" {
		username = "" // поменять тут че-то бы
	} else {
		username = username[:len(username)-1]
	}

	fmt.Println("Введите пароль")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	fmt.Println("Введите описание (если не хотите, нажмите Enter )")
	description, _ := reader.ReadString('\n')
	description = description[:len(description)-1]

	add := pm.CreatePasswordEntry(service, username, password, description)
	if add != nil {
		return
	} else {
		fmt.Println("Данные успешно добавлены!")
	}

	showAllPasswords(pm)
}

func showAllPasswords(pm *modules.PasswordManager) {
	ent, err := pm.GetAllPasswords()
	if err != nil {
		log.Fatal(" Ошибка получения паролей:", err)
	}

	if len(ent) == 0 {
		fmt.Println("Пока что тут ничего нет;(")
	} else {
		fmt.Printf(" Всего записей: %d\n", len(ent))
		for _, entry := range ent {
			fmt.Printf("   ID: %d | Сервис: %s | Логин: %s | Пароль: %s\n",
				entry.ID, entry.Service, entry.Username, entry.Password)
		}
	}
}

func deletePassword(pm *modules.PasswordManager) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите ID пароля, который хотите удалить: ") // как-то переписать немного строчку ( странно звучит )
	id, _ := reader.ReadString('\n')                           // читает строку до нажатия enter
	id = id[:len(id)-1]

	idInt, _ := strconv.Atoi(id)

	for {
		if idInt < 0 {
			fmt.Printf("Введите корректное число!")
		} else {
			rm := pm.DeletePasswordEntry(idInt)
			fmt.Println("Пароль успешно удален!)")
			if rm != nil {
				return
			}
		}

	}

}

func searchPassword(pm *modules.PasswordManager) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\nВведите название сервиса для посика: \n")
	service, _ := reader.ReadString('\n')
	service = strings.TrimSpace(service) // убираем лишние символы по бокам

	if service == "" {
		fmt.Println("Ошибка: введите строку для поиска")
		return
	}

	entries, err := pm.GetAllPasswords()
	if err != nil {
		fmt.Println("ошибка вывода")
		return
	}

	fmt.Println("Поиск по названию: ", service)
	fmt.Println(strings.Repeat("-", 30))

	//ДЛЯ ЛЮБОГО ЦИКЛА, ПЕРВОЕ ЗНАЧЕНИЕ - ИНДЕКС ЭЛЕМЕНТА, ВТОРОЕ - ЕГО ЗНАЧЕНИЕ
	for index, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Service), strings.ToLower(service)) || // поиск по серивсу
			strings.Contains(strings.ToLower(entry.Username), strings.ToLower(service)) || // поиск по юзеру
			strings.Contains(strings.ToLower(entry.Description), strings.ToLower(service)) { // поиск по описанию

			//вывод данных
			fmt.Println()
			fmt.Print("№", index+1)
			fmt.Printf("     ID: %d\n", entry.ID)
			fmt.Printf("     Сервис: %s\n", entry.Service)
			fmt.Printf("     Логин: %s\n", entry.Username)
			fmt.Printf("     Пароль: %s\n", entry.Password)
			if entry.Description != "" {
				fmt.Printf("     Description: %s\n", entry.Description)
			}
			fmt.Println()
			fmt.Println(strings.Repeat("-", 30))
		} else {
			fmt.Println("Ничего не найдено;(")
			fmt.Println()
		}
	}

}

func main() {

	checkDB := createDB()
	if checkDB != nil {
		log.Fatal("ошибка создания БД", checkDB)
	}

	db, err := sql.Open("mysql", "root:SQLpassforCon5@tcp(127.0.0.1:3306)/password") // password - name of database
	if err != nil {
		return
	}

	//очистка таблицы при запуске
	//_, _ = db.Exec("DELETE FROM password_entries")
	//_, _ = db.Exec("ALTER TABLE password_entries AUTO_INCREMENT = 1")
	//fmt.Println("Таблица очищена от старых незашифрованных данных")

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное подключение к БД!")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS password_entries (
		id INT AUTO_INCREMENT PRIMARY KEY,
		service VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		password TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// создание таблицы
	_, errCreateTables := db.Exec(createTableSQL)
	if errCreateTables != nil {
		log.Fatal("ошибка создания таблиц", errCreateTables)
	}

	byteKey := []byte("my-32-byte-super-secret-key-1234") // пока что мастер ключ тут, потом перенесем в функцию
	//fmt.Println("Длина ключа шифрования: ", len(byteKey))

	pm := modules.NewPasswordManager(db, byteKey)

	//делаем паузу перед началом

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Нажмите Enter чтобы продолжить...")
	reader.ReadString('\n')

	// главный цикл в main
	for {
		// приветственное меню
		fmt.Println("\n" + strings.Repeat("═", 50))
		fmt.Println("МЕНЕДЖЕР ПАРОЛЕЙ")
		fmt.Println(strings.Repeat("═", 50))

		fmt.Println("\n ГЛАВНОЕ МЕНЮ:")
		fmt.Println("1.  Показать все пароли")
		fmt.Println("2.  Добавить новый пароль")
		fmt.Println("3.  Изменить существующий пароль")
		fmt.Println("4.  Удалить пароль")
		fmt.Println("5.  Поиск паролей")
		// fmt.Println("6.  Очистить экран") // временно не работает
		fmt.Println("0.  Выход")
		fmt.Print("\nВыберите действие (0-6): ")

		mainReader := bufio.NewReader(os.Stdin)

		choice, _ := mainReader.ReadString('\n') // читаем до Enter
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			showAllPasswords(pm)
		case "2":
			addPassword(pm)
		case "3":
			err := pm.UpdatePasswordInteractive()
			if err != nil {
				fmt.Printf("\n❌ Ошибка: %v\n", err)
			}
		case "4":
			deletePassword(pm)
		case "5":
			searchPassword(pm)
		case "0":
			fmt.Println("До свидания! Ваши пароли в безопасности!!!")
			return
		default:
			fmt.Println("Ошибка! Введите число от 0 до 6!")
			return
		}

		// к каждой команде, кроме case 0, добавляем возможность вернуться в меню
		if choice != "0" {
			fmt.Print("\nНажмите Enter чтобы вернуться в меню...")
			reader.ReadString('\n')
		}

	}
}
