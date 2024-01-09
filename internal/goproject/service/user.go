package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/liangjunmo/gocode"
	"github.com/liangjunmo/gotraceutil"
	"github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/goproject/model"
	"github.com/liangjunmo/goproject/internal/goproject/mutex"
	"github.com/liangjunmo/goproject/internal/goproject/repository"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type UserService interface {
	List(ctx context.Context, cmd ListCommand) (pagination.Pagination, []model.User, error)
	Search(ctx context.Context, cmd SearchCommand) (map[uint32]model.User, error)
	Get(ctx context.Context, cmd GetCommand) (model.User, error)
	GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (model.User, error)
	Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

func NewUserService(mutexProvider mutex.MutexProvider, userRepository repository.UserRepository, userCenterClient usercenterproto.UserCenterClient) UserService {
	return newUserService(mutexProvider, userRepository, userCenterClient)
}

type userService struct {
	log              *logrus.Entry
	mutexProvider    mutex.MutexProvider
	userRepository   repository.UserRepository
	userCenterClient usercenterproto.UserCenterClient
}

func newUserService(mutexProvider mutex.MutexProvider, userRepository repository.UserRepository, userCenterClient usercenterproto.UserCenterClient) *userService {
	return &userService{
		log:              logrus.WithField("tag", "goproject.user_service"),
		mutexProvider:    mutexProvider,
		userRepository:   userRepository,
		userCenterClient: userCenterClient,
	}
}

type ListCommand struct {
	pagination.Request
}

func (service *userService) List(ctx context.Context, cmd ListCommand) (pagination.Pagination, []model.User, error) {
	p, users, err := service.userRepository.List(ctx, repository.UserCriteria{
		Request: cmd.Request,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	uids := make([]uint32, 0, len(users))

	for _, u := range users {
		uids = append(uids, u.UID)
	}

	rep, err := service.userCenterClient.SearchUser(ctx, &usercenterproto.SearchUserRequest{
		Uids: uids,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return nil, nil, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	for i := range users {
		users[i].UserCenterUser = rep.Users[users[i].UID]
	}

	return p, users, nil
}

type SearchCommand struct {
	Uids     []uint32
	Username string
}

func (service *userService) Search(ctx context.Context, cmd SearchCommand) (map[uint32]model.User, error) {
	if len(cmd.Uids) == 0 && cmd.Username == "" {
		return nil, nil
	}

	rep, err := service.userCenterClient.SearchUser(ctx, &usercenterproto.SearchUserRequest{
		Uids:     cmd.Uids,
		Username: cmd.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return nil, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	if len(rep.Users) == 0 {
		return nil, nil
	}

	uids := make([]uint32, 0, len(rep.Users))

	for uid := range rep.Users {
		uids = append(uids, uid)
	}

	users, err := service.userRepository.Search(ctx, repository.UserCriteria{
		Uids: uids,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	for uid, u := range users {
		u.UserCenterUser = rep.Users[uid]
		users[uid] = u
	}

	return users, nil
}

type GetCommand struct {
	UID uint32
}

func (service *userService) Get(ctx context.Context, cmd GetCommand) (model.User, error) {
	user, exist, err := service.userRepository.Get(ctx, cmd.UID)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !exist {
		return model.User{}, codes.UserNotFound
	}

	rep, err := service.userCenterClient.GetUserByUID(ctx, &usercenterproto.GetUserByUIDRequest{
		UID: user.UID,
	})
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return model.User{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	user.UserCenterUser = rep.User

	return user, nil
}

type GetByUsernameCommand struct {
	Username string
}

func (service *userService) GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (model.User, error) {
	rep, err := service.userCenterClient.GetUserByUsername(ctx, &usercenterproto.GetUserByUsernameRequest{
		Username: cmd.Username,
	})
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return model.User{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	user, exist, err := service.userRepository.Get(ctx, rep.User.UID)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !exist {
		return model.User{}, codes.UserNotFound
	}

	user.UserCenterUser = rep.User

	return user, nil
}

type CreateCommand struct {
	Username string
	Password string
}

func (service *userService) Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error) {
	rep, err := service.userCenterClient.CreateUser(ctx, &usercenterproto.CreateUserRequest{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return 0, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	m := service.mutexProvider.ProvideCreateUserMutex(rep.UID)

	err = m.Lock(ctx)
	if err != nil {
		return 0, codes.Timeout
	}
	defer func() {
		ok, err := m.Unlock(ctx)
		if !ok || err != nil {
			service.log.WithContext(ctx).WithError(err).Error("unlock")
		}
	}()

	_, exist, err := service.userRepository.Get(ctx, rep.UID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if exist {
		return 0, codes.UserAlreadyExists
	}

	user := model.User{
		UID: rep.UID,
	}

	err = service.userRepository.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return user.UID, nil
}

type ValidatePasswordCommand struct {
	Username string
	Password string
}

func (service *userService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	rep, err := service.userCenterClient.ValidatePassword(ctx, &usercenterproto.ValidatePasswordRequest{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	return nil
}

func (service *userService) taskToRunExample(ctx context.Context, log *logrus.Entry) {
	log.WithContext(ctx).Info("runExample")
}

func RunUserScheduler(
	ctx context.Context,
	wg *sync.WaitGroup,
	mutexProvider mutex.MutexProvider,
	userRepository repository.UserRepository,
	userCenterClient usercenterproto.UserCenterClient,
) {
	service := newUserService(mutexProvider, userRepository, userCenterClient)

	runUserScheduler(ctx, wg, service)
}

func runUserScheduler(ctx context.Context, wg *sync.WaitGroup, service *userService) {
	wg.Add(1)
	go jobToRunExample(ctx, wg, service)
}

func jobToRunExample(ctx context.Context, wg *sync.WaitGroup, service *userService) {
	log := logrus.WithField("tag", "goproject.job_to_run_example")

	defer func() {
		log.Info("quit")
		wg.Done()
	}()

	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			service.taskToRunExample(gotraceutil.Trace(ctx), log)
		}
	}
}
