package playlist

import (
	"fmt"
	songmodule "gocloudcamp/core/song"
	"log"
	"math/rand"
	"testing"
	"time"
)

func createPlaylist(t *testing.T) *playlist {
	got, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	return got
}

func createPlaylistWithSongs(t *testing.T, songsAmount int) (pl *playlist, songs []*songmodule.Song, ids []uint32) {
	rand.Seed(time.Now().Unix())
	pl = createPlaylist(t)
	for i := 0; i < songsAmount; i++ {
		sng := songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(1*time.Second)))+200*time.Millisecond)
		songs = append(songs, sng)
		added, err := pl.AddSong(*songs[i])
		if err != nil {
			t.Fatalf("Couldn't add a song: %v", err)
		}
		ids = append(ids, added)
	}
	return
}

func validateSong(got, want *songmodule.Song) bool {
	return got != nil && want.Equal(*got)
}

func validateCurrentSong(pl *playlist, want *songmodule.Song) bool {
	return pl != nil && pl.currentSong != nil && validateSong(&pl.currentSong.data, want)
}

func TestNewPlaylist(t *testing.T) {
	got := createPlaylist(t)
	if got.currentSong != nil && got.currentSong.data.IsValid() {
		t.Fatal("Freshly created playlist is initialized with a song")
	}
	if got.timer == nil {
		t.Fatal("Playlist created without a timer")
	}
}

func TestPlaylist_AddSong(t *testing.T) {
	pl := createPlaylist(t)
	song := songmodule.NewSong("Test song", time.Second*30)
	_, err := pl.AddSong(*song)
	if err != nil {
		t.Fatalf("Couldn't add a song: %v", err)
	}
	if pl.currentSong == nil || !song.Equal(pl.currentSong.data) {
		t.Fatal("Data of the song added is not equal to the original")
	}
}

func TestPlaylist_AddNegativeSong(t *testing.T) {
	pl := createPlaylist(t)
	song := songmodule.NewSong("Test song", time.Second*-10)
	_, err := pl.AddSong(*song)
	if err == nil {
		t.Fatal("Added a song with negative length without errors")
	} else {
		t.Logf("Negative-length song couldn't be added due to an error: %v", err)
	}
}

func TestPlaylist_AddSong_many(t *testing.T) {
	pl := createPlaylist(t)
	const amount = 1000000
	var songs []*songmodule.Song
	for i := 0; i < amount; i++ {
		songs = append(songs, songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(time.Hour)))))
	}
	start := time.Now()
	for _, song := range songs {
		_, err := pl.AddSong(*song)
		if err != nil {
			log.Fatalf("Coudln't add song due to an error: %v", err)
		}
	}
	t.Logf("%v elapsed on adding %v songs to playlist", time.Now().Sub(start), amount)
	lastSong := pl.currentSong
	for i, song := range songs {
		if lastSong == nil || !lastSong.data.IsValid() || !song.Equal(lastSong.data) {
			t.Fatalf("Song %v doesn't match the one added to playlist", i)
		}
		lastSong = lastSong.next
	}
	if lastSong != nil && lastSong.data.IsValid() {
		t.Fatalf("Playlist has more songs that were added")
	}
}

func TestPlaylist_Play(t *testing.T) {
	const grace = 50 * time.Millisecond
	pl, songs, _ := createPlaylistWithSongs(t, 4)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Play()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[1]) {
		t.Fatal("Playlist didn't start the next song after the first one should've done playing")
	}
}

func TestPlaylist_Pause(t *testing.T) {
	const grace = 50 * time.Millisecond
	pl, songs, _ := createPlaylistWithSongs(t, 10)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Play()
	pl.Pause()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Song was lost or not paused")
	}
	pl.Play()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[1]) {
		t.Fatal("Playlist didn't start the next song after the first one should've done playing")
	}
}

