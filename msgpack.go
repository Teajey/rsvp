//go:build rsvp_msgpack

package rsvp

import (
	"fmt"
	"io"

	"github.com/Teajey/rsvp/internal/dev"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

const SupportedMediaTypeMsgpack string = "application/vnd.msgpack"

func init() {
	mediaTypeToContentType[SupportedMediaTypeMsgpack] = "application/vnd.msgpack"

	extToProposalMap["msgpack"] = SupportedMediaTypeMsgpack

	mediaTypeExtensionHandlers = append(mediaTypeExtensionHandlers, msgpackHandler)

	extendedMediaTypes = append(extendedMediaTypes, SupportedMediaTypeMsgpack)
}

func msgpackHandler(mediaType string, w io.Writer, res *Body) (bool, error) {
	if mediaType != SupportedMediaTypeMsgpack {
		return false, nil
	}

	dev.Log("Rendering msgpack...")
	err := msgpack.NewEncoder(w).Encode(res.Data)
	if err != nil {
		return true, fmt.Errorf("rendering data as msgpack: %w", err)
	}

	return true, nil
}
