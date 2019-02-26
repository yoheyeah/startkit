package crypt

type OFB struct {
	*Crypt
}

func (d OFB) Encryption() error {
	return nil
}

func (d OFB) Decryption() error {
	return nil
}
