package types

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email     string        `json:"email" bson:"email" validate:"nonnil"`
	LastName  string        `json:"last_name" bson:"last_name" validate:"nonnil"`
	Country   string        `json:"country" bson:"country" validate:"nonnil"`
	City      string        `json:"city" bson:"city" validate:"nonnil"`
	Gender    string        `json:"gender" bson:"gender" validate:"nonnil"`
	BirthDate string        `json:"birch_date" bson:"birch_date" validate:"nonnil"`
}

type Game struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	PointsGained string        `json:"points_gained" bson:"points_gained" validate:"nonnil"`
	WinStatus    string        `json:"win_status" bson:"win_status" validate:"nonnil"`
	GameType     string        `json:"game_type" bson:"game_type" validate:"nonnil"`
	Created      string        `json:"created" bson:"created" validate:"nonnil"`
}
