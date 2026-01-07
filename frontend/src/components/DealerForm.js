import React, { useState, useEffect } from 'react';
import '../styles/App.css';

const DealerForm = ({ dealer, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    name: '',
    city: '',
    address: '',
    area: '',
    rating: '',
  });

  useEffect(() => {
    if (dealer) {
      setFormData({
        name: dealer.name || '',
        city: dealer.city || '',
        address: dealer.address || '',
        area: dealer.area || '',
        rating: dealer.rating || '',
      });
    }
  }, [dealer]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const submitData = {
      ...formData,
      rating: parseFloat(formData.rating),
    };
    
    onSubmit(submitData);
  };

  const isEditMode = !!dealer;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2 className="modal-title">
          {isEditMode ? 'Редактировать дилера' : 'Добавить дилера'}
        </h2>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label required">Название</label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Автоцентр Премиум"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label required">Город</label>
            <input
              type="text"
              name="city"
              value={formData.city}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Москва"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label required">Адрес</label>
            <input
              type="text"
              name="address"
              value={formData.address}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: ул. Ленина, 15"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label required">Район</label>
            <input
              type="text"
              name="area"
              value={formData.area}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Центральный"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label required">Рейтинг (0-5)</label>
            <input
              type="number"
              name="rating"
              value={formData.rating}
              onChange={handleChange}
              className="form-input"
              min="0"
              max="5"
              step="0.1"
              placeholder="Например: 4.5"
              required
            />
          </div>

          <div className="modal-actions">
            <button type="button" onClick={onCancel} className="btn btn-secondary">
              Отмена
            </button>
            <button type="submit" className="btn btn-primary">
              {isEditMode ? 'Обновить' : 'Добавить'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default DealerForm;