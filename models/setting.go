package models

import (
    "errors"
    "fmt"
    log "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "strconv"
)

type Setting struct {
    SettingKey   string `json:"settingKey" gorm:"primaryKey;" extensions:"!x-nullable"`
    SettingValue string `json:"settingValue" gorm:"not null;" extensions:"!x-nullable"`
    SettingType  string `json:"-" gorm:"not null;" extensions:"!x-nullable"`
}

const (
    MinDuration = "min_duration"
    ReqInterval = "req_interval"
)

func InitSettings() error {
    if err := Db.FirstOrCreate(
        &Setting{SettingKey: ReqInterval, SettingValue: "15", SettingType: "int"}).Error; err != nil {
        log.Errorf("[Setting] Init error: %s", err)
        return err
    }

    return nil
}

func GetValue(settingKey string) (interface{}, error) {
    sett := Setting{}

    if err := Db.Table("settings").First(&sett, &Setting{SettingKey: settingKey}).Error; err != nil {
        log.Errorf("[GetValue] Error retreiving setting: %s", err)
        return nil, err
    }

    switch sett.SettingType {
    case "int":
        i, err := strconv.Atoi(sett.SettingValue)
        return i, err
    case "string":
        return sett.SettingValue, nil
    case "bool":
        return sett.SettingValue == "true", nil
    }

    return nil, fmt.Errorf("unknown settings type '%s'", sett.SettingType)
}

func (setting *Setting) Save() error {
    if err := Db.Model(&setting).Where("setting_key = ? ", setting.SettingKey).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            if err2 := Db.Create(&setting).Error; err2 != nil {
                return err2
            }
        } else {
            log.Errorf("[SaveValue] Error retreiving setting: %s", err)
            return err
        }
    }

    return nil
}