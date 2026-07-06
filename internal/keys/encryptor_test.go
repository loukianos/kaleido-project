package keys

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testMasterKey = "56659eba0b040dfd7503c987c0c26428a476b3ee49747a181f6f7361e253e4a3"

func TestAESGCMRoundTrip(t *testing.T) {
	enc, err := NewAESGCM(testMasterKey)
	require.NoError(t, err)

	plaintext := []byte("super secret secp256k1 key material")
	ciphertext, version, err := enc.Encrypt(context.Background(), plaintext)
	require.NoError(t, err)
	require.Equal(t, 1, version)
	require.NotContains(t, string(ciphertext), string(plaintext))

	decrypted, err := enc.Decrypt(context.Background(), ciphertext, version)
	require.NoError(t, err)
	require.True(t, bytes.Equal(plaintext, decrypted))
}

func TestAESGCMNoncesDiffer(t *testing.T) {
	enc, err := NewAESGCM("0x" + testMasterKey)
	require.NoError(t, err)

	first, _, err := enc.Encrypt(context.Background(), []byte("same plaintext"))
	require.NoError(t, err)
	second, _, err := enc.Encrypt(context.Background(), []byte("same plaintext"))
	require.NoError(t, err)
	require.False(t, bytes.Equal(first, second))
}

func TestAESGCMTamperDetected(t *testing.T) {
	enc, err := NewAESGCM(testMasterKey)
	require.NoError(t, err)

	ciphertext, version, err := enc.Encrypt(context.Background(), []byte("payload"))
	require.NoError(t, err)

	ciphertext[len(ciphertext)-1] ^= 0xff
	_, err = enc.Decrypt(context.Background(), ciphertext, version)
	require.Error(t, err)
}

func TestAESGCMWrongKeyFails(t *testing.T) {
	enc, err := NewAESGCM(testMasterKey)
	require.NoError(t, err)
	other, err := NewAESGCM(strings.Repeat("ab", 32))
	require.NoError(t, err)

	ciphertext, version, err := enc.Encrypt(context.Background(), []byte("payload"))
	require.NoError(t, err)

	_, err = other.Decrypt(context.Background(), ciphertext, version)
	require.Error(t, err)
}

func TestAESGCMRejectsBadInputs(t *testing.T) {
	_, err := NewAESGCM("not hex")
	require.ErrorIs(t, err, ErrInvalidMasterKey)
	_, err = NewAESGCM("abcd")
	require.ErrorIs(t, err, ErrInvalidMasterKey)

	enc, err := NewAESGCM(testMasterKey)
	require.NoError(t, err)

	_, err = enc.Decrypt(context.Background(), []byte("short"), 1)
	require.ErrorIs(t, err, ErrCiphertextTooShort)

	ciphertext, _, err := enc.Encrypt(context.Background(), []byte("payload"))
	require.NoError(t, err)
	_, err = enc.Decrypt(context.Background(), ciphertext, 2)
	require.ErrorIs(t, err, ErrUnknownKeyVersion)
}
