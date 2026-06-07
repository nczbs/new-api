package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRelayInfoGetFinalRequestRelayFormatPrefersExplicitFinal(t *testing.T) {
	info := &RelayInfo{
		RelayFormat:             types.RelayFormatOpenAI,
		RequestConversionChain:  []types.RelayFormat{types.RelayFormatOpenAI, types.RelayFormatClaude},
		FinalRequestRelayFormat: types.RelayFormatOpenAIResponses,
	}

	require.Equal(t, types.RelayFormat(types.RelayFormatOpenAIResponses), info.GetFinalRequestRelayFormat())
}

func TestRelayInfoGetFinalRequestRelayFormatFallsBackToConversionChain(t *testing.T) {
	info := &RelayInfo{
		RelayFormat:            types.RelayFormatOpenAI,
		RequestConversionChain: []types.RelayFormat{types.RelayFormatOpenAI, types.RelayFormatClaude},
	}

	require.Equal(t, types.RelayFormat(types.RelayFormatClaude), info.GetFinalRequestRelayFormat())
}

func TestRelayInfoGetFinalRequestRelayFormatFallsBackToRelayFormat(t *testing.T) {
	info := &RelayInfo{
		RelayFormat: types.RelayFormatGemini,
	}

	require.Equal(t, types.RelayFormat(types.RelayFormatGemini), info.GetFinalRequestRelayFormat())
}

func TestRelayInfoGetFinalRequestRelayFormatNilReceiver(t *testing.T) {
	var info *RelayInfo
	require.Equal(t, types.RelayFormat(""), info.GetFinalRequestRelayFormat())
}

func TestGenBaseRelayInfoInitializesClientRequestedStream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		request dto.Request
		want    bool
	}{
		{
			name: "nil request",
			want: false,
		},
		{
			name:    "missing stream",
			request: &dto.GeneralOpenAIRequest{},
			want:    false,
		},
		{
			name: "explicit non stream",
			request: &dto.GeneralOpenAIRequest{
				Stream: boolPtr(false),
			},
			want: false,
		},
		{
			name: "explicit stream",
			request: &dto.GeneralOpenAIRequest{
				Stream: boolPtr(true),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRelayInfoTestContext()

			info := genBaseRelayInfo(c, tt.request)

			require.Equal(t, tt.want, info.IsStream)
			require.Equal(t, tt.want, info.ClientRequestedStream)
			require.Equal(t, tt.want, c.GetBool(string(constant.ContextKeyIsStream)))
		})
	}
}

func TestRelayInfoClientRequestedStreamRemainsStableWhenIsStreamChanges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newRelayInfoTestContext()
	info := genBaseRelayInfo(c, &dto.GeneralOpenAIRequest{
		Stream: boolPtr(false),
	})

	info.IsStream = true

	require.True(t, info.IsStream)
	require.False(t, info.ClientRequestedStream)
}

func newRelayInfoTestContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	return c
}

func boolPtr(v bool) *bool {
	return &v
}
