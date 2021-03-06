package kube

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/jetstack-experimental/cert-manager/pkg/util/errors"
	"github.com/jetstack-experimental/cert-manager/pkg/util/pki"
	api "k8s.io/api/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
)

func GetKeyPair(secretLister corelisters.SecretLister, namespace, name string) (certBytes []byte, keyBytes []byte, err error) {
	secret, err := secretLister.Secrets(namespace).Get(name)

	if err != nil {
		return nil, nil, err
	}

	certBytes = secret.Data[api.TLSCertKey]
	keyBytes = secret.Data[api.TLSPrivateKeyKey]

	return certBytes, keyBytes, err
}

func SecretTLSKey(secretLister corelisters.SecretLister, namespace, name string) (*rsa.PrivateKey, error) {
	secret, err := secretLister.Secrets(namespace).Get(name)

	if err != nil {
		return nil, err
	}

	keyBytes := secret.Data[api.TLSPrivateKeyKey]
	key, err := pki.DecodePKCS1PrivateKeyBytes(keyBytes)

	if err != nil {
		return key, errors.NewInvalidData(err.Error())
	}
	return key, nil
}

func SecretTLSCert(secretLister corelisters.SecretLister, namespace, name string) (*x509.Certificate, error) {
	secret, err := secretLister.Secrets(namespace).Get(name)

	if err != nil {
		return nil, err
	}

	certBytes := secret.Data[api.TLSCertKey]
	cert, err := pki.DecodeX509CertificateBytes(certBytes)

	if err != nil {
		return cert, errors.NewInvalidData(err.Error())
	}
	return cert, nil
}
