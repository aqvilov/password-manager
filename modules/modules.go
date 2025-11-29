package modules

import (
	"database/sql"

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

func (pm *PasswordManager) CreatePasswordEntry(service, username, password, description string) error {
	_, err := pm.db.Exec( // exec не возвращает данные, --> нам не нужно первое значение
		`INSERT INTO password_entries(service, username, password, description) VALUES (?, ?, ?, ?)`,
		service, username, password, description,
	)
	return err
}

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
