package sdk_interface

type WatchEventsStream_Command struct {
	Command string `json:"command"`
}

type WatchEventsStream_WatchCommand struct {
	Command  string   `json:"command"`
	Session  string   `json:"session"`
	Policies []string `json:"policies"`
}

type WatchEventsStream_UnwatchCommand struct {
	Command string `json:"command"`
	Session string `json:"session"`
}

type WatchEventsStream_UnwatchAllCommand struct {
	Command string `json:"command"`
}

type WatchEventsStream_Event struct {
	Session string         `json:"session"`
	Name    string         `json:"name"`
	ID      string         `json:"id"`
	Data    map[string]any `json:"data"`
}
