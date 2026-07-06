// Package keys provides envelope encryption for custodial signing-key material stored in Postgres.
package keys

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Encryptor is the seam between key storage and the encryption backend.
// The local AES-GCM implementation serves dev and demo; the AWS KMS implementation serves cloud deployments behind the same interface without a schema change.
type Encryptor interface {
	// Scheme names the implementation, stored per row so decryption dispatches to the right backend.
	Scheme() string
	// Encrypt seals plaintext and reports the key version it was sealed with.
	Encrypt(ctx context.Context, plaintext []byte) (ciphertext []byte, version int, err error)
	// Decrypt opens ciphertext that was sealed with the given key version.
	Decrypt(ctx context.Context, ciphertext []byte, version int) ([]byte, error)
}

// AESGCMScheme is the value stored in signing_keys.encryptor for rows sealed by AESGCM.
const AESGCMScheme = "local-aes-gcm"

const aesGCMVersion = 1

var (
	ErrInvalidMasterKey   = errors.New("master key must be 32 bytes of hex")
	ErrUnknownKeyVersion  = errors.New("unknown key version")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

// AESGCM encrypts with AES-256-GCM under a single master key, prepending the random nonce to each ciphertext.
type AESGCM struct {
	aead cipher.AEAD
}

func NewAESGCM(masterKeyHex string) (*AESGCM, error) {
	masterKeyHex = strings.TrimPrefix(strings.TrimSpace(masterKeyHex), "0x")
	key, err := hex.DecodeString(masterKeyHex)
	if err != nil || len(key) != 32 {
		return nil, ErrInvalidMasterKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}
	return &AESGCM{aead: aead}, nil
}

func (e *AESGCM) Scheme() string { return AESGCMScheme }

func (e *AESGCM) Encrypt(_ context.Context, plaintext []byte) ([]byte, int, error) {
	nonce := make([]byte, e.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, 0, fmt.Errorf("generate nonce: %w", err)
	}
	return e.aead.Seal(nonce, nonce, plaintext, nil), aesGCMVersion, nil
}

func (e *AESGCM) Decrypt(_ context.Context, ciphertext []byte, version int) ([]byte, error) {
	if version != aesGCMVersion {
		return nil, fmt.Errorf("%w: %d", ErrUnknownKeyVersion, version)
	}
	if len(ciphertext) < e.aead.NonceSize() {
		return nil, ErrCiphertextTooShort
	}
	nonce, sealed := ciphertext[:e.aead.NonceSize()], ciphertext[e.aead.NonceSize():]
	plaintext, err := e.aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt key material: %w", err)
	}
	return plaintext, nil
}
