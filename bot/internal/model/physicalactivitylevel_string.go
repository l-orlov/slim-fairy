// Code generated by "stringer -type=PhysicalActivityLevel"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PhysicalActivityLevelLow-1]
	_ = x[PhysicalActivityLevelMedium-2]
	_ = x[PhysicalActivityLevelHigh-3]
}

const _PhysicalActivityLevel_name = "PhysicalActivityLevelLowPhysicalActivityLevelMediumPhysicalActivityLevelHigh"

var _PhysicalActivityLevel_index = [...]uint8{0, 24, 51, 76}

func (i PhysicalActivityLevel) String() string {
	i -= 1
	if i < 0 || i >= PhysicalActivityLevel(len(_PhysicalActivityLevel_index)-1) {
		return "PhysicalActivityLevel(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _PhysicalActivityLevel_name[_PhysicalActivityLevel_index[i]:_PhysicalActivityLevel_index[i+1]]
}
