package codec

import (
	"bytes"
	"testing"
	"time"

	"github.com/gibram-io/gibram/pkg/types"
	pb "github.com/gibram-io/gibram/proto/gibrampb"
)

// =============================================================================
// Test Type Converters: Document
// =============================================================================

func TestDocumentToProto(t *testing.T) {
	doc := &types.Document{
		ID:         1,
		ExternalID: "ext-123",
		Filename:   "test.txt",
		Status:     types.DocStatusUploaded,
		CreatedAt:  1000,
	}

	pbDoc := DocumentToProto(doc)

	if pbDoc.Id != doc.ID {
		t.Errorf("expected Id %d, got %d", doc.ID, pbDoc.Id)
	}
	if pbDoc.ExternalId != doc.ExternalID {
		t.Errorf("expected ExternalId %s, got %s", doc.ExternalID, pbDoc.ExternalId)
	}
	if pbDoc.Filename != doc.Filename {
		t.Errorf("expected Filename %s, got %s", doc.Filename, pbDoc.Filename)
	}
	if pbDoc.Status != string(doc.Status) {
		t.Errorf("expected Status %s, got %s", doc.Status, pbDoc.Status)
	}
	if pbDoc.CreatedAt != doc.CreatedAt {
		t.Errorf("expected CreatedAt %d, got %d", doc.CreatedAt, pbDoc.CreatedAt)
	}
	// TTL fields removed in v0.1.0 - session-based architecture
}

func TestProtoToDocument(t *testing.T) {
	pbDoc := &pb.Document{
		Id:         1,
		ExternalId: "ext-123",
		Filename:   "test.txt",
		Status:     "pending",
		CreatedAt:  1000,
	}

	doc := ProtoToDocument(pbDoc)

	if doc.ID != pbDoc.Id {
		t.Errorf("expected ID %d, got %d", pbDoc.Id, doc.ID)
	}
	if doc.ExternalID != pbDoc.ExternalId {
		t.Errorf("expected ExternalID %s, got %s", pbDoc.ExternalId, doc.ExternalID)
	}
	if doc.Filename != pbDoc.Filename {
		t.Errorf("expected Filename %s, got %s", pbDoc.Filename, doc.Filename)
	}
	if string(doc.Status) != pbDoc.Status {
		t.Errorf("expected Status %s, got %s", pbDoc.Status, doc.Status)
	}
	// TTL fields removed in v0.1.0
}

// =============================================================================
// Test Type Converters: TextUnit
// =============================================================================

func TestTextUnitToProto(t *testing.T) {
	tu := &types.TextUnit{
		ID:         1,
		ExternalID: "tu-123",
		DocumentID: 10,
		Content:    "Test content",
		TokenCount: 5,
		EntityIDs:  []uint64{1, 2, 3},
		CreatedAt:  1000,
	}

	pbTU := TextUnitToProto(tu)

	if pbTU.Id != tu.ID {
		t.Errorf("expected Id %d, got %d", tu.ID, pbTU.Id)
	}
	if pbTU.ExternalId != tu.ExternalID {
		t.Errorf("expected ExternalId %s, got %s", tu.ExternalID, pbTU.ExternalId)
	}
	if pbTU.DocumentId != tu.DocumentID {
		t.Errorf("expected DocumentId %d, got %d", tu.DocumentID, pbTU.DocumentId)
	}
	if pbTU.Content != tu.Content {
		t.Errorf("expected Content %s, got %s", tu.Content, pbTU.Content)
	}
	if int(pbTU.TokenCount) != tu.TokenCount {
		t.Errorf("expected TokenCount %d, got %d", tu.TokenCount, pbTU.TokenCount)
	}
	if len(pbTU.EntityIds) != len(tu.EntityIDs) {
		t.Errorf("expected %d EntityIds, got %d", len(tu.EntityIDs), len(pbTU.EntityIds))
	}
}

func TestProtoToTextUnit(t *testing.T) {
	pbTU := &pb.TextUnit{
		Id:         1,
		ExternalId: "tu-123",
		DocumentId: 10,
		Content:    "Test content",
		TokenCount: 5,
		EntityIds:  []uint64{1, 2, 3},
		CreatedAt:  1000,
	}

	tu := ProtoToTextUnit(pbTU)

	if tu.ID != pbTU.Id {
		t.Errorf("expected ID %d, got %d", pbTU.Id, tu.ID)
	}
	if tu.DocumentID != pbTU.DocumentId {
		t.Errorf("expected DocumentID %d, got %d", pbTU.DocumentId, tu.DocumentID)
	}
	if tu.TokenCount != int(pbTU.TokenCount) {
		t.Errorf("expected TokenCount %d, got %d", pbTU.TokenCount, tu.TokenCount)
	}
}

