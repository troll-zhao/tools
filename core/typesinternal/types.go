// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package typesinternal provides access to core go/types APIs that are not
// yet exported.
package typesinternal

import (
	"go/token"
	"go/types"
	"reflect"
	"unsafe"

	"github.com/troll-zhao/tools/core/aliases"
)

func SetUsesCgo(conf *types.Config) bool {
	v := reflect.ValueOf(conf).Elem()

	f := v.FieldByName("go115UsesCgo")
	if !f.IsValid() {
		f = v.FieldByName("UsesCgo")
		if !f.IsValid() {
			return false
		}
	}

	addr := unsafe.Pointer(f.UnsafeAddr())
	*(*bool)(addr) = true

	return true
}

// ReadGo116ErrorData extracts additional information from types.Error values
// generated by Go version 1.16 and later: the error code, start position, and
// end position. If all positions are valid, start <= err.Pos <= end.
//
// If the data could not be read, the final result parameter will be false.
func ReadGo116ErrorData(err types.Error) (code ErrorCode, start, end token.Pos, ok bool) {
	var data [3]int
	// By coincidence all of these fields are ints, which simplifies things.
	v := reflect.ValueOf(err)
	for i, name := range []string{"go116code", "go116start", "go116end"} {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return 0, 0, 0, false
		}
		data[i] = int(f.Int())
	}
	return ErrorCode(data[0]), token.Pos(data[1]), token.Pos(data[2]), true
}

// NameRelativeTo returns a types.Qualifier that qualifies members of
// all packages other than pkg, using only the package name.
// (By contrast, [types.RelativeTo] uses the complete package path,
// which is often excessive.)
//
// If pkg is nil, it is equivalent to [*types.Package.Name].
func NameRelativeTo(pkg *types.Package) types.Qualifier {
	return func(other *types.Package) string {
		if pkg != nil && pkg == other {
			return "" // same package; unqualified
		}
		return other.Name()
	}
}

// A NamedOrAlias is a [types.Type] that is named (as
// defined by the spec) and capable of bearing type parameters: it
// abstracts aliases ([types.Alias]) and defined types
// ([types.Named]).
//
// Every type declared by an explicit "type" declaration is a
// NamedOrAlias. (Built-in type symbols may additionally
// have type [types.Basic], which is not a NamedOrAlias,
// though the spec regards them as "named".)
//
// NamedOrAlias cannot expose the Origin method, because
// [types.Alias.Origin] and [types.Named.Origin] have different
// (covariant) result types; use [Origin] instead.
type NamedOrAlias interface {
	types.Type
	Obj() *types.TypeName
}

// TypeParams is a light shim around t.TypeParams().
// (go/types.Alias).TypeParams requires >= 1.23.
func TypeParams(t NamedOrAlias) *types.TypeParamList {
	switch t := t.(type) {
	case *types.Alias:
		return aliases.TypeParams(t)
	case *types.Named:
		return t.TypeParams()
	}
	return nil
}

// TypeArgs is a light shim around t.TypeArgs().
// (go/types.Alias).TypeArgs requires >= 1.23.
func TypeArgs(t NamedOrAlias) *types.TypeList {
	switch t := t.(type) {
	case *types.Alias:
		return aliases.TypeArgs(t)
	case *types.Named:
		return t.TypeArgs()
	}
	return nil
}

// Origin returns the generic type of the Named or Alias type t if it
// is instantiated, otherwise it returns t.
func Origin(t NamedOrAlias) NamedOrAlias {
	switch t := t.(type) {
	case *types.Alias:
		return aliases.Origin(t)
	case *types.Named:
		return t.Origin()
	}
	return t
}