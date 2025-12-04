package oh1

// Out represents a complete checkout sequence for a given score
type Out struct {
	Score   int
	Targets []DartTarget
}

// OutChart contains the official dart checkout chart
type OutChart struct {
	outs map[int]Out
}

// NewOutChart creates and initializes a new OutChart with all official checkouts
func NewOutChart() *OutChart {
	chart := &OutChart{
		outs: make(map[int]Out),
	}
	chart.initializeOuts()
	return chart
}

// GetOut returns the recommended checkout for a given score
// Returns nil if no checkout is available (score too high or odd number below 2)
func (oc *OutChart) GetOut(score int) *Out {
	if out, exists := oc.outs[score]; exists {
		return &out
	}
	return nil
}

// GetNextTarget returns the next target for a given score based on the out chart
// If no out exists, it returns a default target (Triple 20 for high scores, appropriate finishes for lower scores)
func (oc *OutChart) GetNextTarget(score int, preference ScoringPreference) DartTarget {
	if score < 2 {
		panic("Score less than 2, cannot checkout - something is busted")
	}

	// Check if we have a specific out for this score
	if out := oc.GetOut(score); out != nil && len(out.Targets) > 0 {
		return out.Targets[0]
	}

	// if we get here they are over 170
	switch preference {
	case TwentiesScoringPreference:
		// Aim for Triple 20
		return DartTarget{Multiplier: Triple, Number: Twenty}
	case NinteensScoringPreference:
		// Aim for Triple 19
		return DartTarget{Multiplier: Triple, Number: Nineteen}
	default:
		// Default to Triple 20
		return DartTarget{Multiplier: Triple, Number: Twenty}
	}
}

