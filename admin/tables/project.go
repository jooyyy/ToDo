package tables

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"todo/model"
	"todo/storage/mysql"
)

func GetProjectTable(ctx *context.Context) (t table.Table) {
	t = table.NewDefaultTable(table.DefaultConfig())

	info := t.GetInfo()
	info.HideRowSelector()
	info.AddField("ID", "id", db.Int).FieldSortable().FieldHide()
	info.AddField("Name", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator:types.FilterOperatorLike})
	info.AddActionButton("导出表格", action.Ajax("/admin/export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id, err := strconv.ParseInt(ctx.PostForm().Get("id"), 10, 64)
			if err != nil {
				return false, "导出失败", err.Error
			}
			err = exportExel(ctx, id)
			if err != nil {
				return false, "导出失败：" + err.Error(), err.Error
			}
			return true, "导出成功", ""
		}))
	info.SetTable("projects").SetTitle("Projects Manager").SetDescription("")
	info.HideExportButton()

	formList := t.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldHide()
	formList.AddField("Name", "name", db.Varchar, form.Text).FieldMust()
	formList.SetTable("projects").SetTitle("Projects Manager").SetDescription("")

	return
}

func exportExel(ctx *context.Context, projectId int64) error {
	now := time.Now()
	month := now.Format("2006-01")
	if now.Day() < 10 {
		month = now.Add(-10 * time.Hour * 24).Format("2006-01")
	}
	var tasks []model.Task
	err := mysql.DB.Where("project_id = ?", projectId).
		Where("date < ? and date > ?", month + "-31", month).
		Find(&tasks).Error
	if err != nil {
		return err
	}

	outputFile := xlsx.NewFile()
	sheet, err := outputFile.AddSheet(month)
	if err != nil {
		return err
	}

	titleRow := sheet.AddRow()
	for _, one := range []string{"日期", "工作内容", "类型", "时长"} {
		cell := titleRow.AddCell()
		cell.Value = one
	}
	sheet.Cols[0].Width = 15
	sheet.Cols[1].Width = 40

	var totalDuration time.Duration
	for _, one := range tasks {
		taskRow := sheet.AddRow()
		var taskDuration time.Duration
		for _, item := range one.Durations {
			start, error := time.ParseInLocation(TIME_STANDARD, item.StartTime, time.Local)
			if error != nil {
				return errors.New(fmt.Sprintf("时间格式错误: %s %s", one.Date, one.Name))
			}
			end, error := time.ParseInLocation(TIME_STANDARD, item.EndTime, time.Local)
			if error != nil {
				return errors.New(fmt.Sprintf("时间格式错误: %s %s", one.Date, one.Name))
			}
			taskDuration += end.Sub(start)
		}
		totalDuration += taskDuration
		taskRow.AddCell().Value = one.Date
		taskRow.AddCell().Value = one.Name
		taskRow.AddCell().Value = "远程"
		taskRow.AddCell().Value = fmt.Sprintf("%.1f", taskDuration.Hours())
	}
	amountRow := sheet.AddRow()
	amountRow.AddCell().Value = "总计"
	amountRow.AddCell()
	amountRow.AddCell()
	amountRow.AddCell().Value = fmt.Sprintf("%.1f", totalDuration.Hours())

	buf := new(bytes.Buffer)
	err = outputFile.Write(buf)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("joy - %s 工时.xlsx", month)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	return outputFile.Save(filepath.Join(dir, fileName))
}
