package modules

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func (pm *PasswordManager) Encrypt(plaintext string) (string, error) {
	plain := []byte(plaintext)

	block, err := aes.NewCipher(pm.masterKey) // создает специальный алогорит AES
	if err != nil {
		return "", fmt.Errorf("ошибка создания cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block) // GCM - гарантия, что не будет подмены шифровальных данных
	if err != nil {
		return "", fmt.Errorf("ошибка создания gcm: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize()) // просто случайные байты, чтобы одинаковые пароли шифровались по разному

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("ошибка генерации nonce: %v", err)
	}

	// шифруем данные
	ciphertext := gcm.Seal(nonce, nonce, plain, nil)
	//---------------------nonce -куда ложим nonce в результате
	//----------------------------nonce -какой именно nonce используем (тот, что и создали и положили в начало)
	//-----------------------------------plain -что именно шифруем
	//------------------------------------------------------nil -доп данные для проверки(не юз)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

//довести до конца

package modules

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"encoding/base64"
"fmt"
"io"
)

// Encrypt шифрует строку с ключом
func Encrypt(plaintext string, key []byte) (string, error) {
	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("ошибка создания cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка создания GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("ошибка генерации nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintextBytes, nil) // ШИФРУЕМ
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}


// расшифровываем

func Decrypt(encryptedText string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("ошибка создания cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка создания GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()

	if len(ciphertext) < nonceSize {
		return	"", fmt.Errorf("неверная длина ciphertext")
	}

	// извлекаем nonce
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:] // т.к мы знаем, что сначала nonce потом text

	plaintext, err := gcm.Open(nil,nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil

}