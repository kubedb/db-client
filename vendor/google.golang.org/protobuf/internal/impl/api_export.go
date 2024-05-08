// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"reflect"
	"strconv"

	"fmt"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// Export is a zero-length named type that exists only to export a set of
// functions that we do not want to appear in godoc.
type Export struct{}

// NewError formats a string according to the format specifier and arguments and
// returns an error that has a "proto" prefix.
func (Export) NewError(f string, x ...interface{}) error {
	return errors.New(f, x...)
}

// enum is any enum type generated by protoc-gen-go
// and must be a named int32 type.
type enum = interface{}

// EnumOf returns the protoreflect.Enum interface over e.
// It returns nil if e is nil.
func (Export) EnumOf(e enum) protoreflect.Enum {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e
	default:
		return legacyWrapEnum(reflect.ValueOf(e))
	}
}

// EnumDescriptorOf returns the protoreflect.EnumDescriptor for e.
// It returns nil if e is nil.
func (Export) EnumDescriptorOf(e enum) protoreflect.EnumDescriptor {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e.Descriptor()
	default:
		return LegacyLoadEnumDesc(reflect.TypeOf(e))
	}
}

// EnumTypeOf returns the protoreflect.EnumType for e.
// It returns nil if e is nil.
func (Export) EnumTypeOf(e enum) protoreflect.EnumType {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e.Type()
	default:
		return legacyLoadEnumType(reflect.TypeOf(e))
	}
}

// EnumStringOf returns the enum value as a string, either as the name if
// the number is resolvable, or the number formatted as a string.
func (Export) EnumStringOf(ed protoreflect.EnumDescriptor, n protoreflect.EnumNumber) string {
	ev := ed.Values().ByNumber(n)
	if ev != nil {
		return string(ev.Name())
	}
	return strconv.Itoa(int(n))
}

// message is any message type generated by protoc-gen-go
// and must be a pointer to a named struct type.
type message = interface{}

// legacyMessageWrapper wraps a v2 message as a v1 message.
type legacyMessageWrapper struct{ m protoreflect.ProtoMessage }

func (m legacyMessageWrapper) Reset()         { proto.Reset(m.m) }
func (m legacyMessageWrapper) String() string { return Export{}.MessageStringOf(m.m) }
func (m legacyMessageWrapper) ProtoMessage()  {}

// ProtoMessageV1Of converts either a v1 or v2 message to a v1 message.
// It returns nil if m is nil.
func (Export) ProtoMessageV1Of(m message) protoiface.MessageV1 {
	switch mv := m.(type) {
	case nil:
		return nil
	case protoiface.MessageV1:
		return mv
	case unwrapper:
		return Export{}.ProtoMessageV1Of(mv.protoUnwrap())
	case protoreflect.ProtoMessage:
		return legacyMessageWrapper{mv}
	default:
		panic(fmt.Sprintf("message %T is neither a v1 or v2 Message", m))
	}
}

func (Export) protoMessageV2Of(m message) protoreflect.ProtoMessage {
	switch mv := m.(type) {
	case nil:
		return nil
	case protoreflect.ProtoMessage:
		return mv
	case legacyMessageWrapper:
		return mv.m
	case protoiface.MessageV1:
		return nil
	default:
		panic(fmt.Sprintf("message %T is neither a v1 or v2 Message", m))
	}
}

// ProtoMessageV2Of converts either a v1 or v2 message to a v2 message.
// It returns nil if m is nil.
func (Export) ProtoMessageV2Of(m message) protoreflect.ProtoMessage {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv
	}
	return legacyWrapMessage(reflect.ValueOf(m)).Interface()
}

// MessageOf returns the protoreflect.Message interface over m.
// It returns nil if m is nil.
func (Export) MessageOf(m message) protoreflect.Message {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect()
	}
	return legacyWrapMessage(reflect.ValueOf(m))
}

// MessageDescriptorOf returns the protoreflect.MessageDescriptor for m.
// It returns nil if m is nil.
func (Export) MessageDescriptorOf(m message) protoreflect.MessageDescriptor {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect().Descriptor()
	}
	return LegacyLoadMessageDesc(reflect.TypeOf(m))
}

// MessageTypeOf returns the protoreflect.MessageType for m.
// It returns nil if m is nil.
func (Export) MessageTypeOf(m message) protoreflect.MessageType {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect().Type()
	}
	return legacyLoadMessageType(reflect.TypeOf(m), "")
}

// MessageStringOf returns the message value as a string,
// which is the message serialized in the protobuf text format.
func (Export) MessageStringOf(m protoreflect.ProtoMessage) string {
	return prototext.MarshalOptions{Multiline: false}.Format(m)
}
