// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package interfaces_test

import (
	"fmt"

	. "gopkg.in/check.v1"

	. "github.com/ubuntu-core/snappy/interfaces"
	"github.com/ubuntu-core/snappy/snap"
)

type TestInterfaceSuite struct {
	iface Interface
	plug  *Plug
	slot  *Slot
}

var _ = Suite(&TestInterfaceSuite{
	iface: &TestInterface{InterfaceName: "test"},
	plug: &Plug{
		PlugInfo: &snap.PlugInfo{
			Snap:      &snap.Info{Name: "snap"},
			Name:      "name",
			Interface: "test",
		},
	},
	slot: &Slot{
		SlotInfo: &snap.SlotInfo{
			Snap:      &snap.Info{Name: "snap"},
			Name:      "name",
			Interface: "test",
		},
	},
})

// TestInterface has a working Name() function
func (s *TestInterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "test")
}

// TestInterface doesn't do any sanitization by default
func (s *TestInterfaceSuite) TestSanitizePlugOK(c *C) {
	err := s.iface.SanitizePlug(s.plug)
	c.Assert(err, IsNil)
}

// TestInterface has provisions to customize sanitization
func (s *TestInterfaceSuite) TestSanitizePlugError(c *C) {
	iface := &TestInterface{
		InterfaceName: "test",
		SanitizePlugCallback: func(plug *Plug) error {
			return fmt.Errorf("sanitize plug failed")
		},
	}
	err := iface.SanitizePlug(s.plug)
	c.Assert(err, ErrorMatches, "sanitize plug failed")
}

// TestInterface sanitization still checks for interface identity
func (s *TestInterfaceSuite) TestSanitizePlugWrongInterface(c *C) {
	plug := &Plug{
		PlugInfo: &snap.PlugInfo{
			Snap:      &snap.Info{Name: "snap"},
			Name:      "name",
			Interface: "other-interface",
		},
	}
	c.Assert(func() { s.iface.SanitizePlug(plug) }, Panics, "plug is not of interface \"test\"")
}

// TestInterface doesn't do any sanitization by default
func (s *TestInterfaceSuite) TestSanitizeSlotOK(c *C) {
	err := s.iface.SanitizeSlot(s.slot)
	c.Assert(err, IsNil)
}

// TestInterface has provisions to customize sanitization
func (s *TestInterfaceSuite) TestSanitizeSlotError(c *C) {
	iface := &TestInterface{
		InterfaceName: "test",
		SanitizeSlotCallback: func(slot *Slot) error {
			return fmt.Errorf("sanitize slot failed")
		},
	}
	err := iface.SanitizeSlot(s.slot)
	c.Assert(err, ErrorMatches, "sanitize slot failed")
}

// TestInterface sanitization still checks for interface identity
func (s *TestInterfaceSuite) TestSanitizeSlotWrongInterface(c *C) {
	slot := &Slot{
		SlotInfo: &snap.SlotInfo{
			Snap:      &snap.Info{Name: "snap"},
			Name:      "name",
			Interface: "interface",
		},
	}
	c.Assert(func() { s.iface.SanitizeSlot(slot) }, Panics, "slot is not of interface \"test\"")
}

// TestInterface hands out empty plug security snippets
func (s *TestInterfaceSuite) TestPlugSnippet(c *C) {
	snippet, err := s.iface.ConnectedPlugSnippet(s.plug, s.slot, SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedPlugSnippet(s.plug, s.slot, SecuritySecComp)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedPlugSnippet(s.plug, s.slot, SecurityDBus)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedPlugSnippet(s.plug, s.slot, "foo")
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
}

// TestInterface hands out empty slot security snippets
func (s *TestInterfaceSuite) TestSlotSnippet(c *C) {
	snippet, err := s.iface.ConnectedSlotSnippet(s.plug, s.slot, SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedSlotSnippet(s.plug, s.slot, SecuritySecComp)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedSlotSnippet(s.plug, s.slot, SecurityDBus)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedSlotSnippet(s.plug, s.slot, "foo")
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
}

func (s *TestInterfaceSuite) TestAutoConnect(c *C) {
	c.Assert(s.iface.AutoConnect(), Equals, false)
	iface := &TestInterface{AutoConnectCallback: func() bool { return true }}
	c.Assert(iface.AutoConnect(), Equals, true)
}
