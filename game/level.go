package game

const (
	completeBonus        = 50
	pointsPerSavedStroke = 100
	holeInOneBonus       = 100
	basePar              = 3
)

func kindsForLevel(level int) []obstacleKind {
	switch {
	case level <= 1:
		return nil
	case level == 2:
		return []obstacleKind{kindSmall}
	case level == 3:
		return []obstacleKind{kindSmall, kindMedium}
	case level == 4:
		return []obstacleKind{kindMedium, kindLarge}
	case level == 5:
		return []obstacleKind{kindSmall, kindMedium, kindLarge}
	case level == 6:
		return []obstacleKind{kindSmall, kindMedium, kindLarge, kindXLarge}
	default:
		kinds := []obstacleKind{kindSmall, kindMedium, kindLarge, kindXLarge}
		extras := []obstacleKind{kindSmall, kindMedium, kindLarge}
		extra := level - 6
		for i := 0; i < extra && len(kinds) < 8; i++ {
			kinds = append(kinds, extras[i%len(extras)])
		}
		return kinds
	}
}

func parForLevel(level int) int {
	return basePar + (level-1)/3
}

func scoreForHole(level, strokes int) int {
	saved := parForLevel(level) - strokes
	points := completeBonus + saved*pointsPerSavedStroke
	if strokes == 1 {
		points += holeInOneBonus
	}
	if points < completeBonus {
		points = completeBonus
	}
	return points
}
