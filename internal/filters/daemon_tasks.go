package filters

import "github.com/gameap/gameap/internal/domain"

type FindDaemonTask struct {
	IDs                []uint
	DedicatedServerIDs []uint
	ServerIDs          []*uint
	Tasks              []domain.DaemonTaskType
	Statuses           []domain.DaemonTaskStatus
}

func FindDaemonTaskByIDs(ids ...uint) *FindDaemonTask {
	return &FindDaemonTask{
		IDs: ids,
	}
}
