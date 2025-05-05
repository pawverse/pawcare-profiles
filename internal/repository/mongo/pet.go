package mongo

import (
	"context"
	"errors"

	ownerdomain "github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"github.com/pawverse/pawcare-profiles/internal/pet/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	PetCollection = "pets"
)

type pet struct {
	Id          primitive.ObjectID `bson:"_id"`
	OwnerId     primitive.ObjectID `bson:"owner_id"`
	Name        string             `bson:"name"`
	DateOfBirth primitive.DateTime `bson:"date_of_birth"`
	Species     string             `bson:"species"`
	Breed       string             `bson:"breed"`
	Weight      float64            `bson:"weight"`
	Gender      string             `bson:"gender"`
}

func (p *pet) ToModel() *domain.Pet {
	petProfile := domain.NewPetProfile(p.Name, p.DateOfBirth.Time(), p.Species, p.Breed, p.Weight, domain.EGender(p.Gender))
	pet := domain.NewPet(ownerdomain.OwnerId(p.OwnerId.Hex()), petProfile)
	pet.Id = domain.PetId(p.Id.Hex())
	return pet
}

func toPetModel(p pet) *domain.Pet {
	return p.ToModel()
}

func fromPetModel(p domain.Pet) (*pet, error) {
	id, err := primitive.ObjectIDFromHex(string(p.Id))
	if err != nil {
		if !errors.Is(err, primitive.ErrInvalidHex) {
			return nil, err
		}
		id = primitive.NewObjectID()
	}

	ownerId, err := primitive.ObjectIDFromHex(string(p.OwnerId))
	if err != nil {
		return nil, err
	}

	return &pet{
		Id:          id,
		OwnerId:     ownerId,
		Name:        p.Profile.Name,
		DateOfBirth: primitive.NewDateTimeFromTime(p.Profile.DateOfBirth),
		Species:     p.Profile.Species,
		Breed:       p.Profile.Breed,
		Weight:      p.Profile.Weight,
		Gender:      string(p.Profile.Gender),
	}, nil
}

type petRepository struct {
	db *mongo.Database
}

func NewPetRepository(db *mongo.Database, logger *zap.Logger) domain.IPetRepository {
	var repository domain.IPetRepository
	repository = &petRepository{db}
	repository = newPetLoggingMiddleware(logger)(repository)
	return repository
}

func (r *petRepository) Collection() *mongo.Collection {
	return r.db.Collection(PetCollection)
}

func (r *petRepository) FindById(ctx context.Context, id domain.PetId) (*domain.Pet, error) {
	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var result pet
	if err := r.Collection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPetNotFound
		}
		return nil, err
	}

	return result.ToModel(), nil
}

func (r *petRepository) FindByOwnerId(ctx context.Context, id ownerdomain.OwnerId) ([]*domain.Pet, error) {
	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	cursor, err := r.Collection().Find(ctx, bson.M{"owner_id": objectId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPetNotFound
		}
		return nil, err
	}

	var result []pet
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return utils.Map(result, toPetModel), nil
}

func (r *petRepository) Create(ctx context.Context, pet *domain.Pet) error {
	entity, err := fromPetModel(*pet)
	if err != nil {
		return err
	}

	result, err := r.Collection().InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	pet.Id = domain.PetId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (r *petRepository) Update(ctx context.Context, pet *domain.Pet) error {
	entity, err := fromPetModel(*pet)
	if err != nil {
		return err
	}

	if _, err := r.Collection().UpdateOne(ctx, bson.M{"_id": entity.Id}, bson.M{"$set": entity}); err != nil {
		return err
	}

	return nil
}
