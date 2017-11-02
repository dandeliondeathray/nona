package plumber

import "testing"

func TestCodecs_LoadFromPath_KnownSchemaFound(t *testing.T) {
	codecs, err := LoadCodecsFromPath("../schema")
	if err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	_, err = codecs.ByName("Chat")
	if err != nil {
		t.Fatalf("Could not find known schema Chat")
	}

	_, err = codecs.ByName("PuzzleNotification")
	if err != nil {
		t.Fatalf("Could not find known schema PuzzleNotification")
	}
}
