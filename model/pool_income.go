package model

import (
	"github.com/shopspring/decimal"
)

type PoolIncome struct {
	PoolId         int              `gorm:"comment:'矿池ID'" json:"pool_id"`
	PoolSpace      int              `gorm:"type:bigint;comment:'矿池空间大小'" json:"pool_space"`
	PoolName       string           `gorm:"size:255" json:"pool_name"`
	PoolNameEN     string           `gorm:"size:255" json:"pool_name_en"`
	DividendNum    int              `gorm:"comment:'分红期数'" json:"dividend_num"`
	StartTime      int64            `gorm:"type:int" json:"start_time"`
	EndTime        int64            `gorm:"type:int" json:"end_time"`
	CycleIncome    decimal.Decimal  `gorm:"type:decimal(20,8);comment:'周期收益:挖矿时钱包实际收益'" json:"cycle_income"`
	CycleRewards   decimal.Decimal  `gorm:"type:decimal(20,8);comment:'周期爆块奖励'" json:"cycle_rewards"`
	InputIncome    decimal.Decimal  `gorm:"type:decimal(20,8);comment:'管理员指定分配的收益'" json:"input_income"`
	Dividend       decimal.Decimal  `gorm:"type:decimal(20,8);comment:'实际分配收益'" json:"dividend"`
	UserIncome     decimal.Decimal  `gorm:"type:decimal(20,8);comment:'用户本期收益'" json:"user_income"`
	DecimalBalance decimal.Decimal  `gorm:"type:decimal(20,8);comment:'小数结余'" json:"decimal_balance"`
	PlatformIncome decimal.Decimal  `gorm:"type:decimal(20,8);comment:'平台自持空间收益'" json:"platform_income"`
	ServiceCharge  decimal.Decimal  `gorm:"type:decimal(20,8);comment:'服务费'" json:"service_charge"`
	//DividendTime   types.NormalTime `gorm:"comment:'分红时间'" json:"dividend_time"`
	Status         int              `gorm:"type:tinyint;comment:'状态:0-待分配，1-已分配, 2-审核中, 3-驳回, 4-统计中'" json:"status"`
	Reason         string           `gorm:"size:255;commeoont:'人工设置原因'" json:"reason"`
	AuditRemark    string           `gorm:"size:255;comment:'审核备注'" json:"audit_remark"`
	AdminId        int              `gorm:"comment:'最后一次人工设置收益的管理员id'" json:"admin_id"`
	StartHeight    int64            `gorm:"type:int;comment:'统计时间内 开始块高度'" json:"start_height"`
	EndHeight      int64            `gorm:"type:int;comment:'统计时间内 结束块高度'" json:"end_height"`
	SalePower      int              `gorm:"comment:'当期售出的总算力/G'" json:"sale_power"`
}

func (PoolIncome) TableName() string {
	return "pool_income"
}

func FindPoolIncome() (p *PoolIncome) {
	p = new(PoolIncome)
	db.Debug().Model(&PoolIncome{}).Order("id desc").Limit(1).Find(p)
	return
}