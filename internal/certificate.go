package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

// WIP
func GenerateCertificate() {
	// 1. Generate a 4096-bit RSA private key (corresponds to -newkey rsa:4096)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	// 2. Configure the certificate template (corresponds to -days 365 and -subj "/CN=localhost")
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames:              []string{"localhost"}, // Added for compatibility with modern HTTP clients
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // 365 days
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		// Corresponds to -sha256
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	// 3. Create the self-signed certificate bytes
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// 4. Create and write the certificate to "cert.pem" (corresponds to -out cert.pem)
	certFile, err := os.Create("cert.pem")
	if err != nil {
		panic(err)
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	if err != nil {
		panic(err)
	}

	// 5. Create and write the key to "key.pem" (corresponds to -keyout key.pem)
	keyFile, err := os.Create("key.pem")
	if err != nil {
		panic(err)
	}
	defer keyFile.Close()

	err = pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		panic(err)
	}
}