// initializeOuts populates the out chart with all standard checkouts
func (oc *OutChart) initializeOuts() {
	// 2-40: Basic doubles and single+double combinations
	oc.outs[2] = Out{2, []DartTarget{{Double, 1}}}
	oc.outs[3] = Out{3, []DartTarget{{Single, 1}, {Double, 1}}}
	oc.outs[4] = Out{4, []DartTarget{{Double, 2}}}
	oc.outs[5] = Out{5, []DartTarget{{Single, 1}, {Double, 2}}}
	oc.outs[6] = Out{6, []DartTarget{{Double, 3}}}
	oc.outs[7] = Out{7, []DartTarget{{Single, 3}, {Double, 2}}}
	oc.outs[8] = Out{8, []DartTarget{{Double, 4}}}
	oc.outs[9] = Out{9, []DartTarget{{Single, 1}, {Double, 4}}}
	oc.outs[10] = Out{10, []DartTarget{{Double, 5}}}
	oc.outs[11] = Out{11, []DartTarget{{Single, 3}, {Double, 4}}}
	oc.outs[12] = Out{12, []DartTarget{{Double, 6}}}
	oc.outs[13] = Out{13, []DartTarget{{Single, 5}, {Double, 4}}}
	oc.outs[14] = Out{14, []DartTarget{{Double, 7}}}
	oc.outs[15] = Out{15, []DartTarget{{Single, 7}, {Double, 4}}}
	oc.outs[16] = Out{16, []DartTarget{{Double, 8}}}
	oc.outs[17] = Out{17, []DartTarget{{Single, 9}, {Double, 4}}}
	oc.outs[18] = Out{18, []DartTarget{{Double, 9}}}
	oc.outs[19] = Out{19, []DartTarget{{Single, 3}, {Double, 8}}}
	oc.outs[20] = Out{20, []DartTarget{{Double, 10}}}
	oc.outs[21] = Out{21, []DartTarget{{Single, 5}, {Double, 8}}}
	oc.outs[22] = Out{22, []DartTarget{{Double, 11}}}
	oc.outs[23] = Out{23, []DartTarget{{Single, 7}, {Double, 8}}}
	oc.outs[24] = Out{24, []DartTarget{{Double, 12}}}
	oc.outs[25] = Out{25, []DartTarget{{Single, 9}, {Double, 8}}}
	oc.outs[26] = Out{26, []DartTarget{{Double, 13}}}
	oc.outs[27] = Out{27, []DartTarget{{Single, 11}, {Double, 8}}}
	oc.outs[28] = Out{28, []DartTarget{{Double, 14}}}
	oc.outs[29] = Out{29, []DartTarget{{Single, 13}, {Double, 8}}}
	oc.outs[30] = Out{30, []DartTarget{{Double, 15}}}
	oc.outs[31] = Out{31, []DartTarget{{Single, 15}, {Double, 8}}}
	oc.outs[32] = Out{32, []DartTarget{{Double, 16}}}
	oc.outs[33] = Out{33, []DartTarget{{Single, 17}, {Double, 8}}}
	oc.outs[34] = Out{34, []DartTarget{{Double, 17}}}
	oc.outs[35] = Out{35, []DartTarget{{Single, 3}, {Double, 16}}}
	oc.outs[36] = Out{36, []DartTarget{{Double, 18}}}
	oc.outs[37] = Out{37, []DartTarget{{Single, 5}, {Double, 16}}}
	oc.outs[38] = Out{38, []DartTarget{{Double, 19}}}
	oc.outs[39] = Out{39, []DartTarget{{Single, 7}, {Double, 16}}}
	oc.outs[40] = Out{40, []DartTarget{{Double, 20}}}

	// 41-60: Two-dart combinations
	oc.outs[41] = Out{41, []DartTarget{{Single, 9}, {Double, 16}}}
	oc.outs[42] = Out{42, []DartTarget{{Single, 10}, {Double, 16}}}
	oc.outs[43] = Out{43, []DartTarget{{Single, 11}, {Double, 16}}}
	oc.outs[44] = Out{44, []DartTarget{{Single, 12}, {Double, 16}}}
	oc.outs[45] = Out{45, []DartTarget{{Single, 13}, {Double, 16}}}
	oc.outs[46] = Out{46, []DartTarget{{Single, 6}, {Double, 20}}}
	oc.outs[47] = Out{47, []DartTarget{{Single, 7}, {Double, 20}}}
	oc.outs[48] = Out{48, []DartTarget{{Single, 8}, {Double, 20}}}
	oc.outs[49] = Out{49, []DartTarget{{Single, 17}, {Double, 16}}}
	oc.outs[50] = Out{50, []DartTarget{{Single, 18}, {Double, 16}}}
	oc.outs[51] = Out{51, []DartTarget{{Single, 19}, {Double, 16}}}
	oc.outs[52] = Out{52, []DartTarget{{Single, 12}, {Double, 20}}}
	oc.outs[53] = Out{53, []DartTarget{{Single, 13}, {Double, 20}}}
	oc.outs[54] = Out{54, []DartTarget{{Single, 14}, {Double, 20}}}
	oc.outs[55] = Out{55, []DartTarget{{Single, 15}, {Double, 20}}}
	oc.outs[56] = Out{56, []DartTarget{{Single, 16}, {Double, 20}}}
	oc.outs[57] = Out{57, []DartTarget{{Single, 17}, {Double, 20}}}
	oc.outs[58] = Out{58, []DartTarget{{Single, 18}, {Double, 20}}}
	oc.outs[59] = Out{59, []DartTarget{{Single, 19}, {Double, 20}}}
	oc.outs[60] = Out{60, []DartTarget{{Single, 20}, {Double, 20}}}

	// 61-100: Two-dart finishes (treble + double)
	oc.outs[61] = Out{61, []DartTarget{{Triple, 15}, {Double, 8}}}
	oc.outs[62] = Out{62, []DartTarget{{Triple, 10}, {Double, 16}}}
	oc.outs[63] = Out{63, []DartTarget{{Triple, 13}, {Double, 12}}}
	oc.outs[64] = Out{64, []DartTarget{{Triple, 16}, {Double, 8}}}
	oc.outs[65] = Out{65, []DartTarget{{Triple, 19}, {Double, 4}}}
	oc.outs[66] = Out{66, []DartTarget{{Triple, 10}, {Double, 18}}}
	oc.outs[67] = Out{67, []DartTarget{{Triple, 17}, {Double, 8}}}
	oc.outs[68] = Out{68, []DartTarget{{Triple, 20}, {Double, 4}}}
	oc.outs[69] = Out{69, []DartTarget{{Triple, 15}, {Double, 12}}}
	oc.outs[70] = Out{70, []DartTarget{{Triple, 10}, {Double, 20}}}
	oc.outs[71] = Out{71, []DartTarget{{Triple, 13}, {Double, 16}}}
	oc.outs[72] = Out{72, []DartTarget{{Triple, 16}, {Double, 12}}}
	oc.outs[73] = Out{73, []DartTarget{{Triple, 19}, {Double, 8}}}
	oc.outs[74] = Out{74, []DartTarget{{Triple, 14}, {Double, 16}}}
	oc.outs[75] = Out{75, []DartTarget{{Triple, 17}, {Double, 12}}}
	oc.outs[76] = Out{76, []DartTarget{{Triple, 20}, {Double, 8}}}
	oc.outs[77] = Out{77, []DartTarget{{Triple, 19}, {Double, 10}}}
	oc.outs[78] = Out{78, []DartTarget{{Triple, 18}, {Double, 12}}}
	oc.outs[79] = Out{79, []DartTarget{{Triple, 19}, {Double, 11}}}
	oc.outs[80] = Out{80, []DartTarget{{Triple, 20}, {Double, 10}}}
	oc.outs[81] = Out{81, []DartTarget{{Triple, 19}, {Double, 12}}}
	oc.outs[82] = Out{82, []DartTarget{{Triple, 14}, {Double, 20}}}
	oc.outs[83] = Out{83, []DartTarget{{Triple, 17}, {Double, 16}}}
	oc.outs[84] = Out{84, []DartTarget{{Triple, 20}, {Double, 12}}}
	oc.outs[85] = Out{85, []DartTarget{{Triple, 15}, {Double, 20}}}
	oc.outs[86] = Out{86, []DartTarget{{Triple, 18}, {Double, 16}}}
	oc.outs[87] = Out{87, []DartTarget{{Triple, 17}, {Double, 18}}}
	oc.outs[88] = Out{88, []DartTarget{{Triple, 16}, {Double, 20}}}
	oc.outs[89] = Out{89, []DartTarget{{Triple, 19}, {Double, 16}}}
	oc.outs[90] = Out{90, []DartTarget{{Triple, 20}, {Double, 15}}}
	oc.outs[91] = Out{91, []DartTarget{{Triple, 17}, {Double, 20}}}
	oc.outs[92] = Out{92, []DartTarget{{Triple, 20}, {Double, 16}}}
	oc.outs[93] = Out{93, []DartTarget{{Triple, 19}, {Double, 18}}}
	oc.outs[94] = Out{94, []DartTarget{{Triple, 18}, {Double, 20}}}
	oc.outs[95] = Out{95, []DartTarget{{Triple, 19}, {Double, 19}}}
	oc.outs[96] = Out{96, []DartTarget{{Triple, 20}, {Double, 18}}}
	oc.outs[97] = Out{97, []DartTarget{{Triple, 19}, {Double, 20}}}
	oc.outs[98] = Out{98, []DartTarget{{Triple, 20}, {Double, 19}}}
	oc.outs[99] = Out{99, []DartTarget{{Triple, 19}, {Single, 10}, {Double, 16}}}
	oc.outs[100] = Out{100, []DartTarget{{Triple, 20}, {Double, 20}}}

	// 101-130: Three-dart finishes
	oc.outs[101] = Out{101, []DartTarget{{Triple, 20}, {Single, 1}, {Double, 20}}}
	oc.outs[102] = Out{102, []DartTarget{{Triple, 20}, {Single, 10}, {Double, 16}}}
	oc.outs[103] = Out{103, []DartTarget{{Triple, 20}, {Single, 3}, {Double, 20}}}
	oc.outs[104] = Out{104, []DartTarget{{Triple, 18}, {Single, 18}, {Double, 16}}}
	oc.outs[105] = Out{105, []DartTarget{{Triple, 19}, {Single, 16}, {Double, 16}}}
	oc.outs[106] = Out{106, []DartTarget{{Triple, 20}, {Single, 14}, {Double, 16}}}
	oc.outs[107] = Out{107, []DartTarget{{Triple, 19}, {Single, 18}, {Double, 16}}}
	oc.outs[108] = Out{108, []DartTarget{{Triple, 20}, {Single, 16}, {Double, 16}}}
	oc.outs[109] = Out{109, []DartTarget{{Triple, 19}, {Single, 20}, {Double, 16}}}
	oc.outs[110] = Out{110, []DartTarget{{Triple, 20}, {Single, 18}, {Double, 16}}}

	// 111-130: Three-dart finishes (continued)
	oc.outs[111] = Out{111, []DartTarget{{Triple, 20}, {Single, 19}, {Double, 16}}}
	oc.outs[112] = Out{112, []DartTarget{{Triple, 20}, {Single, 12}, {Double, 20}}}
	oc.outs[113] = Out{113, []DartTarget{{Triple, 20}, {Single, 13}, {Double, 20}}}
	oc.outs[114] = Out{114, []DartTarget{{Triple, 20}, {Single, 14}, {Double, 20}}}
	oc.outs[115] = Out{115, []DartTarget{{Triple, 20}, {Single, 15}, {Double, 20}}}
	oc.outs[116] = Out{116, []DartTarget{{Triple, 20}, {Single, 16}, {Double, 20}}}
	oc.outs[117] = Out{117, []DartTarget{{Triple, 20}, {Single, 17}, {Double, 20}}}
	oc.outs[118] = Out{118, []DartTarget{{Triple, 20}, {Single, 18}, {Double, 20}}}
	oc.outs[119] = Out{119, []DartTarget{{Triple, 19}, {Triple, 10}, {Double, 16}}}
	oc.outs[120] = Out{120, []DartTarget{{Triple, 20}, {Single, 20}, {Double, 20}}}
	oc.outs[121] = Out{121, []DartTarget{{Triple, 17}, {Triple, 10}, {Double, 20}}}
	oc.outs[122] = Out{122, []DartTarget{{Triple, 18}, {Triple, 20}, {Double, 4}}}
	oc.outs[123] = Out{123, []DartTarget{{Triple, 19}, {Triple, 16}, {Double, 9}}}
	oc.outs[124] = Out{124, []DartTarget{{Triple, 20}, {Triple, 16}, {Double, 8}}}
	oc.outs[125] = Out{125, []DartTarget{{Single, 25}, {Triple, 20}, {Double, 20}}}
	oc.outs[126] = Out{126, []DartTarget{{Triple, 19}, {Triple, 19}, {Double, 6}}}
	oc.outs[127] = Out{127, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 8}}}
	oc.outs[128] = Out{128, []DartTarget{{Triple, 18}, {Triple, 14}, {Double, 16}}}
	oc.outs[129] = Out{129, []DartTarget{{Triple, 19}, {Triple, 16}, {Double, 12}}}
	oc.outs[130] = Out{130, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 5}}}

	// 131-170: High finishes (treble-treble-double)
	oc.outs[131] = Out{131, []DartTarget{{Triple, 20}, {Triple, 13}, {Double, 16}}}
	oc.outs[132] = Out{132, []DartTarget{{Triple, 20}, {Triple, 16}, {Double, 12}}}
	oc.outs[133] = Out{133, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 8}}}
	oc.outs[134] = Out{134, []DartTarget{{Triple, 20}, {Triple, 14}, {Double, 16}}}
	oc.outs[135] = Out{135, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 12}}}
	oc.outs[136] = Out{136, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 8}}}
	oc.outs[137] = Out{137, []DartTarget{{Triple, 19}, {Triple, 16}, {Double, 16}}}
	oc.outs[138] = Out{138, []DartTarget{{Triple, 20}, {Triple, 18}, {Double, 12}}}
	oc.outs[139] = Out{139, []DartTarget{{Triple, 19}, {Triple, 14}, {Double, 20}}}
	oc.outs[140] = Out{140, []DartTarget{{Triple, 20}, {Triple, 16}, {Double, 16}}}
	oc.outs[141] = Out{141, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 12}}}
	oc.outs[142] = Out{142, []DartTarget{{Triple, 20}, {Triple, 14}, {Double, 20}}}
	oc.outs[143] = Out{143, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 16}}}
	oc.outs[144] = Out{144, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 12}}}
	oc.outs[145] = Out{145, []DartTarget{{Triple, 20}, {Triple, 15}, {Double, 20}}}
	oc.outs[146] = Out{146, []DartTarget{{Triple, 20}, {Triple, 18}, {Double, 16}}}
	oc.outs[147] = Out{147, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 18}}}
	oc.outs[148] = Out{148, []DartTarget{{Triple, 20}, {Triple, 16}, {Double, 20}}}
	oc.outs[149] = Out{149, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 16}}}
	oc.outs[150] = Out{150, []DartTarget{{Triple, 20}, {Triple, 18}, {Double, 18}}}
	oc.outs[151] = Out{151, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 20}}}
	oc.outs[152] = Out{152, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 16}}}
	oc.outs[153] = Out{153, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 18}}}
	oc.outs[154] = Out{154, []DartTarget{{Triple, 20}, {Triple, 18}, {Double, 20}}}
	oc.outs[155] = Out{155, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 19}}}
	oc.outs[156] = Out{156, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 18}}}
	oc.outs[157] = Out{157, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 20}}}
	oc.outs[158] = Out{158, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 19}}}
	// 159 has no out
	oc.outs[160] = Out{160, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 20}}}
	oc.outs[161] = Out{161, []DartTarget{{Triple, 20}, {Triple, 17}, {Double, 25}}}
	// 162 has no out
	// 163 has no out
	oc.outs[164] = Out{164, []DartTarget{{Triple, 20}, {Triple, 18}, {Double, 25}}}
	// 165 has no out
	// 166 has no out
	oc.outs[167] = Out{167, []DartTarget{{Triple, 20}, {Triple, 19}, {Double, 25}}}
	// 168 has no out
	// 169 has no out
	oc.outs[170] = Out{170, []DartTarget{{Triple, 20}, {Triple, 20}, {Double, 25}}}
}
