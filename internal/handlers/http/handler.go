package httphandler

import (
	service "github.com/DarRo9/Tenders/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service interface {
	service.TenderService
	service.BidService
}

type Handler struct {
	srv Service
	log *logrus.Logger
}

func New(srv Service, log *logrus.Logger) *Handler {
	return &Handler{srv: srv, log: log}
}

func (h *Handler) CreateRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/ping", h.PingStatus) 

		tenders := api.Group("/tenders")
		{
			tenders.GET("", h.GetAllTenders)        
			tenders.POST("/new", h.BuildTender) 
			tenders.GET("/my", h.GetOnesTenders)   

			tenders.GET("/:tenderId/status", h.GetStatusOfTender)               
			tenders.PUT("/:tenderId/status", h.RefreshTenderStatus)             
			tenders.PATCH("/:tenderId/edit", h.ChangeTender)                       
			tenders.PUT("/:tenderId/rollback/:version", h.RefreshTenderVersion) 
		}

		bids := api.Group("/bids")
		{
			bids.POST("/new", h.ConstructBid) 
			bids.GET("/my", h.GetOnesBids)   

		
			{
				bids.GET("/:id/status", h.GetStatusOfBids)                
				bids.PUT("/:id/status", h.RenewStatusOfBid)             
				bids.PATCH("/:id/edit", h.ChangeBid)                       
				bids.PUT("/:id/submit_decision", h.ApplyDecision)     
				bids.PUT("/:id/feedback", h.ApplyFeedback)             
				bids.PUT("/:id/rollback/:version", h.ReturnBidVersion) 
			}

			
			{
				bids.GET("/:id/list", h.GetBidsOfTender) 
				bids.GET("/:id/reviews", h.GetCommentsOfBid) 
			}
		}
	}

	return r
}
