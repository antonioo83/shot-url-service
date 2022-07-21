package interfaces

type ServerCertificate interface {
	//SaveCertificateAndPrivateKeyToFiles save certificate and private key to the file.
	SaveCertificateAndPrivateKeyToFiles(certFileName string, privateKeyFileName string) error
}
