package tines

type Action struct {
	Id       int              `json:"id,omitempty"`
	Type     string           `json:"type,omitempty"`
	UserID   int              `json:"user_id,omitempty"`
	Options  ActionOptions    `json:"options,omitempty"`
	Name     string           `json:"name,omitempty"`
	Schedule []ActionSchedule `json:"schedule,omitempty"`
}

type ActionOptions struct {
	Mode string `json:"mode,omitempty"`
}

type ActionSchedule struct {
	Cron     string `json:"cron,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}
