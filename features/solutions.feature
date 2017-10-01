Feature: Solve a puzzle
	A user can solve a puzzle by finding a Swedish word that matches the letters in the puzzle.

	Scenario: Solving a puzzle
		Given that the current puzzle is VANTRIVSA
		 When the user attempts the solution VANTRIVAS
		 Then the solution is correct

	Scenario: Getting a puzzle
		 When Erik has requested the puzzle
		 Then the response is a nine letter word

	Scenario: All users have the same chain of words
		Given Erik has gotten five puzzles
		 When Johan has gotten five puzzles
		 Then Erik and Johan have received the same puzzles

	Scenario: Next puzzle
		Given Erik has requested the puzzle
		  And he has solved the puzzle
		 When Erik has requested a new puzzle
		 Then the response is a different puzzle
