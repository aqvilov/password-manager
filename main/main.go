package main

import (
	"database/sql"
	"fmt"
	"log"
	"password-manager/modules"
)

func main() {
	db, err := sql.Open("mysql", "root:SQLpassforCon5@tcp(127.0.0.1:3306)/password") // password - name of database
	if err != nil {
		return
	}

	_, _ = db.Exec("TRUNCATE TABLE password_entries") // сбрасывает auto_increment

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное подключение к БД!")

	defer db.Close()

	pm := modules.NewPasswordManager(db, []byte("123456"))

	fmt.Println("Добавляем тестовый пароль")
	err1 := pm.CreatePasswordEntry("Telegram", "Aqvi", "123456", "test")
	err2 := pm.CreatePasswordEntry("Telegram2", "Aqvi", "123456", "test")
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
		log.Fatal("❌ Ошибка получения паролей:", err)
	}

	if len(ent) == 0 {
		fmt.Println("📭 База данных пустая")
	} else {
		fmt.Printf("📊 Всего записей: %d\n", len(ent))
		for _, entry := range ent {
			fmt.Printf("   ID: %d | Сервис: %s | Логин: %s | Пароль: %s\n",
				entry.ID, entry.Service, entry.Username, entry.Password)
		}
	}
}
