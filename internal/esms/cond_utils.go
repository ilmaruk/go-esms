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
