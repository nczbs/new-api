package dto

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/QuantumNous/new-api/common"
)

func TestChannelResponseFormatSettingsIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		settings ChannelResponseFormatSettings
		want     bool
	}{
		{
			name:     "missing config is disabled",
			settings: ChannelResponseFormatSettings{},
			want:     false,
		},
		{
			name: "disabled config is disabled",
			settings: ChannelResponseFormatSettings{
				Enabled: false,
				Mode:    ChannelResponseFormatModeClientStream,
			},
			want: false,
		},
		{
			name: "enabled with empty mode is enabled",
			settings: ChannelResponseFormatSettings{
				Enabled: true,
			},
			want: true,
		},
		{
			name: "enabled with client stream mode is enabled",
			settings: ChannelResponseFormatSettings{
				Enabled: true,
				Mode:    ChannelResponseFormatModeClientStream,
			},
			want: true,
		},
		{
			name: "enabled with unknown mode is disabled",
			settings: ChannelResponseFormatSettings{
				Enabled: true,
				Mode:    "upstream_content_type",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.settings.IsEnabled())
		})
	}
}

func TestChannelResponseFormatRulesRoundTripUnknownFields(t *testing.T) {
	raw := []byte(`{
		"response_format": {
			"enabled": true,
			"mode": "client_stream",
			"rules": [
				{
					"name": "future-rule",
					"priority": 7,
					"nested": {
						"content_type": "text/event-stream",
						"preserve": true
					}
				}
			]
		}
	}`)

	var settings ChannelSettings
	require.NoError(t, common.Unmarshal(raw, &settings))

	require.True(t, settings.ResponseFormat.IsEnabled())
	require.Len(t, settings.ResponseFormat.Rules, 1)
	require.Equal(t, "future-rule", settings.ResponseFormat.Rules[0]["name"])
	require.Equal(t, float64(7), settings.ResponseFormat.Rules[0]["priority"])
	require.Equal(t, map[string]any{
		"content_type": "text/event-stream",
		"preserve":     true,
	}, settings.ResponseFormat.Rules[0]["nested"])

	encoded, err := common.Marshal(settings)
	require.NoError(t, err)

	var roundTripped map[string]any
	require.NoError(t, common.Unmarshal(encoded, &roundTripped))

	responseFormat, ok := roundTripped["response_format"].(map[string]any)
	require.True(t, ok)
	rules, ok := responseFormat["rules"].([]any)
	require.True(t, ok)
	require.Len(t, rules, 1)
	rule, ok := rules[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "future-rule", rule["name"])
	require.Equal(t, float64(7), rule["priority"])
	require.Equal(t, map[string]any{
		"content_type": "text/event-stream",
		"preserve":     true,
	}, rule["nested"])
}
