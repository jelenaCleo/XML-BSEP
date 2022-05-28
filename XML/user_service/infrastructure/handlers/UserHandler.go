package handlers

import (
	pb "common/module/proto/user_service"
	"context"
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"strconv"
	"strings"
	"user/module/application/helpers"
	"user/module/application/services"
	"user/module/infrastructure/api"
)

type UserHandler struct {
	l            *log.Logger
	service      *services.UserService
	jsonConv     *helpers.JsonConverters
	validator    *validator.Validate
	passwordUtil *helpers.PasswordUtil
}

func NewUserHandler(l *log.Logger, service *services.UserService, jsonConv *helpers.JsonConverters, validator *validator.Validate,
	passwordUtil *helpers.PasswordUtil) *UserHandler {
	return &UserHandler{l, service, jsonConv, validator, passwordUtil}
}

func (u UserHandler) MustEmbedUnimplementedUserServiceServer() {
	u.l.Println("Handling MustEmbedUnimplementedUserServiceServer Users")
}

func (u UserHandler) ActivateUserAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	u.l.Println("Handling ActivateUserAccount ")
	//TODO:mzd dodati provjeru da li se uspelo ok mapirati?
	requstDto := api.MapPbToUserActivateRequest(request)
	//TODO:dodati validaciju u obliku regexa, spreciti injection napad
	err := u.validator.Struct(requstDto)
	if err != nil {

		u.l.Println("1111111111111")
		u.l.Println(err)
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, err
		//http.Error(rw, "New user dto fields aren't entered in valid format! error:"+err.Error(), http.StatusExpectationFailed) //400
	}
	u.l.Println("222222222222222")
	policy := bluemonday.UGCPolicy()
	//sanitize everything
	requstDto.Username = strings.TrimSpace(policy.Sanitize(requstDto.Username))
	requstDto.Code = strings.TrimSpace(policy.Sanitize(requstDto.Code))
	if requstDto.Username == "" || requstDto.Code == "" {
		u.l.Println("fields are empty or xss")
		//http.Error(rw, "Fields are empty or xss attack happened! error:"+err.Error(), http.StatusExpectationFailed) //400
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, errors.New("fields are empty or xss happened")
	}
	u.l.Println("3333333333")
	existsErr := u.service.UserExists(requstDto.Username)
	if existsErr != nil {
		u.l.Println(existsErr)
		//http.Error(rw, "User with entered username already exists!", http.StatusConflict) //409
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, existsErr
	}
	u.l.Println("444444444444")
	var code int
	code, convertError := strconv.Atoi(requstDto.Code)
	if convertError != nil {
		u.l.Println(convertError)
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, convertError
		//http.Error(rw, "Error converting code from string to int! error:"+convertError.Error(), http.StatusConflict) //409
	}
	u.l.Println("555555")
	activated, e := u.service.ActivateUserAccount(requstDto.Username, code)
	if e != nil {
		u.l.Println(e)
		//http.Error(rw, e.Error(), http.StatusConflict) //409
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, e
	}
	u.l.Println("6666666")
	if !activated {
		u.l.Println("account activation failed")
		//http.Error(rw, "Account activation failed!", http.StatusConflict) //409
		return &pb.ActivationResponse{Activated: false, Username: requstDto.Username}, errors.New("account activation failed")
	}
	u.l.Println("skoro pa kraj")
	return &pb.ActivationResponse{Activated: activated, Username: requstDto.Username}, nil
}
func (u UserHandler) GetAll(ctx context.Context, request *pb.EmptyRequest) (*pb.GetAllResponse, error) {
	u.l.Println("Handling GetAll Users")
	users, err := u.service.GetUsers()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := api.MapProduct(&user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (u UserHandler) UpdateUser(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateUserResponse, error) {
	u.l.Println("Handling UpdateUser Users")

	return &pb.UpdateUserResponse{UpdatedUser: nil}, nil
}
func (u UserHandler) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	//	fmt.Println(request.UserRequest.Email)
	u.l.Println("Handling RegisterUser")
	//TODO:mzd dodati provjeru da li se uspelo ok mapirati?
	newUser := api.MapPbUserToNewUserDto(request)
	//TODO:dodati validaciju u obliku regexa, spreciti injection napad
	if err := u.validator.Struct(newUser); err != nil {
		fmt.Println(err)
		return nil, err
		//http.Error(rw, "New user dto fields aren't entered in valid format! error:"+err.Error(), http.StatusExpectationFailed) //400
	}
	policy := bluemonday.UGCPolicy()
	//sanitize everything
	newUser.Username = strings.TrimSpace(policy.Sanitize(newUser.Username))
	newUser.FirstName = strings.TrimSpace(policy.Sanitize(newUser.FirstName))
	newUser.LastName = strings.TrimSpace(policy.Sanitize(newUser.LastName))
	newUser.Email = strings.TrimSpace(policy.Sanitize(newUser.Email))
	newUser.Password = strings.TrimSpace(policy.Sanitize(newUser.Password))
	newUser.Gender = strings.TrimSpace(policy.Sanitize(newUser.Gender))
	newUser.DateOfBirth = strings.TrimSpace(policy.Sanitize(newUser.DateOfBirth))
	newUser.PhoneNumber = strings.TrimSpace(policy.Sanitize(newUser.PhoneNumber))
	newUser.RecoveryEmail = strings.TrimSpace(policy.Sanitize(newUser.RecoveryEmail))

	if newUser.Username == "" || newUser.FirstName == "" || newUser.LastName == "" ||
		newUser.Gender == "" || newUser.DateOfBirth == "" || newUser.PhoneNumber == "" ||
		newUser.Password == "" || newUser.Email == "" || newUser.RecoveryEmail == "" {
		fmt.Println("fields are empty or xss")
		//http.Error(rw, "Fields are empty or xss attack happened! error:"+err.Error(), http.StatusExpectationFailed) //400
		return nil, errors.New("fields are empty or xss happened")
	}

	err := u.service.UserExists(newUser.Username)
	if err == nil {
		fmt.Println(err)
		//http.Error(rw, "User with entered username already exists!", http.StatusConflict) //409
		return nil, err
	}

	var hashedSaltedPassword = ""
	validPassword := u.passwordUtil.IsValidPassword(newUser.Password)

	if validPassword {
		pass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		hashedSaltedPassword = string(pass)
	} else {
		fmt.Println("Password format is not valid!")
		//http.Error(rw, "Password format is not valid! error:"+err.Error(), http.StatusBadRequest) //400
		return nil, errors.New("password format is not valid")
	}
	newUser.Password = hashedSaltedPassword
	registeredUser, er := u.service.CreateRegisteredUser(api.MapDtoToUser(newUser))

	if er != nil {
		fmt.Println(er)
		//http.Error(rw, "Failed creating registered user! error:"+er.Error(), http.StatusExpectationFailed) //
		return nil, er
	}

	return &pb.RegisterUserResponse{RegisteredUser: api.MapUserToPbResponseUser(registeredUser)}, nil
}
