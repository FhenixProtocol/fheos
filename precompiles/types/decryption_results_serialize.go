package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"
)

// GetSerializedDecryptionResult returns a byte-serialization of a decryption result.
func (dr *DecryptionResults) GetSerializedDecryptionResult(key PendingDecryption) ([]byte, error) {
	// The structure the encoded message is:
	//	encoded_result = key | result_req OR result_seal OR result_decrypt
	// see DecryptionRecord.Serialize for the format of each result type
	result, ok := dr.Get(key)
	if !ok {
		return nil, errors.New("tried to serialize result of unknown decryption")
	}

	if result.Value == nil {
		return nil, errors.New("tried to serialize result of decryption which is still pending")
	}

	serializedKey, err := key.Serialize()
	if err != nil {
		return nil, err
	}
	serializedResult, err := result.Serialize(key.Type)
	if err != nil {
		return nil, err
	}

	return append(serializedKey, serializedResult...), nil
}

func (dr *DecryptionResults) LoadResolvedDecryption(reader io.Reader) error {
	var pendingDecryptionKey PendingDecryption
	err := pendingDecryptionKey.Deserialize(reader)
	if err != nil {
		return err
	}

	var record DecryptionRecord
	err = record.Deserialize(reader, pendingDecryptionKey.Type)
	if err != nil {
		return err
	}

	return dr.SetRecord(pendingDecryptionKey, record)
}

// Serialize the struct into binary
func (p *PendingDecryption) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write the Hash (32 bytes)
	if err := binary.Write(buf, binary.LittleEndian, p.Hash); err != nil {
		return nil, err
	}

	// Write the Type (int)
	if err := binary.Write(buf, binary.LittleEndian, int32(p.Type)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Deserialize binary data into the struct
func (p *PendingDecryption) Deserialize(reader io.Reader) error {
	// Read the Hash (32 bytes)
	if err := binary.Read(reader, binary.LittleEndian, &p.Hash); err != nil {
		return err
	}

	// Read the Type (int32)
	var typeVal int32
	if err := binary.Read(reader, binary.LittleEndian, &typeVal); err != nil {
		return err
	}
	p.Type = PrecompileName(typeVal)

	return nil
}

// Serialize a decryptionRecord into binary, based on the resultType
func (d *DecryptionRecord) Serialize(resultType PrecompileName) ([]byte, error) {
	// The structure the encoded message is:
	//    encoded_result = result | timestamp
	// where result is:
	// if resultType == SealOutput:
	//    len(result) | result (byte slice)
	// if resultType == Require:
	//    result (bool)
	// if resultType == Decrypt:
	//    len(result) | result (byte slice)
	buf := new(bytes.Buffer)

	// Serialize the Value based on resultType
	switch resultType {
	case SealOutput:
		value, ok := d.Value.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected []byte for SealOutput")
		}
		// Write the length of the byte slice, then the slice itself
		if err := binary.Write(buf, binary.LittleEndian, int32(len(value))); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case Require:
		value, ok := d.Value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool for Require")
		}
		// Write the boolean value
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case Decrypt:
		value, ok := d.Value.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("expected *big.Int for Decrypt")
		}
		// Write the big.Int as a byte slice
		bigIntBytes := value.Bytes()
		if err := binary.Write(buf, binary.LittleEndian, int32(len(bigIntBytes))); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, bigIntBytes); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("tried to serialize unsupported result type")
	}

	// Serialize the Timestamp as int64 (UnixNano)
	timestamp := d.Timestamp.UnixNano()
	if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Deserialize binary data into the struct based on the known resultType
func (d *DecryptionRecord) Deserialize(reader io.Reader, resultType PrecompileName) error {
	// Deserialize the Value based on resultType
	switch resultType {
	case SealOutput:
		var length int32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			return err
		}
		byteSlice := make([]byte, length)
		if err := binary.Read(reader, binary.LittleEndian, &byteSlice); err != nil {
			return err
		}
		d.Value = byteSlice
	case Require:
		var value bool
		if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
			return err
		}
		d.Value = value
	case Decrypt:
		var length int32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			return err
		}
		bigIntBytes := make([]byte, length)
		if err := binary.Read(reader, binary.LittleEndian, &bigIntBytes); err != nil {
			return err
		}
		d.Value = new(big.Int).SetBytes(bigIntBytes)
	default:
		return fmt.Errorf("tried to deserialize unsupported type of result")
	}

	// Deserialize the Timestamp as int64 and convert to time.Time
	var timestamp int64
	if err := binary.Read(reader, binary.LittleEndian, &timestamp); err != nil {
		return err
	}
	d.Timestamp = time.Unix(0, timestamp)

	return nil
}
