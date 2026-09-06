package v1

import "testing"

func TestLocalMediaItemValidRejectsEmptyIdentity(t *testing.T) {
	item := LocalMediaItem{ID:"", Title:"title", Type:"video", Source:"local"}
	if item.Valid() { t.Fatal("empty media ID must be invalid") }
}
