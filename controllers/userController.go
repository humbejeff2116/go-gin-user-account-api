package controllers

import (
	"context"
	"fmt"
	"go-gin-user-account-api/configs"
	"go-gin-user-account-api/models"
	"go-gin-user-account-api/responses"
    "go-gin-user-account-api/lib"
	"log"
	"net/http"
	"time"
    "path/filepath"


	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "github.com/google/uuid" // To generate random profile image names
)


type UserController interface {
  
	SignupUser()
    LoginUser()
    UpdateUser()
	RemoveUser()
    GetUser()
	GetUsers()
	
}

var userCollection *mongo.Collection = configs.GetCollection(configs.MongodbClient, "golang", "users");

var validate = validator.New()

func SignupUser(c *gin.Context) {

    var response responses.UserResponse;

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    var user models.UserModel

    defer cancel()

    err := c.ShouldBind(&user)

    if  err != nil {

        log.Println(err)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse{
            Status: http.StatusBadRequest,
            Error: true,
            Message: "bad request format",
            ErrorData: map[string]interface{}{"data": err.Error()},
        }

        c.JSON(http.StatusBadRequest, response)

        return

    }

    //validate required fields using validator library
    validationErr := validate.Struct(&user);

    if  validationErr != nil {

        log.Println(validationErr)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse {
            Status: http.StatusBadRequest,
            Error: true,
            Message: "Required user fields validation failed",
            ErrorData: map[string]interface{}{"data": validationErr.Error()},
        }

        c.JSON(http.StatusBadRequest, response)

        return

    }

    file, err := c.FormFile("file")

    if err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusBadRequest, 
            Error: true,
            Message: "error occured while reading form file", 
        }

        c.AbortWithStatusJSON(http.StatusBadRequest, response)

        return

    }
    
    filename := filepath.Base(file.Filename)

    // Generate a random file name for file so it doesn't override any file that has already been uploaded with the same name
    newFileName := uuid.New().String() + filename
    
    // TODO... upload file to cloud storage/CDN and save image url in data base

    err = c.SaveUploadedFile(file, "../public/uploads/product_images/" + newFileName);

    if  err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusBadRequest, 
            Error: true,
            Message: "error occured while uploading file", 
        }

        c.AbortWithStatusJSON(http.StatusBadRequest, response)

        return

    }

    // hash user password
    user.Password = lib.GeneratePasswordHash(user.Password)

    newUser := models.UserModel {
        Id: primitive.NewObjectID(),
        FullName: user.FullName,
        UserName: user.UserName,
        UserEmail: user.UserEmail,
        Password: user.Password,
        ProfileImage: "../public/uploads/user_images/" + newFileName,
    }
  
    result, err := userCollection.InsertOne(ctx, newUser)

    // NOTE: return ErrorData to client only when in development
    // TODO... log error and return a custom error to client
    if err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusInternalServerError, 
            Error: true,
            Message: "failed to create user", 
            ErrorData: map[string]interface{}{ "data": err.Error() },
        }

        c.JSON(http.StatusInternalServerError, response)

        return

    }

    response = responses.UserResponse{
        Status: http.StatusCreated, 
        Message: "user created successfully", 
        Data: map[string]interface{}{"data": result },
    }

    c.JSON(http.StatusCreated, response)

}

func LoginUser(c *gin.Context) {

    var response responses.UserResponse;

    var user models.UserModel

    var savedUser models.UserModel

    err := c.BindJSON(&user);
    
    if  err != nil {

        log.Println(err)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse{
            Status: http.StatusBadRequest, 
            Error: true,
            Message: "JSON format is incorrect", 
            ErrorData: map[string]interface{}{"data": err.Error()},
        } 

        c.JSON(http.StatusBadRequest, response)

        return

    }

    //validate required fields using validator library
    validationErr := validate.Struct(&user);

    if  validationErr != nil {

        log.Println(validationErr)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse{
            Status: http.StatusBadRequest,
            Error: true, 
            Message: "JSON fields validation failed", 
            ErrorData: map[string]interface{}{"data": validationErr.Error()},
        }

        c.JSON(http.StatusBadRequest, response)

        return

    }


    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()

    err = userCollection.FindOne(ctx, bson.M{"userEmail": user.UserEmail}).Decode(&savedUser)

    // TODO... return no user found response to client if err == ErrNoDocuments

    if err != nil {

        log.Println(err);

        response = responses.UserResponse{
            Status: http.StatusInternalServerError,
            Error: true, 
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "An error occured while getting user",   
        }

        c.JSON(http.StatusInternalServerError, response)

        return

    }

    userPasswordGuess := user.Password

    hashedUserPassword:= savedUser.Password
  
    passwordErr := lib.CheckPassword(hashedUserPassword, userPasswordGuess)
  
    if passwordErr != nil {

        log.Println(passwordErr)

        response = responses.UserResponse {
            Status: http.StatusBadRequest,
            Error: true, 
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "incorrect password",   
        }

        c.JSON(http.StatusBadRequest, response)

        return

    }

    jwtToken, err := lib.GenerateJWT()

    if err != nil {

        response = responses.UserResponse{
            Status: http.StatusInternalServerError,
            Error: true, 
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "error occured while generating token",   
        }

        c.JSON(http.StatusInternalServerError, response)

        return

    }

    response = responses.UserResponse {
        Status: http.StatusOK, 
        Message: "token generated succesfully",
        Data: map[string]interface{}{"token": jwtToken },
    }

    c.JSON(http.StatusOK, response);
    
}

