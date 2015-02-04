package db

import (
	"../model"
)

const ()

type ()

func Crowdfunds(offset, limit uint) (list []model.Crowdfund) {

	rows, err := db.Query(`SELECT "id", "owner_id", "ctime", "utime", "goal", "collected", "name", "cover" `+
		`FROM "public"."crowdfunds" order by "id" LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		logger.Println(err)
		return
	}
	defer rows.Close()

	cf := model.Crowdfund{}
	list = make([]model.Crowdfund, 0, 32)

	for rows.Next() {

		if err = rows.Scan(&cf.Id, &cf.OwnerId, &cf.CreateTime, &cf.UpdateTime, &cf.Goal,
			&cf.Сollected, &cf.Name, &cf.Cover); err != nil {

			logger.Println(err)
			return
		}

		list = append(list, cf)
	}

	return
}

func SaveCrowdfund(crowdfund *model.Crowdfund) (err error) {

	if crowdfund.Id == 0 {

		// insert

	} else {

		// update
	}

	return
}
