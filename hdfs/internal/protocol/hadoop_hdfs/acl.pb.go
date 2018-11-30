// Code generated by protoc-gen-go.
// source: acl.proto
// DO NOT EDIT!

package hadoop_hdfs

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type AclEntryProto_AclEntryScopeProto int32

const (
	AclEntryProto_ACCESS  AclEntryProto_AclEntryScopeProto = 0
	AclEntryProto_DEFAULT AclEntryProto_AclEntryScopeProto = 1
)

var AclEntryProto_AclEntryScopeProto_name = map[int32]string{
	0: "ACCESS",
	1: "DEFAULT",
}
var AclEntryProto_AclEntryScopeProto_value = map[string]int32{
	"ACCESS":  0,
	"DEFAULT": 1,
}

func (x AclEntryProto_AclEntryScopeProto) Enum() *AclEntryProto_AclEntryScopeProto {
	p := new(AclEntryProto_AclEntryScopeProto)
	*p = x
	return p
}
func (x AclEntryProto_AclEntryScopeProto) String() string {
	return proto.EnumName(AclEntryProto_AclEntryScopeProto_name, int32(x))
}
func (x *AclEntryProto_AclEntryScopeProto) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AclEntryProto_AclEntryScopeProto_value, data, "AclEntryProto_AclEntryScopeProto")
	if err != nil {
		return err
	}
	*x = AclEntryProto_AclEntryScopeProto(value)
	return nil
}

type AclEntryProto_AclEntryTypeProto int32

const (
	AclEntryProto_USER  AclEntryProto_AclEntryTypeProto = 0
	AclEntryProto_GROUP AclEntryProto_AclEntryTypeProto = 1
	AclEntryProto_MASK  AclEntryProto_AclEntryTypeProto = 2
	AclEntryProto_OTHER AclEntryProto_AclEntryTypeProto = 3
)

var AclEntryProto_AclEntryTypeProto_name = map[int32]string{
	0: "USER",
	1: "GROUP",
	2: "MASK",
	3: "OTHER",
}
var AclEntryProto_AclEntryTypeProto_value = map[string]int32{
	"USER":  0,
	"GROUP": 1,
	"MASK":  2,
	"OTHER": 3,
}

func (x AclEntryProto_AclEntryTypeProto) Enum() *AclEntryProto_AclEntryTypeProto {
	p := new(AclEntryProto_AclEntryTypeProto)
	*p = x
	return p
}
func (x AclEntryProto_AclEntryTypeProto) String() string {
	return proto.EnumName(AclEntryProto_AclEntryTypeProto_name, int32(x))
}
func (x *AclEntryProto_AclEntryTypeProto) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AclEntryProto_AclEntryTypeProto_value, data, "AclEntryProto_AclEntryTypeProto")
	if err != nil {
		return err
	}
	*x = AclEntryProto_AclEntryTypeProto(value)
	return nil
}

type AclEntryProto_FsActionProto int32

const (
	AclEntryProto_NONE          AclEntryProto_FsActionProto = 0
	AclEntryProto_EXECUTE       AclEntryProto_FsActionProto = 1
	AclEntryProto_WRITE         AclEntryProto_FsActionProto = 2
	AclEntryProto_WRITE_EXECUTE AclEntryProto_FsActionProto = 3
	AclEntryProto_READ          AclEntryProto_FsActionProto = 4
	AclEntryProto_READ_EXECUTE  AclEntryProto_FsActionProto = 5
	AclEntryProto_READ_WRITE    AclEntryProto_FsActionProto = 6
	AclEntryProto_PERM_ALL      AclEntryProto_FsActionProto = 7
)

var AclEntryProto_FsActionProto_name = map[int32]string{
	0: "NONE",
	1: "EXECUTE",
	2: "WRITE",
	3: "WRITE_EXECUTE",
	4: "READ",
	5: "READ_EXECUTE",
	6: "READ_WRITE",
	7: "PERM_ALL",
}
var AclEntryProto_FsActionProto_value = map[string]int32{
	"NONE":          0,
	"EXECUTE":       1,
	"WRITE":         2,
	"WRITE_EXECUTE": 3,
	"READ":          4,
	"READ_EXECUTE":  5,
	"READ_WRITE":    6,
	"PERM_ALL":      7,
}