func GetUsers(c *gin.Context) {

    var response responses.UserResponse;

    var users []bson.D;

    // create a custom mongoDB context 
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);

    defer cancel();

    // access users using a cursor (allows us to iterate over db while holding only a subset of them in memory at a given time)
    cursor, err := userCollection.Find(ctx, bson.D{});
   
    if err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusInternalServerError, 
            Error: true,
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "error occured while getting users from database",
        }

        c.JSON(http.StatusInternalServerError, response);

        return

    }

    // close cursor to free resources it consumes in both the client application and the MongoDB server
    defer cursor.Close(ctx);

    //populate users array with all users query results
    err = cursor.All(ctx, &users);

    if  err != nil {

        response = responses.UserResponse{
            Status: http.StatusInternalServerError,
            Error: true, 
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "error",   
        }

        log.Println(err);

        c.JSON(http.StatusInternalServerError, response);

        return

    }

    response = responses.UserResponse{
        Status: http.StatusOK, 
        Message: "users gotten successfully", 
        Data: map[string]interface{}{"data": users },
    }

    c.JSON(http.StatusOK, response);

}

func GetUser(c *gin.Context) {

    var response responses.UserResponse;

    var user models.UserModel

    userId := c.Param("userId")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

    if err != nil {

        log.Println(err);

        response = responses.UserResponse{
            Status: http.StatusInternalServerError,
            Error: true, 
            ErrorData: map[string]interface{}{ "data": err.Error() },
            Message: "An error occured while getting user",   
        }

        c.JSON(http.StatusInternalServerError, response)

        return

    }

    response = responses.UserResponse{
        Status: http.StatusOK, 
        Message: "user gotten successfully", 
        Data: map[string]interface{}{"data": user },
    }

    c.JSON(http.StatusOK, response)
    
}


// find product with id and update product 
func UpdateUser(c *gin.Context) {
   
    var response responses.UserResponse;

    var updateUser models.UpdateUserModel;

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()
    
    //validate the request body
    err := c.BindJSON(&updateUser);
    
    if  err != nil {

        log.Println(err)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse{
            Status: http.StatusBadRequest, 
            Error: true,
            Message: "JSON format is incorrect", 
            ErrorData: map[string]interface{}{"data": err.Error()},
        } 

        c.JSON(http.StatusBadRequest, response)

        return

    }

    //validate required fields using validator library
    validationErr := validate.Struct(&updateUser);

    if  validationErr != nil {

        log.Println(validationErr)

        // NOTE: return ErrorData to client only when in development
        response = responses.UserResponse{
            Status: http.StatusBadRequest,
            Error: true, 
            Message: "JSON fields validation failed", 
            ErrorData: map[string]interface{}{"data": validationErr.Error()},
        }

        c.JSON(http.StatusBadRequest, response)

        return

    }

    objId, _ := primitive.ObjectIDFromHex(updateUser.Id)

    filter := bson.M{"_id" : objId}

    updateQuery := bson.M{ "$set": bson.M{updateUser.Key: updateUser.Value}}

    result, err := userCollection.UpdateOne(ctx, filter, updateQuery)

    if err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusBadRequest,
            Error: true, 
            Message: "user update failed", 
            ErrorData: map[string]interface{}{"data": err.Error()},
        } 

        c.JSON(http.StatusBadRequest, response)

        return

    }

    response = responses.UserResponse{
        Status: http.StatusOK, 
        Message: "user Updated sucessfully", 
        ErrorData: map[string]interface{}{"data": result},
    } 

    fmt.Printf("Documents matched: %v\n", result.MatchedCount)

    fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

    c.JSON(http.StatusOK, response)

}


func RemoveUser(c *gin.Context) {

    var response responses.UserResponse;

    userId := c.Param("userId")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})

    if err != nil {

        log.Println(err)

        response = responses.UserResponse{
            Status: http.StatusBadRequest, 
            Error: true,
            Message: "failed to delete user", 
            ErrorData: map[string]interface{}{"data": err.Error()},
        } 

        c.JSON(http.StatusInternalServerError, response)

        return

    }

    if result.DeletedCount < 1 {

        response =  responses.UserResponse{
            Status: http.StatusNotFound,
            Error: true, 
            Message: "user with specified id not found",   
        }

        c.JSON(http.StatusNotFound, response)

        return
    }

    response =  responses.UserResponse{
        Status: http.StatusOK, 
        Message: "user deleted successfully", 
    }

    c.JSON(http.StatusOK, response)
}