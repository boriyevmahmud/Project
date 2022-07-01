package casbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mahmud3253/Project/Api-Gateway/api/auth"
	"github.com/mahmud3253/Project/Api-Gateway/api/model"
	"github.com/mahmud3253/Project/Api-Gateway/config"
)

type JwtRoleStruct struct {
	enforce    *casbin.Enforcer
	config     config.Config
	jwtHandler auth.JwtHandler
}

// NewAuthorizer is a middleware for gin to get role and
// allow or deny to endpoints
func NewJwtRoleStruct(e *casbin.Enforcer, c config.Config, jwtHandler auth.JwtHandler) gin.HandlerFunc {
	conf := &JwtRoleStruct{
		enforce:    e,
		config:     c,
		jwtHandler: jwtHandler,
	}

	return func(c *gin.Context) {

		allow, err := conf.CheckPermission(c.Request)
		if err != nil {
			v,_:=err.(*jwt.ValidationError)
			if v.Errors==jwt.ValidationErrorExpired{
				conf.RequireRefresh(c)
			}else{
				conf.RequirePermission(c)
			}
		}else if !allow{
			conf.RequirePermission(c)
		}

	}

}

// CheckPermission checks whether role is allowed to use certain endpoint
func (a *JwtRoleStruct) CheckPermission(r *http.Request) (bool, error) {
	role, err := a.GetRole(r)
	if err != nil {
		return false, err
	}

	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforce.Enforce(role, path, method)
	fmt.Println(allowed,err,"---------------------")
	if err != nil {
		panic(err)
	}

	return allowed, nil
}

// GetRole gets role from Authorization header if there
//parsed and in role got from role claim. If there is
//unathorized
func (a *JwtRoleStruct) GetRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwt.MapClaims
		err    error
	)

	jwToken := r.Header.Get("Authorization")
	if jwToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwToken, "Basic") {
		return "unauthorized", nil
	}

	a.jwtHandler.Token = jwToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "", err
	}

	if claims["role"].(string) == "authorized" {
		role = "authorized"
	} else {
		role = "unknown"
	}
	return role, nil
}


// RequireRefresh aborts request with 403 status
func (a *JwtRoleStruct) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, model.ResponseError{
		Error: model.ServerError{
			Status:  "UNAUTHORIZED",
			Message: "Token is expired",
		},
	})
	c.AbortWithStatus(401)
}

func (a *JwtRoleStruct) RequirePermission(c *gin.Context){
	c.AbortWithStatus(401)
}