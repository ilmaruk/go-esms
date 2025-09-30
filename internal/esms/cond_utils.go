package esms

// Given a full position (like DML), get only
// the position (DM)
func fullpos2position(fullpos string) string {
	// assert(fullpos.size() == 3);
	return fullpos[0:2]
}

// Given full position (like DML), get only
// the side (L)
func fullpos2side(fullpos string) string {
	// assert(fullpos.size() == 3);
	return fullpos[2:3]
}

// Given a position (DM) and a side (L), returns the
// full position (DML)
func posAndSide2fullpos(pos, side string) string {
	fullpos := pos
	fullpos += side
	return fullpos
}

// Position: 3 letters (DFL, AMC, etc.) or GK
func is_legal_position(position string) bool {
	if position == "GK" {
		return true
	}

	if len(position) != 3 {
		return false
	}

	raw_position := position[0:2]
	side := position[2:3]

	return tact_manager.position_exists(raw_position) && is_legal_side(side)
}

func is_legal_side(side string) bool {
	return side == "L" || side == "R" || side == "C"
}
