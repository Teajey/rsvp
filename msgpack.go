//go:build rsvp_msgpack

package rsvp

import (
	"fmt"
	"net/http"

	"github.com/Teajey/rsvp/internal/log"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

const SupportedMediaTypeMsgpack string = "application/vnd.msgpack"

func init() {
	mediaTypeToContentType[SupportedMediaTypeMsgpack] = "application/vnd.msgpack"

	extToProposalMap["msgpack"] = SupportedMediaTypeMsgpack

	mediaTypeExtensionHandlers = append(mediaTypeExtensionHandlers, msgpackHandler)

	extendedMediaTypes = append(extendedMediaTypes, SupportedMediaTypeMsgpack)
}

func msgpackHandler(mediaType string, w http.ResponseWriter, res *Response) (bool, error) {
	if mediaType != SupportedMediaTypeMsgpack {
		return false, nil
	}

	log.Dev("Rendering msgpack...")
	err := msgpack.NewEncoder(w).Encode(res.Body)
	if err != nil {
		return true, fmt.Errorf("failed to render body as msgpack: %w", err)
	}

	return true, nil
}
