// ovirt-imageio
// Copyright (C) 2021 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package nbd

import (
	"encoding/json"
	"reflect"
	"testing"

	"ovirt.org/imageio"
	"ovirt.org/imageio/units"
)

func TestNbdSize(t *testing.T) {
	// TODO: Requires empty 6g image exposed via qemu-nbd
	b, err := Connect("nbd+unix://?socket=/tmp/nbd.sock")
	if err != nil {
		t.Fatalf("Connect failed: %s", err)
	}
	defer b.Close()

	imageSize := 6 * units.GiB
	size, err := b.Size()
	if err != nil {
		t.Fatalf("Size() failed: %s", err)
	} else if size != imageSize {
		t.Errorf("backend.Size() = %v, expected %v", size, imageSize)
	}
}

func TestNbdExtents(t *testing.T) {
	// TODO: Requires empty 6g image exposed via qemu-nbd
	b, err := Connect("nbd+unix://?socket=/tmp/nbd.sock")
	if err != nil {
		t.Fatalf("Connect failed: %s", err)
	}
	defer b.Close()

	res, err := b.Extents()
	if err != nil {
		t.Fatalf("Extents() failed: %s", err)
	}

	var extents []*imageio.Extent
	for res.Next() {
		extents = append(extents, res.Value())
	}

	expected := []*imageio.Extent{
		{Start: 0 * units.GiB, Length: 6 * units.GiB, Zero: true},
	}
	if !reflect.DeepEqual(extents, expected) {
		t.Fatalf("extents:\n%s\nexpected:\n%s\n", dump(extents), dump(expected))
	}
}

func dump(v interface{}) []byte {
	res, _ := json.MarshalIndent(v, "", " ")
	return res
}
