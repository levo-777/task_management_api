import React, { createContext, useContext, useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';

const UserContext = createContext({
  user: null,
  setUser: () => {},
  login: () => {},
  logout: () => {},
});

export const useUser = () => useContext(UserContext);

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const access_token = localStorage.getItem('access_token');
    const storedUser = localStorage.getItem('user');
    
    if (access_token && storedUser) {
      try {
        // Use stored user data directly (already combined in LoginPage)
        const userData = JSON.parse(storedUser);
        setUser(userData);
      } catch (error) {
        console.error('Failed to parse user data', error);
        setUser(null);
        localStorage.removeItem('access_token');
        localStorage.removeItem('user');
      }
    } else {
      setUser(null);
    }
  }, []);

  const login = (userData) => {
    setUser(userData);
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
  };

  return (
    <UserContext.Provider value={{ user, setUser, login, logout }}>
      {children}
    </UserContext.Provider>
  );
};
