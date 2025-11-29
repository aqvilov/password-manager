package main

import (
	"database/sql"
	"fmt"
	"log"
	"password/modules"
)

func main() {
	db, err := sql.Open("mysql", "root:SQLpassforCon5@tcp(127.0.0.1:3306)/password") // password - name of database
	if err != nil {
		return
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное подключение к БД!")

	defer db.Close()

	pm := modules.NewPasswordManager(db, []byte("123456"))

	fmt.Println("Добавляем тестовый пароль")
	err1 := pm.CreatePasswordEntry("Telegram", "Aqvi", "123456", "test")
	if err1 != nil {
		log.Fatal(err1)
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
}
