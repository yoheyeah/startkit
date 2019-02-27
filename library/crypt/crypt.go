package crypt

type Crypter interface {
	Encryption() error
	Decryption() error
}

type Crypt struct {
	ToExt    string
	Key      string
	Dir      string
	SavePath string
}

var (
	_ Crypter = &GCM{}
	_ Crypter = &CBC{}
	_ Crypter = &CFB{}
	_ Crypter = &CTR{}
	_ Crypter = &OFB{}
	_ Crypter = &Stream{}
)

func GetCrypter(c *Crypt, mode string) Crypter {
	switch mode {
	case "GCM":
		return GCM{c}
	case "CBC":
		return CBC{c}
	case "CFB":
		return CFB{c}
	case "CTR":
		return CTR{c}
	case "OFB":
		return OFB{c}
	default:
		return Stream{c}
	}
}

func FileDycryption(c *Crypt, mode string) error {
	crypter := GetCrypter(c, mode)
	return crypter.Decryption()
}

func FileEncryption(c *Crypt, mode string) error {
	crypter := GetCrypter(c, mode)
	return crypter.Encryption()
}
