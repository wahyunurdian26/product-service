package kit

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wahyunurdian26/util/model"
	"github.com/wahyunurdian26/util/requestid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapGRPCCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.FailedPrecondition:
		return http.StatusBadRequest // Or 422 if preferred
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// Context wraps standard context with request and path variables
type Context struct {
	context.Context
	request      *http.Request
	pathVariable map[string]string
}

func (c Context) Request() *http.Request {
	return c.request
}

func (c Context) GetPathVariable(name string) string {
	return c.pathVariable[name]
}

func (c Context) BindJSON(v interface{}) error {
	return json.NewDecoder(c.request.Body).Decode(v)
}

// HandlerFunc is the signature for gateway handlers
type HandlerFunc func(Context) (interface{}, error)

// RawResponse tells the gateway to return the result without the standard wrapper
type RawResponse struct {
	Data interface{}
}

// Router wraps mux.Router to provide a cleaner API
type Router struct {
	*mux.Router
}

func NewRouter(r *mux.Router) *Router {
	return &Router{Router: r}
}

var (
	rgbb *regexp.Regexp
)

func getBBPattern() *regexp.Regexp {
	r, _ := regexp.Compile(`^(?:https?:\/\/)?(?:[^.]+\.)?bluebird\.id(?::\d{1,5})?(\/.*)?`)
	return r
}

// HealthCheckHandler returns http 200 for root path
func HealthCheckHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

// CORSHandler enables cross-origin resource sharing
func CORSHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			if rgbb == nil {
				rgbb = getBBPattern()
			}

			if rgbb.MatchString(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					headers := []string{"Content-Type", "Accept", "Authorization"}
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
					methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"}
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
					w.Header().Set("Access-Control-Max-Age", "86400")
					return
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}

// DefaultHTTPHandler specifies default http handler
func DefaultHTTPHandler(handler http.Handler) http.Handler {
	handler = handlers.CompressHandler(handler)
	handler = CORSHandler(handler)
	handler = HealthCheckHandler(handler)
	return handler
}

func (r *Router) Get(path string, h HandlerFunc) {
	r.HandleFunc(path, r.makeHTTPHandler(h)).Methods(http.MethodGet)
}

func (r *Router) Post(path string, h HandlerFunc) {
	r.HandleFunc(path, r.makeHTTPHandler(h)).Methods(http.MethodPost)
}

func (r *Router) makeHTTPHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Initialize context with Request ID and other metadata
		newCtx := requestid.MiddlewareRequestId(req.Context(), nil)

		ctx := Context{
			Context:      newCtx,
			request:      req,
			pathVariable: mux.Vars(req),
		}

		w.Header().Set("Content-Type", "application/json")
		resp, err := h(ctx)

		reqID := requestid.GetRequestId(ctx)

		if err != nil {
			statusCode := http.StatusInternalServerError
			message := err.Error()

			// Extract gRPC status if available
			if st, ok := status.FromError(err); ok {
				statusCode = mapGRPCCodeToHTTP(st.Code())
				message = st.Message()
			}

			w.WriteHeader(statusCode)
			errorResp := model.Response{
				Code:      statusCode,
				Message:   message,
				RequestID: reqID,
			}
			json.NewEncoder(w).Encode(errorResp)
			return
		}

		if raw, ok := resp.(RawResponse); ok {
			json.NewEncoder(w).Encode(raw.Data)
			return
		}

		w.WriteHeader(http.StatusOK)
		finalResp := model.Response{
			Code:      http.StatusOK,
			Message:   "Success",
			RequestID: reqID,
			Result:    resp,
		}
		json.NewEncoder(w).Encode(finalResp)
	}
}
