import React, { useState, useEffect } from 'react';
import './styles/App.css';
import { carApi, dealerApi } from './services/api';
import CarList from './components/CarList';
import CarForm from './components/CarForm';
import DealerList from './components/DealerList';

function App() {
  const [activeTab, setActiveTab] = useState('cars');
  const [cars, setCars] = useState([]);
  const [dealers, setDealers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [showCarForm, setShowCarForm] = useState(false);
  const [editingCar, setEditingCar] = useState(null);
  const [showDealerForm, setShowDealerForm] = useState(false);
  const [editingDealer, setEditingDealer] = useState(null);
  const [searchId, setSearchId] = useState('');

  // Загрузить все автомобили
  const fetchCars = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await carApi.getAll();
      setCars(response.data);
    } catch (err) {
      setError('Не удалось загрузить автомобили. Убедитесь, что сервер запущен.');
    } finally {
      setLoading(false);
    }
  };

  // Загрузить всех дилеров
  const fetchDealers = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await dealerApi.getAll();
      setDealers(response.data);
    } catch (err) {
      setError('Не удалось загрузить дилеров.');
    } finally {
      setLoading(false);
    }
  };

  // Первоначальная загрузка данных
  useEffect(() => {
    fetchCars();
    fetchDealers();
  }, []);

  // Поиск автомобиля по ID
  const searchCarById = async (id) => {
    if (!id) return;
    
    setLoading(true);
    setError('');
    try {
      const response = await carApi.getById(id);
      setCars([response.data]);
    } catch (err) {
      setError('Автомобиль не найден. Проверьте ID.');
      setCars([]);
    } finally {
      setLoading(false);
    }
  };

  // Поиск дилера по ID
  const searchDealerById = async (id) => {
    if (!id) return;
    
    setLoading(true);
    setError('');
    try {
      const response = await dealerApi.getById(id);
      setDealers([response.data]);
    } catch (err) {
      setError('Дилер не найден. Проверьте ID.');
      setDealers([]);
    } finally {
      setLoading(false);
    }
  };

  // Обработка отправки формы автомобиля
  const handleCarSubmit = async (carData) => {
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      if (editingCar) {
        await carApi.update(editingCar.id, carData);
        setSuccess('Автомобиль успешно обновлен!');
      } else {
        await carApi.create(carData);
        setSuccess('Автомобиль успешно добавлен!');
      }
      
      setShowCarForm(false);
      setEditingCar(null);
      fetchCars();
    } catch (err) {
      setError(err.response?.data?.message || 'Не удалось сохранить автомобиль');
    } finally {
      setLoading(false);
    }
  };

  // Обработка отправки формы дилера
  const handleDealerSubmit = async (dealerData) => {
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      if (editingDealer) {
        await dealerApi.update(editingDealer.id, dealerData);
        setSuccess('Дилер успешно обновлен!');
      } else {
        await dealerApi.create(dealerData);
        setSuccess('Дилер успешно добавлен!');
      }
      
      setShowDealerForm(false);
      setEditingDealer(null);
      fetchDealers();
    } catch (err) {
      setError(err.response?.data?.message || 'Не удалось сохранить дилера');
    } finally {
      setLoading(false);
    }
  };

  // Удалить автомобиль
  const deleteCar = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот автомобиль?')) return;
    
    setLoading(true);
    setError('');
    
    try {
      await carApi.delete(id);
      setSuccess('Автомобиль успешно удален!');
      fetchCars();
    } catch (err) {
      setError('Не удалось удалить автомобиль');
    } finally {
      setLoading(false);
    }
  };

  // Удалить дилера
  const deleteDealer = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить этого дилера?')) return;
    
    setLoading(true);
    setError('');
    
    try {
      await dealerApi.delete(id);
      setSuccess('Дилер успешно удален!');
      fetchDealers();
    } catch (err) {
      setError('Не удалось удалить дилера');
    } finally {
      setLoading(false);
    }
  };

  // Редактировать автомобиль
  const editCar = async (id) => {
    const car = cars.find(c => c.id === id);
    if (car) {
      setEditingCar(car);
      setShowCarForm(true);
    }
  };

  // Редактировать дилера
  const editDealer = async (id) => {
    const dealer = dealers.find(d => d.id === id);
    if (dealer) {
      setEditingDealer(dealer);
      setShowDealerForm(true);
    }
  };

  // Очистить поиск и показать все записи
  const clearSearch = () => {
    setSearchId('');
    if (activeTab === 'cars') {
      fetchCars();
    } else {
      fetchDealers();
    }
  };

  // Обработка поиска
  const handleSearch = () => {
    if (!searchId.trim()) {
      clearSearch();
      return;
    }
    
    if (activeTab === 'cars') {
      searchCarById(searchId);
    } else {
      searchDealerById(searchId);
    }
  };

  // Очистка сообщений через 5 секунд
  useEffect(() => {
    const timer = setTimeout(() => {
      setError('');
      setSuccess('');
    }, 5000);
    return () => clearTimeout(timer);
  }, [error, success]);

  // Рендеринг формы дилера
  const renderDealerForm = () => {
    return (
      <div className="modal-overlay">
        <div className="modal-content">
          <h2 className="modal-title">
            {editingDealer ? 'Редактировать дилера' : 'Добавить дилера'}
          </h2>
          
          <form onSubmit={(e) => {
            e.preventDefault();
            const formData = {
              name: e.target.name.value,
              city: e.target.city.value,
              address: e.target.address.value,
              area: e.target.area.value,
              rating: parseFloat(e.target.rating.value),
            };
            handleDealerSubmit(formData);
          }}>
            <div className="form-group">
              <label className="form-label">Название *</label>
              <input
                type="text"
                name="name"
                defaultValue={editingDealer?.name || ''}
                className="form-input"
                placeholder="Например: Автоцентр Премиум"
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label">Город *</label>
              <input
                type="text"
                name="city"
                defaultValue={editingDealer?.city || ''}
                className="form-input"
                placeholder="Например: Москва"
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label">Адрес *</label>
              <input
                type="text"
                name="address"
                defaultValue={editingDealer?.address || ''}
                className="form-input"
                placeholder="Например: ул. Ленина, 15"
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label">Район *</label>
              <input
                type="text"
                name="area"
                defaultValue={editingDealer?.area || ''}
                className="form-input"
                placeholder="Например: Центральный"
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label">Рейтинг (0-5) *</label>
              <input
                type="number"
                name="rating"
                defaultValue={editingDealer?.rating || ''}
                className="form-input"
                min="0"
                max="5"
                step="0.1"
                placeholder="Например: 4.5"
                required
              />
            </div>

            <div className="modal-actions">
              <button 
                type="button" 
                onClick={() => {
                  setShowDealerForm(false);
                  setEditingDealer(null);
                }} 
                className="btn btn-secondary"
              >
                Отмена
              </button>
              <button type="submit" className="btn btn-primary">
                {editingDealer ? 'Обновить' : 'Добавить'}
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  };

  return (
    <div className="container">
      {/* Шапка */}
      <header className="header">
        <h1>Система управления автосалоном</h1>
        <p className="header-subtitle">Управление автомобилями и дилерами</p>
      </header>

      {/* Вкладки */}
      <div className="tabs-container">
        <div className="tabs">
          <button
            className={`tab-button ${activeTab === 'cars' ? 'active' : ''}`}
            onClick={() => setActiveTab('cars')}
          >
            Автомобили
          </button>
          <button
            className={`tab-button ${activeTab === 'dealers' ? 'active' : ''}`}
            onClick={() => setActiveTab('dealers')}
          >
            Дилеры
          </button>
        </div>
      </div>

      {/* Сообщения */}
      {error && (
        <div className="message message-error">
          <span>{error}</span>
          <button 
            className="message-close"
            onClick={() => setError('')}
          >
            ×
          </button>
        </div>
      )}
      
      {success && (
        <div className="message message-success">
          <span>{success}</span>
          <button 
            className="message-close"
            onClick={() => setSuccess('')}
          >
            ×
          </button>
        </div>
      )}

      {/* Панель управления автомобилями */}
      {activeTab === 'cars' && (
        <>
          <div className="controls-panel">
            <div className="search-section">
              <input
                type="text"
                value={searchId}
                onChange={(e) => setSearchId(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Введите ID автомобиля..."
                className="search-input"
              />
              <div className="controls-actions">
                <button onClick={handleSearch} className="btn btn-warning">
                  Найти
                </button>
                <button onClick={clearSearch} className="btn btn-secondary">
                  Очистить
                </button>
                <button onClick={fetchCars} className="btn btn-success">
                  Обновить список
                </button>
                <button 
                  onClick={() => {
                    setEditingCar(null);
                    setShowCarForm(true);
                  }} 
                  className="btn btn-primary"
                >
                  Добавить автомобиль
                </button>
              </div>
            </div>
          </div>

          {loading ? (
            <div className="loading-container">
              <div className="spinner"></div>
              <p>Загрузка автомобилей...</p>
            </div>
          ) : (
            <CarList
              cars={cars}
              onEdit={editCar}
              onDelete={deleteCar}
            />
          )}
        </>
      )}

      {/* Панель управления дилерами */}
      {activeTab === 'dealers' && (
        <>
          <div className="controls-panel">
            <div className="search-section">
              <input
                type="text"
                value={searchId}
                onChange={(e) => setSearchId(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Введите ID дилера..."
                className="search-input"
              />
              <div className="controls-actions">
                <button onClick={handleSearch} className="btn btn-warning">
                  Найти
                </button>
                <button onClick={clearSearch} className="btn btn-secondary">
                  Очистить
                </button>
                <button onClick={fetchDealers} className="btn btn-success">
                  Обновить список
                </button>
                <button 
                  onClick={() => {
                    setEditingDealer(null);
                    setShowDealerForm(true);
                  }} 
                  className="btn btn-primary"
                >
                  Добавить дилера
                </button>
              </div>
            </div>
          </div>

          {loading ? (
            <div className="loading-container">
              <div className="spinner"></div>
              <p>Загрузка дилеров...</p>
            </div>
          ) : (
            <DealerList
              dealers={dealers}
              onEdit={editDealer}
              onDelete={deleteDealer}
            />
          )}
        </>
      )}

      {/* Модальное окно формы автомобиля */}
      {showCarForm && (
        <CarForm
          car={editingCar}
          onSubmit={handleCarSubmit}
          onCancel={() => {
            setShowCarForm(false);
            setEditingCar(null);
          }}
        />
      )}

      {/* Модальное окно формы дилера */}
      {showDealerForm && renderDealerForm()}
    </div>
  );
}

export default App;