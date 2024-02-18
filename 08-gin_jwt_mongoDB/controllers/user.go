package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go_gin_jwt/databases"
	"go_gin_jwt/helpers"
	"go_gin_jwt/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 创建用户集合
var userCollection = databases.OpenCollection(databases.Client, "user")

// validator 校验
var validate = validator.New()

func HashPassword(password string) string {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(fromPassword)
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	// 密码比较 第一个参数应该是数据库密码，第二个参数应该是用户输入提供的密码
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("邮箱密码不正确")
		check = false
	}
	return check, msg
}

func Signup(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		fmt.Printf("Validation error: %+v\n", validationErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查电子邮件时发生错误"})
		log.Panic(err)
	}
	// defer cancel()

	// 在处理手机号之前处理密码(邮箱校验通过后可以保存密码进数据库然后让用户绑定手机号)
	password := HashPassword(*user.Password)
	user.Password = &password

	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查手机号码时发生错误"})
		log.Panic(err)
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "电子邮件或手机号码已存在,请勿重复注册"})
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserId, *user.UserType)
	user.Token = &token
	user.RefreshToken = &refreshToken

	result, insterErr := userCollection.InsertOne(ctx, user)
	if insterErr != nil {
		msg := fmt.Sprintf("创建用户失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var foundUser models.User
	defer cancel()
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "电子邮件或密码不正确"})
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()

	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	// 加盐
	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "找不到该用户"})
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId, *foundUser.UserType)
	helpers.UpdateAllToken(token, refreshToken, foundUser.UserId)
	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, foundUser)
}

func GetUsers(c *gin.Context) {
	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{"match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"_id", "null"}}},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}},
	}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}},
	}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表时发生错误"})
	}

	var allUsers []bson.M

	if err = result.All(ctx, &allUsers); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allUsers[0])
}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")
	err := helpers.MatchUserTypeToUid(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// 查询结果一般只使用普通的结构体变量即可,因为解码的过程只是将数据库返回的数据按照指定的结构体字段进行填充，并不会涉及到额外的内存分配或复制操作
	// 在需要修改结构体字段,或者函数参数传递,或者处理大型结构体,可以考虑使用指针
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
