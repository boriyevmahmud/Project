package v1

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gomodule/redigo/redis"
	"github.com/mahmud3253/Project/Api-Gateway/api/auth"
	pb "github.com/mahmud3253/Project/Api-Gateway/genproto"
	l "github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/protojson"
	gomail "gopkg.in/mail.v2"
)

type CreateUserRequestBody struct {
	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName string  `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName  string  `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Posts     []*Post `protobuf:"bytes,4,rep,name=posts,proto3" json:"posts"`
}

type RegisterUserAuthReqBody struct {
	//	Id          string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id"`
	FirstName   string `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName"`
	Username    string `protobuf:"bytes,3,opt,name=Username,proto3" json:"Username"`
	PhoneNumber string `protobuf:"bytes,4,opt,name=PhoneNumber,proto3" json:"PhoneNumber"`
	Email       string `protobuf:"bytes,5,opt,name=Email,proto3" json:"Email"`
	Code        string `protobuf:"bytes,6,opt,name=Code,proto3" json:"Code"`
	Password    string `protobuf:"bytes,7,opt,name=Password,proto3" json:"Password"`
}

type RegisterResponse struct {
	UserID       string
	Accesstoken  string
	Refreshtoken string
}

type Emailver struct {
	Email string `json:"Email"`
	Code  string `json:"Code"`
}

type Post struct {
	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name        string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId      string   `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias      []*Media `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}

type LoginResponse struct {
	Id          string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id"`
	FirstName   string `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName"`
	Username    string `protobuf:"bytes,3,opt,name=Username,proto3" json:"Username"`
	PhoneNumber string `protobuf:"bytes,4,opt,name=PhoneNumber,proto3" json:"PhoneNumber"`
	Email       string `protobuf:"bytes,5,opt,name=Email,proto3" json:"Email"`
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
	bodyByte, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to set to redis", l.Error(err))
		return
	}

	err = h.redisStorage.Set(body.FirstName, string(bodyByte))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to set redis", l.Error(err))
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

// RegisterUser register user
// @Summary Register user summary
// @Description This api is using for registering user
// @Tags user
// @Accept json
// @Produce json
// @Param user body RegisterUserAuthReqBody true "user_body"
// @Succes 200 {string} Succes!
// @Router /v1/users/register [post]
func (h *handlerV1) RegisterUser(c *gin.Context) {
	var (
		body        RegisterUserAuthReqBody
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

	// Check password
	err = verifyPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		h.log.Error("your password doesn't respond to requests", l.Error(err))
		return
	}

	//Hashing password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), len(body.Password))
	fmt.Println(string(hashedPassword))

	body.Password = string(hashedPassword)
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniquess", l.Error(err))
		return
	}
	if exists.Check {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email",
		})
		h.log.Error("failed to check email uniquess", l.Error(err))
		return
	}

	exists, err = h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "username",
		Value: body.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check username uniquess", l.Error(err))
		return
	}
	if exists.Check {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This username already in use, please use another username",
		})
		h.log.Error("failed to check username uniquess", l.Error(err))
		return
	}
	// generate 6 digits code for sending gmail
	min := 99999
	max := 1000000
	rand.Seed(time.Now().UnixNano())
	gen := rand.Intn((max - min) + min)
	code := strconv.Itoa(gen)

	body.Code = code

	bodyByte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to set to redis", l.Error(err))
		return
	}

	// writing to redis
	err = h.redisStorage.SetWithTTL(body.Email, string(bodyByte), int64(5*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to set redis", l.Error(err))
		return
	}
	fmt.Printf("%T", code)
	fmt.Println(code)
	SendEmail(body.Email, code)

}

// VerifyUser verify user
// @Description This api using for verifying registered user
// @Tags user
// @Accept json
// @Produce json
// @Param user body Emailver true "user body"
// @Succes 200 {string} success
// @Router /v1/users/verfication [post]
func (h *handlerV1) VerifyUser(c *gin.Context) {
	var dataemail Emailver
	var jspbMarshal protojson.MarshalOptions

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&dataemail)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	dataemail.Email = strings.TrimSpace(dataemail.Email)
	dataemail.Email = strings.ToLower(dataemail.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	getRedis, err := redis.String(h.redisStorage.Get(dataemail.Email))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to getting redis for write db", l.Error(err))
		return
	}

	var redisBody *pb.CreateUserAuthReqBody
	_ = json.Unmarshal([]byte(getRedis), &redisBody)

	fmt.Println(redisBody)

	if dataemail.Code == redisBody.Code {
		_, err := h.serviceManager.UserService().RegisterUser(ctx, redisBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to writing db", l.Error(err))
			return
		} else {
			h.log.Error("failed to writing db", l.Error(err))
		}
	} else if dataemail.Code != redisBody.Code {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("your code is wrong", l.Error(err))
			return
		}
	}
	id, _ := uuid.NewUUID()

	h.jwtHandler = auth.JwtHandler{
		Sub:  id.String(),
		Iss:  "client",
		Role: "authorized",
		Log:  h.log,
	}
	access, refresh, err := h.jwtHandler.GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating jwt",
		})
		h.log.Error("your code is wrong", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, &RegisterResponse{
		UserID:       id.String(),
		Accesstoken:  access,
		Refreshtoken: refresh,
	})
}

var email string
var password string

// Login login user
// @Description This api using for logging registered user
// @Tags user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param password path string true "Password"
// @Succes 200 {string} LoginResponse
// @Router /v1/users/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	email = c.Param("email")
	password = c.Param("password")
	fmt.Println(password)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	userData, err := h.serviceManager.UserService().LoginUser(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to getting datas", l.Error(err))
		return
	}
	userData.Password = ""
	c.JSON(http.StatusOK, userData)

}

func SendEmail(email, code string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "testapigomail@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("code:", "dfsdfdsf")

	m.SetBody("text/plain", code)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "testapigomail@gmail.com", "cpebajsbmuddenig")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func verifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 32
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}
