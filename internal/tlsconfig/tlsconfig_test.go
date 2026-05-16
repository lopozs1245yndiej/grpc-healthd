package tlsconfig_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/example/grpc-healthd/internal/tlsconfig"
)

func TestDefaultConfig_Disabled(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_TLS_CERT", "")
	t.Setenv("GRPC_HEALTHD_TLS_KEY", "")
	cfg := tlsconfig.DefaultConfig()
	if cfg.Enabled {
		t.Fatal("expected TLS to be disabled when env vars are empty")
	}
}

func TestDefaultConfig_Enabled(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_TLS_CERT", "/tmp/cert.pem")
	t.Setenv("GRPC_HEALTHD_TLS_KEY", "/tmp/key.pem")
	cfg := tlsconfig.DefaultConfig()
	if !cfg.Enabled {
		t.Fatal("expected TLS to be enabled when both env vars are set")
	}
}

func TestLoad_Disabled(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: false}
	tc, err := tlsconfig.Load(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tc != nil {
		t.Fatal("expected nil tls.Config when disabled")
	}
}

func TestLoad_ValidCertAndKey(t *testing.T) {
	certFile, keyFile := writeSelfSigned(t)
	cfg := tlsconfig.Config{Enabled: true, CertFile: certFile, KeyFile: keyFile}
	tc, err := tlsconfig.Load(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tc == nil {
		t.Fatal("expected non-nil tls.Config")
	}
	if len(tc.Certificates) != 1 {
		t.Fatalf("expected 1 certificate, got %d", len(tc.Certificates))
	}
}

func TestLoad_InvalidFiles(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: true, CertFile: "/nonexistent/cert.pem", KeyFile: "/nonexistent/key.pem"}
	_, err := tlsconfig.Load(cfg)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}

// writeSelfSigned creates a temporary self-signed cert+key pair and returns
// the file paths. Files are removed when the test completes.
func writeSelfSigned(t *testing.T) (certPath, keyPath string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}
	certF, _ := os.CreateTemp(t.TempDir(), "cert*.pem")
	_ = pem.Encode(certF, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	certF.Close()
	keyBytes, _ := x509.MarshalECPrivateKey(priv)
	keyF, _ := os.CreateTemp(t.TempDir(), "key*.pem")
	_ = pem.Encode(keyF, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	keyF.Close()
	return certF.Name(), keyF.Name()
}
