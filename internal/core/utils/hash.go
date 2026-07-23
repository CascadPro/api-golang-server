package core_utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Параметры хеширования
var (
	saltLength         = 16
	keyLength   uint32 = 32
	memory      uint32 = 64 * 1024
	iterations  uint32 = 3
	parallelism uint8  = 2
)

func GenerateHash(password string) (string, error) {
	// 1. Генерируем случайную соль
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Вычисляем хеш с помощью Argon2id (наиболее безопасный режим)
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// 3. Кодируем соль и хеш в base64 для хранения в БД
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 4. Собираем итоговую строку в формате PHP
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	// 1. Разбираем строку на части
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("неверный формат хеша")
	}

	// 2. Парсим параметры версии и памяти/итераций/параллелизма
	var version int
	var memoryParam, iterationsParam, parallelismParam int

	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, fmt.Errorf("несовместимая версия argon2")
	}

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memoryParam, &iterationsParam, &parallelismParam)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	existingHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// 3. Вычисляем хеш из введенного пароля с теми же параметрами
	computedHash := argon2.IDKey([]byte(password), salt, uint32(iterationsParam), uint32(memoryParam), uint8(parallelismParam), uint32(keyLength))

	// 4. Сравниваем хеши. Используем subtle.ConstantTimeCompare для защиты от атак по времени.
	if subtle.ConstantTimeCompare(existingHash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}
