package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"password/modules"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func masterKey() []byte { //temporary
	masterKey := []byte("my-32-byte-super-secret-key-1234")
	return masterKey
}

func createDB() error {
	connStr := "host=localhost port=5432 user=postgres password=SQLpassforCon5 dbname=password sslmode=disable"
	rootDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}
	defer rootDB.Close()

	if err := rootDB.Ping(); err != nil {
		return fmt.Errorf("PostgreSQL не запущен: %v", err)
	}

	var exists bool
	err = rootDB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'password')",
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка проверки базы данных: %v", err)
	}

	if !exists {
		_, err = rootDB.Exec("CREATE DATABASE password")
		if err != nil {
			return fmt.Errorf("не могу создать БД: %v", err)
		}
		fmt.Println(" База данных создана")
	} else {
		fmt.Println(" База данных уже существует")
	}

	return nil
}

func addPassword(pm *modules.PasswordManager) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Добавляем что-то новое...")

	fmt.Println("Введите название сервиса: ")
	service, _ := reader.ReadString('\n')
	service = strings.TrimSpace(service) // обрезаем переход на новую строку

	fmt.Println("А теперь введите логин ( если не хотите, нажмите Enter )")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Введите пароль")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Println("Введите описание (если не хотите, нажмите Enter )")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	add := pm.CreatePasswordEntry(service, username, password, description)
	if add != nil {
		return
	} else {
		fmt.Println("Данные успешно добавлены!")
	}

	//showAllPasswords(pm)
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
		fmt.Println()
		for i, entry := range ent {
			fmt.Printf("   ID: %d | Сервис: %s | Логин: %s | Пароль: %s | Описание: %s\n",
				i+1, entry.Service, entry.Username, entry.Password, entry.Description)
			fmt.Println()
		}
	}
}

