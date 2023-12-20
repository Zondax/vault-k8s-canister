package admctrl_injector

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// mutate is an HTTP handler for the admission controller. It processes incoming admission requests,
// performs mutations on the Pod object, and returns an admission response with the required patches.
func mutate(w http.ResponseWriter, r *http.Request) {
	// Read the request body into a byte slice
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer func() {
		// Close the request body, handling any error if it occurs
		if closeErr := r.Body.Close(); closeErr != nil {
			// You can log the error here, or handle it as required.
			zap.S().Errorf("[ADM_CTRL] Failed to close request body: %v\n", closeErr)
		}
	}()

	// Log receipt of the admission request
	zap.S().Infof("[ADM CTRL] Admission Controller received request")

	// Initialize an empty AdmissionReview object
	admissionReview := v1beta1.AdmissionReview{}

	// Unmarshal the request body into the AdmissionReview object
	err = json.Unmarshal(body, &admissionReview)
	if err != nil {
		http.Error(w, "[ADM CTRL] Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Initialize an empty Pod object
	pod := corev1.Pod{}

	// Unmarshal the raw object from the AdmissionReview into the Pod object
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
	if err != nil {
		http.Error(w, "[ADM CTRL] Failed to decode raw object", http.StatusBadRequest)
		return
	}

	// Generate patches based on the Pod object
	patch := getPatchObject(context.TODO(), pod)

	// Marshal the patches into JSON format
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		http.Error(w, "[ADM CTRL] Failed to encode patch", http.StatusInternalServerError)
		return
	}

	// Create an AdmissionResponse with the required information
	admissionResponse := v1beta1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: func() *v1beta1.PatchType { pt := v1beta1.PatchTypeJSONPatch; return &pt }(),
	}

	// Set the response in the AdmissionReview object
	admissionReview.Response = &admissionResponse

	// Marshal the AdmissionReview response into JSON format
	responseBytes, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, "[ADM CTRL] Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Set the content type header and write the response
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBytes)
	if err != nil {
		zap.S().Error(err)
	}
}