// =============================================================================
// Test Type Converters: Entity
// =============================================================================

func TestEntityToProto(t *testing.T) {
	ent := &types.Entity{
		ID:          1,
		ExternalID:  "ent-123",
		Title:       "Test Entity",
		Type:        "PERSON",
		Description: "A test entity",
		TextUnitIDs: []uint64{1, 2},
		CreatedAt:   1000,
	}

	pbEnt := EntityToProto(ent)

	if pbEnt.Id != ent.ID {
		t.Errorf("expected Id %d, got %d", ent.ID, pbEnt.Id)
	}
	if pbEnt.ExternalId != ent.ExternalID {
		t.Errorf("expected ExternalId %s, got %s", ent.ExternalID, pbEnt.ExternalId)
	}
	if pbEnt.Title != ent.Title {
		t.Errorf("expected Title %s, got %s", ent.Title, pbEnt.Title)
	}
	if pbEnt.Type != ent.Type {
		t.Errorf("expected Type %s, got %s", ent.Type, pbEnt.Type)
	}
	if pbEnt.Description != ent.Description {
		t.Errorf("expected Description %s, got %s", ent.Description, pbEnt.Description)
	}
}

func TestProtoToEntity(t *testing.T) {
	pbEnt := &pb.Entity{
		Id:          1,
		ExternalId:  "ent-123",
		Title:       "Test Entity",
		Type:        "PERSON",
		Description: "A test entity",
		TextunitIds: []uint64{1, 2},
		CreatedAt:   1000,
	}

	ent := ProtoToEntity(pbEnt)

	if ent.ID != pbEnt.Id {
		t.Errorf("expected ID %d, got %d", pbEnt.Id, ent.ID)
	}
	if ent.Title != pbEnt.Title {
		t.Errorf("expected Title %s, got %s", pbEnt.Title, ent.Title)
	}
	if len(ent.TextUnitIDs) != len(pbEnt.TextunitIds) {
		t.Errorf("expected %d TextUnitIDs, got %d", len(pbEnt.TextunitIds), len(ent.TextUnitIDs))
	}
}

// =============================================================================
// Test Type Converters: Relationship
// =============================================================================

func TestRelationshipToProto(t *testing.T) {
	rel := &types.Relationship{
		ID:          1,
		ExternalID:  "rel-123",
		SourceID:    10,
		TargetID:    20,
		Type:        "KNOWS",
		Description: "Test relationship",
		Weight:      0.8,
		CreatedAt:   1000,
	}

	pbRel := RelationshipToProto(rel)

	if pbRel.Id != rel.ID {
		t.Errorf("expected Id %d, got %d", rel.ID, pbRel.Id)
	}
	if pbRel.SourceId != rel.SourceID {
		t.Errorf("expected SourceId %d, got %d", rel.SourceID, pbRel.SourceId)
	}
	if pbRel.TargetId != rel.TargetID {
		t.Errorf("expected TargetId %d, got %d", rel.TargetID, pbRel.TargetId)
	}
	if pbRel.Type != rel.Type {
		t.Errorf("expected Type %s, got %s", rel.Type, pbRel.Type)
	}
	if pbRel.Weight != rel.Weight {
		t.Errorf("expected Weight %f, got %f", rel.Weight, pbRel.Weight)
	}
}

func TestProtoToRelationship(t *testing.T) {
	pbRel := &pb.Relationship{
		Id:          1,
		ExternalId:  "rel-123",
		SourceId:    10,
		TargetId:    20,
		Type:        "KNOWS",
		Description: "Test relationship",
		Weight:      0.8,
		CreatedAt:   1000,
	}

	rel := ProtoToRelationship(pbRel)

	if rel.ID != pbRel.Id {
		t.Errorf("expected ID %d, got %d", pbRel.Id, rel.ID)
	}
	if rel.SourceID != pbRel.SourceId {
		t.Errorf("expected SourceID %d, got %d", pbRel.SourceId, rel.SourceID)
	}
	if rel.TargetID != pbRel.TargetId {
		t.Errorf("expected TargetID %d, got %d", pbRel.TargetId, rel.TargetID)
	}
}

// =============================================================================
// Test Type Converters: Community
// =============================================================================

