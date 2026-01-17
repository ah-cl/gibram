// Package codec provides encoding/decoding for GibRAM protocol
package codec

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"

	"github.com/gibram-io/gibram/pkg/types"
	pb "github.com/gibram-io/gibram/proto/gibrampb"
	"google.golang.org/protobuf/proto"
)

// CodecType represents the encoding type
type CodecType byte

const (
	CodecJSON     CodecType = 0x00 // JSON encoding (legacy, default)
	CodecProtobuf CodecType = 0x01 // Protobuf encoding (new)
)

// Frame represents a wire frame
type Frame struct {
	CodecType CodecType
	CmdType   pb.CommandType
	Payload   []byte
}

// =============================================================================
// Protobuf Encoder
// =============================================================================

// EncodeEnvelope encodes an envelope to wire format
func EncodeEnvelope(env *pb.Envelope) ([]byte, error) {
	data, err := proto.Marshal(env)
	if err != nil {
		return nil, err
	}

	// Frame: [1 byte codec][4 bytes length][payload]
	frame := make([]byte, 1+4+len(data))
	frame[0] = byte(CodecProtobuf)
	binary.BigEndian.PutUint32(frame[1:5], uint32(len(data)))
	copy(frame[5:], data)

	return frame, nil
}

// DecodeEnvelope decodes an envelope from wire format
func DecodeEnvelope(r io.Reader) (*pb.Envelope, CodecType, error) {
	// Read codec type
	var codecByte [1]byte
	if _, err := io.ReadFull(r, codecByte[:]); err != nil {
		return nil, 0, err
	}

	codecType := CodecType(codecByte[0])

	// Read length
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, codecType, err
	}

	if length > 64*1024*1024 { // 64MB max
		return nil, codecType, errors.New("frame too large")
	}

	// Read payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, codecType, err
	}

	if codecType == CodecProtobuf {
		var env pb.Envelope
		if err := proto.Unmarshal(payload, &env); err != nil {
			return nil, codecType, err
		}
		return &env, codecType, nil
	}

	// Legacy JSON - convert to envelope
	return nil, codecType, errors.New("legacy JSON codec not supported in this decoder")
}

// =============================================================================
// Type Converters: types.* <-> pb.*
// =============================================================================

// DocumentToProto converts types.Document to pb.Document
func DocumentToProto(doc *types.Document) *pb.Document {
	return &pb.Document{
		Id:         doc.ID,
		ExternalId: doc.ExternalID,
		Filename:   doc.Filename,
		Status:     string(doc.Status),
		CreatedAt:  doc.CreatedAt,
	}
}

// ProtoToDocument converts pb.Document to types.Document
func ProtoToDocument(doc *pb.Document) *types.Document {
	return &types.Document{
		ID:         doc.Id,
		ExternalID: doc.ExternalId,
		Filename:   doc.Filename,
		Status:     types.DocumentStatus(doc.Status),
		CreatedAt:  doc.CreatedAt,
	}
}

// TextUnitToProto converts types.TextUnit to pb.TextUnit
func TextUnitToProto(tu *types.TextUnit) *pb.TextUnit {
	return &pb.TextUnit{
		Id:         tu.ID,
		ExternalId: tu.ExternalID,
		DocumentId: tu.DocumentID,
		Content:    tu.Content,
		TokenCount: int32(tu.TokenCount),
		EntityIds:  tu.EntityIDs,
		CreatedAt:  tu.CreatedAt,
	}
}

// ProtoToTextUnit converts pb.TextUnit to types.TextUnit
func ProtoToTextUnit(tu *pb.TextUnit) *types.TextUnit {
	return &types.TextUnit{
		ID:         tu.Id,
		ExternalID: tu.ExternalId,
		DocumentID: tu.DocumentId,
		Content:    tu.Content,
		TokenCount: int(tu.TokenCount),
		EntityIDs:  tu.EntityIds,
		CreatedAt:  tu.CreatedAt,
	}
}

// EntityToProto converts types.Entity to pb.Entity
func EntityToProto(ent *types.Entity) *pb.Entity {
	return &pb.Entity{
		Id:          ent.ID,
		ExternalId:  ent.ExternalID,
		Title:       ent.Title,
		Type:        ent.Type,
		Description: ent.Description,
		TextunitIds: ent.TextUnitIDs,
		CreatedAt:   ent.CreatedAt,
	}
}

// ProtoToEntity converts pb.Entity to types.Entity
func ProtoToEntity(ent *pb.Entity) *types.Entity {
	return &types.Entity{
		ID:          ent.Id,
		ExternalID:  ent.ExternalId,
		Title:       ent.Title,
		Type:        ent.Type,
		Description: ent.Description,
		TextUnitIDs: ent.TextunitIds,
		CreatedAt:   ent.CreatedAt,
	}
}

// RelationshipToProto converts types.Relationship to pb.Relationship
func RelationshipToProto(rel *types.Relationship) *pb.Relationship {
	return &pb.Relationship{
		Id:          rel.ID,
		ExternalId:  rel.ExternalID,
		SourceId:    rel.SourceID,
		TargetId:    rel.TargetID,
		Type:        rel.Type,
		Description: rel.Description,
		Weight:      rel.Weight,
		CreatedAt:   rel.CreatedAt,
	}
}

