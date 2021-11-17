package common

import "sync/atomic"

type SerialIdService struct {
	nextId uint32
}

func (s *SerialIdService) AllocId() uint32 {
	return atomic.AddUint32(&s.nextId, 1)
}

func NewSerialIdService() *SerialIdService {
	return &SerialIdService{
		nextId: 1,
	}
}


type SerialIdService64 struct {
    nextId uint64
}

func (s *SerialIdService64) AllocId() uint64 {
    return atomic.AddUint64(&s.nextId, 1)
}

func NewSerialIdService64() *SerialIdService64 {
    return &SerialIdService64{
        nextId: 1,
    }
}
