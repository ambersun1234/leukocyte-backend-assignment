package types

type CallbackFunc func(string) error

type JobObject struct {
	Namespace     string
	Name          string
	Image         string
	RestartPolicy string
	Commands      []string
}

type RoutingKey = string
