// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: consensus/ibft/proto/ibft_operator.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IbftStatusResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *IbftStatusResp) Reset() {
	*x = IbftStatusResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IbftStatusResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IbftStatusResp) ProtoMessage() {}

func (x *IbftStatusResp) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IbftStatusResp.ProtoReflect.Descriptor instead.
func (*IbftStatusResp) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{0}
}

func (x *IbftStatusResp) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type SnapshotReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latest bool   `protobuf:"varint,1,opt,name=latest,proto3" json:"latest,omitempty"`
	Number uint64 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *SnapshotReq) Reset() {
	*x = SnapshotReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SnapshotReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SnapshotReq) ProtoMessage() {}

func (x *SnapshotReq) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SnapshotReq.ProtoReflect.Descriptor instead.
func (*SnapshotReq) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{1}
}

func (x *SnapshotReq) GetLatest() bool {
	if x != nil {
		return x.Latest
	}
	return false
}

func (x *SnapshotReq) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type Snapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Validators []*Snapshot_Validator `protobuf:"bytes,1,rep,name=validators,proto3" json:"validators,omitempty"`
	Number     uint64                `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
	Hash       string                `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Votes      []*Snapshot_Vote      `protobuf:"bytes,4,rep,name=votes,proto3" json:"votes,omitempty"`
}

func (x *Snapshot) Reset() {
	*x = Snapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot) ProtoMessage() {}

func (x *Snapshot) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot.ProtoReflect.Descriptor instead.
func (*Snapshot) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{2}
}

func (x *Snapshot) GetValidators() []*Snapshot_Validator {
	if x != nil {
		return x.Validators
	}
	return nil
}

func (x *Snapshot) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Snapshot) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Snapshot) GetVotes() []*Snapshot_Vote {
	if x != nil {
		return x.Votes
	}
	return nil
}

type ProposeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Auth    bool   `protobuf:"varint,2,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *ProposeReq) Reset() {
	*x = ProposeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeReq) ProtoMessage() {}

func (x *ProposeReq) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeReq.ProtoReflect.Descriptor instead.
func (*ProposeReq) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{3}
}

func (x *ProposeReq) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ProposeReq) GetAuth() bool {
	if x != nil {
		return x.Auth
	}
	return false
}

type CandidatesResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Candidates []*Candidate `protobuf:"bytes,1,rep,name=candidates,proto3" json:"candidates,omitempty"`
}

func (x *CandidatesResp) Reset() {
	*x = CandidatesResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CandidatesResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CandidatesResp) ProtoMessage() {}

func (x *CandidatesResp) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CandidatesResp.ProtoReflect.Descriptor instead.
func (*CandidatesResp) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{4}
}

func (x *CandidatesResp) GetCandidates() []*Candidate {
	if x != nil {
		return x.Candidates
	}
	return nil
}

type Candidate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address   string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	BlsPubkey []byte `protobuf:"bytes,2,opt,name=bls_pubkey,json=blsPubkey,proto3" json:"bls_pubkey,omitempty"`
	Auth      bool   `protobuf:"varint,3,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Candidate) Reset() {
	*x = Candidate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Candidate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Candidate) ProtoMessage() {}

func (x *Candidate) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Candidate.ProtoReflect.Descriptor instead.
func (*Candidate) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{5}
}

func (x *Candidate) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Candidate) GetBlsPubkey() []byte {
	if x != nil {
		return x.BlsPubkey
	}
	return nil
}

func (x *Candidate) GetAuth() bool {
	if x != nil {
		return x.Auth
	}
	return false
}

type Snapshot_Validator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Data    []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Snapshot_Validator) Reset() {
	*x = Snapshot_Validator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot_Validator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot_Validator) ProtoMessage() {}

