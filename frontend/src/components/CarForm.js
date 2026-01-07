import React, { useState, useEffect } from 'react';
import '../styles/App.css';

const CarForm = ({ car, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    firm: '',
    model: '',
    year: '',
    power: '',
    color: '',
    price: '',
    dealer_id: '',
  });

  useEffect(() => {
    if (car) {
      setFormData({
        firm: car.firm || '',
        model: car.model || '',
        year: car.year || '',
        power: car.power || '',
        color: car.color || '',
        price: car.price || '',
        dealer_id: car.dealer_id || '',
      });
    }
  }, [car]);

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
      year: parseInt(formData.year),
      power: parseInt(formData.power),
      price: parseInt(formData.price),
      dealer_id: parseInt(formData.dealer_id),
    };
    
    onSubmit(submitData);
  };

  const isEditMode = !!car;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2 className="modal-title">
          {isEditMode ? 'Редактировать автомобиль' : 'Добавить автомобиль'}
        </h2>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Марка *</label>
            <input
              type="text"
              name="firm"
              value={formData.firm}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Toyota"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Модель *</label>
            <input
              type="text"
              name="model"
              value={formData.model}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Camry"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Год выпуска *</label>
            <input
              type="number"
              name="year"
              value={formData.year}
              onChange={handleChange}
              className="form-input"
              min="1900"
              max="2024"
              placeholder="Например: 2023"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Мощность (л.с.) *</label>
            <input
              type="number"
              name="power"
              value={formData.power}
              onChange={handleChange}
              className="form-input"
              min="1"
              placeholder="Например: 200"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Цвет *</label>
            <input
              type="text"
              name="color"
              value={formData.color}
              onChange={handleChange}
              className="form-input"
              placeholder="Например: Красный"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Цена ($) *</label>
            <input
              type="number"
              name="price"
              value={formData.price}
              onChange={handleChange}
              className="form-input"
              min="1"
              placeholder="Например: 25000"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">ID дилера *</label>
            <input
              type="number"
              name="dealer_id"
              value={formData.dealer_id}
              onChange={handleChange}
              className="form-input"
              min="1"
              placeholder="Например: 1"
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

export default CarForm;