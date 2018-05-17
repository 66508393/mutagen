// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/havoc-io/mutagen/pkg/session/state.proto

package session

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import rsync "github.com/havoc-io/mutagen/pkg/rsync"
import sync "github.com/havoc-io/mutagen/pkg/sync"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Status int32

const (
	Status_Disconnected         Status = 0
	Status_HaltedOnRootDeletion Status = 1
	Status_ConnectingAlpha      Status = 2
	Status_ConnectingBeta       Status = 3
	Status_Watching             Status = 4
	Status_ScanningAlpha        Status = 5
	Status_ScanningBeta         Status = 6
	Status_WaitingForRescan     Status = 7
	Status_Reconciling          Status = 8
	Status_StagingAlpha         Status = 9
	Status_StagingBeta          Status = 10
	Status_TransitioningAlpha   Status = 11
	Status_TransitioningBeta    Status = 12
	Status_Saving               Status = 13
)

var Status_name = map[int32]string{
	0:  "Disconnected",
	1:  "HaltedOnRootDeletion",
	2:  "ConnectingAlpha",
	3:  "ConnectingBeta",
	4:  "Watching",
	5:  "ScanningAlpha",
	6:  "ScanningBeta",
	7:  "WaitingForRescan",
	8:  "Reconciling",
	9:  "StagingAlpha",
	10: "StagingBeta",
	11: "TransitioningAlpha",
	12: "TransitioningBeta",
	13: "Saving",
}
var Status_value = map[string]int32{
	"Disconnected":         0,
	"HaltedOnRootDeletion": 1,
	"ConnectingAlpha":      2,
	"ConnectingBeta":       3,
	"Watching":             4,
	"ScanningAlpha":        5,
	"ScanningBeta":         6,
	"WaitingForRescan":     7,
	"Reconciling":          8,
	"StagingAlpha":         9,
	"StagingBeta":          10,
	"TransitioningAlpha":   11,
	"TransitioningBeta":    12,
	"Saving":               13,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_state_71685f93597e686f, []int{0}
}

