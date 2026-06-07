package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestApplyResponseContentTypeNormalizationFlag(t *testing.T) {
	tests := []struct {
		name                   string
		relayMode              int
		responseFormatSettings dto.ChannelResponseFormatSettings
		clientRequestedStream  bool
		isStream               bool
		statusCode             int
		wantFlag               bool
	}{
		{
			name:                   "chat completions enabled non stream ok",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               true,
		},
		{
			name:                   "responses enabled non stream ok",
			relayMode:              relayconstant.RelayModeResponses,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               true,
		},
		{
			name:                   "disabled response format",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: dto.ChannelResponseFormatSettings{Enabled: false, Mode: dto.ChannelResponseFormatModeClientStream},
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "non 200",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusBadRequest,
			wantFlag:               false,
		},
		{
			name:                   "client stream",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			clientRequestedStream:  true,
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "upstream stream",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			isStream:               true,
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "responses compact",
			relayMode:              relayconstant.RelayModeResponsesCompact,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "gemini native",
			relayMode:              relayconstant.RelayModeGemini,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "images",
			relayMode:              relayconstant.RelayModeImagesGenerations,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
		{
			name:                   "messages or unknown",
			relayMode:              relayconstant.RelayModeUnknown,
			responseFormatSettings: enabledResponseFormatSettings(),
			statusCode:             http.StatusOK,
			wantFlag:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			info := streamDetectionRelayInfo(tt.relayMode, tt.responseFormatSettings, tt.clientRequestedStream)
			info.IsStream = tt.isStream

			ApplyResponseContentTypeNormalizationFlag(c, info, tt.statusCode)

			require.Equal(t, tt.wantFlag, common.GetContextKeyBool(c, constant.ContextKeyNormalizeResponseContentType))
		})
	}
}
