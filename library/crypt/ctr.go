package crypt

type CTR struct {
	*Crypt
}

func (d CTR) Encryption() error {
	return nil
}

func (d CTR) Decryption() error {
	return nil
}
