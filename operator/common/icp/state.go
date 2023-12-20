package icp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/certificate"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/fxamacker/cbor/v2"
	"github.com/gibson042/canonicaljson-go"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type TororuSecretResource struct {
	Name        string
	Status      string
	ROConsumers []string
	RWConsumer  string
}

type State struct {
	Cert    Certificate
	Secrets []TororuSecretResource
}

type Certificate struct{}

func GetCanisterStatus() (*CertifiedStatus, error) {
	client, err := GetBackendClient()
	if err != nil {
		return nil, err
	}

	state, err := client.GetCertifiedStatus()
	if err != nil {
		return nil, err
	}

	// Validate cert here
	valid, err := validateCertificate(client.canisterId, state)
	if err != nil {
		zap.S().Errorf("certificate validation failed with err: %v", err)
		return nil, err
	}
	if !valid {
		err := fmt.Errorf("certificate not valid")
		zap.S().Errorf("%s", err)
		return nil, err
	}

	return state, nil
}

func validateCertificate(canisterId principal.Principal, status *CertifiedStatus) (bool, error) {
	var cert map[string]any
	if err := cbor.Unmarshal(*status.Certificate, &cert); err != nil {
		fmt.Printf("Error: %+v \n", err)
	}

	node, err := certificate.DeserializeNode(cert["tree"].([]any))
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
	}

	// Get time from certificate
	timeLeaf := certificate.Lookup(certificate.LookupPath("time"), node)
	c := concat([]byte{'D', 'I', 'D', 'L', 0x00, 0x01, 0x7d}, timeLeaf)
	_, certificateTimeDecoded, err := idl.Decode(c)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
	}

	// Check 2: The diff between decoded time and local time is within 5s.
	certificateTimeNum := certificateTimeDecoded[0].(idl.Nat)
	certificateTimeSecs := big.NewInt(0).Div(certificateTimeNum.BigInt(), big.NewInt(1e9))
	certificateTimeNano := big.NewInt(0).Mod(certificateTimeNum.BigInt(), big.NewInt(1e9))

	zap.S().Debugf("Certificate Time Secs: %d", certificateTimeSecs.Int64())
	zap.S().Debugf("Certificate Time Nano: %d", certificateTimeNano.Int64())

	now := time.Now()
	certificateTime := time.Unix(certificateTimeSecs.Int64(), certificateTimeNano.Int64())
	zap.S().Debugf("Certificate Time Parsed: %s", certificateTime.UTC())
	zap.S().Debugf("Current Time: %s", now.UTC())

	diff := now.Sub(certificateTime)
	if diff > 5*time.Second {
		return false, fmt.Errorf("certificate is expired")
	} else {
		zap.S().Debugf("Certificate is not expired: diff %f seconds", diff.Seconds())
	}

	// Verify Certificate Data
	certifiedData := certificate.Lookup(certificate.LookupPath("canister", string(canisterId.Raw), "certified_data"), node)

	consumers, err := serializeAndHash(status.Consumers)
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	pendingConsumerReqs, err := serializeAndHash(status.PendingConsumerReqs)
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	pendingSecretReqs, err := serializeAndHash(status.PendingSecretReqs)
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	secrets, err := serializeAndHash(status.Secrets)
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	toHash := fmt.Sprintf("%s%s%s%s", consumers, pendingConsumerReqs, pendingSecretReqs, secrets)

	h := sha256.New()
	h.Write([]byte(toHash))
	bs := h.Sum(nil)

	return hex.EncodeToString(certifiedData) == hex.EncodeToString(bs), nil
}

func serializeAndHash(data interface{}) (string, error) {
	canonicalDataJSON, err := canonicaljson.Marshal(data)
	if err != nil {
		fmt.Printf("Error %s \n", err)
		return "", err
	}

	h := sha256.New()
	h.Write(canonicalDataJSON)
	bs := h.Sum(nil)

	return hex.EncodeToString(bs), nil
}

func concat(bs ...[]byte) []byte {
	var c []byte
	for _, b := range bs {
		c = append(c, b...)
	}
	return c
}
