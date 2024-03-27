package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type advertisementRoutes struct {
	advertisementService service.Advertisement
	AuthService          service.Auth
}

type createAdvertisementRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Pictures    []string `json:"pictures" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
}

type getAdvertisementResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"name"`
	Description string    `json:"description"`
	Pictures    []string  `json:"pictures"`
	Price       float64   `json:"price"`
	Owner       string    `json:"owner"`
	CreatedAt   string    `json:"created_at"`
}

type createAdvertisementResponse struct {
	ID   uuid.UUID `json:"id"`
	Code int       `json:"code"`
}

func newAdvertisementRoutes(c *gin.RouterGroup, advertisementService service.Advertisement, authService service.Auth) {
	r := &advertisementRoutes{
		advertisementService: advertisementService,
		AuthService:          authService,
	}

	c.POST("/create", r.create)
	c.GET("/:id", r.get)
	c.GET("/", r.getAll)
	c.GET("/user/:id", r.getByUserID)
}

// @Summary Get advertisements by user ID
// @Description Get advertisements created by a user
// @Tags advertisements
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Success 200 {array} getAdvertisementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/advertisements/user/{id} [get]
func (r *advertisementRoutes) getByUserID(ctx *gin.Context) {

	userID := ctx.Param("id")
	IDParsed, err := uuid.Parse(userID)
	if err != nil {
		log.Errorf("failed to parse user ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	advertisements, err := r.advertisementService.GetAdvertisementsByUserID(ctx.Request.Context(), IDParsed)
	if err != nil {
		log.Errorf("failed to get advertisements by user ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := r.AuthService.GetUserByID(ctx.Request.Context(), IDParsed)

	if err != nil {
		log.Errorf("AdvertisementService.GetUserByID: cannot get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]getAdvertisementResponse, len(advertisements))
	for i, ad := range advertisements {
		response[i] = getAdvertisementResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Description: ad.Description,
			Pictures:    ad.Pictures,
			Price:       ad.Price,
			Owner:       user.Email,
			CreatedAt:   ad.CreatedAt.Format("January 2, 2006 15:04:05"),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Create Advertisement
// @Description Create advertisement
// @Tags advertisements
// @Accept json
// @Produce json
// @Param request body createAdvertisementRequest true "Advertisement Request"
// @Success 201 {object} createAdvertisementResponse
// @Security JWT
// @Failure 400 {object} createAdvertisementResponse
// @Failure 500 {object} createAdvertisementResponse
// @Failure 403 {object} createAdvertisementResponse
// @Failure default {object} createAdvertisementResponse
// @Router /api/v1/advertisements/create [post]
func (r *advertisementRoutes) create(ctx *gin.Context) {
	token, ok := bearerToken(ctx.Request)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := r.AuthService.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request createAdvertisementRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	advertisement := &entity.Advertisement{
		Title:       request.Title,
		Description: request.Description,
		Pictures:    request.Pictures,
		Price:       float64(request.Price),
		UserID:      userID,
	}

	id, err := r.advertisementService.CreateAdvertisement(*advertisement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := createAdvertisementResponse{
		ID:   id,
		Code: http.StatusCreated,
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary Get advertisement
// @Description Get advertisement by ID
// @Tags advertisements
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Advertisement ID"
// @Success 200 {object} getAdvertisementResponse
// @Failure 400 {object} getAdvertisementResponse
// @Failure 403 {object} getAdvertisementResponse
// @Failure 404 {object} getAdvertisementResponse
// @Failure 500 {object} getAdvertisementResponse
// @Router /api/v1/advertisements/{id} [get]
func (r *advertisementRoutes) get(ctx *gin.Context) {
	id := ctx.Param("id")

	idParsed, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	advertisement, err := r.advertisementService.GetAdvertisementByID(ctx.Request.Context(), idParsed)
	fmt.Println(advertisement)

	if err != nil {
		log.Errorf("failed to get advertisement by id: %v", err.Error())
		if errors.Is(err, service.ErrAdvertisementNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "advertisement not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	user, err := r.AuthService.GetUserByID(ctx.Request.Context(), advertisement.UserID)

	if err != nil {
		log.Errorf("failed to get user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := getAdvertisementResponse{
		ID:          advertisement.ID,
		Title:       advertisement.Title,
		Description: advertisement.Description,
		Pictures:    advertisement.Pictures,
		Price:       float64(advertisement.Price),
		Owner:       user.Email,
		CreatedAt:   advertisement.CreatedAt.Format("January 2, 2006 15:04:05"),
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get advertisement
// @Description Get all advertisements
// @Tags advertisements
// @Produce json
// @Success 200 {object} v1.advertisementRoutes.getAll.response
// @Router /api/v1/advertisements/ [get]
func (r *advertisementRoutes) getAll(ctx *gin.Context) {

	page := parsePage(ctx.DefaultQuery("page", "1"))
	limit := parseLimit(ctx.DefaultQuery("limit", "10"))
	sortBy := ctx.DefaultQuery("sortBy", "createdAt")
	sortOrder := ctx.DefaultQuery("sortOrder", "desc")
	minPrice := parsePrice(ctx.DefaultQuery("minPrice", "0"))
	maxPrice := parsePrice(ctx.DefaultQuery("maxPrice", "9999999"))

	advertisements, err := r.advertisementService.GetAllAdvertisements(ctx.Request.Context(), page, limit, sortBy, sortOrder, minPrice, maxPrice)
	if err != nil {
		log.Errorf("failed to get advertisements: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]getAdvertisementResponse, len(advertisements))
	for i, ad := range advertisements {
		user, err := r.AuthService.GetUserByID(ctx.Request.Context(), ad.UserID)
		if err != nil {
			log.Errorf("failed to get user details for advertisement ID %v: %v", ad.ID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		response[i] = getAdvertisementResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Description: ad.Description,
			Pictures:    ad.Pictures,
			Price:       ad.Price,
			Owner:       user.Email,
			CreatedAt:   ad.CreatedAt.Format("January 2, 2006 15:04:05"),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func parsePage(pageStr string) int {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func parseLimit(limitStr string) int {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 10
	}
	return limit
}

func parsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		return 0
	}
	return price
}
