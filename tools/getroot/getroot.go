package getroot

import (
	"fmt"
	"strings"
	"github.com/ziutek/mymysql/autorc"
)

type UserRelation struct {
	UserId   int `json:"user_id"`
	ParentId int `json:"parent_id"`
	GrandId  int `json:"grand_id"`
	RootId   int `json:"root_id"`
}

func GetUserRoot(db *autorc.Conn) error {

	query := ` SELECT user_id,parent_id,grand_id,root_id FROM tmm.invite_codes `
	var list []*UserRelation
	rows, _, err := db.Query(query)
	if err != nil {
		return err
	}
	for _, row := range rows {
		list = append(list, &UserRelation{
			UserId:   row.Int(0),
			ParentId: row.Int(1),
			GrandId:  row.Int(2),
			RootId:   row.Int(3),
		})
	}
	for _, d := range list {
		if d.ParentId == 0 || d.RootId != 0 {
			continue
		}
		d.RootId = getRoot(*d, list)
	}
	query = `INSERT INTO tmm.invite_codes(user_id,root_id)VALUES %s
	ON DUPLICATE KEY UPDATE root_id=VALUES(root_id)`
	var value []string
	for _, user := range list {
		value = append(value, fmt.Sprintf("(%d,%d)", user.UserId, user.RootId))
		if len(value) > 1000 {
			if _, _, err := db.Query(query, strings.Join(value, `,`)); err != nil {
				return err
			}
			value = []string{}
		}
	}
	if len(value) > 0 {
		if _, _, err := db.Query(query, strings.Join(value, `,`)); err != nil {
			return err
		}
	}
	return nil
}
func getRoot(inv UserRelation, lists []*UserRelation) int {
	if inv.ParentId != 0 {
		for _, _inv := range lists {
			if _inv.UserId == inv.ParentId {
				return getRoot(*_inv, lists)
			}
		}
	}
	return inv.UserId
}
