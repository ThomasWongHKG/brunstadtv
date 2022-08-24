package sqlc

import (
	"context"
	"encoding/json"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/samber/lo"
	"time"
)

func mapToPages(pages []getPagesRow) []common.Page {
	return lo.Map(pages, func(p getPagesRow, _ int) common.Page {
		var title common.LocaleString
		var description common.LocaleString

		_ = json.Unmarshal(p.Title.RawMessage, &title)
		_ = json.Unmarshal(p.Title.RawMessage, &description)

		return common.Page{
			ID:          int(p.ID),
			Title:       title,
			Description: description,
		}
	})
}

// GetPages returns a list of pages retrieved by ids
func (q *Queries) GetPages(ctx context.Context, ids []int) ([]common.Page, error) {
	pages, err := q.getPages(ctx, intToInt32(ids))
	if err != nil {
		return nil, err
	}
	return mapToPages(pages), err
}

// ListPages returns a list of pages
func (q *Queries) ListPages(ctx context.Context) ([]common.Page, error) {
	pages, err := q.listPages(ctx)
	if err != nil {
		return nil, err
	}
	return mapToPages(lo.Map(pages, func(p listPagesRow, _ int) getPagesRow {
		return getPagesRow(p)
	})), nil
}

// GetPermissionsForPages returns permissions for pages
func (q *Queries) GetPermissionsForPages(ctx context.Context, ids []int) ([]common.Permissions[int], error) {
	items, err := q.getPermissionsForPages(ctx, intToInt32(ids))
	if err != nil {
		return nil, err
	}

	from, _ := time.Parse("2006-01-02", "1900-01-01")
	to, _ := time.Parse("2006-01-02", "2100-01-01")

	return lo.Map(items, func(row getPermissionsForPagesRow, _ int) common.Permissions[int] {
		return common.Permissions[int]{
			ItemID: int(row.ID),
			Type:   common.TypePage,
			Availability: common.Availability{
				Published: row.Published,
				From:      from,
				To:        to,
			},
			Roles: common.Roles{
				Access: row.Roles,
			},
		}
	}), nil
}

// GetOriginal returns the requested string
func (row getPageIDsForCodesRow) GetOriginal() string {
	return row.Code.String
}

// GetResult returns the id from the query
func (row getPageIDsForCodesRow) GetResult() int {
	return int(row.ID)
}

// GetPageIDsForCodes returns ids for the requested codes
func (q *Queries) GetPageIDsForCodes(ctx context.Context, codes []string) ([]common.Conversion[string, int], error) {
	rows, err := q.getPageIDsForCodes(ctx, codes)
	if err != nil {
		return nil, err
	}
	return lo.Map(rows, func(i getPageIDsForCodesRow, _ int) common.Conversion[string, int] {
		return i
	}), nil
}