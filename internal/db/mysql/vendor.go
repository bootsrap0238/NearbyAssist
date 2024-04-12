package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) CountVendor(filter models.VendorStatus) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Vendor"

	switch filter {
	case models.VENDOR_STATUS_RESTRICTED:
		query += " WHERE restricted = 1"
	case models.VENDOR_STATUS_UNRESTRICTED:
		query += " WHERE restricted = 0"
	}

	count := -1
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return count, nil
}

func (m *Mysql) FindVendorById(id int) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, vendorId, rating, job, restricted FROM Vendor WHERE id = ?"

	vendor := models.NewVendorModel()
	err := m.Conn.GetContext(ctx, vendor, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (m *Mysql) RestrictVendor(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Vendor SET restricted = 1 WHERE vendorId = ?"

	_, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (m *Mysql) UnrestrictVendor(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Vendor SET restricted = 0 WHERE vendorId = ?"

	_, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
