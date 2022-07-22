package main

import (
	"crypto/rand"
	"encoding/binary"
	"errors"

	"github.com/shimin/pow/internal/pow"
	"github.com/shimin/pow/internal/proto"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.SugaredLogger
}

func NewHandler(log *zap.SugaredLogger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) AuthFlow(stream proto.AuthService_AuthFlowServer) error {
	data := make([]byte, keySize)
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

	ok := pow.Validate(data, targetBits, binary.LittleEndian.Uint64(solution))
	if ok {
		h.log.Infof("Validation passed")
		stream.Send(&proto.Packet{
			Data: []byte("Here is your word of wisdom"),
		})
		return nil
	}

	h.log.Infof("Validation failed")
	stream.Send(&proto.Packet{
		Data: []byte("Access denied"),
	})
	return nil
}
