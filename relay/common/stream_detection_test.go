package common

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/stretchr/testify/require"
)

func TestApplyUpstreamContentTypeStreamDetectionTargetsResponseFormat(t *testing.T) {
	tests := []struct {
		name                   string
		relayMode              int
		responseFormatSettings dto.ChannelResponseFormatSettings
		clientRequestedStream  bool
		wantStream             bool
	}{
		{
			name:                   "chat completions enabled non stream ignores upstream event stream",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			clientRequestedStream:  false,
			wantStream:             false,
		},
		{
			name:                   "responses enabled non stream ignores upstream event stream",
			relayMode:              relayconstant.RelayModeResponses,
			responseFormatSettings: enabledResponseFormatSettings(),
			clientRequestedStream:  false,
			wantStream:             false,
		},
		{
			name:                   "chat completions enabled client stream upgrades",
			relayMode:              relayconstant.RelayModeChatCompletions,
			responseFormatSettings: enabledResponseFormatSettings(),
			clientRequestedStream:  true,
			wantStream:             true,
		},
		{
			name:      "chat completions disabled upgrades",
			relayMode: relayconstant.RelayModeChatCompletions,
			responseFormatSettings: dto.ChannelResponseFormatSettings{
				Enabled: false,
				Mode:    dto.ChannelResponseFormatModeClientStream,
			},
			clientRequestedStream: false,
			wantStream:            true,
		},
		{
			name:      "chat completions unknown mode upgrades",
			relayMode: relayconstant.RelayModeChatCompletions,
			responseFormatSettings: dto.ChannelResponseFormatSettings{
				Enabled: true,
				Mode:    "upstream_content_type",
			},
			clientRequestedStream: false,
			wantStream:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := streamDetectionRelayInfo(tt.relayMode, tt.responseFormatSettings, tt.clientRequestedStream)

			ApplyUpstreamContentTypeStreamDetection(info, "text/event-stream")

			require.Equal(t, tt.wantStream, info.IsStream)
			require.Equal(t, tt.clientRequestedStream, info.ClientRequestedStream)
		})
	}
}

func TestApplyUpstreamContentTypeStreamDetectionKeepsNonTargetModes(t *testing.T) {
	tests := []struct {
		name      string
		relayMode int
	}{
		{name: "responses compact", relayMode: relayconstant.RelayModeResponsesCompact},
		{name: "completions", relayMode: relayconstant.RelayModeCompletions},
		{name: "gemini native", relayMode: relayconstant.RelayModeGemini},
		{name: "images generations", relayMode: relayconstant.RelayModeImagesGenerations},
		{name: "images edits", relayMode: relayconstant.RelayModeImagesEdits},
		{name: "messages or unknown", relayMode: relayconstant.RelayModeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := streamDetectionRelayInfo(tt.relayMode, enabledResponseFormatSettings(), false)

			ApplyUpstreamContentTypeStreamDetection(info, "text/event-stream")

			require.True(t, info.IsStream)
			require.False(t, info.ClientRequestedStream)
		})
	}
}

func TestApplyUpstreamContentTypeStreamDetectionParsesEventStreamMediaType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		wantStream  bool
	}{
		{
			name:        "plain event stream",
			contentType: "text/event-stream",
			wantStream:  true,
		},
		{
			name:        "case insensitive",
			contentType: "TEXT/EVENT-STREAM",
			wantStream:  true,
		},
		{
			name:        "trimmed with charset",
			contentType: " text/event-stream; charset=utf-8 ",
			wantStream:  true,
		},
		{
			name:        "json does not upgrade",
			contentType: "application/json; charset=utf-8",
			wantStream:  false,
		},
		{
			name:        "invalid media type does not upgrade",
			contentType: "text/event-stream; charset",
			wantStream:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := streamDetectionRelayInfo(relayconstant.RelayModeCompletions, enabledResponseFormatSettings(), false)

			ApplyUpstreamContentTypeStreamDetection(info, tt.contentType)

			require.Equal(t, tt.wantStream, info.IsStream)
		})
	}
}

func streamDetectionRelayInfo(relayMode int, responseFormatSettings dto.ChannelResponseFormatSettings, clientRequestedStream bool) *RelayInfo {
	return &RelayInfo{
		RelayMode:             relayMode,
		ClientRequestedStream: clientRequestedStream,
		ChannelMeta: &ChannelMeta{
			ChannelSetting: dto.ChannelSettings{
				ResponseFormat: responseFormatSettings,
			},
		},
	}
}

func enabledResponseFormatSettings() dto.ChannelResponseFormatSettings {
	return dto.ChannelResponseFormatSettings{
		Enabled: true,
		Mode:    dto.ChannelResponseFormatModeClientStream,
	}
}
