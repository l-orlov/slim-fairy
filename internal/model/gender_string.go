// Code generated by "stringer -type=Gender"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GenderMan-1]
	_ = x[GenderWoman-2]
}

const _Gender_name = "GenderManGenderWoman"

var _Gender_index = [...]uint8{0, 9, 20}

func (i Gender) String() string {
	i -= 1
	if i < 0 || i >= Gender(len(_Gender_index)-1) {
		return "Gender(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Gender_name[_Gender_index[i]:_Gender_index[i+1]]
}
