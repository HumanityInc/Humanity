package db

import (
	"../model"
	"database/sql"
	"time"
)

const ()

type ()

func UserCrowdfunds(prof_id int64, offset, limit uint) (list []model.Crowdfund) {

	rows, err := db.Query(`SELECT "id", "owner_id", "ctime", "utime", "goal", "collected", "name", "cover" `+
		`FROM "public"."crowdfunds" `+
		`WHERE "owner_id"=$3 `+
		`OR "owner_id" IN (SELECT "prof_id1" FROM "public"."follows" WHERE "prof_id0"=$3) `+
		`OR "id" IN (SELECT "cf_id" FROM "public"."favorits" WHERE "prof_id"=$3) `+
		`ORDER BY "id" DESC LIMIT $1 OFFSET $2`, limit, offset, prof_id)

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

func Follow(prof_id0, prof_id1 int64, check bool) {

	var err error

	if check {

		_, err = db.Exec(`INSERT INTO "public"."follows" ("prof_id0", "prof_id1") VALUES ($1, $2)`,
			prof_id0, prof_id1)

	} else {

		_, err = db.Exec(`DELETE FROM "public"."follows" WHERE "prof_id0"=$1 AND "prof_id1"=$2`,
			prof_id0, prof_id1)
	}

	if err != nil {
		logger.Println(err)
		return
	}
}

func Favorit(prof_id, crowdfund_id int64, check bool) {

	var err error

	if check {

		_, err = db.Exec(`INSERT INTO "public"."favorits" ("prof_id", "cf_id") VALUES ($1, $2)`,
			prof_id, crowdfund_id)

	} else {

		_, err = db.Exec(`DELETE FROM "public"."favorits" WHERE "prof_id"=$1 AND "cf_id"=$2`,
			prof_id, crowdfund_id)
	}

	if err != nil {
		logger.Println(err)
		return
	}
}

func Crowdfunds(prof_id int64, offset, limit uint) (list []model.Crowdfund) {

	// rows, err := db.Query(`SELECT "id", "owner_id", "ctime", "utime", "goal", "collected", "name", "cover" `+
	// `FROM "public"."crowdfunds" ORDER BY "id" DESC LIMIT $1 OFFSET $2`, limit, offset)

	rows, err := db.Query(`SELECT "id", "owner_id", "ctime", "utime", "goal", "collected", "name", "cover", "prof_id" `+
		`FROM "public".crowdfunds `+
		`LEFT JOIN "public"."favorits" ON (id = cf_id AND prof_id=$3) ORDER BY "id" DESC LIMIT $1 OFFSET $2`,
		limit, offset, prof_id)

	if err != nil {
		logger.Println(err)
		return
	}
	defer rows.Close()

	var flag sql.NullInt64
	cf := model.Crowdfund{}
	list = make([]model.Crowdfund, 0, 32)

	for rows.Next() {

		if err = rows.Scan(&cf.Id, &cf.OwnerId, &cf.CreateTime, &cf.UpdateTime, &cf.Goal,
			&cf.Сollected, &cf.Name, &cf.Cover, &flag); err != nil {

			logger.Println(err)
			return
		}

		if flag.Valid && flag.Int64 != 0 {
			cf.Favorit = true
		} else {
			cf.Favorit = false
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

	db.QueryRow(`SELECT "id" FROM "public"."crowdfunds" WHERE "id">$1 ORDER BY "id" LIMIT 1`, id).Scan(&crowdfund_ptr.NextId)
	db.QueryRow(`SELECT "id" FROM "public"."crowdfunds" WHERE "id"<$1 ORDER BY "id" DESC LIMIT 1`, id).Scan(&crowdfund_ptr.PrevId)

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
