package base_postgres

type FindAll func(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error
type FindAllDeleted func(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error
type FindOne func(one HasId, s Scope) error
type FindOneDeleted func(one HasId, s Scope) error
type Create func(one HasId) error
type Update func(one HasId) error
type Delete func(one HasId) error

type FindAllInterface interface {
	FindAll(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error
}
type FindAllDeletedInterface interface {
	FindAllDeleted(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error
}

type GetFullInterface interface {
	GetFull(s Scope, all interface{}) error
}

type FindOneInterface interface {
	FindOne(one HasId, s Scope) error
}
type FindOneDeletedInterface interface {
	FindOneDeleted(one HasId, s Scope) error
}

type CreateInterface interface {
	Create(one HasId) error
}

type UpdateInterface interface {
	Update(one HasId) error
	PartialUpdate(one HasId) error
	Recover(one HasId) error
}

type DeleteInterface interface {
	Delete(one HasId) error
}

type CrudServiceInterface interface {
	GetFullInterface
	FindAllInterface
	FindAllDeletedInterface
	FindOneInterface
	FindOneDeletedInterface
	CreateInterface
	UpdateInterface
	DeleteInterface
}

type CrudService struct {
	repo CrudRepository
}

func NewCrudService() *CrudService {
	return &CrudService{
		repo: New(),
	}
}

func (c *CrudService) GetFull(s Scope, all interface{}) error {
	return c.repo.GetFull(s, all)
}

func (c *CrudService) FindAll(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error {
	return c.repo.FindAll(p, o, s, total, all, se)
}
func (c *CrudService) FindAllDeleted(all interface{}, s Scope, p Pager, o OrderFilter, total *int64, se Searcher) error {
	return c.repo.FindAllDeleted(p, o, s, total, all, se)
}

func (c *CrudService) FindOne(one HasId, s Scope) error {
	return c.repo.FindOne(one.GetId(), s, one)
}
func (c *CrudService) FindOneDeleted(one HasId, s Scope) error {
	return c.repo.FindOneDeleted(one.GetId(), s, one)
}

func (c *CrudService) Create(one HasId) error {
	return c.repo.Save(one)
}

func (c *CrudService) PartialUpdate(one HasId) error {
	return c.repo.PartialUpdate(one)
}

func (c *CrudService) Update(one HasId) error {
	return c.repo.Save(one)
}

func (c *CrudService) Recover(one HasId) error {
	return c.repo.Recover(one)
}

func (c *CrudService) Delete(one HasId) error {
	return c.repo.Delete(one)
}
