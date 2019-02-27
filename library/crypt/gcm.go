package crypt

type GCM struct {
	*Crypt
}

func (d GCM) Encryption() error {
	return nil
}

func (d GCM) Decryption() error {
	return nil
}
