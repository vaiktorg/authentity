package gwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

var Spice encData

// include "_ github.com/vaiktorg/gwt" to initialize the default global Spice encryption additives
func init() {
	initEncData()
}

func initEncData() {
	Spice.Salt = []byte(uuid.New().String())
	Spice.Salt = []byte(uuid.New().String())
}

type (
	encData struct {
		Salt   []byte
		Pepper []byte
	}

	GWT struct {
		repo *gorm.DB
	}

	Header struct {
		ID        string
		Issuer    string
		Timestamp time.Time
	}
	Payload struct {
		Data interface{}
	}

	Token struct {
		Header    Header
		Payload   Payload
		Signature string // Token Signature
		hash      string // Whole token.
		salt      string
	}
)

const (
	TokenName = "gwt"
)

const (
	ErrorFailedToDecode    = "failed to decode token signature"
	ErrorSignatureNotMatch = "signatures do not match"
)

func NewToken(id string, issuer string, salt, pepper []byte, payload interface{}) (tkn, sig string, _ error) {
	return EncodeGWT(&Header{
		ID:        id,
		Issuer:    issuer,
		Timestamp: time.Now(),
	}, &Payload{Data: payload},
		salt, pepper,
	)
}

func EncodeGWT(header *Header, payload *Payload, salt, pepper []byte) (token, signature string, err error) {
	if header == nil || payload == nil {
		return "", "", errors.New("arguments payload or header are nil")
	}

	// Payload JSON string
	payloadBuffer := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuffer).Encode(payload)
	if err != nil {
		return "", "", err
	}

	// Header JSON string
	headerBuffer := new(bytes.Buffer)
	err = json.NewEncoder(headerBuffer).Encode(header)
	if err != nil {
		return "", "", err
	}

	// ----------------------------------------------------------------------------------------------
	// Signature Encoding
	hash := hmac.New(sha256.New, salt)
	_, err = hash.Write(payloadBuffer.Bytes())
	if err != nil {
		return "", "", err
	}

	hashSignature := hash.Sum(pepper)

	b64signature := make([]byte, base64.URLEncoding.EncodedLen(len(hashSignature)))
	base64.URLEncoding.Encode(b64signature, hashSignature)

	// ----------------------------------------------------------------------------------------------
	// Payload Encoding
	b64payload := make([]byte, base64.URLEncoding.EncodedLen(payloadBuffer.Len()))
	base64.URLEncoding.Encode(b64payload, payloadBuffer.Bytes())

	//----------------------------------------------------------------------------------------------
	// Header Encoding
	b64header := make([]byte, base64.URLEncoding.EncodedLen(headerBuffer.Len()))
	base64.URLEncoding.Encode(b64header, headerBuffer.Bytes())
	//-----------

	return strings.Join([]string{string(b64header), string(b64payload), string(b64signature)}, "."), string(b64signature), nil
}

func DecodeGWT(tkn string) (gwt *Token, err error) {
	if tkn == "" {
		return nil, errors.New("no token provided")
	}

	// Split Message
	gwtParts := strings.Split(tkn, ".")
	if len(gwtParts) > 3 || len(gwtParts) < 3 {
		return nil, errors.New("parts of token are missing")
	}

	// ------------------------------------------------------------------------------------------------
	// B64 Decoding
	// Signature
	signatureBuff, err := base64.URLEncoding.DecodeString(gwtParts[2])
	if err != nil {
		return nil, err
	}

	// Payload
	payloadBuff, err := base64.URLEncoding.DecodeString(gwtParts[1])
	if err != nil {
		return nil, err
	}

	// Header
	headerBuff, err := base64.URLEncoding.DecodeString(gwtParts[0])
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------------------------------------------
	// Decode data
	header := new(Header)
	err = json.NewDecoder(bytes.NewReader(headerBuff)).Decode(header)
	if err != nil {
		return nil, errors.New(ErrorFailedToDecode)
	}

	payload := new(Payload)
	err = json.NewDecoder(bytes.NewReader(payloadBuff)).Decode(payload)
	if err != nil {
		return nil, errors.New(ErrorFailedToDecode)
	}

	// ------------------------------------------------------------------------------------------------
	// Signature Validation

	//Generate gwt Signature from decoded payload
	hash := hmac.New(sha256.New, Spice.Salt)
	_, err = hash.Write(payloadBuff)
	if err != nil {
		return nil, errors.New(ErrorFailedToDecode)
	}

	hashSignature := hash.Sum(Spice.Pepper)

	// Validate
	if !bytes.Equal(hashSignature, signatureBuff) {
		return nil, errors.New(ErrorSignatureNotMatch)
	}

	gwt = &Token{
		Header:    *header,
		Payload:   *payload,
		Signature: gwtParts[2],
	}

	return gwt, nil
}
