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
	Type     int    `json:"type" bson:"type" map:"type" form:"type" validate:"oneof=0 1"`
	Nation   string `json:"nation" bson:"nation" map:"nation" form:"nation" validate:"required"`
	Province string `json:"province" bson:"province" map:"province" form:"province" validate:"required"`
	City     string `json:"city" bson:"city" map:"city" form:"city" validate:"required"`
	District string `json:"district" bson:"district" map:"district" form:"district" validate:"required"`
	Address  string `json:"address" bson:"address" map:"address" form:"address" validate:"required"`
}

type ProfileInfoEducation struct {
	HighestDegree    string `json:"highest_degree" bson:"highest_degree" map:"highest_degree" form:"highest_degree" validate:"required"`
	GraduatedCollege string `json:"graduated_college" bson:"graduated_college" map:"graduated_college" form:"graduated_college" validate:"required"`
}

type ProfileInfoProfession struct {
	Company      string `json:"company" bson:"company" map:"company" form:"company" validate:"required"`
	Position     string `json:"position" bson:"position" map:"position" form:"position" validate:"required"`
	AnnualIncome int    `json:"annual_income" bson:"annual_income" map:"annual_income" form:"annual_income" validate:"oneof=0 1 2 3 4 5"`
}

type ProfileInfoContact struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number" map:"phone_number" form:"position" validate:"required,chinamobile"`
	Email       string `json:"email" bson:"email" map:"email" form:"position" validate:"required,email"`
}

type ProfileInfoBasic struct {
	Sex           int              `json:"sex" bson:"sex" map:"sex" form:"sex" validate:"oneof=0 1 2"`
	Birth         ProfileInfoBirth `json:"birth" bson:"birth" map:"birth" form:"birth" validate:"required,dive"`
	Height        int              `json:"height" bson:"height" map:"height" form:"height" validate:"required,gte=140,lte=250"`
	Weight        int              `json:"weight" bson:"weight" map:"weight" form:"weight" validate:"required,gte=50,lte=250"`
	Marriage      int              `json:"marriage" bson:"marriage" map:"marriage" form:"marriage" validate:"oneof=0 1 2"`
	Constellation int              `json:"constellation" bson:"constellation" map:"constellation" form:"constellation" validate:"gte=0,lte=11"`
}

type ProfileInfoBirth struct {
	Year  int `json:"year" bson:"year" map:"year" form:"year" validate:"numeric"`
	Month int `json:"month" bson:"month" map:"month" form:"month" validate:"required,gte=1,lte=12"`
	Day   int `json:"day" bson:"day" map:"day" form:"day" validate:"required,gte=1,lte=31"`
}

type ProfileInfoIntroduction struct {
	Type         int      `json:"type" bson:"type" map:"type" form:"type" validate:"oneof=0 1 2"`
	Introduction string   `json:"introduction" bson:"introduction" map:"introduction" form:"introduction" validate:"required"`
	Photos       []string `json:"photos" bson:"photos,omitempty" map:"photos,omitempty" form:"photos" validate:"gt=0,dive,required"`
}
type ProfileProofEducation struct {
	Type        int      `json:"type" bson:"type" map:"type" form:"type" validate:"oneof=0 1 2 3 4"`
	ID          string   `json:"id" bson:"id" map:"id" form:"id" validate:"required"`
	CollegeName string   `json:"college_name" bson:"college_name" map:"college_name" form:"college_name" validate:"required"`
	Degree      string   `json:"degree" bson:"degree" map:"degree" form:"degree" validate:"required"`
	Place       string   `json:"place" bson:"place" map:"place" form:"place" validate:"required_if=type 4"`
	Year        string   `json:"year" bson:"year" map:"year" form:"year" validate:"required_if=type 4"`
	Images      []string `json:"images" bson:"images,omitempty" map:"images,omitempty" form:"images" validate:"required_if=type 1,gt=0,dive,required"`
	Verified    bool     `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProofStudent struct {
	SchoolName string   `json:"school_name" bson:"school_name" map:"school_name" form:"school_name" validate:"required"`
	Degree     string   `json:"degree" bson:"degree" map:"degree" form:"degree" validate:"required"`
	Images     []string `json:"images" bson:"images,omitempty" map:"images,omitempty" form:"images" validate:"gt=0,dive,required"`
}

type ProfileProofCompany struct {
	CompanyName string   `json:"company_name" bson:"company_name" map:"company_name" form:"company_name" validate:"required"`
	ShowName    string   `json:"show_name" bson:"show_name" map:"show_name" form:"degree" show_name:"required"`
	Images      []string `json:"images" bson:"images,omitempty" map:"images,omitempty" form:"images" validate:"gt=0,dive,required"`
}

type ProfileProofProfession struct {
	Type     int                 `json:"type" bson:"type" map:"type" form:"type" validate:"oneof=0 1 2 3 4 5 6"`
	Company  ProfileProofCompany `json:"company" bson:"company" map:"company" form:"company" validate:"required_without=Student,dive"`
	Student  ProfileProofStudent `json:"student" bson:"student" map:"student" form:"student" validate:"required_if=type 6,required_without=Company,dive"`
	Verified bool                `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProofIdentityCard struct {
	ID       string   `json:"id" bson:"id" map:"id" form:"id" validate:"required"`
	RealName string   `json:"real_name" bson:"real_name" map:"real_name" form:"real_name" validate:"required"`
	Images   []string `json:"images" bson:"images,omitempty" map:"images,omitempty" form:"images" validate:"gt=0,dive,required"`
	Verified bool     `json:"verified" bson:"verified" map:"verified"`
}

type ProfileProof struct {
	IdentityCard ProfileProofIdentityCard `json:"identity_card" bson:"identity_card" map:"identity_card"`
	Profession   ProfileProofProfession   `json:"profession" bson:"profession" map:"profession"`
	Education    ProfileProofEducation    `json:"education" bson:"education" map:"education"`
}

type ProfileQuestion struct {
	Type     int    `json:"type" bson:"type" map:"type" form:"type" validate:"numeric,oneof=0"`
	Question string `json:"question" bson:"question" map:"question" form:"question" validate:"required"`
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
