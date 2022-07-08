package interfaces

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
)

type ServerCertificate interface {
	GeneratePrivateKey(privateKey *rsa.PrivateKey) (bytes.Buffer, error)
	GenerateKey(bits int) (*rsa.PrivateKey, error)
	GenerateCertificate(cert *x509.Certificate, privateKey *rsa.PrivateKey) (bytes.Buffer, error)
	CreateTemplate() *x509.Certificate
}
