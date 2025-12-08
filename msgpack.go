//go:build rsvp_msgpack

package rsvp

import (
	"fmt"
	"net/http"

	"github.com/Teajey/rsvp/internal/log"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

const mMsgpack supportedType = "application/vnd.msgpack"

func init() {
	mediaTypeToContentType[mMsgpack] = "application/vnd.msgpack"

	defaultExtToProposalMap["msgpack"] = string(mMsgpack)

	mediaTypeExtensionHandlers = append(mediaTypeExtensionHandlers, msgpackHandler)

	extendedMediaTypes = append(extendedMediaTypes, mMsgpack)
}

func msgpackHandler(mediaType supportedType, w http.ResponseWriter, res *Rsvp) (bool, error) {
	if mediaType != mMsgpack {
		return false, nil
	}

	log.Dev("Rendering msgpack...")
	err := msgpack.NewEncoder(w).Encode(res.Body)
	if err != nil {
		return true, fmt.Errorf("failed to render body as msgpack: %w", err)
	}

	return true, nil
}
