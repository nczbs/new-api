package common

import (
	"mime"
	"strings"

	relayconstant "github.com/QuantumNous/new-api/relay/constant"
)

const upstreamEventStreamMediaType = "text/event-stream"

func ApplyUpstreamContentTypeStreamDetection(info *RelayInfo, contentType string) {
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil || !strings.EqualFold(mediaType, upstreamEventStreamMediaType) {
		return
	}

	if shouldIgnoreUpstreamStreamContentType(info) {
		return
	}

	if info != nil {
		info.IsStream = true
	}
}

func shouldIgnoreUpstreamStreamContentType(info *RelayInfo) bool {
	if info == nil || info.ChannelMeta == nil {
		return false
	}
	if !info.ChannelMeta.ChannelSetting.ResponseFormat.IsEnabled() {
		return false
	}
	if info.ClientRequestedStream {
		return false
	}

	return info.RelayMode == relayconstant.RelayModeChatCompletions ||
		info.RelayMode == relayconstant.RelayModeResponses
}