func (x AclEntryProto_FsActionProto) Enum() *AclEntryProto_FsActionProto {
	p := new(AclEntryProto_FsActionProto)
	*p = x
	return p
}
func (x AclEntryProto_FsActionProto) String() string {
	return proto.EnumName(AclEntryProto_FsActionProto_name, int32(x))
}
func (x *AclEntryProto_FsActionProto) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AclEntryProto_FsActionProto_value, data, "AclEntryProto_FsActionProto")
	if err != nil {
		return err
	}
	*x = AclEntryProto_FsActionProto(value)
	return nil
}

type AclEntryProto struct {
	Type             *AclEntryProto_AclEntryTypeProto  `protobuf:"varint,1,req,name=type,enum=hadoop.hdfs.AclEntryProto_AclEntryTypeProto" json:"type,omitempty"`
	Scope            *AclEntryProto_AclEntryScopeProto `protobuf:"varint,2,req,name=scope,enum=hadoop.hdfs.AclEntryProto_AclEntryScopeProto" json:"scope,omitempty"`
	Permissions      *AclEntryProto_FsActionProto      `protobuf:"varint,3,req,name=permissions,enum=hadoop.hdfs.AclEntryProto_FsActionProto" json:"permissions,omitempty"`
	Name             *string                           `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte                            `json:"-"`
}

func (m *AclEntryProto) Reset()         { *m = AclEntryProto{} }
func (m *AclEntryProto) String() string { return proto.CompactTextString(m) }
func (*AclEntryProto) ProtoMessage()    {}

func (m *AclEntryProto) GetType() AclEntryProto_AclEntryTypeProto {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return AclEntryProto_USER
}

func (m *AclEntryProto) GetScope() AclEntryProto_AclEntryScopeProto {
	if m != nil && m.Scope != nil {
		return *m.Scope
	}
	return AclEntryProto_ACCESS
}

func (m *AclEntryProto) GetPermissions() AclEntryProto_FsActionProto {
	if m != nil && m.Permissions != nil {
		return *m.Permissions
	}
	return AclEntryProto_NONE
}

func (m *AclEntryProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type AclStatusProto struct {
	Owner            *string            `protobuf:"bytes,1,req,name=owner" json:"owner,omitempty"`
	Group            *string            `protobuf:"bytes,2,req,name=group" json:"group,omitempty"`
	Sticky           *bool              `protobuf:"varint,3,req,name=sticky" json:"sticky,omitempty"`
	Entries          []*AclEntryProto   `protobuf:"bytes,4,rep,name=entries" json:"entries,omitempty"`
	Permission       *FsPermissionProto `protobuf:"bytes,5,opt,name=permission" json:"permission,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *AclStatusProto) Reset()         { *m = AclStatusProto{} }
func (m *AclStatusProto) String() string { return proto.CompactTextString(m) }
func (*AclStatusProto) ProtoMessage()    {}

func (m *AclStatusProto) GetOwner() string {
	if m != nil && m.Owner != nil {
		return *m.Owner
	}
	return ""
}

func (m *AclStatusProto) GetGroup() string {
	if m != nil && m.Group != nil {
		return *m.Group
	}
	return ""
}

func (m *AclStatusProto) GetSticky() bool {
	if m != nil && m.Sticky != nil {
		return *m.Sticky
	}
	return false
}

func (m *AclStatusProto) GetEntries() []*AclEntryProto {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AclStatusProto) GetPermission() *FsPermissionProto {
	if m != nil {
		return m.Permission
	}
	return nil
}

type ModifyAclEntriesRequestProto struct {
	Src              *string          `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	AclSpec          []*AclEntryProto `protobuf:"bytes,2,rep,name=aclSpec" json:"aclSpec,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *ModifyAclEntriesRequestProto) Reset()         { *m = ModifyAclEntriesRequestProto{} }
func (m *ModifyAclEntriesRequestProto) String() string { return proto.CompactTextString(m) }
func (*ModifyAclEntriesRequestProto) ProtoMessage()    {}

func (m *ModifyAclEntriesRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

func (m *ModifyAclEntriesRequestProto) GetAclSpec() []*AclEntryProto {
	if m != nil {
		return m.AclSpec
	}
	return nil
}

type ModifyAclEntriesResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ModifyAclEntriesResponseProto) Reset()         { *m = ModifyAclEntriesResponseProto{} }
func (m *ModifyAclEntriesResponseProto) String() string { return proto.CompactTextString(m) }
func (*ModifyAclEntriesResponseProto) ProtoMessage()    {}

type RemoveAclRequestProto struct {
	Src              *string `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RemoveAclRequestProto) Reset()         { *m = RemoveAclRequestProto{} }
func (m *RemoveAclRequestProto) String() string { return proto.CompactTextString(m) }
func (*RemoveAclRequestProto) ProtoMessage()    {}

func (m *RemoveAclRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

type RemoveAclResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RemoveAclResponseProto) Reset()         { *m = RemoveAclResponseProto{} }
func (m *RemoveAclResponseProto) String() string { return proto.CompactTextString(m) }
func (*RemoveAclResponseProto) ProtoMessage()    {}

type RemoveAclEntriesRequestProto struct {
	Src              *string          `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	AclSpec          []*AclEntryProto `protobuf:"bytes,2,rep,name=aclSpec" json:"aclSpec,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *RemoveAclEntriesRequestProto) Reset()         { *m = RemoveAclEntriesRequestProto{} }
func (m *RemoveAclEntriesRequestProto) String() string { return proto.CompactTextString(m) }
func (*RemoveAclEntriesRequestProto) ProtoMessage()    {}

func (m *RemoveAclEntriesRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

func (m *RemoveAclEntriesRequestProto) GetAclSpec() []*AclEntryProto {
	if m != nil {
		return m.AclSpec
	}
	return nil
}

type RemoveAclEntriesResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RemoveAclEntriesResponseProto) Reset()         { *m = RemoveAclEntriesResponseProto{} }
func (m *RemoveAclEntriesResponseProto) String() string { return proto.CompactTextString(m) }
func (*RemoveAclEntriesResponseProto) ProtoMessage()    {}

type RemoveDefaultAclRequestProto struct {
	Src              *string `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RemoveDefaultAclRequestProto) Reset()         { *m = RemoveDefaultAclRequestProto{} }
func (m *RemoveDefaultAclRequestProto) String() string { return proto.CompactTextString(m) }
func (*RemoveDefaultAclRequestProto) ProtoMessage()    {}

func (m *RemoveDefaultAclRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

type RemoveDefaultAclResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RemoveDefaultAclResponseProto) Reset()         { *m = RemoveDefaultAclResponseProto{} }
func (m *RemoveDefaultAclResponseProto) String() string { return proto.CompactTextString(m) }
func (*RemoveDefaultAclResponseProto) ProtoMessage()    {}

type SetAclRequestProto struct {
	Src              *string          `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	AclSpec          []*AclEntryProto `protobuf:"bytes,2,rep,name=aclSpec" json:"aclSpec,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *SetAclRequestProto) Reset()         { *m = SetAclRequestProto{} }
func (m *SetAclRequestProto) String() string { return proto.CompactTextString(m) }
func (*SetAclRequestProto) ProtoMessage()    {}

func (m *SetAclRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

func (m *SetAclRequestProto) GetAclSpec() []*AclEntryProto {
	if m != nil {
		return m.AclSpec
	}
	return nil
}

type SetAclResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *SetAclResponseProto) Reset()         { *m = SetAclResponseProto{} }
func (m *SetAclResponseProto) String() string { return proto.CompactTextString(m) }
func (*SetAclResponseProto) ProtoMessage()    {}

type GetAclStatusRequestProto struct {
	Src              *string `protobuf:"bytes,1,req,name=src" json:"src,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GetAclStatusRequestProto) Reset()         { *m = GetAclStatusRequestProto{} }
func (m *GetAclStatusRequestProto) String() string { return proto.CompactTextString(m) }
func (*GetAclStatusRequestProto) ProtoMessage()    {}

func (m *GetAclStatusRequestProto) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

type GetAclStatusResponseProto struct {
	Result           *AclStatusProto `protobuf:"bytes,1,req,name=result" json:"result,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *GetAclStatusResponseProto) Reset()         { *m = GetAclStatusResponseProto{} }
func (m *GetAclStatusResponseProto) String() string { return proto.CompactTextString(m) }
func (*GetAclStatusResponseProto) ProtoMessage()    {}

func (m *GetAclStatusResponseProto) GetResult() *AclStatusProto {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterEnum("hadoop.hdfs.AclEntryProto_AclEntryScopeProto", AclEntryProto_AclEntryScopeProto_name, AclEntryProto_AclEntryScopeProto_value)
	proto.RegisterEnum("hadoop.hdfs.AclEntryProto_AclEntryTypeProto", AclEntryProto_AclEntryTypeProto_name, AclEntryProto_AclEntryTypeProto_value)
	proto.RegisterEnum("hadoop.hdfs.AclEntryProto_FsActionProto", AclEntryProto_FsActionProto_name, AclEntryProto_FsActionProto_value)
}
