package keys

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSScheme is the value stored in signing_keys.encryptor for rows sealed by AWS KMS.
const KMSScheme = "aws-kms"

const kmsVersion = 1

// KMSClient is the sliver of the AWS KMS API the encryptor uses; *kms.Client satisfies it.
type KMSClient interface {
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// KMS seals key material with AWS KMS directly: plaintexts are 32-byte curve keys, far under the 4KiB direct-encrypt limit, so no data-key envelope is needed.
// The master key never leaves KMS, which is the custody upgrade over the local AES-GCM encryptor.
type KMS struct {
	client KMSClient
	keyID  string
}

func NewKMS(client KMSClient, keyID string) (*KMS, error) {
	if client == nil || keyID == "" {
		return nil, errors.New("kms client and key id are required")
	}
	return &KMS{client: client, keyID: keyID}, nil
}

func (e *KMS) Scheme() string { return KMSScheme }

func (e *KMS) Encrypt(ctx context.Context, plaintext []byte) ([]byte, int, error) {
	out, err := e.client.Encrypt(ctx, &kms.EncryptInput{
		KeyId:     aws.String(e.keyID),
		Plaintext: plaintext,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("kms encrypt: %w", err)
	}
	return out.CiphertextBlob, kmsVersion, nil
}

func (e *KMS) Decrypt(ctx context.Context, ciphertext []byte, version int) ([]byte, error) {
	if version != kmsVersion {
		return nil, fmt.Errorf("%w: %d", ErrUnknownKeyVersion, version)
	}
	// The ciphertext blob embeds the key id; passing ours pins decryption to the expected key.
	out, err := e.client.Decrypt(ctx, &kms.DecryptInput{
		KeyId:          aws.String(e.keyID),
		CiphertextBlob: ciphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("kms decrypt: %w", err)
	}
	return out.Plaintext, nil
}
