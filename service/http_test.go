package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestIOCopyBytesGracefullyNormalizesContentTypeWhenFlagged(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	common.SetContextKey(c, constant.ContextKeyNormalizeResponseContentType, true)

	src := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   {"text/event-stream; charset=utf-8"},
			"Content-Length": {"999"},
			"X-Custom":       {"kept"},
		},
	}
	body := []byte(`{"ok":true}`)

	IOCopyBytesGracefully(c, src, body)

	result := recorder.Result()
	defer result.Body.Close()
	writtenBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))
	require.Equal(t, "11", result.Header.Get("Content-Length"))
	require.Equal(t, "kept", result.Header.Get("X-Custom"))
	require.Equal(t, body, writtenBody)
}

func TestIOCopyBytesGracefullyPreservesUpstreamHeadersWhenNotFlagged(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	src := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": {"text/event-stream; charset=utf-8"},
		},
	}

	IOCopyBytesGracefully(c, src, []byte(`{}`))

	require.Equal(t, "text/event-stream; charset=utf-8", recorder.Result().Header.Get("Content-Type"))
}

func TestIOCopyBytesGracefullyRequestIdHandlingUnchanged(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	src := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			common.RequestIdKey: {"upstream-request-id"},
		},
	}

	IOCopyBytesGracefully(c, src, []byte(`{}`))

	_, copied := recorder.Result().Header[common.RequestIdKey]
	require.False(t, copied)
	require.Equal(t, "upstream-request-id", c.GetString(common.UpstreamRequestIdKey))
}
