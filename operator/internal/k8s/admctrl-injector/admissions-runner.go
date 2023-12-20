package admctrl_injector

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"github.com/zondax/vault-k8s-canister/operator/common"
	"github.com/zondax/vault-k8s-canister/operator/internal/conf"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type AdmCtrlInjector struct {
	name   string
	config *conf.Config
}

func NewAdmCtrlInjector(config *conf.Config) *AdmCtrlInjector {
	kubeconfig := ""
	if config.Dev != nil {
		kubeconfig = config.Dev.Kubeconfig
	}
	common.CreateCommonKubernetesClient(kubeconfig)
	return &AdmCtrlInjector{
		name:   "admctrl-injector",
		config: config,
	}
}

func (a AdmCtrlInjector) Name() string {
	return a.name
}

func (a AdmCtrlInjector) Start() error {
	// Chi
	r := chi.NewRouter()

	r.Post("/", mutate)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello!"))
		if err != nil {
			zap.S().Warnf("[ADM CTRL] Failed to write response: %v", err)
			return
		}
	})

	port := "8282"
	if a.config.Dev != nil {
		port = a.config.Dev.AdmController.Port
	}

	var tlsConfig *tls.Config
	var err error
	admControllerCert := []byte(os.Getenv("ADM_CONTROLLER_CERT"))
	admControllerKey := []byte(os.Getenv("ADM_CONTROLLER_KEY"))

	admControllerCertBase64 := os.Getenv("ADM_CONTROLLER_CERT_BASE64")
	admControllerKeyBase64 := os.Getenv("ADM_CONTROLLER_KEY_BASE64")

	if admControllerCertBase64 != "" {
		admControllerCert, err = base64.StdEncoding.DecodeString(admControllerCertBase64)
		if err != nil {
			zap.S().Fatalf("[ADM CTRL] Failed to read cert file from base64: %v", err)
			return err
		}
	}

	if admControllerKeyBase64 != "" {
		admControllerKey, err = base64.StdEncoding.DecodeString(admControllerKeyBase64)
		if err != nil {
			zap.S().Fatalf("[ADM CTRL] Failed to read key file from base64: %v", err)
			return err
		}
	}

	if len(admControllerKey) > 0 && len(admControllerCert) > 0 {
		certificates, err := tls.X509KeyPair(admControllerCert, admControllerKey)
		if err != nil {
			zap.S().Fatalf("[ADM CTRL] Failed to read cert and key files: %v", err)
			return err
		}
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{certificates}, MinVersion: tls.VersionTLS12}
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         tlsConfig,
	}

	if tlsConfig != nil {
		zap.S().Infof("[ADM CTRL] Admission Controller starting as HTTPS on port %s", port)
		err = server.ListenAndServeTLS("", "")
	} else {
		zap.S().Infof("[ADM CTRL] Admission Controller starting as HTTP on port: %s", port)
		err = server.ListenAndServe()
	}

	if err != nil {
		zap.S().Errorf("[ADB_CRTL] Prometheus server error: %v", err)
	} else {
		zap.S().Infof("[ADB_CRTL] Prometheus server serving at port %s", port)
	}
	return err
}

func (a AdmCtrlInjector) Stop() error {
	// TODO implement me
	return nil
}
