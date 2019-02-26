package crypt

type CFB struct {
	*Crypt
}

func (d CFB) Encryption() error {
	return nil
}

func (d CFB) Decryption() error {
	return nil
}