func TestPlaylist_Seek(t *testing.T) {
	const amount = 10
	pl, songs, _ := createPlaylistWithSongs(t, amount)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	for i := 1; i < amount; i++ {
		pl.Next()
		if !validateCurrentSong(pl, songs[i]) {
			t.Fatalf("Playlist should play %v-th song after %v next() method calls", i+1, i)
		}
	}
	if pl.currentSong.next != nil && pl.currentSong.next.data.IsValid() {
		t.Fatal("Current song in playlist is not last, but should be")
	}
	for i := amount - 2; i > -1; i-- {
		pl.Prev()
		if !validateCurrentSong(pl, songs[i]) {
			t.Fatalf("Playlist should play %v-th song after %v Prev() method calls", i+1, amount-i+1)
		}
	}
	if pl.currentSong.previous != nil && pl.currentSong.previous.data.IsValid() {
		t.Fatal("Current song in playlist is not first, but should be")
	}
}

func TestPlaylist_Index(t *testing.T) {
	const amount = 10
	pl, songs, ids := createPlaylistWithSongs(t, amount)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	for i := 0; i < amount; i++ {
		want := songs[i]
		got, exists := pl.GetSong(ids[i])
		if !exists || !want.Equal(got) {
			t.Fatalf("Song %v wasn't indexed correctly", i)
		}
	}
}

func TestPlaylist_Delete(t *testing.T) {
	const amount = 10
	const amountToDelete = amount / 3
	pl, songs, ids := createPlaylistWithSongs(t, amount)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Play()
	_, err := pl.RemoveSong(ids[0])
	if err == nil {
		t.Fatal("Playlist allowed to remove the song which is now playing")
	} else {
		t.Logf("Got an error while trying to remove the now-playing song: %v", err.Error())
	}
	_, err = pl.RemoveSong(0)
	if err == nil {
		t.Fatal("Playlist didn't create an error when trying to remove a non-existent song (or you're VERY lucky)")
	} else {
		t.Logf("Got an error while trying to remove a non-existent song: %v", err.Error())
	}
	pl.Pause()
	rand.Shuffle(amount-1, func(i, j int) {
		i++
		j++
		ids[i], ids[j] = ids[j], ids[i]
		songs[i], songs[j] = songs[j], songs[i]
	})
	for i := 1; i <= amountToDelete; i++ {
		removed, err := pl.RemoveSong(ids[i])
		if err != nil {
			t.Fatalf("Playlist couldn't remove a song that is allowed to remove: %v", err)
		}
		if !songs[i].Equal(removed) {
			t.Fatal("Playlist removed wrong song")
		}
	}
}

func TestPlaylist_DeleteEdges(t *testing.T) {
	const amount = 10
	pl, songs, ids := createPlaylistWithSongs(t, amount)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Next()
	removed, err := pl.RemoveSong(ids[0])
	if err != nil {
		t.Fatalf("Playlist couldn't remove the first song, but it's allowed to remove: %v", err)
	}
	if !songs[0].Equal(removed) {
		t.Fatal("Playlist removed wrong song")
	}
	removed, err = pl.RemoveSong(ids[amount-1])
	if err != nil {
		t.Fatalf("Playlist couldn't remove the first song, but it's allowed to remove: %v", err)
	}
	if !songs[amount-1].Equal(removed) {
		t.Fatal("Playlist removed wrong song")
	}
	songs = songs[1 : amount-1]
	for i := 0; i < amount; i++ {
		song := songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(time.Hour))))
		songs = append(songs, song)
		_, err := pl.AddSong(*song)
		if err != nil {
			log.Fatalf("Couldn't add song due to an error: %v", err)
		}
	}
	lastSong := pl.currentSong
	for i, song := range songs {
		if lastSong == nil || !lastSong.data.IsValid() || !song.Equal(lastSong.data) {
			t.Fatalf("Song %v doesn't match the one added to playlist", i)
		}
		lastSong = lastSong.next
	}
	if lastSong != nil && lastSong.data.IsValid() {
		t.Fatalf("Playlist has more songs that were added")
	}
}
