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
	"github.com/antonioo83/shot-url-service/internal/utils"
	"math/big"
	"net"
	"os"
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

//SaveCertificateAndPrivateKeyToFiles save certificate and private key to the file.
func (s serverCertificateService) SaveCertificateAndPrivateKeyToFiles(certFileName string, privateKeyFileName string) error {
	certificatePEM, privateKeyPEM, err := s.generateCertificateAndPrivateKey()
	if err != nil {
		return fmt.Errorf("i can't open a file: %w", err)
	}

	err = s.saveToFile(certFileName, certificatePEM.Bytes())
	if err != nil {
		return fmt.Errorf("i can't open a file: %w", err)
	}

	err = s.saveToFile(privateKeyFileName, privateKeyPEM.Bytes())
	if err != nil {
		return fmt.Errorf("i can't open a file: %w", err)
	}

	return nil
}

// generateCertificateAndPrivateKey encode certificate to PEM format.
func (s serverCertificateService) generateCertificateAndPrivateKey() (certificatePEM bytes.Buffer, privatePEM bytes.Buffer, error error) {
	var certPEM bytes.Buffer
	var privateKeyPEM bytes.Buffer

	key, err := s.generateKey(4096)
	if err != nil {
		return certPEM, privateKeyPEM, fmt.Errorf("i can't generate a rsa key: %w", err)
	}

	privateKeyPEM, err = s.generatePrivateKey(key)
	if err != nil {
		return certPEM, privateKeyPEM, fmt.Errorf("i can't generate a private key: %w", err)
	}

	template := s.createTemplate()
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return certPEM, privateKeyPEM, fmt.Errorf("i can't create a certificate: %w", err)
	}

	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return certPEM, privateKeyPEM, nil
}

//saveToFile save array of byte to the file.
func (s serverCertificateService) saveToFile(fileName string, data []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("i can't open a file: %w", err)
	}
	defer utils.ResourceClose(file)

	err = utils.LogErr(file.Write(data))
	if err != nil {
		return fmt.Errorf("i can't write to file: %w", err)
	}

	return nil
}

//createTemplate create new template for a certificate.
func (s serverCertificateService) createTemplate() *x509.Certificate {
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

//generatePrivateKey generate new private key.
func (s serverCertificateService) generatePrivateKey(privateKey *rsa.PrivateKey) (bytes.Buffer, error) {
	var privateKeyPEM bytes.Buffer
	err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return privateKeyPEM, fmt.Errorf("can't generate private key: %w", err)
	}

	return privateKeyPEM, nil
}

// generateKey create new private RSA key.
func (s serverCertificateService) generateKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("i can't generate private key: %w", err)
	}

	return privateKey, nil
}
