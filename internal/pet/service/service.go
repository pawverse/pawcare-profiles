package service

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	ownerservice "github.com/pawverse/pawcare-profiles/internal/owner/service"
	petdomain "github.com/pawverse/pawcare-profiles/internal/pet/domain"
	"github.com/pawverse/pawcare-profiles/pkg/events"
	"github.com/pawverse/pawcare-core/pkg/common"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type IPetService interface {
	GetAll(ctx context.Context) ([]*petdomain.Pet, error)
	GetById(ctx context.Context, id petdomain.PetId) (*petdomain.Pet, error)
	Create(ctx context.Context, petProfile petdomain.PetProfile) (*petdomain.Pet, error)
}

type petService struct {
	petRepository petdomain.IPetRepository
	ownerService  ownerservice.IOwnerService
	eventBus      *cqrs.EventBus
}

func NewPetService(petRepository petdomain.IPetRepository, ownerService ownerservice.IOwnerService, eventBus *cqrs.EventBus, logger *zap.Logger) IPetService {
	var svc IPetService
	svc = &petService{petRepository: petRepository, ownerService: ownerService, eventBus: eventBus}
	svc = newLoggingMiddleware(logger)(svc)

	return svc
}

func (svc *petService) GetAll(ctx context.Context) ([]*petdomain.Pet, error) {
	owner, err := svc.ownerService.Get(ctx)
	if err != nil {
		return nil, err
	}

	return svc.petRepository.FindByOwnerId(ctx, owner.Id)
}

func (svc *petService) GetById(ctx context.Context, id petdomain.PetId) (*petdomain.Pet, error) {
	owner, err := svc.ownerService.Get(ctx)
	if err != nil {
		return nil, err
	}

	pet, err := svc.petRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if pet.OwnerId != owner.Id {
		return nil, common.ErrUnauthorized
	}

	return pet, nil
}

func (svc *petService) Create(ctx context.Context, petProfile petdomain.PetProfile) (*petdomain.Pet, error) {
	owner, err := svc.ownerService.Get(ctx)
	if err != nil {
		return nil, err
	}

	pet := petdomain.NewPet(owner.Id, petProfile)
	if err := svc.petRepository.Create(ctx, pet); err != nil {
		return nil, err
	}

	if err := svc.eventBus.Publish(ctx, events.PetCreated{
		Id:          string(pet.Id),
		UserId:      string(owner.UserId),
		Name:        pet.Profile.Name,
		Weight:      pet.Profile.Weight,
		Species:     pet.Profile.Species,
		DateOfBirth: pet.Profile.DateOfBirth,
		Breed:       pet.Profile.Breed,
		Gender:      string(pet.Profile.Gender),
	}); err != nil {
		return nil, err
	}

	return pet, nil
}
