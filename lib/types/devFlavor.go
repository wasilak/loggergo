package types

import (
	"fmt"
	"log/slog"

	"github.com/xybor-x/enum"
)

// DevFlavor represents the flavor of the development environment.
type devFlavor int
type DevFlavor struct{ enum.SafeEnum[devFlavor] }

var (
	DevFlavorTint    = enum.NewExtended[DevFlavor]("tint")
	DevFlavorSlogor  = enum.NewExtended[DevFlavor]("slogor")
	DevFlavorDevslog = enum.NewExtended[DevFlavor]("devslog")

	_ = enum.Finalize[DevFlavor]() // still required internally
)

// AllDevFlavors returns all defined DevFlavor values.
func AllDevFlavors() []DevFlavor {
	return enum.All[DevFlavor]()
}

// DevFlavorFromString parses a string to a DevFlavor, returning a fallback if not found.
func DevFlavorFromString(name string) DevFlavor {
	if v, ok := enum.FromString[DevFlavor](name); ok {
		return v
	}
	slog.Warn(fmt.Sprintf("Invalid dev flavor: %q, defaulting to %s", name, DevFlavorTint))
	return DevFlavorTint
}
