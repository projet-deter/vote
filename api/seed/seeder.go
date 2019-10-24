package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/m2fof/vote/api/models"
)

var users = []models.User{
	models.User{
		Uuid:        "69cda7b3-9a67-47ce-beac-16ce1dc177e0",
		AccessLevel: 1,
		First_name:  "Steven ",
		Last_name:   " victor",
		Email:       "steven@gmail.com",
		Password:    "password",
		Birth_date:  "11/05/2000",
	},
	models.User{
		Uuid:        "33ert7b3-9a67-47ce-cgrt-16ce1dc177e0",
		AccessLevel: 0,
		First_name:  "Martin ",
		Last_name:   "Luther",
		Email:       "luther@gmail.com",
		Password:    "password",
		Birth_date:  "24/04/2005",
	},
}

var votes = []models.Vote{
	models.Vote{
		Title: "Title 1",
		Desc:  "Hello world 1",
	},
	models.Vote{
		Title: "Title 2",
		Desc:  "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Vote{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Vote{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Vote{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		votes[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Vote{}).Create(&votes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed votes table: %v", err)
		}
	}
}
