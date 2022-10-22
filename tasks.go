package snapattack

type Task struct {
	Creation string      `json:"creation"`
	Modified string      `json:"modified"`
	ID       string      `json:"task_id"`
	Status   string      `json:"status"`
	Output   interface{} `json:"output"`
}
