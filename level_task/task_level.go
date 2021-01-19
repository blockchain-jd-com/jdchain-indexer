package level_task

import "strconv"

type TaskLevel int

const (
	// 已完成
	MetaInfoLevel0 TaskLevel = 0
	// 索引node数据
	MetaInfoLevel1 TaskLevel = 1
	// 索引link数据
	MetaInfoLevel2 TaskLevel = 2
)

func (level TaskLevel) String() string {
	return strconv.Itoa(int(level))
}

func (level TaskLevel) Upgrade() TaskLevel {
	switch level {
	case MetaInfoLevel1:
		return MetaInfoLevel2
	case MetaInfoLevel2:
		return MetaInfoLevel0
	}
	return MetaInfoLevel0
}
