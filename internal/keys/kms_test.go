package keys

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/stretchr/testify/require"
)

// fakeKMS reverses plaintext as its "encryption" so round-trips are observable without AWS.
type fakeKMS struct {
	encryptKeyID string
	decryptKeyID string
	err          error
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i, v := range b {
		out[len(b)-1-i] = v
	}
	return out
}

func (f *fakeKMS) Encrypt(_ context.Context, params *kms.EncryptInput, _ ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.encryptKeyID = aws.ToString(params.KeyId)
	return &kms.EncryptOutput{CiphertextBlob: reverse(params.Plaintext)}, nil
}

func (f *fakeKMS) Decrypt(_ context.Context, params *kms.DecryptInput, _ ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.decryptKeyID = aws.ToString(params.KeyId)
	return &kms.DecryptOutput{Plaintext: reverse(params.CiphertextBlob)}, nil
}

func TestKMSRoundTrip(t *testing.T) {
	fake := &fakeKMS{}
	encryptor, err := NewKMS(fake, "alias/test-key")
	require.NoError(t, err)
	require.Equal(t, KMSScheme, encryptor.Scheme())

	plaintext := []byte("super secret key material")
	ciphertext, version, err := encryptor.Encrypt(context.Background(), plaintext)
	require.NoError(t, err)
	require.Equal(t, 1, version)
	require.Equal(t, "alias/test-key", fake.encryptKeyID)
	require.False(t, bytes.Equal(plaintext, ciphertext))

	decrypted, err := encryptor.Decrypt(context.Background(), ciphertext, version)
	require.NoError(t, err)
	require.True(t, bytes.Equal(plaintext, decrypted))
	require.Equal(t, "alias/test-key", fake.decryptKeyID)
}

func TestKMSRejectsUnknownVersion(t *testing.T) {
	encryptor, err := NewKMS(&fakeKMS{}, "alias/test-key")
	require.NoError(t, err)

	_, err = encryptor.Decrypt(context.Background(), []byte("blob"), 2)
	require.ErrorIs(t, err, ErrUnknownKeyVersion)
}

func TestKMSPropagatesServiceErrors(t *testing.T) {
	encryptor, err := NewKMS(&fakeKMS{err: errors.New("throttled")}, "alias/test-key")
	require.NoError(t, err)

	_, _, err = encryptor.Encrypt(context.Background(), []byte("payload"))
	require.ErrorContains(t, err, "throttled")
}

func TestNewKMSRequiresClientAndKey(t *testing.T) {
	_, err := NewKMS(nil, "alias/test-key")
	require.Error(t, err)
	_, err = NewKMS(&fakeKMS{}, "")
	require.Error(t, err)
}
