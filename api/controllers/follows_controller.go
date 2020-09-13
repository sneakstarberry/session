package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sneakstarberry/session/api/auth"
	"github.com/sneakstarberry/session/api/models"
	"github.com/sneakstarberry/session/api/utils/formaterror"
)

func (server *Server) FollowUser(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	userID := c.Param("id")
	pid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// check if the post exist:
	user2 := models.User{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&user2).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	follow := models.Follow{}
	follow.UserAID = user.ID
	follow.UserBID = user2.ID

	followCreated, err := follow.SaveFollow(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": followCreated,
	})
}

func (server *Server) GetFollows(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	userID := c.Param("id")

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		fmt.Println("this is the error: ", err)
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Check if the post exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", pid).Take(&user).Error
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	follow := models.Follow{}

	follows, err := follow.GetFollowsInfo(server.DB, pid)
	if err != nil {
		errList["No_follows"] = "No Follows found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": follows,
	})
}

func (server *Server) UnFollowUser(c *gin.Context) {

	followID := c.Param("id")
	// Is a valid like id given to us?
	fid, err := strconv.ParseUint(followID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the post exist
	follow := models.Follow{}
	err = server.DB.Debug().Model(models.Follow{}).Where("id = ?", fid).Take(&follow).Error
	if err != nil {
		errList["No_follow"] = "No Follow Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this post?
	if uid != follow.UserAID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If all the conditions are met, delete the post
	_, err = follow.DeleteFollow(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Follow deleted",
	})
}
