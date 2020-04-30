package pagination

import (
	"github.com/vfw4g/base/errors"
	"math"
)

// Pagination
type Pagination struct {
	//查询结果列表
	Items []interface{} `json:"items"`
	//是否为最后一页
	IsLastPage bool `json:"isLastPage"`
	//总页数
	Total int `json:"total"`
	//总记录数
	Count int64 `json:"count"`
	//每页记录数
	Limit int `json:"limit"`
	//当前为第几页
	PageNo int `json:"pageNo"`
	//是否为第一页
	IsFirstPage bool `json:"isFirstPage"`
	//当前页记录数
	PageCount int64 `json:"pageCount"`
}

//实例化Pageable
func New(limit, pageNo int) (p *Pagination, err error) {
	if pageNo == 0 || limit == 0 {
		return nil, errors.New("limit or pageNo must not be zero")
	}
	p = &Pagination{
		Limit:  limit,
		PageNo: pageNo,
	}
	if pageNo <= 0 {
		p.PageNo = 1
	}
	if p.PageNo == 1 {
		p.IsFirstPage = true
	}
	return p, err
}

func (p *Pagination) SetCount(count int64) {
	p.Count = count
	p.Total = int(math.Ceil(float64(count) / float64(p.Limit)))
	if p.Total == 0 {
		p.IsFirstPage = true
		p.IsLastPage = true
	}
	if p.PageNo > p.Total {
		p.PageNo = p.Total
	}
	if p.PageNo == p.Total {
		p.IsLastPage = true
	}
	if p.PageNo == 1 {
		p.IsFirstPage = true
	}
}

func (p *Pagination) GetStart() int {
	return p.Limit * (p.PageNo - 1)
}

func (p *Pagination) AddItem(item interface{}) {
	p.Items = append(p.Items, item)
	p.PageCount = int64(len(p.Items))
}
