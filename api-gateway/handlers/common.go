package handlers

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcToHTTP(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "internal server error"
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()
	case codes.NotFound:
		return http.StatusNotFound, st.Message()
	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()
	case codes.Unauthenticated:
		return http.StatusUnauthorized, st.Message()
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "upstream timeout"
	default:
		return http.StatusInternalServerError, st.Message()
	}
}
