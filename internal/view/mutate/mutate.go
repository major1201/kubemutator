package mutate

import (
	"context"
	"encoding/json"
	"github.com/major1201/k8s-mutator/pkg/httputils"
	"github.com/major1201/k8s-mutator/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch
)

// ServeMutate serve the /mutate path
func ServeMutate(w http.ResponseWriter, r *http.Request) {
	// set logger
	r = r.WithContext(context.WithValue(r.Context(), httputils.CtxLogger, zap.L().Named("mutator").With(zap.String("requestId", httputils.RequestID(r)))))

	// serve
	serve(w, r, mutatePods)
}

func getLogger(r *http.Request) *zap.Logger {
	ctx := r.Context()
	if ctx != nil {
		return ctx.Value(httputils.CtxLogger).(*zap.Logger)
	}
	return nil
}

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(*http.Request, v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		err = errors.Wrap(err, "deserializer decode error")
		getLogger(r).Error("deserializer decoding error", log.Error(err))
		responseAdmissionReview.Response = toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(r, requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	response := responseAdmissionReview.Response
	getLogger(r).Info("sending response",
		zap.String("uid", string(response.UID)),
		zap.Bool("allowed", response.Allowed),
		zap.Any("auditAnnotations", response.AuditAnnotations),
		zap.String("patch", string(response.Patch)),
		zap.Any("result", response.Result),
	)

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		err = errors.Wrap(err, "json marshal error")
		getLogger(r).Error("json marshal error", log.Error(err))
	}
	if _, err := w.Write(respBytes); err != nil {
		err = errors.Wrap(err, "write response error")
		getLogger(r).Error("write response error", zap.Error(err))
	}
}