func TestCommunityToProto(t *testing.T) {
	comm := &types.Community{
		ID:              1,
		ExternalID:      "comm-123",
		Title:           "Test Community",
		Summary:         "A test community",
		FullContent:     "Full content here",
		Level:           2,
		EntityIDs:       []uint64{1, 2, 3},
		RelationshipIDs: []uint64{10, 20},
		CreatedAt:       1000,
	}

	pbComm := CommunityToProto(comm)

	if pbComm.Id != comm.ID {
		t.Errorf("expected Id %d, got %d", comm.ID, pbComm.Id)
	}
	if pbComm.Title != comm.Title {
		t.Errorf("expected Title %s, got %s", comm.Title, pbComm.Title)
	}
	if pbComm.Summary != comm.Summary {
		t.Errorf("expected Summary %s, got %s", comm.Summary, pbComm.Summary)
	}
	if int(pbComm.Level) != comm.Level {
		t.Errorf("expected Level %d, got %d", comm.Level, pbComm.Level)
	}
	if len(pbComm.EntityIds) != len(comm.EntityIDs) {
		t.Errorf("expected %d EntityIds, got %d", len(comm.EntityIDs), len(pbComm.EntityIds))
	}
}

func TestProtoToCommunity(t *testing.T) {
	pbComm := &pb.Community{
		Id:              1,
		ExternalId:      "comm-123",
		Title:           "Test Community",
		Summary:         "A test community",
		FullContent:     "Full content here",
		Level:           2,
		EntityIds:       []uint64{1, 2, 3},
		RelationshipIds: []uint64{10, 20},
		CreatedAt:       1000,
	}

	comm := ProtoToCommunity(pbComm)

	if comm.ID != pbComm.Id {
		t.Errorf("expected ID %d, got %d", pbComm.Id, comm.ID)
	}
	if comm.Level != int(pbComm.Level) {
		t.Errorf("expected Level %d, got %d", pbComm.Level, comm.Level)
	}
	if len(comm.RelationshipIDs) != len(pbComm.RelationshipIds) {
		t.Errorf("expected %d RelationshipIDs, got %d", len(pbComm.RelationshipIds), len(comm.RelationshipIDs))
	}
}

// =============================================================================
// Test Envelope Encoding/Decoding
// =============================================================================

func TestEncodeDecodeEnvelope(t *testing.T) {
	env := &pb.Envelope{
		RequestId: 123,
		Version:   2,
	}

	// Encode
	data, err := EncodeEnvelope(env)
	if err != nil {
		t.Fatalf("failed to encode envelope: %v", err)
	}

	// Verify frame structure
	if data[0] != byte(CodecProtobuf) {
		t.Errorf("expected codec type %d, got %d", CodecProtobuf, data[0])
	}

	// Decode
	reader := bytes.NewReader(data)
	decoded, codecType, err := DecodeEnvelope(reader)
	if err != nil {
		t.Fatalf("failed to decode envelope: %v", err)
	}

	if codecType != CodecProtobuf {
		t.Errorf("expected codec type %d, got %d", CodecProtobuf, codecType)
	}

	if decoded.RequestId != env.RequestId {
		t.Errorf("expected RequestId %d, got %d", env.RequestId, decoded.RequestId)
	}
}

func TestDecodeEnvelope_FrameTooLarge(t *testing.T) {
	// Create frame header with length > 64MB
	data := make([]byte, 5)
	data[0] = byte(CodecProtobuf)
	// Set length to 70MB (> 64MB limit)
	data[1] = 0x04
	data[2] = 0x30
	data[3] = 0x00
	data[4] = 0x00

	reader := bytes.NewReader(data)
	_, _, err := DecodeEnvelope(reader)
	if err == nil {
		t.Error("expected error for frame too large")
	}
}

func TestDecodeEnvelope_LegacyJSON(t *testing.T) {
	// Create frame with JSON codec type
	data := make([]byte, 9)
	data[0] = byte(CodecJSON)
	// Length = 4
	data[1] = 0x00
	data[2] = 0x00
	data[3] = 0x00
	data[4] = 0x04
	// Payload
	data[5] = '{'
	data[6] = '}'
	data[7] = '\n'
	data[8] = '\n'

	reader := bytes.NewReader(data)
	_, _, err := DecodeEnvelope(reader)
	if err == nil {
		t.Error("expected error for legacy JSON codec")
	}
}

// =============================================================================
// Test WAL Entry Encoding/Decoding
// =============================================================================

func TestEncodeDecodeWALEntry(t *testing.T) {
	payload := map[string]interface{}{
		"id":    1,
		"name":  "test",
		"value": 42.5,
	}

	// Encode
	data, err := EncodeWALEntry(0x01, payload)
	if err != nil {
		t.Fatalf("failed to encode WAL entry: %v", err)
	}

	// Verify structure
	if data[0] != 0x01 {
		t.Errorf("expected op 0x01, got %x", data[0])
	}

	// Decode
	reader := bytes.NewReader(data)
	entry, err := DecodeWALEntry(reader)
	if err != nil {
		t.Fatalf("failed to decode WAL entry: %v", err)
	}

	if entry.Op != 0x01 {
		t.Errorf("expected Op 0x01, got %x", entry.Op)
	}

	if len(entry.Payload) == 0 {
		t.Error("expected non-empty payload")
	}
}

