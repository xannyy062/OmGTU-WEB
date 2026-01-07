import React from 'react';
import '../styles/App.css';

const SearchBar = ({ onSearch, onClear, placeholder = "Введите ID..." }) => {
  const [searchId, setSearchId] = React.useState('');

  const handleSearch = () => {
    if (searchId.trim()) {
      onSearch(searchId);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  const handleClear = () => {
    setSearchId('');
    onClear();
  };

  return (
    <div className="search-section">
      <input
        type="text"
        value={searchId}
        onChange={(e) => setSearchId(e.target.value)}
        onKeyPress={handleKeyPress}
        placeholder={placeholder}
        className="search-input"
      />
      <div className="controls-actions">
        <button onClick={handleSearch} className="btn btn-warning">
          Найти
        </button>
        <button onClick={handleClear} className="btn btn-secondary">
          Очистить
        </button>
      </div>
    </div>
  );
};

export default SearchBar;