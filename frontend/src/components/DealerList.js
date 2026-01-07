import React from 'react';
import '../styles/App.css';

const DealerList = ({ dealers, onEdit, onDelete }) => {
  if (!dealers || dealers.length === 0) {
    return (
      <div className="empty-state">
        <div className="empty-icon">🏢</div>
        <h3>Дилеры не найдены</h3>
        <p>Добавьте нового дилера</p>
      </div>
    );
  }

  return (
    <div className="cards-grid">
      {dealers.map((dealer) => (
        <div key={dealer.id} className="card">
          <div className="card-header">
            <h3 className="card-title">{dealer.name}</h3>
            <div className="card-badge">ID: {dealer.id}</div>
          </div>
          
          <div className="card-details">
            <div className="detail-row">
              <span className="detail-label">Город:</span>
              <span className="detail-value">{dealer.city}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Адрес:</span>
              <span className="detail-value">{dealer.address}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Район:</span>
              <span className="detail-value">{dealer.area}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Рейтинг:</span>
              <span className="detail-value">
                {dealer.rating}/5
              </span>
            </div>
          </div>
          
          <div className="card-actions">
            <button 
              onClick={() => onEdit(dealer.id)} 
              className="btn btn-warning"
            >
              Редактировать
            </button>
            <button 
              onClick={() => onDelete(dealer.id)} 
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

export default DealerList;