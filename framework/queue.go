package framework

type SongQueue struct {
	List    []Song
	Current *Song
	Running bool
}

// Gets the full queue list
func (queue SongQueue) Get() []Song {
	return queue.List
}

// Sets the full queue list
func (queue *SongQueue) Set(list []Song) {
	queue.List = list
}

// Appends a song to the end of the queue list
func (queue *SongQueue) Add(song Song) {
	queue.List = append(queue.List, song)
}

// Checks to see if the queue has another song
func (queue SongQueue) HasNext() bool {
	return len(queue.List) > 0
}

// Returns the next song in the queue
func (queue *SongQueue) Next() Song {
	song := queue.List[0]
	queue.List = queue.List[1:]
	queue.Current = &song
	return song
}

// Clears the queue, regardless of size
func (queue *SongQueue) Clear() {
	queue.List = make([]Song, 0)
	queue.Running = false
	queue.Current = nil
}

// Starts a new session and queue
func (queue *SongQueue) Start(sess *Session, callback func(string)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		callback("Now playing `" + song.Title + "`.")
		sess.Play(song)
	}

	if !queue.Running {
		callback("Stopped playing.")
	} else {
		callback("Finished queue.")
	}
}

// Returns the song that's currently playing
func (queue *SongQueue) CurrentSong() *Song {
	return queue.Current
}

// Pauses the queue list, pause happens after the song ends
func (queue *SongQueue) Pause() {
	queue.Running = false
}

// Creates a new queue list
func newSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.List = make([]Song, 0)
	return queue
}
