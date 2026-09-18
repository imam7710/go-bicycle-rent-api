package repository

import "bicycle-rent-api/entity"

type VoucherRepository interface {
	Create(voucher entity.Voucher) (entity.Voucher, error)
	GetAll() ([]entity.Voucher, error)
	GetByID(voucherID int64) (entity.Voucher, error)
	GetByCode(code string) (entity.Voucher, error)
	Update(voucher entity.Voucher) (entity.Voucher, error)
	Delete(voucherID int64) error
}
