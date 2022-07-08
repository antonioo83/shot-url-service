package services

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/services/interfaces"
	"math/big"
	"net"
	"time"
)

type serverCertificateService struct {
	serialNumber int64
	organization string
	country      string
}

func NewServerCertificate509Service(serialNumber int64, organization string, country string) interfaces.ServerCertificate {
	return &serverCertificateService{serialNumber, organization, country}
}

//CreateTemplate create new template for a certificate.
func (s serverCertificateService) CreateTemplate() *x509.Certificate {
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(s.serialNumber),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{s.organization},
			Country:      []string{s.country},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	return cert
}

//GeneratePrivateKey generate new private key.
func (s serverCertificateService) GeneratePrivateKey(privateKey *rsa.PrivateKey) (bytes.Buffer, error) {
	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return privateKeyPEM, nil
}

// GenerateKey create new private RSA key.
func (s serverCertificateService) GenerateKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("i can't generate private key: %w", err)
	}

	return privateKey, nil
}

// GenerateCertificate encode certificate to PEM format.
func (s serverCertificateService) GenerateCertificate(cert *x509.Certificate, privateKey *rsa.PrivateKey) (bytes.Buffer, error) {
	var certPEM bytes.Buffer
	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return certPEM, fmt.Errorf("i can't certificate: %w", err)
	}

	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return certPEM, nil
}
