package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/certificate"
	"github.com/fxamacker/cbor/v2"
	"github.com/gibson042/canonicaljson-go"
	"icp_vault_client/backend"
	"net/url"
	"os"

	"github.com/aviate-labs/agent-go/principal"
)

func main() {
	ledgerID, _ := principal.Decode("d6g4o-amaaa-aaaaa-qaaoq-cai")
	localReplica, _ := url.Parse("http://127.0.0.1:4943")

	ccfg := &agent.ClientConfig{Host: localReplica}
	a, err := backend.NewAgent(ledgerID, agent.Config{ClientConfig: ccfg, FetchRootKey: true})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	res, err := a.GetCertifiedStatus()
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		return
	}

	if res.Certificate != nil {
		var cert map[string]any
		if err := cbor.Unmarshal(*res.Certificate, &cert); err != nil {
			fmt.Printf("Error: %+v \n", err)
		}

		node, err := certificate.DeserializeNode(cert["tree"].([]any))
		if err != nil {
			fmt.Printf("Error: %+v \n", err)
		}

		// Get time from certificate
		time := certificate.Lookup(certificate.LookupPath("time"), node)
		c := concat([]byte{'D', 'I', 'D', 'L', 0x00, 0x01, 0x7d}, time)
		_, certificateTime, err := idl.Decode(c)
		if err != nil {
			fmt.Printf("Error: %+v \n", err)
		}
		fmt.Println("Certificate Time:", certificateTime[0])

		// Verify Certificate Data
		certifiedData := certificate.Lookup(certificate.LookupPath("canister", string(ledgerID.Raw), "certified_data"), node)

		// https://github.com/aviate-labs/agent-go/issues/9
		// FIXME this is required as icp client is setting emtpy arrays incorrectly, and they are marshalled to nil instead of []
		if len(res.Consumers) == 0 {
			res.Consumers = make([]backend.Consumer, 0)
		}
		consumers, err := serializeAndHash(res.Consumers)
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

		// https://github.com/aviate-labs/agent-go/issues/9
		// FIXME this is required as icp client is setting emtpy arrays incorrectly, and they are marshalled to nil instead of []
		if len(res.PendingConsumerReqs) == 0 {
			res.PendingConsumerReqs = make([]backend.Consumer, 0)
		}
		pendingConsumerReqs, err := serializeAndHash(res.PendingConsumerReqs)
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

		// https://github.com/aviate-labs/agent-go/issues/9
		// FIXME this is required as icp client is setting emtpy arrays incorrectly, and they are marshalled to nil instead of []
		if len(res.PendingSecretReqs) == 0 {
			res.PendingSecretReqs = make([]backend.Secret, 0)
		}
		pendingSecretReqs, err := serializeAndHash(res.PendingSecretReqs)
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

		// https://github.com/aviate-labs/agent-go/issues/9
		// FIXME this is required as icp client is setting emtpy arrays incorrectly, and they are marshalled to nil instead of []
		if len(res.Secrets) == 0 {
			res.Secrets = make([]backend.Secret, 0)
		}
		secrets, err := serializeAndHash(res.Secrets)
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

		toHash := fmt.Sprintf("%s%s%s%s", consumers, pendingConsumerReqs, pendingSecretReqs, secrets)

		h := sha256.New()
		h.Write([]byte(toHash))
		bs := h.Sum(nil)

		fmt.Println("Is certified data valid:", hex.EncodeToString(certifiedData) == hex.EncodeToString(bs))
	}
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
