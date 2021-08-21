package upload

import (
	"context"
	"net/http"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"

	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_Index(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileModel := upload.NewModel(self.M)
	filter := bson.D{}

	types := r.Form["type[]"]
	if len(types) > 0 {
		filter = append(filter, bson.E{
			"type", bson.D{
				{"$in", types},
			},
		})
	}

	storages := r.Form["storage[]"]
	if len(storages) > 0 {
		filter = append(filter, bson.E{
			"storage", bson.D{
				{"$in", storages},
			},
		})
	}

	if userIDHex := r.FormValue("user_id"); userIDHex != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			filter = append(filter, bson.E{"user_id", userID})
		}
	}

	if contentType := r.FormValue("content_type"); contentType != "" {
		filter = append(filter, bson.E{"content_type", contentType})
	}

	if originalFileName := r.FormValue("original_file_name"); originalFileName != "" {
		filter = append(filter, bson.E{"original_file_name", originalFileName})
	}

	if fileName := r.FormValue("file_name"); fileName != "" {
		filter = append(filter, bson.E{"file_name", fileName})
	}

	if fileExt := r.FormValue("file_ext"); fileExt != "" {
		filter = append(filter, bson.E{"file_ext", fileExt})
	}

	cnt, _ := fileModel.CountDocuments(context.Background(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	items := []model.Upload{}
	if cur, err := fileModel.Find(context.Background(), filter, opt); err == nil {
		cur.All(context.TODO(), &items)
	}

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: items,
	}

	response.JSON(rw, 0, rpr, "")
}

//单独给editor_upload设置无需权限即可使用
// func (p *Service) Route_POST_Editor_upload(rw http.ResponseWriter, r *http.Request) bool {
// 	returnData := map[string]interface{}{
// 		"code":    "001",
// 		"message": "上传失败",
// 		"data":    utils.M{},
// 	}

// 	fileType := r.FormValue("fileType")
// 	sBase64 := r.FormValue("base64")
// 	if sBase64 == "1" {
// 		fileFormName := "img_base64_data"
// 		if fileType == predefined.UploadFileTypeImage {
// 			if fm, err := p.SaveEncodeFile(r, fileFormName); err != nil {
// 				log.Error(err)
// 			} else {
// 				returnData["code"] = "000"
// 				returnData["message"] = "上传成功"
// 				returnData["data"].(utils.M)["url"] = fm.URL
// 			}
// 		}
// 	} else {
// 		fileFormName := "imgFile"
// 		switch fileType {
// 		case predefined.UploadFileTypeImage, predefined.UploadFileTypeMedia, predefined.UploadFileTypeFile, "":
// 			if fm, err := p.SaveFormFile(r, fileFormName); err != nil {
// 				log.Error(err)
// 			} else {
// 				returnData["code"] = "000"
// 				returnData["message"] = "上传成功"
// 				returnData["data"].(utils.M)["url"] = fm.URL
// 			}
// 		}
// 	}

// 	httputils.ResponseJson(rw, returnData)

// 	return true
// }

//单独给editor_manager设置无需权限即可使用
// func (p *Service) Init_Editor_manager(rw http.ResponseWriter, r *http.Request) {
// 	p.Service.Init(rw, r)
// 	p.OffRight = true
// }

// func (p *Service) Route_GET_Editor_manager(rw http.ResponseWriter, r *http.Request) bool {
// 	returnData := map[string]interface{}{
// 		"code":    "001",
// 		"message": "操作失败",
// 		"data":    []utils.M{},
// 	}

// 	fileType := r.URL.Query().Get("fileType")

// 	fileModel := file.NewModel(p.M, schema.File, p.Cache)

// 	filter := bson.D{}
// 	switch fileType {
// 	case predefined.FileTypeImage, predefined.FileTypeMedia:
// 		filter = append(filter, bson.E{"type", fileType})
// 	default:
// 		filter = append(filter,
// 			bson.E{
// 				"type", bson.D{
// 					{
// 						"$in", bson.A{
// 							predefined.FileTypeFile,
// 							predefined.FileTypeDocument,
// 							predefined.FileTypeCompression,
// 						},
// 					},
// 				},
// 			},
// 		)
// 	}

// 	var prePageNum int64 = 20
// 	var showPageNum int64 = 5

// 	cnt, _ := fileModel.CountDocuments(context.Background(), filter)
// 	pagination := p.Pagination(rw, r, cnt, prePageNum, showPageNum)
// 	skipNum := pagination["skipNum"].(int64)

// 	opt := options.Find().SetSort(bson.D{
// 		{"create_time", -1},
// 	}).SetSkip(skipNum).SetLimit(prePageNum)

// 	cols := []model.File{}
// 	if cur, err := fileModel.Find(context.Background(), filter, opt); err == nil {
// 		if err := cur.All(context.TODO(), &cols); err == nil {
// 			returnData["code"] = "000"
// 			returnData["message"] = "操作成功"
// 			returnData["count"] = cnt
// 			returnData["page"] = pagination["currentPage"]
// 			returnData["pagesize"] = pagination["pageNum"]
// 			data := []utils.M{}
// 			for _, col := range cols {
// 				data = append(data, utils.M{
// 					"thumbURL": col.URL,
// 					"oriURL":   col.URL,
// 					"filesize": col.Size,
// 					"width":    0,
// 					"height":   0,
// 				})
// 			}

// 			if len(data) > 0 {
// 				returnData["data"] = data
// 			}
// 		}
// 	}

// 	httputils.ResponseJson(rw, returnData)

// 	return true
// }
