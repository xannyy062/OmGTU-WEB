import React from 'react';
import '../styles/App.css';

const CarList = ({ cars, onEdit, onDelete }) => {
  if (!cars || cars.length === 0) {
    return (
      <div className="empty-state">
        <div className="empty-icon">🚗</div>
        <h3>Автомобили не найдены</h3>
        <p>Добавьте новый автомобиль или измените критерии поиска</p>
      </div>
    );
  }

  return (
    <div className="cards-grid">
      {cars.map((car) => (
        <div key={car.id} className="card">
          <div className="card-header">
            <h3 className="card-title">
              {car.firm} {car.model}
            </h3>
            <div className="card-badge">ID: {car.id}</div>
          </div>
          
          <div className="card-details">
            <div className="detail-row">
              <span className="detail-label">Год выпуска:</span>
              <span className="detail-value">{car.year}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Мощность:</span>
              <span className="detail-value">{car.power} л.с.</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Цвет:</span>
              <span className="detail-value">{car.color}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Цена:</span>
              <span className="detail-value">${car.price.toLocaleString('ru-RU')}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">ID дилера:</span>
              <span className="detail-value">{car.dealer_id}</span>
            </div>
          </div>
          
          <div className="card-actions">
            <button 
              onClick={() => onEdit(car.id)} 
              className="btn btn-warning"
            >
              Редактировать
            </button>
            <button 
              onClick={() => onDelete(car.id)} 
              className="btn btn-danger"
            >
              Удалить
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default CarList;