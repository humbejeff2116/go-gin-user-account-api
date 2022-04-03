
package responses

import (
    "go-gin-user-account-api/models"
)

type UserResponse struct {
    Status  int `json:"status"`
    Error bool  `json:"error"`
    ErrorData map[string]interface{} `json:"erroData"`
    Message string `json:"message"`
    Data map[string]interface{} `json:"data"`
}

type Response struct {
    Status  int `json:"status"`
    Error bool  `json:"error"`
    Message string `json:"message"`
    Data []models.UserModel `json:"data"`
}