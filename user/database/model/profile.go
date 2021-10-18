package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileCover struct {
	Image string `json:"image" bson:"image" map:"image"` //图片
	Video string `json:"video" bson:"video" map:"video"` //视频
	Voice string `json:"voice" bson:"voice" map:"voice"` //语音
}

type ProfileInfo struct {
	Avatar       string                    `json:"avatar" bson:"avatar" map:"avatar"`
	NickName     string                    `json:"nickname" bson:"nickname" map:"nickname"`
	Cover        ProfileCover              `json:"cover" bson:"cover" map:"cover"`
	Basic        ProfileInfoBasic          `json:"basic" bson:"basic" map:"basic"`
	Photos       []string                  `json:"photos" bson:"photos,omitempty" map:"photos,omitempty"`
	Introduction []ProfileInfoIntroduction `json:"introduction" bson:"introduction,omitempty" map:"introduction,omitempty"`
	Education    ProfileInfoEducation      `json:"education" bson:"education" map:"education"`
	Profession   ProfileInfoProfession     `json:"profession" bson:"profession" map:"profession"`
	Contact      ProfileInfoContact        `json:"contact" bson:"contact" map:"contact"`
	Address      []ProfileInfoAddress      `json:"address" bson:"address,omitempty" map:"address,omitempty"`
	Tags         []string                  `json:"tags" bson:"tags,omitempty" map:"tags,omitempty"`
}

type ProfileInfoAddress struct {
	Type     int    `json:"type" bson:"type" map:"type"`
	Nation   string `json:"nation" bson:"nation" map:"nation"`
	Province string `json:"province" bson:"province" map:"province"`
	City     string `json:"city" bson:"city" map:"city"`
	District string `json:"district" bson:"district" map:"district"`
	Address  string `json:"address" bson:"address" map:"address"`
}

type ProfileInfoEducation struct {
	HighestDegree    string `json:"highest_degree" bson:"highest_degree" map:"highest_degree"`
	GraduatedCollege string `json:"graduated_college" bson:"graduated_college" map:"graduated_college"`
}

type ProfileInfoProfession struct {
	Company      string `json:"company" bson:"company" map:"company"`
	Position     string `json:"position" bson:"position" map:"position"`
	AnnualIncome int    `json:"annual_income" bson:"annual_income" map:"annual_income"`
}

type ProfileInfoContact struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number" map:"phone_number"`
	Email       string `json:"email" bson:"email" map:"email"`
}

type ProfileInfoBasic struct {
	Sex           int              `json:"sex" bson:"sex" map:"sex"`
	Birth         ProfileInfoBirth `json:"birth" bson:"birth" map:"birth"`
	Height        int              `json:"height" bson:"height" map:"height"`
	Weight        int              `json:"weight" bson:"weight" map:"weight"`
	Marriage      int              `json:"marriage" bson:"marriage" map:"marriage"`
	Constellation int              `json:"constellation" bson:"constellation" map:"constellation"`
}

type ProfileInfoBirth struct {
	Year  int `json:"year" bson:"year" map:"year"`
	Month int `json:"month" bson:"month" map:"month"`
	Day   int `json:"day" bson:"day" map:"day"`
}

type ProfileInfoIntroduction struct {
	Type         int      `json:"type" bson:"type" map:"type"`
	Introduction string   `json:"introduction" bson:"introduction" map:"introduction"`
	Photos       []string `json:"photos" bson:"photos,omitempty" map:"photos,omitempty"`
}
type ProfileProofEducation struct {
	Type        int      `json:"type" bson:"type" map:"type"`
	ID          string   `json:"id" bson:"id" map:"id"`
	CollegeName string   `json:"college_name" bson:"college_name" map:"college_name"`
	Degree      string   `json:"degree" bson:"degree" map:"degree"`
	Place       string   `json:"place" bson:"place" map:"place"`
	Year        string   `json:"year" bson:"year" map:"year"`
	Images      []string `json:"images" bson:"images,omitempty" map:"images,omitempty"`
	Verified    bool     `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProofStudent struct {
	SchoolName string   `json:"school_name" bson:"school_name" map:"school_name"`
	Degree     string   `json:"degree" bson:"degree" map:"degree"`
	Images     []string `json:"images" bson:"images,omitempty" map:"images,omitempty"`
}

type ProfileProofCompany struct {
	CompanyName string   `json:"company_name" bson:"company_name" map:"company_name"`
	ShowName    string   `json:"show_name" bson:"show_name" map:"show_name"`
	Images      []string `json:"images" bson:"images,omitempty" map:"images,omitempty"`
}

type ProfileProofProfession struct {
	Type     int                 `json:"type" bson:"type" map:"type"`
	Company  ProfileProofCompany `json:"company" bson:"company" map:"company"`
	Student  ProfileProofStudent `json:"student" bson:"student" map:"student"`
	Verified bool                `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProofIdentityCard struct {
	ID       string   `json:"id" bson:"id" map:"id"`
	RealName string   `json:"real_name" bson:"real_name" map:"real_name"`
	Images   []string `json:"images" bson:"images,omitempty" map:"images,omitempty"`
	Verified bool     `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProof struct {
	IdentityCard ProfileProofIdentityCard `json:"identity_card" bson:"identity_card" map:"identity_card"`
	Profession   ProfileProofProfession   `json:"profession" bson:"profession" map:"profession"`
	Education    ProfileProofEducation    `json:"education" bson:"education" map:"education"`
}

type ProfileQuestion struct {
	Type     int    `json:"type" bson:"type" map:"type"`
	Question string `json:"question" bson:"question" map:"question"`
}

type Profile struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"_id" map:"_id"`
	Proof      ProfileProof       `json:"proof" bson:"proof" map:"proof"`
	Info       ProfileInfo        `json:"info" bson:"info" map:"info"`
	Questions  []ProfileQuestion  `json:"questions" bson:"questions,omitempty" map:"questions,omitempty"`
	Enable     bool               `json:"enable" bson:"enable" map:"enable"`
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
