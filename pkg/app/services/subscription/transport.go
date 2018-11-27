// Code generated by genservices. DO NOT EDIT.
package subscription

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golangci/golangci-api/pkg/endpoint/apierrors"
	"github.com/golangci/golangci-api/pkg/transportutil"
	"github.com/pkg/errors"
)

func RegisterHandlers(svc Service, regCtx *transportutil.HandlerRegContext) {

	hList := httptransport.NewServer(
		makeListEndpoint(svc, regCtx.Log),
		decodeListRequest,
		encodeListResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAuthorizedRequestContext(regCtx.Log,
			regCtx.ErrTracker, regCtx.DB, regCtx.AuthSessFactory)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("GET").Path("/v1/orgs/{org_id}/subs").Handler(hList)

	hGet := httptransport.NewServer(
		makeGetEndpoint(svc, regCtx.Log),
		decodeGetRequest,
		encodeGetResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAuthorizedRequestContext(regCtx.Log,
			regCtx.ErrTracker, regCtx.DB, regCtx.AuthSessFactory)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("GET").Path("/v1/orgs/{org_id}/subs/{sub_id}").Handler(hGet)

	hCreate := httptransport.NewServer(
		makeCreateEndpoint(svc, regCtx.Log),
		decodeCreateRequest,
		encodeCreateResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAuthorizedRequestContext(regCtx.Log,
			regCtx.ErrTracker, regCtx.DB, regCtx.AuthSessFactory)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("POST").Path("/v1/orgs/{org_id}/subs").Handler(hCreate)

	hUpdate := httptransport.NewServer(
		makeUpdateEndpoint(svc, regCtx.Log),
		decodeUpdateRequest,
		encodeUpdateResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAuthorizedRequestContext(regCtx.Log,
			regCtx.ErrTracker, regCtx.DB, regCtx.AuthSessFactory)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("PUT").Path("/v1/orgs/{org_id}/subs/{sub_id}").Handler(hUpdate)

	hDelete := httptransport.NewServer(
		makeDeleteEndpoint(svc, regCtx.Log),
		decodeDeleteRequest,
		encodeDeleteResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAuthorizedRequestContext(regCtx.Log,
			regCtx.ErrTracker, regCtx.DB, regCtx.AuthSessFactory)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("DELETE").Path("/v1/orgs/{org_id}/subs/{sub_id}").Handler(hDelete)

	hEventCreate := httptransport.NewServer(
		makeEventCreateEndpoint(svc, regCtx.Log),
		decodeEventCreateRequest,
		encodeEventCreateResponse,
		httptransport.ServerBefore(transportutil.StoreHTTPRequestToContext),
		httptransport.ServerAfter(transportutil.FinalizeSession),

		httptransport.ServerBefore(transportutil.MakeStoreAnonymousRequestContext(
			regCtx.Log, regCtx.ErrTracker, regCtx.DB)),

		httptransport.ServerFinalizer(transportutil.FinalizeRequest),
		httptransport.ServerErrorEncoder(transportutil.EncodeError),
		httptransport.ServerErrorLogger(transportutil.AdaptErrorLogger(regCtx.Log)),
	)
	regCtx.Router.Methods("POST").Path("/v1/payments/{provider}/events").Handler(hEventCreate)

}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ListRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(ListResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		ListResponse
	}{
		ListResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeGetResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(GetResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		GetResponse
	}{
		GetResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(CreateResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		CreateResponse
	}{
		CreateResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UpdateRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeUpdateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(UpdateResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		UpdateResponse
	}{
		UpdateResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request DeleteRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeDeleteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(DeleteResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		DeleteResponse
	}{
		DeleteResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}

func decodeEventCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request EventCreateRequest
	if err := transportutil.DecodeRequest(&request, r); err != nil {
		return nil, errors.Wrap(err, "can't decode request")
	}

	return request, nil
}

func encodeEventCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	if err := transportutil.GetContextError(ctx); err != nil {
		wrappedResp := struct {
			Error *transportutil.Error
		}{
			Error: transportutil.MakeError(err),
		}
		w.WriteHeader(wrappedResp.Error.HTTPCode)
		return json.NewEncoder(w).Encode(wrappedResp)
	}

	resp := response.(EventCreateResponse)
	wrappedResp := struct {
		transportutil.ErrorResponse
		EventCreateResponse
	}{
		EventCreateResponse: resp,
	}

	if resp.err != nil {
		if apierrors.IsErrorLikeResult(resp.err) {
			return transportutil.HandleErrorLikeResult(ctx, w, resp.err)
		}

		terr := transportutil.MakeError(resp.err)
		wrappedResp.Error = terr
		w.WriteHeader(terr.HTTPCode)
	}

	return json.NewEncoder(w).Encode(wrappedResp)
}
