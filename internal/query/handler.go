package query

import (
	"github.com/Omotolani98/insightly/internal/storage"
	querypb "github.com/Omotolani98/insightly/proto/query"
	"gorm.io/gorm"
)

type QueryServer struct {
	querypb.UnimplementedLogQueryServer
	db *gorm.DB
}

func NewQueryServer(db *gorm.DB) *QueryServer {
	return &QueryServer{db: db}
}

func (s *QueryServer) GetSummaries(req *querypb.GetReq, stream querypb.LogQuery_GetSummariesServer) error {
	var summaries []storage.Summary

	if err := s.db.
		Where("stream = ?", req.Stream).
		Order("window_start desc").
		Limit(int(req.Limit)).
		Find(&summaries).Error; err != nil {
		return err
	}

	for _, summary := range summaries {
		resp := &querypb.SummaryResp{
			Stream:      summary.Stream,
			WindowStart: summary.WindowStart.Format("2006-01-02T15:04:05Z"),
			WindowEnd:   summary.WindowEnd.Format("2006-01-02T15:04:05Z"),
			Text:        summary.Text,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}
