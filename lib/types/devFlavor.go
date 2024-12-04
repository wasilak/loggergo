package types

// DevFlavor represents the flavor of the development environment.
type DevFlavor int

const (
	// DevFlavorTint represents the "tint" development flavor.
	DevFlavorTint DevFlavor = iota
	// DevFlavorSlogor represents the "slogor" development flavor.
	DevFlavorSlogor
	// DevFlavorDevslog represents the production "devslog" flavor.
	DevFlavorDevslog
)

func (f DevFlavor) String() string {
	switch f {
	case DevFlavorTint:
		return "tint"
	case DevFlavorSlogor:
		return "slogor"
	case DevFlavorDevslog:
		return "devslog"
	default:
		return "unknown"
	}
}
