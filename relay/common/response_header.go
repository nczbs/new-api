package common

import (
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"

	"github.com/gin-gonic/gin"
)

func ApplyResponseContentTypeNormalizationFlag(c *gin.Context, info *RelayInfo, upstreamStatusCode int) {
	if c == nil {
		return
	}
	if shouldNormalizeResponseContentType(info, upstreamStatusCode) {
		common.SetContextKey(c, constant.ContextKeyNormalizeResponseContentType, true)
	}
}

func shouldNormalizeResponseContentType(info *RelayInfo, upstreamStatusCode int) bool {
	if info == nil || info.ChannelMeta == nil {
		return false
	}
	if upstreamStatusCode != http.StatusOK {
		return false
	}
	if info.ClientRequestedStream || info.IsStream {
		return false
	}
	if !info.ChannelMeta.ChannelSetting.ResponseFormat.IsEnabled() {
		return false
	}

	return info.RelayMode == relayconstant.RelayModeChatCompletions ||
		info.RelayMode == relayconstant.RelayModeResponses
}
