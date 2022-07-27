package server

import (
	"crypto/rand"
	"encoding/binary"
	"errors"

	"github.com/shimin/pow/internal/pow"
	"github.com/shimin/pow/internal/wisdom"
	"github.com/shimin/pow/proto"
	"go.uber.org/zap"
)

type Handler struct {
	log        *zap.SugaredLogger
	keySize    uint16
	targetBits uint16
	quotes     *wisdom.Set
}

func NewHandler(log *zap.SugaredLogger, keySize, targetBits uint16, quotes *wisdom.Set) *Handler {
	return &Handler{
		log:        log,
		keySize:    keySize,
		targetBits: targetBits,
		quotes:     quotes,
	}
}

func (h *Handler) AuthFlow(stream proto.AuthService_AuthFlowServer) error {
	data := make([]byte, h.keySize)
	rand.Read(data)

	stream.Send(&proto.Packet{
		Data: data,
	})

	answer, err := stream.Recv()
	if err != nil {
		h.log.Error(err)
		return err
	}
	solution := answer.GetData()
	if len(solution) != 8 {
		h.log.Error("answer size is wrong")
		return errors.New("answer size is wrong")
	}

	ok := pow.Validate(data, h.targetBits, binary.LittleEndian.Uint64(solution))
	if ok {
		h.log.Infof("Validation passed")
		quote := h.quotes.GetRandQuote()
		stream.Send(&proto.Packet{
			Data: []byte(quote),
		})
		return nil
	}

	h.log.Infof("Validation failed")
	stream.Send(&proto.Packet{
		Data: []byte("Access denied"),
	})
	return nil
}
