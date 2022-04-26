package paginator

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

const MaxResults = 3

// OffsetPaginater is a paginator that uses offset and limit to paginate
// TODO: we should favor cursor pagination over offset pagination
type OffsetPaginater struct {
	Total int `json:"total"`
	// length of items per page
	ItemsThisPage int `json:"items_this_page"`
	// NextCursor is the next page cursor.
	// we are returning a base64 encoded string to prevent the client from
	// knowing the offset value.
	NextCursor string        `json:"next"`
	Items      []interface{} `json:"items"`
	offset     int
}

func (p *OffsetPaginater) GetOffset() int {
	return p.offset
}

func (*OffsetPaginater) GetLimit() int {
	return MaxResults
}

func (p *OffsetPaginater) Paginate(items []interface{}, total int) *OffsetPaginater {
	p.Items = items
	p.Total = total
	p.ItemsThisPage = len(items)
	if p.Total == 0 || p.Total < MaxResults {
		p.NextCursor = ""
	}
	return p
}

func NewOffsetPaginater(cursor string) (*OffsetPaginater, error) {
	if cursor == "" {
		return &OffsetPaginater{offset: 0, NextCursor: strToBase64("2")}, nil
	}
	pager := &OffsetPaginater{}
	offset, err := strconv.Atoi(base64ToStr(cursor))
	if err != nil {
		return nil, err
	}
	if offset <= 1 {
		pager.offset = 0
	} else {
		pager.offset = (offset - 1) * MaxResults
	}
	nextPage := strconv.Itoa(offset + 1)
	pager.NextCursor = strToBase64(nextPage)
	return pager, nil
}

func strToBase64(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func base64ToStr(data string) string {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("I got an error:", err)
		return ""
	}
	return string(decoded)
}
