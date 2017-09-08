package packet

import (
	"bytes"
	"encoding/binary"
	"io"
)

type SubAck struct {
	ID      uint16
	RetCode []byte
}

var SubackReturnCodes = map[uint8]string{
	0:    "Success - Maximum QoS 0",
	1:    "Success - Maximum QoS 1",
	2:    "Success - Maximum QoS 2",
	0x80: "Failure",
}

func createSubAck(r io.Reader, remainingLen int, fixFlags byte) (interface{}, error) {
	// TODO: move to ReadPacket, read full expect for Publish

	buf := make([]byte, remainingLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	m := &SubAck{}
	m.ID = binary.BigEndian.Uint16(buf[:2])
	m.RetCode = buf[2:]
	return m, nil
}

func (p *SubAck) Write(w io.Writer) error {
	buf := bytes.NewBuffer(nil)
	remainingLength := 2 + len(p.RetCode)
	buf.WriteByte(CtrlTypeSUBACK << 4)
	buf.Write(encodeLength(remainingLength))
	buf.Write(encodeUint16(p.ID))
	buf.Write(p.RetCode)
	_, err := buf.WriteTo(w)
	return err
}
