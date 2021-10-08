package tables

import (
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"strconv"
	"time"
	"todo/model"
	"todo/storage/mysql"
)

const TIME_STANDARD = "2006-01-02 15:04:05"
const TIME_DATE = "2006-01-02"

func GetTaskTable(ctx *context.Context) (t table.Table) {
	t = table.NewDefaultTable(table.DefaultConfig())
	info := t.GetInfo()
	info.HideRowSelector().HideDetailButton()
	info.AddActionIconButton(icon.Check, action.Ajax("/admin/done", func(ctx *context.Context) (success bool, msg string, data interface{}) {
		return true, "success", ""
	}), "")
	info.AddField("ID", "id", db.Int).FieldSortable().FieldHide()
	info.AddField("Date", "date", db.Varchar)
	info.AddField("Status", "status", db.Varchar).FieldEditAble(editType.Select).FieldEditOptions(taskStatusOpts).FieldDisplay(func(value types.FieldModel) interface{} {
		switch value.Value {
		case model.TaskStatusPause:
			return "暂停"
		case model.TaskStatusDoing:
			return "进行中"
		case model.TaskStatusDone:
			return "结束"
		default:
			return "未知状态"
		}
	})
	info.AddField("Name", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Duration", "durations", db.Text).FieldDisplay(func(value types.FieldModel) interface{} {
		var durations model.Durations
		err := json.Unmarshal([]byte(value.Value), &durations)
		if err != nil {
			return err.Error()
		}
		var totalDuration time.Duration
		for _, one := range durations {
			start, _ := time.ParseInLocation(TIME_STANDARD, one.StartTime, time.Local)
			end, _ := time.ParseInLocation(TIME_STANDARD, one.EndTime, time.Local)
			if one.EndTime == "" {
				end = time.Now()
			}
			totalDuration += end.Sub(start)
		}
		return fmt.Sprintf("%.1f", totalDuration.Hours())
	})
	info.SetTable("tasks").SetTitle("Tasks Manager").SetDescription("")

	formList := t.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldHide()
	formList.AddField("Date", "date", db.Varchar, form.Text).
		FieldDefault(time.Now().Format("2006-01-02")).FieldMust()
	formList.AddField("Project", "project_id", db.Int, form.SelectSingle).
		FieldOptions(getProjectsOptions()).
		FieldMust()
	formList.AddField("Name", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Status", "status", db.Text, form.SelectSingle).
		FieldOptions(taskStatusOpts).
		FieldDefault("pause").
		FieldMust()
	formList.AddField("Description", "description", db.Text, form.TextArea)
	formList.AddField("Duration", "durations", db.Text, form.TextArea).FieldDefault("[]")
	formList.SetTable("tasks").SetTitle("Tasks Manager").SetDescription("")
	formList.SetPreProcessFn(updateFunc())

	return
}

var taskStatusOpts = types.FieldOptions{
	types.FieldOption{
		Text:  "暂停",
		Value: "pause",
	},
	types.FieldOption{
		Text:  "进行中",
		Value: "doing",
	},
	types.FieldOption{
		Text:  "完成",
		Value: "done",
	},
}

func getProjectsOptions() types.FieldOptions {
	var list []model.Project
	mysql.DB.Find(&list)
	var out types.FieldOptions
	for _, one := range list {
		out = append(out, types.FieldOption{
			Text:  one.Name,
			Value: strconv.FormatInt(int64(one.ID), 10),
		})
	}
	return out
}

func updateFunc() types.FormPreProcessFn {
	return func(values form2.Values) form2.Values {
		now := time.Now().Format("2006-01-02 15:04:05")

		var durations model.Durations
		if value, ok := values["durations"]; ok {
			err := json.Unmarshal([]byte(value[0]), &durations)
			if err != nil {
				return values
			}
		}else {
			var task model.Task
			err := mysql.DB.Where("id = ?", values.Get("id")).First(&task).Error
			if err != nil {
				return values
			}
			durations = task.Durations
		}

		switch values.Get("status") {
		case model.TaskStatusPause:
			if len(durations) == 0 {
				break
			}
			if durations[len(durations)-1].EndTime == "" {
				durations[len(durations)-1].EndTime = now
			}
			break
		case model.TaskStatusDoing:
			if len(durations) == 0 {
				durations = append(durations, model.Duration{
					StartTime: now,
				})
			} else if durations[len(durations)-1].EndTime != "" {
				durations = append(durations, model.Duration{
					StartTime: now,
				})
			}
		case model.TaskStatusDone:
			if len(durations) == 0 {
				break
			} else if durations[len(durations)-1].EndTime == "" {
				durations[len(durations)-1].EndTime = now
			}
			break
		default:
			return values
		}

		values.Delete("durations")
		newValue, _ := json.Marshal(durations)
		values.Add("durations", string(newValue))

		return values
	}
}
