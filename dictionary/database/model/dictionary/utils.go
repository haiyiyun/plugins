package dictionary

import (
	"context"
	"sort"

	"github.com/haiyiyun/plugins/dictionary/database/model"
	"github.com/haiyiyun/plugins/dictionary/predefined"
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Model) FilterByEnable() bson.D {
	filter := bson.D{
		{"enable", true},
		{"delete", false},
	}

	return filter
}

func (self *Model) FilterByKey(key string) bson.D {
	filter := bson.D{
		{"key", key},
	}

	return filter
}

func (self *Model) sort(values []*model.DictionaryValue) []*model.DictionaryValue {
	sort.Slice(values, func(i, j int) bool {
		return values[i].Order < values[j].Order
	})

	return values
}

func (self *Model) FindValues(filter bson.D) []*model.DictionaryValue {
	filter = append(filter, self.FilterByEnable()...)
	opt := self.ProjectionOne(bson.D{
		{"_id", 0},
		{"values", 1},
	})

	if sr := self.FindOne(context.TODO(), filter, opt); sr.Err() == nil {
		dict := model.Dictionary{}
		if err := sr.Decode(&dict); err == nil {
			return self.sort(dict.Values)
		}
	}

	return nil
}

func (self *Model) FindSliceLable(key string) []string {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		ss := []string{}
		for _, v := range vals {
			if v.Enable {
				ss = append(ss, v.Lable)
			}
		}

		return ss
	}

	return nil
}

func (self *Model) FindMapLable(key string) map[string]string {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		m := map[string]string{}
		for _, v := range vals {
			if v.Enable {
				m[v.Key] = v.Lable
			}
		}

		return m
	}

	return nil
}

func (self *Model) FindSliceValue(key string) []int {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		is := []int{}
		for _, v := range vals {
			if v.Enable {
				is = append(is, v.Value)
			}
		}

		return is
	}

	return nil
}

func (self *Model) FindMapValue(key string) map[string]int {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		m := map[string]int{}
		for _, v := range vals {
			if v.Enable {
				m[v.Key] = v.Value
			}
		}

		return m
	}

	return nil
}

func (self *Model) FindSliceLableValue(key string) []help.M {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		ms := []help.M{}
		for _, v := range vals {
			if v.Enable {
				m := help.M{
					"lable": v.Lable,
					"value": v.Value,
				}

				ms = append(ms, m)
			}
		}

		return ms
	}

	return nil
}

func (self *Model) FindMapLableValue(key string) help.M {
	if vals := self.FindValues(self.FilterByKey(key)); vals != nil {
		m := help.M{}
		for _, v := range vals {
			if v.Enable {
				m[v.Key] = help.M{
					"lable": v.Lable,
					"value": v.Value,
				}
			}
		}

		return m
	}

	return nil
}

func (self *Model) Lookup(key, structrue, want string) interface{} {
	var dict interface{}

	switch structrue {
	case predefined.DictionaryLookupStructureSlice:
		switch want {
		case predefined.DictionaryLookupWantLable:
			dict = self.FindSliceLable(key)
		case predefined.DictionaryLookupWantValue:
			dict = self.FindSliceValue(key)
		case predefined.DictionaryLookupWantLableAndValue:
			dict = self.FindSliceLableValue(key)
		}
	case predefined.DictionaryLookupStructureMap:
		switch want {
		case predefined.DictionaryLookupWantLable:
			dict = self.FindMapLable(key)
		case predefined.DictionaryLookupWantValue:
			dict = self.FindMapValue(key)
		case predefined.DictionaryLookupWantLableAndValue:
			dict = self.FindMapLableValue(key)
		}
	}

	return dict
}
