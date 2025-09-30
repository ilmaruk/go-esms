package esms

type TacticManager struct{}

func (tm *TacticManager) get_mult(tactic_a, tactic_b, pos, skill string) float64 {
	// Dummy implementation, replace with actual logic
	return 1.0
}

func (tm *TacticManager) get_positions_names() []string {
	// Dummy implementation, replace with actual logic
	return []string{"GK", "DF", "MF", "FW"}
}

func (tm *TacticManager) tactic_exists(tactic string) bool {
	// Dummy implementation, replace with actual logic
	return true
}

func (tm *TacticManager) position_exists(position string) bool {
	// Dummy implementation, replace with actual logic
	return true
}
