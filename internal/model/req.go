package model

type PageReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func (p *PageReq) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = MaxInt
	}
}