func TestDecodeWALEntry_EmptyReader(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	_, err := DecodeWALEntry(reader)
	if err == nil {
		t.Error("expected error for empty reader")
	}
}

func TestDecodeWALEntry_TruncatedLength(t *testing.T) {
	// Only op byte, no length
	reader := bytes.NewReader([]byte{0x01})
	_, err := DecodeWALEntry(reader)
	if err == nil {
		t.Error("expected error for truncated length")
	}
}

func TestDecodeWALEntry_TruncatedPayload(t *testing.T) {
	// Op + length header but not enough payload
	data := []byte{0x01, 0x10, 0x00, 0x00, 0x00} // length = 16, but no payload
	reader := bytes.NewReader(data)
	_, err := DecodeWALEntry(reader)
	if err == nil {
		t.Error("expected error for truncated payload")
	}
}

// =============================================================================
// Test Snapshot Header Encoding/Decoding
// =============================================================================

func TestEncodeDecodeSnapshotHeader(t *testing.T) {
	header := &SnapshotHeader{
		Magic:      SnapshotMagic,
		Version:    1,
		VectorDim:  1536,
		DocCount:   100,
		TUCount:    500,
		EntCount:   200,
		RelCount:   300,
		CommCount:  50,
		QueryCount: 1000,
		CreatedAt:  time.Now().Unix(),
	}

	// Encode
	data := EncodeSnapshotHeader(header)
	if len(data) != 64 {
		t.Errorf("expected header size 64, got %d", len(data))
	}

	// Decode
	reader := bytes.NewReader(data)
	decoded, err := DecodeSnapshotHeader(reader)
	if err != nil {
		t.Fatalf("failed to decode snapshot header: %v", err)
	}

	if decoded.Magic != header.Magic {
		t.Errorf("expected Magic %x, got %x", header.Magic, decoded.Magic)
	}
	if decoded.Version != header.Version {
		t.Errorf("expected Version %d, got %d", header.Version, decoded.Version)
	}
	if decoded.VectorDim != header.VectorDim {
		t.Errorf("expected VectorDim %d, got %d", header.VectorDim, decoded.VectorDim)
	}
	if decoded.DocCount != header.DocCount {
		t.Errorf("expected DocCount %d, got %d", header.DocCount, decoded.DocCount)
	}
	if decoded.TUCount != header.TUCount {
		t.Errorf("expected TUCount %d, got %d", header.TUCount, decoded.TUCount)
	}
	if decoded.EntCount != header.EntCount {
		t.Errorf("expected EntCount %d, got %d", header.EntCount, decoded.EntCount)
	}
	if decoded.RelCount != header.RelCount {
		t.Errorf("expected RelCount %d, got %d", header.RelCount, decoded.RelCount)
	}
	if decoded.CommCount != header.CommCount {
		t.Errorf("expected CommCount %d, got %d", header.CommCount, decoded.CommCount)
	}
	if decoded.QueryCount != header.QueryCount {
		t.Errorf("expected QueryCount %d, got %d", header.QueryCount, decoded.QueryCount)
	}
	if decoded.CreatedAt != header.CreatedAt {
		t.Errorf("expected CreatedAt %d, got %d", header.CreatedAt, decoded.CreatedAt)
	}
}

func TestDecodeSnapshotHeader_InvalidMagic(t *testing.T) {
	header := &SnapshotHeader{
		Magic:   0x12345678, // Invalid magic
		Version: 1,
	}

	data := EncodeSnapshotHeader(header)
	reader := bytes.NewReader(data)
	_, err := DecodeSnapshotHeader(reader)
	if err == nil {
		t.Error("expected error for invalid magic")
	}
}

func TestDecodeSnapshotHeader_TruncatedData(t *testing.T) {
	data := make([]byte, 32) // Only half the header size
	reader := bytes.NewReader(data)
	_, err := DecodeSnapshotHeader(reader)
	if err == nil {
		t.Error("expected error for truncated data")
	}
}

// =============================================================================
// Test Codec Constants
// =============================================================================

func TestCodecConstants(t *testing.T) {
	if CodecJSON != 0x00 {
		t.Errorf("expected CodecJSON = 0x00, got %x", CodecJSON)
	}
	if CodecProtobuf != 0x01 {
		t.Errorf("expected CodecProtobuf = 0x01, got %x", CodecProtobuf)
	}
	if SnapshotMagic != 0x47494232 {
		t.Errorf("expected SnapshotMagic = 0x47494232 (GIB2), got %x", SnapshotMagic)
	}
}
