package db

import (
	"../model"
	"time"
)

const ()

type ()

func Crowdfunds(offset, limit uint) (list []model.Crowdfund) {

	rows, err := db.Query(`SELECT "id", "owner_id", "ctime", "utime", "goal", "collected", "name", "cover" `+
		`FROM "public"."crowdfunds" ORDER BY "id" DESC LIMIT $1 OFFSET $2`, limit, offset)

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

func GetCrowdfund(id int64) (crowdfund_ptr *model.Crowdfund, err error) {

	var crowdfund model.Crowdfund

	err = db.QueryRow(`SELECT "goal", "name", "cover", "video" `+
		`FROM "public"."crowdfunds" `+
		`WHERE "id"=$1`, id).Scan(&crowdfund.Goal, &crowdfund.Name, &crowdfund.Cover, &crowdfund.Video)

	if err != nil {
		logger.Println(err)
		return
	}

	crowdfund_ptr = &crowdfund
	return
}

func SaveCrowdfund(crowdfund *model.Crowdfund) (err error) {

	unix_time := time.Now().Unix()

	if crowdfund.Id == 0 {

		err = db.QueryRow(`INSERT INTO "public"."crowdfunds" `+
			`("owner_id", "ctime", "goal", "name", "cover", "utime", "video") `+
			`VALUES `+
			`($1, $2, $3, $4, $5, $2, $6) `+
			`RETURNING id`,
			crowdfund.OwnerId, unix_time, crowdfund.Goal, crowdfund.Name,
			crowdfund.Cover, crowdfund.Video).Scan(&crowdfund.Id)

		if err != nil {
			logger.Println(err)
			return
		}

	} else {

		_, err = db.Exec(`UPDATE "public"."crowdfunds" `+
			`SET "goal"=$3, "name"=$4, "cover"=$5, "utime"=$6, "video"=$7 `+
			`WHERE "id"=$1 AND "owner_id"=$2`,
			crowdfund.Id, crowdfund.OwnerId, crowdfund.Goal, crowdfund.Name,
			crowdfund.Cover, unix_time, crowdfund.Video)

		if err != nil {
			logger.Println(err)
			return
		}
	}

	return
}
