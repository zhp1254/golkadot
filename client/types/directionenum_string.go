// Code generated by "stringer -type=directionEnum"; DO NOT EDIT.

package clienttypes

import "strconv"

const _directionEnum_name = "AscendingDescending"

var _directionEnum_index = [...]uint8{0, 9, 19}

func (i directionEnum) String() string {
	if i < 0 || i >= directionEnum(len(_directionEnum_index)-1) {
		return "directionEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _directionEnum_name[_directionEnum_index[i]:_directionEnum_index[i+1]]
}