// ProtoToRelationship converts pb.Relationship to types.Relationship
func ProtoToRelationship(rel *pb.Relationship) *types.Relationship {
	return &types.Relationship{
		ID:          rel.Id,
		ExternalID:  rel.ExternalId,
		SourceID:    rel.SourceId,
		TargetID:    rel.TargetId,
		Type:        rel.Type,
		Description: rel.Description,
		Weight:      rel.Weight,
		CreatedAt:   rel.CreatedAt,
	}
}

// CommunityToProto converts types.Community to pb.Community
func CommunityToProto(comm *types.Community) *pb.Community {
	return &pb.Community{
		Id:              comm.ID,
		ExternalId:      comm.ExternalID,
		Title:           comm.Title,
		Summary:         comm.Summary,
		FullContent:     comm.FullContent,
		Level:           int32(comm.Level),
		EntityIds:       comm.EntityIDs,
		RelationshipIds: comm.RelationshipIDs,
		CreatedAt:       comm.CreatedAt,
	}
}

// ProtoToCommunity converts pb.Community to types.Community
func ProtoToCommunity(comm *pb.Community) *types.Community {
	return &types.Community{
		ID:              comm.Id,
		ExternalID:      comm.ExternalId,
		Title:           comm.Title,
		Summary:         comm.Summary,
		FullContent:     comm.FullContent,
		Level:           int(comm.Level),
		EntityIDs:       comm.EntityIds,
		RelationshipIDs: comm.RelationshipIds,
		CreatedAt:       comm.CreatedAt,
	}
}

// =============================================================================
// Binary WAL Encoding (more compact than JSON)
// =============================================================================

// WALEntry represents a binary WAL entry
type WALEntry struct {
	Op      byte
	Payload []byte
}

// EncodeWALEntry encodes a WAL entry to binary
func EncodeWALEntry(op byte, payload interface{}) ([]byte, error) {
	// For now, use JSON for payload but binary framing
	// This can be optimized later to full binary
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// [1 byte op][4 bytes length][payload]
	entry := make([]byte, 1+4+len(payloadBytes))
	entry[0] = op
	binary.LittleEndian.PutUint32(entry[1:5], uint32(len(payloadBytes)))
	copy(entry[5:], payloadBytes)

	return entry, nil
}

// DecodeWALEntry decodes a binary WAL entry
func DecodeWALEntry(r io.Reader) (*WALEntry, error) {
	// Read op
	var op [1]byte
	if _, err := io.ReadFull(r, op[:]); err != nil {
		return nil, err
	}

	// Read length
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	// Read payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return &WALEntry{Op: op[0], Payload: payload}, nil
}

// =============================================================================
// Snapshot Binary Encoding
// =============================================================================

// SnapshotHeader for binary snapshots
type SnapshotHeader struct {
	Magic      uint32 // "GIB2" = 0x47494232
	Version    uint16
	VectorDim  uint16
	DocCount   uint64
	TUCount    uint64
	EntCount   uint64
	RelCount   uint64
	CommCount  uint64
	QueryCount uint64
	CreatedAt  int64
}

const SnapshotMagic = 0x47494232 // "GIB2"

// EncodeSnapshotHeader encodes header to binary
func EncodeSnapshotHeader(h *SnapshotHeader) []byte {
	buf := make([]byte, 64) // Fixed header size
	binary.LittleEndian.PutUint32(buf[0:4], h.Magic)
	binary.LittleEndian.PutUint16(buf[4:6], h.Version)
	binary.LittleEndian.PutUint16(buf[6:8], h.VectorDim)
	binary.LittleEndian.PutUint64(buf[8:16], h.DocCount)
	binary.LittleEndian.PutUint64(buf[16:24], h.TUCount)
	binary.LittleEndian.PutUint64(buf[24:32], h.EntCount)
	binary.LittleEndian.PutUint64(buf[32:40], h.RelCount)
	binary.LittleEndian.PutUint64(buf[40:48], h.CommCount)
	binary.LittleEndian.PutUint64(buf[48:56], h.QueryCount)
	binary.LittleEndian.PutUint64(buf[56:64], uint64(h.CreatedAt))
	return buf
}

// DecodeSnapshotHeader decodes header from binary
func DecodeSnapshotHeader(r io.Reader) (*SnapshotHeader, error) {
	buf := make([]byte, 64)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	h := &SnapshotHeader{
		Magic:      binary.LittleEndian.Uint32(buf[0:4]),
		Version:    binary.LittleEndian.Uint16(buf[4:6]),
		VectorDim:  binary.LittleEndian.Uint16(buf[6:8]),
		DocCount:   binary.LittleEndian.Uint64(buf[8:16]),
		TUCount:    binary.LittleEndian.Uint64(buf[16:24]),
		EntCount:   binary.LittleEndian.Uint64(buf[24:32]),
		RelCount:   binary.LittleEndian.Uint64(buf[32:40]),
		CommCount:  binary.LittleEndian.Uint64(buf[40:48]),
		QueryCount: binary.LittleEndian.Uint64(buf[48:56]),
		CreatedAt:  int64(binary.LittleEndian.Uint64(buf[56:64])),
	}

	if h.Magic != SnapshotMagic {
		return nil, errors.New("invalid snapshot magic")
	}

	return h, nil
}
