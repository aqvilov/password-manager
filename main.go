package main

import (
	"database/sql"
	"fmt"
	"log"
	"password/modules"
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
	_, errCreateTables := db.Exec(createTableSQL)
	if errCreateTables != nil {
		log.Fatal("ошибка создания таблиц", errCreateTables)
	}

	byteKey := []byte("my-32-byte-super-secret-key-1234")
	fmt.Println(len(byteKey))

	pm := modules.NewPasswordManager(db, byteKey)

	fmt.Println("Добавляем тестовый пароль")
	err1 := pm.CreatePasswordEntry("Telegram", "Aqvi", "12345678", "test")
	err2 := pm.CreatePasswordEntry("Telegram2", "Aqvi", "12345678", "test")
	if err1 != nil {
		log.Fatal(err1)
	}
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("Пароль добавлен!")

	ent, err := pm.GetAllPasswords()
	if err != nil {
	}

	if len(ent) == 0 {
		log.Println("БАЗА ДАННЫХ ПУСТАЯ")
	} else {
		for _, entry := range ent {
			fmt.Printf("%d. Сервис: %s | Логин: %s | Пароль: %s |\nОписание: %s\n", entry.ID, entry.Service, entry.Username, entry.Password, entry.Description)
		}
	}
	showAllPasswords(pm)

	testDeleting := pm.DeletePasswordEntry(1)
	if testDeleting != nil {
		return
	}
	showAllPasswords(pm)

	err15 := pm.UpdatePasswordInteractive()
	if err15 != nil {
		return
	}
	showAllPasswords(pm)

}

func showAllPasswords(pm *modules.PasswordManager) {
	ent, err := pm.GetAllPasswords()
	if err != nil {
		log.Fatal(" Ошибка получения паролей:", err)
	}

	if len(ent) == 0 {
		fmt.Println(" База данных пустая")
	} else {
		fmt.Printf(" Всего записей: %d\n", len(ent))
		for _, entry := range ent {
			fmt.Printf("   ID: %d | Сервис: %s | Логин: %s | Пароль: %s\n",
				entry.ID, entry.Service, entry.Username, entry.Password)
		}
	}
}
