package stats

import "fmt"

func GetCacheKey(str string) string {
	return fmt.Sprintf("info-stats-%s", str)
}

type StatsList struct {
	Yesterday StatsData `json:"yesterday"`
	Today     StatsData `json:"today"`
}
type StatsData struct {
	UserCountByUcDrawCash    uint64 `json:"point_exchange_number,omitempty"`
	UserCountByPointDrawCash uint64 `json:"ucoin_exchange_number,omitempty"`
	Cash                     string `json:"cash,omitempty"`
	PointSupply              string `json:"point_supply,omitempty"`
	UcSupply                 string `json:"uc_supply,omitempty"`
	TotalTaskUser            uint64 `json:"total_task_user,omitempty"`
	TotalFinishTask          uint64 `json:"total_finish_task,omitempty"`
	InviteNumber             uint64 `json:"invite_number,omitempty"`
	Active                   uint64 `json:"active,omitempty"`
	NewUsers                 uint64 `json:"new_users,omitempty"`
	AllActiveUsers           uint64 `json:"all_active_users,omitempty"`
}
