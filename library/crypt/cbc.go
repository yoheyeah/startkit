package crypt

type CBC struct {
	*Crypt
}

func (d CBC) Encryption() error {
	return nil
}

func (d CBC) Decryption() error {
	return nil
}