func (x *Snapshot_Validator) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot_Validator.ProtoReflect.Descriptor instead.
func (*Snapshot_Validator) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Snapshot_Validator) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Snapshot_Validator) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Snapshot_Validator) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Snapshot_Vote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Validator string `protobuf:"bytes,1,opt,name=validator,proto3" json:"validator,omitempty"`
	Proposed  string `protobuf:"bytes,2,opt,name=proposed,proto3" json:"proposed,omitempty"`
	Auth      bool   `protobuf:"varint,3,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Snapshot_Vote) Reset() {
	*x = Snapshot_Vote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot_Vote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot_Vote) ProtoMessage() {}

func (x *Snapshot_Vote) ProtoReflect() protoreflect.Message {
	mi := &file_consensus_ibft_proto_ibft_operator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot_Vote.ProtoReflect.Descriptor instead.
func (*Snapshot_Vote) Descriptor() ([]byte, []int) {
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Snapshot_Vote) GetValidator() string {
	if x != nil {
		return x.Validator
	}
	return ""
}

func (x *Snapshot_Vote) GetProposed() string {
	if x != nil {
		return x.Proposed
	}
	return ""
}

func (x *Snapshot_Vote) GetAuth() bool {
	if x != nil {
		return x.Auth
	}
	return false
}

var File_consensus_ibft_proto_ibft_operator_proto protoreflect.FileDescriptor

var file_consensus_ibft_proto_ibft_operator_proto_rawDesc = []byte{
	0x0a, 0x28, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2f, 0x69, 0x62, 0x66, 0x74,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x62, 0x66, 0x74, 0x5f, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x0e, 0x49,
	0x62, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22,
	0x3d, 0x0a, 0x0b, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0xbc,
	0x02, 0x0a, 0x08, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x36, 0x0a, 0x0a, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x27, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x1a, 0x4d, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x54, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x75, 0x74,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x3a, 0x0a,
	0x0a, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x3f, 0x0a, 0x0e, 0x43, 0x61, 0x6e,
	0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2d, 0x0a, 0x0a, 0x63,
	0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x0a,
	0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x73, 0x22, 0x58, 0x0a, 0x09, 0x43, 0x61,
	0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x73, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x61, 0x75, 0x74, 0x68, 0x32, 0xde, 0x01, 0x0a, 0x0c, 0x49, 0x62, 0x66, 0x74, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x2c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x12, 0x0f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68,
	0x6f, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x12, 0x30, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x12, 0x0d,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x0a, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x34, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x62, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x42, 0x17, 0x5a, 0x15, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e,
	0x73, 0x75, 0x73, 0x2f, 0x69, 0x62, 0x66, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_consensus_ibft_proto_ibft_operator_proto_rawDescOnce sync.Once
	file_consensus_ibft_proto_ibft_operator_proto_rawDescData = file_consensus_ibft_proto_ibft_operator_proto_rawDesc
)

func file_consensus_ibft_proto_ibft_operator_proto_rawDescGZIP() []byte {
	file_consensus_ibft_proto_ibft_operator_proto_rawDescOnce.Do(func() {
		file_consensus_ibft_proto_ibft_operator_proto_rawDescData = protoimpl.X.CompressGZIP(file_consensus_ibft_proto_ibft_operator_proto_rawDescData)
	})
	return file_consensus_ibft_proto_ibft_operator_proto_rawDescData
}

var file_consensus_ibft_proto_ibft_operator_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_consensus_ibft_proto_ibft_operator_proto_goTypes = []interface{}{
	(*IbftStatusResp)(nil),     // 0: v1.IbftStatusResp
	(*SnapshotReq)(nil),        // 1: v1.SnapshotReq
	(*Snapshot)(nil),           // 2: v1.Snapshot
	(*ProposeReq)(nil),         // 3: v1.ProposeReq
	(*CandidatesResp)(nil),     // 4: v1.CandidatesResp
	(*Candidate)(nil),          // 5: v1.Candidate
	(*Snapshot_Validator)(nil), // 6: v1.Snapshot.Validator
	(*Snapshot_Vote)(nil),      // 7: v1.Snapshot.Vote
	(*emptypb.Empty)(nil),      // 8: google.protobuf.Empty
}
var file_consensus_ibft_proto_ibft_operator_proto_depIdxs = []int32{
	6, // 0: v1.Snapshot.validators:type_name -> v1.Snapshot.Validator
	7, // 1: v1.Snapshot.votes:type_name -> v1.Snapshot.Vote
	5, // 2: v1.CandidatesResp.candidates:type_name -> v1.Candidate
	1, // 3: v1.IbftOperator.GetSnapshot:input_type -> v1.SnapshotReq
	5, // 4: v1.IbftOperator.Propose:input_type -> v1.Candidate
	8, // 5: v1.IbftOperator.Candidates:input_type -> google.protobuf.Empty
	8, // 6: v1.IbftOperator.Status:input_type -> google.protobuf.Empty
	2, // 7: v1.IbftOperator.GetSnapshot:output_type -> v1.Snapshot
	8, // 8: v1.IbftOperator.Propose:output_type -> google.protobuf.Empty
	4, // 9: v1.IbftOperator.Candidates:output_type -> v1.CandidatesResp
	0, // 10: v1.IbftOperator.Status:output_type -> v1.IbftStatusResp
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_consensus_ibft_proto_ibft_operator_proto_init() }
func file_consensus_ibft_proto_ibft_operator_proto_init() {
	if File_consensus_ibft_proto_ibft_operator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IbftStatusResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SnapshotReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProposeReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CandidatesResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Candidate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot_Validator); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_consensus_ibft_proto_ibft_operator_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot_Vote); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_consensus_ibft_proto_ibft_operator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_consensus_ibft_proto_ibft_operator_proto_goTypes,
		DependencyIndexes: file_consensus_ibft_proto_ibft_operator_proto_depIdxs,
		MessageInfos:      file_consensus_ibft_proto_ibft_operator_proto_msgTypes,
	}.Build()
	File_consensus_ibft_proto_ibft_operator_proto = out.File
	file_consensus_ibft_proto_ibft_operator_proto_rawDesc = nil
	file_consensus_ibft_proto_ibft_operator_proto_goTypes = nil
	file_consensus_ibft_proto_ibft_operator_proto_depIdxs = nil
}
