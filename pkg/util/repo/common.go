package repoutil

import (
	"context"
	requestutil "runar-himmel/pkg/util/request"

	"gorm.io/gorm"
)

// NewRepo creates new Repo instance
func NewRepo[T any](db *gorm.DB) *Repo[T] {
	return &Repo[T]{db}
}

// Repo represents the client for common usages
type Repo[T any] struct {
	GDB *gorm.DB
}

// Create creates a new record
func (d *Repo[T]) Create(ctx context.Context, input *T) error {
	return d.GDB.WithContext(ctx).Create(input).Error
}

// CreateInBatches creates multiple records in batches
func (d *Repo[T]) CreateInBatches(ctx context.Context, input []T, batchSize int) error {
	return d.GDB.WithContext(ctx).CreateInBatches(input, batchSize).Error
}

// Read get a record by primary key
func (d *Repo[T]) Read(ctx context.Context, output *T, conds ...any) error {
	return d.GDB.WithContext(ctx).First(output, parseConds(conds)...).Error
}

// ReadByID gets a record by primary key
func (d *Repo[T]) ReadByID(ctx context.Context, output *T, id string) error {
	return d.GDB.WithContext(ctx).Where(`id = ?`, id).Take(output).Error
}

// List gets all records that match given conditions
func (d *Repo[T]) List(ctx context.Context, output interface{}, conds ...any) error {
	return d.GDB.WithContext(ctx).Find(output, parseConds(conds)...).Error
}

// Update updates a record by conditions
func (d *Repo[T]) Update(ctx context.Context, updates any, conds ...any) error {
	db := d.GDB.WithContext(ctx).Model(new(T))
	conds = parseConds(conds)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}
	db = db.Omit("id").Updates(updates)
	return db.Error
}

// Delete deletes a record by conditions
func (d *Repo[T]) Delete(ctx context.Context, conds ...any) error {
	return d.GDB.WithContext(ctx).Delete(new(T), parseConds(conds)...).Error
}

// Count counts records that match given conditions
func (d *Repo[T]) Count(ctx context.Context, count *int64, conds ...any) error {
	db := d.GDB.WithContext(ctx).Model(new(T))
	if len(conds) > 0 {
		conds = parseConds(conds)
		db = db.Where(conds[0], conds[1:]...)
	}
	return db.Count(count).Error
}

// Existed checks if a record exists by conditions
func (d *Repo[T]) Existed(ctx context.Context, conds ...any) (bool, error) {
	var count int64
	if err := d.Count(ctx, &count, conds...); err != nil {
		return false, err
	}
	return count > 0, nil
}

// ReadAllByCondition retrieves a list of entities based on the provided query conditions.
func (d *Repo[T]) ReadAllByCondition(ctx context.Context, output interface{}, count *int64, lqc *requestutil.ListQueryCondition) error {
	db := d.GDB.WithContext(ctx)

	if lqc != nil {
		// Parse and apply filter conditions
		lqc.Filter = parseConds(lqc.Filter)
		if len(lqc.Filter) > 0 {
			db = d.GDB.Where(lqc.Filter[0], lqc.Filter[1:]...)
		}

		// Apply pagination and sorting
		db = withPaging(db, lqc.Page, lqc.PerPage)
		db = withSorting(db, lqc.Sort, d.quoteCol)

		// Retrieve data
		if err := db.Find(output).Error; err != nil {
			return err
		}
	}

	// Count total records if requested
	if lqc != nil && lqc.Count {
		if err := db.Limit(-1).Offset(-1).Count(count).Error; err != nil {
			return err
		}
	}

	return nil
}