func deletePassword(pm *modules.PasswordManager) {
	// выводим текущие пароли для удобства пользователя
	fmt.Println("Все текущие пароли: ")
	show, err := pm.GetAllPasswords()
	if err != nil {
		return
	}

	if len(show) == 0 {
		fmt.Println("У вас пока нет никаих паролей")
	}

	for i, entry := range show {
		fmt.Printf("ID: %d | Сервис: %s | Логин: %s\n",
			i+1, entry.Service, entry.Username) // i + 1 это айди, который отображается пользователю (не айди из БД)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите ID пароля, который хотите удалить: ") // как-то переписать немного строчку ( странно звучит )
	id, _ := reader.ReadString('\n')                           // читает строку до нажатия enter
	id = strings.TrimSpace(id)

	idInt, _ := strconv.Atoi(id)

	realID := show[idInt-1].ID // восстанавливаем id из БДшки, чтобы программе работать с ним

	for {
		if idInt < 0 {
			fmt.Printf("Введите корректное число!")
		} else {
			rm := pm.DeletePasswordEntry(realID)
			if rm != nil {
				return
			}
			fmt.Println("Пароль успешно удален!)")
			return
		}

	}

}

func searchPassword(pm *modules.PasswordManager) {
	reader := bufio.NewReader(os.Stdin) // читаем всю строку целиком

	fmt.Printf("\nВведите ключевые слова для поиска (Сервис, Логин, Описание): \n")
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

	fmt.Println("\nПоиск по названию: ", service)
	fmt.Println(strings.Repeat("-", 30))

	matchesCount := 0

	for index, entry := range entries {

		searchByService := strings.Contains(strings.ToLower(entry.Service), strings.ToLower(service))
		searchByUser := strings.Contains(strings.ToLower(entry.Username), strings.ToLower(service))
		searchByDescription := strings.Contains(strings.ToLower(entry.Description), strings.ToLower(service))

		if searchByUser || searchByDescription || searchByService { // поиск по описанию
			matchesCount++
			//вывод данных
			fmt.Println()
			fmt.Printf("   ID: %d\n", index+1)
			fmt.Printf("     Сервис: %s\n", entry.Service)
			fmt.Printf("     Логин: %s\n", entry.Username)
			fmt.Printf("     Пароль: %s\n", entry.Password)
			fmt.Printf("     Description: %s\n", entry.Description)
			fmt.Println()
			fmt.Println(strings.Repeat("-", 30))
		}
	}

	if matchesCount != 0 {
		fmt.Printf("Найдено совпадений: %d\n", matchesCount)
	} else {
		fmt.Println("Совпадений не найдено ;(")
	}

}

func GreetingMenu() {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("МЕНЕДЖЕР ПАРОЛЕЙ")
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n ГЛАВНОЕ МЕНЮ:")
	fmt.Println("1.  Показать все пароли")
	fmt.Println("2.  Добавить новый пароль")
	fmt.Println("3.  Изменить существующий пароль")
	fmt.Println("4.  Удалить пароль")
	fmt.Println("5.  Поиск паролей")
	fmt.Println("6.  Очистить экран") // временно не работает
	fmt.Println("0.  Выход")
	fmt.Print("\nВыберите действие (0-6): ")
}

func clearScreen() {

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout // подключаем вывод к стандартному выводу нашей программы
	cmd.Run()

}

func main() {

	checkDB := createDB()
	if checkDB != nil {
		log.Fatal("ошибка создания БД", checkDB)
	}

	connStr := "host=localhost port=5432 user=postgres password=SQLpassforCon5 dbname=password sslmode=disable"
	db, err := sql.Open("postgres", connStr) // password - name of database
	if err != nil {
		return
	}

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
		id SERIAL PRIMARY KEY,
		service VARCHAR(255),
		username VARCHAR(255),
		password TEXT NOT NULL,
		description TEXT
	)`

	// создание таблицы
	_, errCreateTables := db.Exec(createTableSQL)
	if errCreateTables != nil {
		log.Fatal("ошибка создания таблиц", errCreateTables)
	}

	byteKey := masterKey() // пока что мастер ключ тут, потом перенесем в функцию
	//fmt.Println("Длина ключа шифрования: ", len(byteKey))

	pm := modules.NewPasswordManager(db, byteKey)

	//делаем паузу перед началом

	reader := bufio.NewReader(os.Stdin)

	// анимация ожидания ввода текста/символа
	var animation = []string{"...", ".. ", ".  ", "   "}
	done := make(chan bool)

	go func() {
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			default:
				index := i % len(animation)
				fmt.Printf("\r%s%s", "Нажмите Enter чтобы продолжить", animation[index])
				time.Sleep(350 * time.Millisecond)
			}
		}

	}()

	readString, _ := reader.ReadString('\n')
	readString = strings.TrimSpace(readString)

	close(done) // закрываем канал --> передается сигнал --> return из цикла

	if readString == "void" {
		fmt.Println("Вы ввели секретное слово\nМожете посмотреть это видео: https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=RDdQw4w9WgXcQ&start_radio=1")
	}

	// главный цикл в main
	for {
		// приветственное меню
		GreetingMenu()

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
				fmt.Printf("\n Ошибка: %v\n", err)
			}
		case "4":
			deletePassword(pm)
		case "5":
			searchPassword(pm)
		case "0":
			fmt.Println()
			fmt.Println("До свидания!\n Ваши пароли в безопасности!!!")
			fmt.Println()
			return
		case "6":
			clearScreen()

		default:
			fmt.Println("Ошибка! Введите число от 0 до 6!")
			fmt.Println()
			fmt.Println()
			continue // НА ГЕНИАЛЫЧАХ!!!
		}

		// к каждой команде, кроме case 0, добавляем возможность вернуться в меню
		if choice != "0" {
			fmt.Print("\nНажмите Enter чтобы вернуться в меню...")
			reader.ReadString('\n')
		}

	}
}
