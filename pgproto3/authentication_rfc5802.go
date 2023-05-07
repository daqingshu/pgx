package pgproto3

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/internal/pgio"
)

// AuthenticationSASL is a message sent from the backend indicating that SASL authentication is required.
type AuthenticationRFC5802 struct {
	PwdMethod       int32
	Rnd64           [64]byte
	ServerSignature [64]byte
	Token           [8]byte
	Iterstep        int
	Data            [4]byte
}

// Backend identifies this message as sendable by the PostgreSQL backend.
func (*AuthenticationRFC5802) Backend() {}

// Backend identifies this message as an authentication response.
func (*AuthenticationRFC5802) AuthenticationResponse() {}

// Decode decodes src into dst. src must contain the complete message with the exception of the initial 1 byte message
// type identifier and 4 byte message length.
func (dst *AuthenticationRFC5802) Decode(src []byte) error {
	if len(src) < 4 {
		return errors.New("authentication message too short")
	}

	authType := binary.BigEndian.Uint32(src)

	if authType != AuthTypeSASL {
		return errors.New("bad auth type")
	}

	dst.PwdMethod = int32(binary.BigEndian.Uint32(src[4:]))

	if dst.PwdMethod == 0 || dst.PwdMethod == 2 {
		copy(dst.Rnd64[:], src[8:72])
		copy(dst.Token[:], src[72:80])
		dst.Iterstep = int(int32(binary.BigEndian.Uint32(src[80:])))
	} else {
		copy(dst.Data[:], src[8:12])
	}

	return nil
}

// Encode encodes src into dst. dst will include the 1 byte message type identifier and the 4 byte message length.
func (src *AuthenticationRFC5802) Encode(dst []byte) []byte {
	dst = append(dst, 'R')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)
	dst = pgio.AppendUint32(dst, AuthTypeSASL)
	dst = pgio.AppendUint32(dst, uint32(src.PwdMethod))
	if src.PwdMethod == 0 || src.PwdMethod == 2 {
		dst = append(dst, src.Rnd64[:]...)
		dst = append(dst, src.Token[:]...)
		dst = pgio.AppendUint32(dst, uint32(src.Iterstep))
	} else {
		dst = append(dst, src.Data[:]...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

// MarshalJSON implements encoding/json.Marshaler.
func (src AuthenticationRFC5802) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type            string
		PwdMethod       int32
		Rnd64           string
		ServerSignature string
		Token           string
		Iterstep        int
		Data            string
	}{
		Type:            "AuthenticationRFC5802",
		PwdMethod:       src.PwdMethod,
		Rnd64:           string(src.Rnd64[:]),
		ServerSignature: "",
		Token:           string(src.Token[:]),
		Iterstep:        src.Iterstep,
		Data:            string(src.Data[:]),
	})
}

func (src AuthenticationRFC5802) String() string {
	b, _ := src.MarshalJSON()
	return string(b)
}
