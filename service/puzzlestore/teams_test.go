package puzzlestore

import "testing"

func TestTeam_Team_CanGetFirstPuzzle(t *testing.T) {
	teams := NewTeams(dictionary)
	teams.NewRound("myteam", 0)

	_, err := teams.GetPuzzle("myteam", 0)
	if err != nil {
		t.Fatalf("unexpected error when getting puzzle for myteam at index 0: %s", err)
	}
}

func TestTeam_TwoTeamsSameSeed_OnlySeedDeterminesPuzzleChain(t *testing.T) {
	teams := NewTeams(dictionary)

	seed := int64(0)
	teams.NewRound("team1", seed)
	teams.NewRound("team2", seed)

	puzzle1, _ := teams.GetPuzzle("team1", 0)
	puzzle2, _ := teams.GetPuzzle("team2", 0)

	if puzzle1 != puzzle2 {
		t.Fatalf("Both teams have the same seed, so they should both have the same puzzles at"+
			"index 0, but they were %s and %s", puzzle1, puzzle2)
	}
}

func TestTeam_GetPuzzleBeforeNewRound_Fails(t *testing.T) {
	teams := NewTeams(dictionary)

	_, err := teams.GetPuzzle("nosuchteam", 0)

	if err == nil {
		t.Fatalf("Get puzzle did not fail, despite not calling NewRound before GetPuzzle")
	}
}

func TestTeam_GetPuzzle_TeamBackByPuzzles(t *testing.T) {
	teams := NewTeams(dictionary)
	seed := int64(0)
	teams.NewRound("myteam", seed)

	// Create a puzzles with the same dictionary and seed.
	// The chain should be the same as the team "myteam"
	puzzles := NewPuzzles(dictionary, seed)

	for i := 0; i < 10; i++ {
		teamPuzzle, _ := teams.GetPuzzle("myteam", i)
		puzzle := puzzles.Get(i)

		if teamPuzzle != puzzle {
			t.Fatalf("Team puzzle %s should be identical to puzzle chain %s at index %d", teamPuzzle, puzzle, i)
		}
	}
}
