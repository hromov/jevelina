package cdb

import "gorm.io/gorm/clause"

func (c *CDB) Create(i interface{}) (interface{}, error) {
	if err := c.DB.Omit(clause.Associations).Create(i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func (c *CDB) Update(i interface{}) (interface{}, error) {
	if err := c.DB.Save(i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func (c *CDB) Delete(i interface{}) (interface{}, error) {
	if err := c.DB.Delete(i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func (c *CDB) List(list []interface{}, limit, offset int) ([]interface{}, error) {
	if result := c.DB.Limit(limit).Offset(offset).Find(&list); result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}
