package modules

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type PasswordEntry struct {
	ID          int    // для БД
	Service     string // приложение, название и тд
	Username    string // логин
	Password    string // сам пароль
	Description string // возможное описание

}

type PasswordManager struct {
	db        *sql.DB
	masterKey []byte /*  ключ для шифрования (то есть шифрует (меняет вид пароля в БД))
	просто пароль (зашифрованный) */
}

// получаем новый пароль в БД

func NewPasswordManager(db *sql.DB, masterKey []byte) *PasswordManager {
	return &PasswordManager{
		db:        db,
		masterKey: []byte(masterKey),
	}
}

//Создания пароля

func (pm *PasswordManager) CreatePasswordEntry(service, username, password, description string) error {
	_, err := pm.db.Exec( // exec не возвращает данные, --> нам не нужно первое значение
		`INSERT INTO password_entries(service, username, password, description) VALUES (?, ?, ?, ?)`,
		service, username, password, description,
	)
	return err
}

// Удаление пароля

func (pm *PasswordManager) DeletePasswordEntry(id int) error {
	query := `DELETE FROM password_entries WHERE ID = ?`

	_, err := pm.db.Exec(query, id)
	if err != nil {
		log.Printf("Ошибка удаления пароля %v", err.Error())
	}

	return err
}

// Получение всех паролей

func (pm *PasswordManager) GetAllPasswords() ([]PasswordEntry, error) {
	command := `SELECT id, service, username, password, description FROM password_entries ORDER BY id`

	// QUERY -- ДЕЛАЕТ ЗАПРОС К БД И ВОЗВРАЩАЕТ ДАННЫЕ
	rows, err := pm.db.Query(command) // аналог Exec, работает с функциями, которые возвращают данные
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []PasswordEntry // слайс структур, в который добавляем все данные из БД, чтобы показать их

	for rows.Next() { // проходимся по каждой строке, которая содержит все (проще говоря проходимся по таблице)
		structForCopy := PasswordEntry{}

		oneStringFromColumn := rows.Scan(&structForCopy.ID, &structForCopy.Service, &structForCopy.Username,
			&structForCopy.Password, &structForCopy.Description) // КОПИРУЕМ ЦЕЛУЮ СТРОКУ В СОЗДАННУЮ НАМИ СТРУКТУРУ

		if oneStringFromColumn != nil {
			return entries, oneStringFromColumn
		}

		entries = append(entries, structForCopy)

	}
	return entries, nil

}

//ЧАСТИЧНОЕ ИЗМЕНЕНИЕ КАКОГО-ТО ПАРОЛЯ

func (pm *PasswordManager) UpdatePasswordInteractive() error {
	fmt.Println("Все текущие пароли: ")
	show, err := pm.GetAllPasswords()
	if err != nil {
		return err
	}

	if len(show) == 0 {
		return fmt.Errorf("Ошибка, у вас нет никаких паролей")
	}

	for _, entry := range show {
		fmt.Printf("ID: %d | Сервис: %s | Логин: %s\n",
			entry.ID, entry.Service, entry.Username)
	}

	// спрашиваем id, который надо менять
	var id int
	fmt.Print("\n🎯 Введите ID записи для изменения и нажмите Enter: ")
	_, err12 := fmt.Scanln(&id) // Scnaln - в отличии от Scan, читает ДО НАЖАТИЯ ENTER!!!
	if err12 != nil {
		return fmt.Errorf("ошибка ввода выбора: %v", err)
	}

	//Теперь мы должны найти все данные, которые есть в введенном id
	var AllInfo *PasswordEntry
	for _, entry := range show {
		if entry.ID == id {
			AllInfo = &entry
			break
		}
	}

	if AllInfo == nil { // если значение нулевое
		return fmt.Errorf("Запись с ID %d не найдена!", id)
	}

	fmt.Println("\n Текущие данные:")
	fmt.Printf("1. Сервис: %s\n", AllInfo.Service)
	fmt.Printf("2. Логин: %s\n", AllInfo.Username)
	fmt.Printf("3. Пароль: %s\n", AllInfo.Password)
	fmt.Printf("4. Описание: %s\n", AllInfo.Description)

	fmt.Println("\n Что вы хотите изменить?")
	fmt.Println("   1 - Сервис")    // 1
	fmt.Println("   2 - Логин")     // 2
	fmt.Println("   3 - Пароль")    // 3
	fmt.Println("   4 - Описание")  // 4
	fmt.Println("   5 - Всё сразу") // 5
	fmt.Println("   0 - Отмена")    // 0

	var choice int
	fmt.Println("Выберите номер: ")
	fmt.Scanln(&choice)

	switch choice {
	case 0:
		fmt.Println("❌ Отмена операции")
		return nil
	case 1: // меняем название сервиса
		var newService string
		fmt.Printf("Текущий сервис: %s\n", AllInfo.Service)
		fmt.Print("Новый серис: ")
		fmt.Scanln(&newService) //--------- -> ОБНОВИТ SERVICE ТОЛЬКО ГДЕ ID == ID
		_, err1 := pm.db.Exec("UPDATE password_entries SET service=? WHERE id = ?", newService, id)
		if err1 != nil {
			return fmt.Errorf("ошибка обновления сервиса: %v", err)
		}

	case 2:
		var newUsername string
		fmt.Printf("Текущий логин: %s\n", AllInfo.Username)
		fmt.Print("Новый логин: ")
		fmt.Scanln(&newUsername)
		_, err2 := pm.db.Exec("UPDATE password_entries SET username=? WHERE id=?", newUsername, id)
		if err2 != nil {
			return fmt.Errorf("ошибка обновления логина: %v", err)
		}
	case 3:
		var newPassword string
		fmt.Printf("Текущий пароль: %s\n", AllInfo.Password)
		fmt.Print("Новый пароль: ")
		fmt.Scanln(&newPassword)
		_, err3 := pm.db.Exec("UPDATE password_entries SET password=? WHERE id=?", newPassword, id)
		if err3 != nil {
			return fmt.Errorf("ошибка обновления пароля: %v", err3)
		}
	case 4:
		var newDescription string
		fmt.Printf("Текущее описание: %s\n", AllInfo.Description)
		fmt.Print("Новое описание: ")
		fmt.Scanln(&newDescription)
		_, err4 := pm.db.Exec("UPDATE	password_entries SET description=? WHERE id=?", newDescription, id)
		if err4 != nil {
			return fmt.Errorf("ошибка обновления описания: %v", err4)
		}
	case 5:
		var Service, Username, Password, Description string
		fmt.Print("Новый сервис: ")
		fmt.Scanln(&Service)
		fmt.Print("Новый логин: ")
		fmt.Scanln(&Username)
		fmt.Print("Новый пароль: ")
		fmt.Scanln(&Password)
		fmt.Print("Новое описание: ")
		fmt.Scanln(&Description)

		_, err5 := pm.db.Exec("UPDATE password_entries SET (service, username, password, description)=(?, ?, ?, ?) WHERE id=?",
			Service, Username, Password, Description, id)
		if err5 != nil {
			return fmt.Errorf("ошибка обновления данных: %v", err5)
		}
	default:
		return fmt.Errorf("неверный выбор")
	}

	return nil

}