type State struct {
	Session              *Session              `protobuf:"bytes,1,opt,name=session" json:"session,omitempty"`
	Status               Status                `protobuf:"varint,2,opt,name=status,enum=session.Status" json:"status,omitempty"`
	AlphaConnected       bool                  `protobuf:"varint,3,opt,name=alphaConnected" json:"alphaConnected,omitempty"`
	BetaConnected        bool                  `protobuf:"varint,4,opt,name=betaConnected" json:"betaConnected,omitempty"`
	LastError            string                `protobuf:"bytes,5,opt,name=lastError" json:"lastError,omitempty"`
	StagingStatus        *rsync.ReceiverStatus `protobuf:"bytes,6,opt,name=stagingStatus" json:"stagingStatus,omitempty"`
	Conflicts            []*sync.Conflict      `protobuf:"bytes,7,rep,name=conflicts" json:"conflicts,omitempty"`
	AlphaProblems        []*sync.Problem       `protobuf:"bytes,8,rep,name=alphaProblems" json:"alphaProblems,omitempty"`
	BetaProblems         []*sync.Problem       `protobuf:"bytes,9,rep,name=betaProblems" json:"betaProblems,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_state_71685f93597e686f, []int{0}
}
func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (dst *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(dst, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

func (m *State) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_Disconnected
}

func (m *State) GetAlphaConnected() bool {
	if m != nil {
		return m.AlphaConnected
	}
	return false
}

func (m *State) GetBetaConnected() bool {
	if m != nil {
		return m.BetaConnected
	}
	return false
}

func (m *State) GetLastError() string {
	if m != nil {
		return m.LastError
	}
	return ""
}

func (m *State) GetStagingStatus() *rsync.ReceiverStatus {
	if m != nil {
		return m.StagingStatus
	}
	return nil
}

func (m *State) GetConflicts() []*sync.Conflict {
	if m != nil {
		return m.Conflicts
	}
	return nil
}

func (m *State) GetAlphaProblems() []*sync.Problem {
	if m != nil {
		return m.AlphaProblems
	}
	return nil
}

func (m *State) GetBetaProblems() []*sync.Problem {
	if m != nil {
		return m.BetaProblems
	}
	return nil
}

func init() {
	proto.RegisterType((*State)(nil), "session.State")
	proto.RegisterEnum("session.Status", Status_name, Status_value)
}

func init() {
	proto.RegisterFile("github.com/havoc-io/mutagen/pkg/session/state.proto", fileDescriptor_state_71685f93597e686f)
}

var fileDescriptor_state_71685f93597e686f = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x6f, 0xd3, 0x30,
	0x14, 0xc6, 0xc9, 0xba, 0xa6, 0xed, 0x6b, 0xd3, 0x66, 0x8f, 0x0d, 0x45, 0x13, 0x87, 0x0a, 0x21,
	0xa8, 0x26, 0x48, 0x44, 0x2b, 0x4e, 0x9c, 0x60, 0x03, 0x71, 0x03, 0xb9, 0x48, 0x3b, 0xbb, 0x9e,
	0x49, 0x2d, 0x52, 0xbb, 0xb2, 0xdd, 0x4a, 0xfc, 0x21, 0x9c, 0xf8, 0x67, 0x91, 0x1d, 0xa7, 0x25,
	0x08, 0x69, 0x3b, 0x59, 0xfa, 0xf2, 0xfd, 0x9e, 0xbf, 0xf7, 0x39, 0xb0, 0x28, 0x85, 0x5d, 0xef,
	0x56, 0x39, 0x53, 0x9b, 0x62, 0x4d, 0xf7, 0x8a, 0xbd, 0x16, 0xaa, 0xd8, 0xec, 0x2c, 0x2d, 0xb9,
	0x2c, 0xb6, 0x3f, 0xca, 0xc2, 0x70, 0x63, 0x84, 0x92, 0x85, 0xb1, 0xd4, 0xf2, 0x7c, 0xab, 0x95,
	0x55, 0xd8, 0x0b, 0xe2, 0xe5, 0xbd, 0xb4, 0x36, 0x3f, 0x25, 0x2b, 0x34, 0x67, 0x5c, 0xec, 0x03,
	0x7d, 0xf9, 0xf6, 0xc1, 0x57, 0xd6, 0x67, 0xc0, 0xee, 0x4f, 0xea, 0xae, 0x62, 0x4a, 0x7e, 0xaf,
	0x04, 0xb3, 0x01, 0x9a, 0x3f, 0x08, 0xda, 0x6a, 0xb5, 0xaa, 0xf8, 0xa6, 0x66, 0x9e, 0xfd, 0xee,
	0x40, 0x77, 0xe9, 0xb6, 0xc5, 0x2b, 0x68, 0x36, 0xcd, 0xa2, 0x69, 0x34, 0x1b, 0xce, 0xd3, 0xbc,
	0xc9, 0xb4, 0xac, 0x4f, 0xd2, 0x18, 0xf0, 0x25, 0xc4, 0xae, 0xa2, 0x9d, 0xc9, 0x4e, 0xa6, 0xd1,
	0x6c, 0x3c, 0x9f, 0x1c, 0xad, 0x5e, 0x26, 0xe1, 0x33, 0xbe, 0x80, 0x31, 0xad, 0xb6, 0x6b, 0x7a,
	0xad, 0xa4, 0xe4, 0xcc, 0xf2, 0xbb, 0xac, 0x33, 0x8d, 0x66, 0x7d, 0xf2, 0x8f, 0x8a, 0xcf, 0x21,
	0x59, 0x71, 0xfb, 0x97, 0xed, 0xd4, 0xdb, 0xda, 0x22, 0x3e, 0x85, 0x41, 0x45, 0x8d, 0xfd, 0xa8,
	0xb5, 0xd2, 0x59, 0x77, 0x1a, 0xcd, 0x06, 0xe4, 0x28, 0xe0, 0x3b, 0x48, 0x8c, 0xa5, 0xa5, 0x90,
	0x65, 0x1d, 0x22, 0x8b, 0xfd, 0x1a, 0x17, 0xb9, 0x7f, 0x97, 0x9c, 0xd4, 0xef, 0xa2, 0x43, 0xc2,
	0xb6, 0x17, 0x5f, 0xc1, 0xa0, 0x69, 0xd3, 0x64, 0xbd, 0x69, 0x67, 0x36, 0x9c, 0x8f, 0x73, 0xcf,
	0x5d, 0x07, 0x99, 0x1c, 0x0d, 0xb8, 0x80, 0xc4, 0x2f, 0xf0, 0xb5, 0xee, 0xd2, 0x64, 0x7d, 0x4f,
	0x24, 0x35, 0x11, 0x54, 0xd2, 0xf6, 0xe0, 0x1b, 0x18, 0xb9, 0x75, 0x0e, 0xcc, 0xe0, 0x7f, 0x4c,
	0xcb, 0x72, 0xf5, 0xeb, 0x04, 0xe2, 0x10, 0x30, 0x85, 0xd1, 0x8d, 0x30, 0xac, 0xe9, 0x22, 0x7d,
	0x84, 0x19, 0x9c, 0x7f, 0xa6, 0x95, 0xe5, 0x77, 0x5f, 0x24, 0x51, 0xca, 0xde, 0xf0, 0x8a, 0x5b,
	0xa1, 0x64, 0x1a, 0xe1, 0x63, 0x98, 0x84, 0xd2, 0x84, 0x2c, 0xdf, 0xbb, 0x10, 0xe9, 0x09, 0x22,
	0x8c, 0x8f, 0xe2, 0x07, 0x6e, 0x69, 0xda, 0xc1, 0x11, 0xf4, 0x6f, 0xa9, 0x65, 0x6b, 0x21, 0xcb,
	0xf4, 0x14, 0xcf, 0x20, 0x59, 0x32, 0x2a, 0xe5, 0x01, 0xea, 0xba, 0x5b, 0x1b, 0xc9, 0x23, 0x31,
	0x9e, 0x43, 0x7a, 0x4b, 0x85, 0x9b, 0xf1, 0x49, 0x69, 0xc2, 0x0d, 0xa3, 0x32, 0xed, 0xe1, 0x04,
	0x86, 0x84, 0x33, 0x25, 0x99, 0xa8, 0xdc, 0xac, 0xbe, 0x07, 0xeb, 0x82, 0xeb, 0x51, 0x03, 0x67,
	0x09, 0x8a, 0x9f, 0x04, 0xf8, 0x04, 0xf0, 0x9b, 0xa6, 0xd2, 0x08, 0x97, 0xfa, 0x60, 0x1c, 0xe2,
	0x05, 0x9c, 0xb5, 0x74, 0x6f, 0x1f, 0x21, 0x40, 0xbc, 0xa4, 0x7b, 0x37, 0x3d, 0x59, 0xc5, 0xfe,
	0xe7, 0x5d, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x01, 0xaa, 0x22, 0x6d, 0xd1, 0x03, 0x00, 0x00,
}
