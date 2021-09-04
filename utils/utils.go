package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"io"
	"math/rand"
	"strconv"
	"thxy/setting"
	"time"
	crand "crypto/rand"
)

const (
	digestLength       = 4
	defaultNonceLength = 12
	secretKeyLength    = 16
	SecretKey          = `Uy8&9@iL186BvNcc`
)

var (
	HashStrength = bcrypt.MinCost + 1

	ErrInvalidParam      = errors.New("invalid param")
	ErrInvalidCiphertext = errors.New("invalid encrypted")
)

func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

func CurrentTimestr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(max - min)
	return i + min
}

// GenRandNum 生成随机数字
func GenRandNumStr(size int) string {
	if size <= 0 {
		size = 6
	}
	str := "0123456789"
	bytes := []byte(str)
	result := make([]byte, size)

	result[0] = bytes[RandInt(1, 9)]
	for i := 1; i < size; i++ {
		result[i] = bytes[RandInt(0, 9)]
	}
	return string(result)
}

func GenUserCode() string {
	return fmt.Sprintf("%v%v", CurrentTimestamp(), GenRandNumStr(6))
}



// GenRandStr 生成随机字符串 数字+大写+小写
func GenRandStr(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetTimeZeroLastSecond(timeDuration string) string {
	now := time.Now()
	duration, _ := time.ParseDuration(timeDuration)
	newTime := now.Add(duration)
	newTimeStr := newTime.Format("2006-01-02") + " 00:00:00"
	return newTimeStr
}

func GetTimeStrFromSecond(seconds int) string {
	minute := strconv.Itoa(seconds / 60)
	second := strconv.Itoa(seconds % 60)

	if second == "0" || second == "" {
		return minute + ":00"
	} else {
		return minute + ":" + second
	}
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func GenToken(size int) (nonce []byte, err error) {
	nonce = make([]byte, size)

	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}


func Encrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) == 0 || len(plaintext) == 0 {
		return nil, ErrInvalidParam
	}

	aesgcm, err := newAEAD(key)
	if err != nil {
		return nil, err
	}

	size := defaultNonceLength
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce, err := GenToken(size)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return digest(ciphertext, nonce), nil
}

// EncryptPassword 对入库的用户密码加密
func EncryptPassword(phrase string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(phrase), HashStrength)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Decrypt(key, encrypted []byte) ([]byte, error) {
	length := len(encrypted)
	size := defaultNonceLength
	ciphertext := encrypted[:length-size-digestLength]
	nonce := encrypted[length-size-digestLength : length-digestLength]
	if !verify(encrypted) {
		return nil, ErrInvalidCiphertext
	}

	aesgcm, err := newAEAD(key)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func verify(encrypted []byte) bool {
	length := len(encrypted)
	cksm := encrypted[length-digestLength:]
	hash := blake2b.Sum256(encrypted[:length-digestLength])
	return bytes.Equal(cksm, hash[:digestLength])
}


func digest(ciphertext, nonce []byte) []byte {
	l1 := len(ciphertext)
	l2 := l1 + len(nonce)

	output := make([]byte, l2+4)

	copy(output, ciphertext)
	copy(output[l1:], nonce)
	digest := blake2b.Sum256(output[:l2])
	copy(output[l2:], digest[:4])

	return output
}

func GenSid() string {
	//return fmt.Sprintf("%x", string(bson.NewObjectId()))
	s, _ := GenerateRandomString(32)
	return fmt.Sprintf("%s", s)
}

func ComparePassword(phrase, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(phrase)) == nil
}

func Password(str string) string {
	str = setting.PasswordPrefix + str
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, _ := GenerateRandomBytes(s)
	return hex.EncodeToString(b), nil
}