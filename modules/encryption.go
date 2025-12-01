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

func (pm *PasswordManager) Decrypt(ciphertext string) (string, error) {
	cipher, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("ошибка декодировки: %v", err)
	}

}
