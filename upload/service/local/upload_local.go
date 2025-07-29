package local

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"bufio"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 全局随机数生成器，线程安全
var (
	globalRand *rand.Rand
	randMutex  sync.Mutex
)

func init() {
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 全局魔数映射表，按文件扩展名组织
var magicNumbers = map[string][][]byte{
	// 图片格式
	"jpg":  {{0xFF, 0xD8, 0xFF}},
	"jpeg": {{0xFF, 0xD8, 0xFF}},
	"png":  {{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
	"gif":  {{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, {0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
	"bmp":  {{0x42, 0x4D}},
	"webp": {{0x52, 0x49, 0x46, 0x46, 0, 0, 0, 0, 0x57, 0x45, 0x42, 0x50}},
	"tiff": {{0x49, 0x49, 0x2A, 0x00}, {0x4D, 0x4D, 0x00, 0x2A}},
	"ico":  {{0x00, 0x00, 0x01, 0x00}},
	"psd":  {{0x38, 0x42, 0x50, 0x53}},                                                 // Photoshop格式
	"heic": {{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63}}, // HEIC格式
	"svg":  {{0x3C, 0x73, 0x76, 0x67}},                                                 // SVG格式
	"psb":  {{0x38, 0x42, 0x50, 0x53}},                                                 // Photoshop Large Document

	// 文档格式
	"pdf":  {{0x25, 0x50, 0x44, 0x46}},
	"doc":  {{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	"docx": {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}},
	"xls":  {{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	"xlsx": {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}},
	"ppt":  {{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	"pptx": {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}},
	"rtf":  {{0x7B, 0x5C, 0x72, 0x74, 0x66, 0x31}},
	"odt":  {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}}, // OpenDocument Text
	"ods":  {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}}, // OpenDocument Spreadsheet
	"odp":  {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}}, // OpenDocument Presentation
	"epub": {{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x06, 0x00}}, // EPUB电子书
	"mobi": {{0x42, 0x4F, 0x4F, 0x4B, 0x4D, 0x4F, 0x42, 0x49}}, // MOBI电子书
	"chm":  {{0x49, 0x54, 0x53, 0x46}},                         // CHM帮助文件

	// 压缩格式
	"zip": {{0x50, 0x4B, 0x03, 0x04}, {0x50, 0x4B, 0x05, 0x06}, {0x50, 0x4B, 0x07, 0x08}},
	"rar": {{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}, {0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01, 0x00}},
	"7z":  {{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}},
	"gz":  {{0x1F, 0x8B}},
	"bz2": {{0x42, 0x5A, 0x68}},
	"xz":  {{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}},
	"tar": {{0x75, 0x73, 0x74, 0x61, 0x72}},
	"dmg": {{0x78, 0x01}},             // Apple Disk Image
	"cab": {{0x4D, 0x53, 0x43, 0x46}}, // Windows Cabinet文件

	// 音视频格式
	"mp3":  {{0x49, 0x44, 0x33}, {0xFF, 0xFB}},
	"wav":  {{0x52, 0x49, 0x46, 0x46}},
	"flac": {{0x66, 0x4C, 0x61, 0x43}},
	"aac":  {{0xFF, 0xF1}, {0xFF, 0xF9}},
	"ogg":  {{0x4F, 0x67, 0x67, 0x53}},
	"mp4":  {{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x6D, 0x70, 0x34, 0x32}, {0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}},
	"avi":  {{0x52, 0x49, 0x46, 0x46}},
	"mov":  {{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x71, 0x74, 0x20, 0x20}},
	"wmv":  {{0x30, 0x26, 0xB2, 0x75, 0x8E, 0x66, 0xCF, 0x11, 0xA6, 0xD9}},
	"flv":  {{0x46, 0x4C, 0x56, 0x01}},
	"mkv":  {{0x1A, 0x45, 0xDF, 0xA3}},
	"webm": {{0x1A, 0x45, 0xDF, 0xA3}},
	"m4a":  {{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x4D, 0x34, 0x41, 0x20}}, // M4A音频
	"aiff": {{0x46, 0x4F, 0x52, 0x4D}},                                                 // AIFF音频

	// 其他安全格式
	"swf":    {{0x46, 0x57, 0x53}},
	"ps":     {{0x25, 0x21, 0x50, 0x53}}, // PostScript
	"eps":    {{0x25, 0x21, 0x50, 0x53, 0x2D, 0x41, 0x64, 0x6F, 0x62, 0x65}},
	"sqlite": {{0x53, 0x51, 0x4C, 0x69, 0x74, 0x65, 0x20, 0x66, 0x6F, 0x72, 0x6D, 0x61, 0x74, 0x20, 0x33, 0x00}},
	"ttf":    {{0x00, 0x01, 0x00, 0x00, 0x00}},
	"otf":    {{0x4F, 0x54, 0x54, 0x4F}}, // OpenType字体
	"woff":   {{0x77, 0x4F, 0x46, 0x46}},
	"woff2":  {{0x77, 0x4F, 0x46, 0x32}},
	"eot":    {{0x4C, 0x50}},                                                                               // Embedded OpenType
	"ics":    {{0x42, 0x45, 0x47, 0x49, 0x4E, 0x3A, 0x56, 0x43, 0x41, 0x4C, 0x45, 0x4E, 0x44, 0x41, 0x52}}, // iCalendar
	"vcf":    {{0x42, 0x45, 0x47, 0x49, 0x4E, 0x3A, 0x56, 0x43, 0x41, 0x52, 0x44}},                         // vCard

	// 文本文件不验证
	"txt":  nil,
	"csv":  nil,
	"xml":  nil,
	"html": nil,
	"htm":  nil,
	"js":   nil,
	"css":  nil,
	"json": nil,
	"md":   nil,
	"conf": nil,
	"ini":  nil,
	"log":  nil,
	"sh":   nil,
	"bat":  nil,
	"php":  nil,
	"py":   nil,
	"rb":   nil,
	"go":   nil,
	"yaml": nil,
	"yml":  nil,
	"toml": nil,
	"pem":  nil, // 证书文件
	"crt":  nil, // 证书文件
	"key":  nil, // 密钥文件
}

// 魔数验证辅助函数
func (self *Service) validateMagicNumber(fileExt string, magicNum []byte) error {
	// 转换为小写并去掉点
	normalizedExt := strings.ToLower(strings.TrimPrefix(fileExt, "."))

	if expectedMagics, ok := magicNumbers[normalizedExt]; ok {
		// 如果魔数映射中该扩展名对应的魔数为nil，表示不验证（如文本文件）
		if expectedMagics == nil {
			return nil
		}

		// 检查所有可能的魔数序列
		for _, expectedMagic := range expectedMagics {
			if len(magicNum) < len(expectedMagic) {
				continue
			}

			match := true
			for i, b := range expectedMagic {
				if magicNum[i] != b {
					match = false
					break
				}
			}

			if match {
				return nil
			}
		}

		return errors.New("file content does not match the file extension")
	}

	// 如果该扩展名不在魔数映射中，则不验证（避免对未知类型造成阻碍）
	return nil
}

func (self *Service) generateUploadName(contentType, originalFileName string) (fileType, fileDir, fileName, fileExt string, err error) {
	if contentTypes := strings.Split(contentType, "/"); len(contentTypes) > 1 {
		if originalFileName == "" {
			fileExt = "." + contentTypes[1]
		} else {
			fileExt = filepath.Ext(originalFileName)
		}

		// 创建后缀过滤map
		allowedMap := make(map[string]bool)
		for _, ext := range self.Config.AllowedFileExtensions {
			allowedMap[strings.ToLower(strings.TrimSpace(ext))] = true
		}

		disallowedMap := make(map[string]bool)
		for _, ext := range self.Config.DisallowedFileExtensions {
			disallowedMap[strings.ToLower(strings.TrimSpace(ext))] = true
		}

		// 后缀名过滤
		normalizedExt := strings.ToLower(strings.TrimSpace(fileExt))

		// 1. 黑名单检查
		if len(disallowedMap) > 0 {
			if disallowedMap[normalizedExt] {
				return "", "", "", "", fmt.Errorf(predefined.ErrorFileExtensionNotAllowed)
			}
		}

		// 2. 白名单检查
		if len(allowedMap) > 0 {
			if !allowedMap[normalizedExt] {
				return "", "", "", "", fmt.Errorf(predefined.ErrorFileExtensionNotAllowed)
			}
		}

		// 类型过滤
		fileType = predefined.UploadTypeFile
		pathType := self.Config.UploadFileDirectory
		switch contentTypes[0] {
		case "image":
			fileType = predefined.UploadTypeImage
			pathType = self.Config.UploadImageDirectory
		case "video", "audio":
			switch contentTypes[1] {
			case "x-ms-asf", "avi", "x-ivf", "x-mpeg", "mp4", "mpeg4", "x-sgi-movie", "mpeg", "x-mpg", "mpg", "vnd.rn-realvideo", "x-ms-wm", "x-ms-wmv", "x-ms-wmx", "x-mei-aac",
				"aiff", "basic", "x-liquid-file", "x-liquid-secure", "x-la-lms", "mpegurl", "mid", "x-musicnet-download", "x-musicnet-stream", "mp1", "mp2", "mp3", "rn-mpeg", "scpls",
				"vnd.rn-realaudio", "x-pn-realaudio", "x-pn-realaudio-plugin", "wav", "x-ms-wax", "x-ms-wma":
				fileType = predefined.UploadTypeMedia
				pathType = self.Config.UploadMediaDirectory
			}
		case "application":
			switch contentTypes[1] {
			case "msword",
				"vnd.ms-excel",
				"vnd.ms-powerpoint",
				"vnd.ms-visio.drawing",
				"vnd.openxmlformats-officedocument.wordprocessingml.document",
				"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"vnd.openxmlformats-officedocument.presentationml.presentation",
				"pdf":
				fileType = predefined.UploadTypeDocument
				pathType = self.Config.UploadDocumentDirectory
			case "x-zip-compressed":
				fileType = predefined.UploadTypeCompression
				pathType = self.Config.UploadCompressionDirectory
			case "octet-stream":
				switch fileExt {
				case ".xmind", ".rp":
					fileType = predefined.UploadTypeDocument
					pathType = self.Config.UploadDocumentDirectory
				case ".rar", ".gz", ".bz2":
					fileType = predefined.UploadTypeCompression
					pathType = self.Config.UploadCompressionDirectory
				}
			}
		}

		now := time.Now()
		fileDir = filepath.Clean(fmt.Sprintf("%s%s/%d/%d/%d", self.Config.UploadDirectory, pathType, now.Year(), now.Month(), now.Day()))
		if _, err = os.Stat(fileDir); err != nil {
			err = os.MkdirAll(fileDir, 0755)
		}

		// 使用线程安全的随机数生成器
		randMutex.Lock()
		fileName = fmt.Sprintf("%v", globalRand.Uint64())
		randMutex.Unlock()
	} else {
		err = errors.New("parse Content-Type failed")
	}

	return
}

func (self *Service) relativePath(sPath string) (newPath string) {
	uploadDir := filepath.Clean(self.Config.UploadDirectory)
	newPath = sPath

	// 防止路径遍历攻击 - 增强安全检查
	if strings.Contains(sPath, "..") || strings.HasPrefix(sPath, "/") ||
		strings.Contains(sPath, "\\") || strings.Contains(sPath, "//") {
		return ""
	}

	// 确保路径在允许的目录内
	cleanPath := filepath.Clean(sPath)
	if !strings.HasPrefix(cleanPath, uploadDir) {
		return ""
	}

	if iPos := strings.Index(sPath, uploadDir); iPos > -1 {
		newPath = sPath[iPos+len(uploadDir):]
	}

	return strings.TrimLeft(newPath, "/")
}

func (self *Service) saveFormFile(r *http.Request, fileFormName string, bEncode bool, remark string) (fm *model.Upload, err error) {
	if bEncode {
		form_value := r.FormValue(fileFormName)
		if form_value != "" {
			if datas := strings.Split(form_value, ":"); len(datas) > 1 {
				fm, err = self.encodeDataToFile(datas[1], self.Config.AppendFileExt)
			}
		} else {
			err = errors.New(predefined.ErrorNotFoundFormData)
		}
	} else {
		if f, fh, fErr := r.FormFile(fileFormName); fErr != nil {
			err = fErr
		} else {
			defer f.Close()

			// 生成文件路径
			contentType := fh.Header.Get("Content-Type")
			originalFileName := fh.Filename
			fileType, fileDir, fileName, fileExt, genErr := self.generateUploadName(contentType, originalFileName)
			if genErr != nil {
				return nil, genErr
			}

			if self.Config.AppendFileExt {
				fileName = fileName + fileExt
			}

			filePath := fileDir + "/" + fileName

			// 创建目标文件
			dst, createErr := os.Create(filePath)
			if createErr != nil {
				return nil, createErr
			}
			defer func() {
				dst.Close()
				// 如果发生错误，删除可能创建的文件
				if err != nil {
					os.Remove(filePath)
				}
			}()

			// 使用缓冲写入提升大文件性能
			bufferedWriter := bufio.NewWriterSize(dst, 4*1024*1024) // 4MB buffer
			defer bufferedWriter.Flush()

			// 读取文件头（最多512字节）用于魔数验证
			headerBuffer := make([]byte, 512)
			n, errRead := f.Read(headerBuffer)
			if errRead != nil && errRead != io.EOF {
				return nil, errRead
			}

			// 验证魔数
			if err = self.validateMagicNumber(fileExt, headerBuffer[:n]); err != nil {
				os.Remove(filePath) // 删除无效文件
				return nil, fmt.Errorf("invalid file content: %w", err)
			}

			// 将文件头写入目标文件
			if _, err = bufferedWriter.Write(headerBuffer[:n]); err != nil {
				return nil, err
			}

			// 流式复制剩余文件内容
			size, copyErr := io.Copy(bufferedWriter, f)
			if copyErr != nil {
				return nil, copyErr
			}

			// 确保缓冲区数据写入磁盘
			if err = bufferedWriter.Flush(); err != nil {
				return nil, err
			}

			// 总文件大小 = 头部大小 + 剩余内容大小
			totalSize := int64(n) + size

			// 构建文件信息
			fm = &model.Upload{
				ID:               primitive.NewObjectID(),
				Type:             fileType,
				Storage:          predefined.UploadStorageLocal,
				UserID:           self.userID,
				ContentType:      contentType,
				OriginalFileName: originalFileName,
				FileName:         fileName,
				FileExt:          fileExt,
				Path:             self.relativePath(filePath),
				URL:              self.relativePath(filePath),
				Size:             totalSize,
			}
		}
	}

	if err == nil && fm != nil {
		uploadModel := upload.NewModel(self.M)
		if remark != "" {
			fm.Remark = remark
		}

		// 保存到数据库
		_, err = uploadModel.Create(context.Background(), fm)

		// 数据库失败时清理文件
		if err != nil {
			// 构建完整文件路径
			fullPath := filepath.Join(self.Config.UploadDirectory, fm.Path)
			// 删除已上传的文件
			if removeErr := os.Remove(fullPath); removeErr != nil {
				// 记录清理失败日志
				return nil, fmt.Errorf("Failed to clean upload file path: %w, error: %w", fullPath, removeErr)
			}
			return nil, fmt.Errorf("database error: %w, file cleaned", err)
		}
	}

	return
}
func (self *Service) encodeDataToFile(encode_string string, appendFileExt bool) (*model.Upload, error) {
	err := errors.New(predefined.ErrorNotFoundEncodeData)
	if encode_strings := strings.Split(encode_string, ";"); len(encode_strings) > 1 {
		if encode_datas := strings.Split(encode_strings[1], ","); len(encode_datas) > 1 {
			encode_type := encode_datas[0]
			encode_data := encode_datas[1]
			switch encode_type {
			case "base64":
				var file_data []byte
				if file_data, err = base64.StdEncoding.DecodeString(encode_data); err == nil {
					contentType := encode_strings[0]
					originalFileName := ""

					// 直接生成文件信息
					fileType, fileDir, fileName, fileExt, genErr := self.generateUploadName(contentType, originalFileName)
					if genErr != nil {
						return nil, genErr
					}

					if appendFileExt {
						fileName = fileName + fileExt
					}

					filePath := fileDir + "/" + fileName

					// 魔数验证
					magicNum := file_data
					if len(magicNum) > 512 {
						magicNum = magicNum[:512]
					}
					if err := self.validateMagicNumber(fileExt, magicNum); err != nil {
						return nil, err
					}

					// 写入文件
					if err = os.WriteFile(filePath, file_data, 0755); err != nil {
						return nil, err
					}

					// 获取文件信息
					fileinfo, err := os.Stat(filePath)
					if err != nil {
						return nil, err
					}

					// 返回文件模型
					return &model.Upload{
						ID:               primitive.NewObjectID(),
						Type:             fileType,
						Storage:          predefined.UploadStorageLocal,
						UserID:           self.userID,
						ContentType:      contentType,
						OriginalFileName: originalFileName,
						FileName:         fileName,
						FileExt:          fileExt,
						Path:             self.relativePath(filePath),
						URL:              self.relativePath(filePath),
						Size:             fileinfo.Size(),
					}, nil
				}
			}
		}
	}
	return nil, err
}
