package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/mahmud3253/Project/Api-Gateway/genproto"
	l "github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

type CreateUserRequestBody struct {
	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName string  `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName  string  `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Posts     []*Post `protobuf:"bytes,4,rep,name=posts,proto3" json:"posts"`
}

type Post struct {
	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name        string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId      string   `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias      []*Media `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}

type Media struct {
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type"`
	Link string `protobuf:"bytes,3,opt,name=link,proto3" json:"link"`
}

// CreateUser creates user
// @Summary Create user summary
// @Description This Api is using for creating new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body CreateUserRequestBody true "user body"
// @Success 200 {string} Succes!
// @Router /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	fmt.Println(&body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser get user by id
// @Summary Get user user summary
// @Description This api is using for getting user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Succes 200 {string} CreateUserRequestBody
// @Router /v1/users/getbyid/{id} [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))

	defer cancel()

	response, err := h.serviceManager.UserService().GetByIdUser(
		ctx, &pb.GetByIdRequest{
			UserId: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, response)
}

// DeleteUser deletes user
// @Summary Delete user summary
// @Description This Api is using for deleting user
// @Tags user
// @Accecpt json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} Succes!
// @Router /v1/users/delete/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))

	defer cancel()

	response, err := h.serviceManager.UserService().DeleteById(
		ctx, &pb.DeleteByIdReq{
			UserId: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delte user", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}

// UpdateUser update user
// @Summary Update user
// @Description This Api is using for updating user with posts
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Param user body CreateUserRequestBody true "user body"
// @Success 200 {string} Succes!
// @Router /v1/users/update/{id}  [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	id := c.Param("id")

	var body pb.User

	c.ShouldBindJSON(&body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateById(
		ctx, &pb.UpdateByIdReq{
			UserId: id,
			Users:  &body,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}

// ListUser list user
// @Summary ListUser user
// @Description This Api is using for listing users
// @Tags user
// @Accept json
// @Produce json
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Success 200 {string} []CreateUserRequestBody
// @Router /v1/users/listuser  [get]
func (h *handlerV1) ListUser(c *gin.Context) {
	p := c.Query("page")
	l := c.Query("limit")
	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parsing limit or page to conv")
		return
	}
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().ListUser(
		ctx, &pb.ListUserReq{
			Page:  page,
			Limit: limit,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to getting user listing")
		return
	}
	c.JSON(http.StatusOK, response)
}
