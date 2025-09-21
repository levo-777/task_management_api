import React from 'react';
import { useForm } from 'react-hook-form';
import { TextField, Button, Typography, Box, Container } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

const RegistrationPage = () => {
  const { register, handleSubmit, formState: { errors } } = useForm();
  const navigate = useNavigate();

  const onSubmit = async (data) => {
    try {
      console.log('Registration data:', data); // Debug log
      const response = await api.post('/auth/register', data);
      console.log('Registration response:', response.data); // Debug log
      if (response.data.message) {
        alert('Registration successful! Please login.');
        navigate('/'); 
      }
    } catch (error) {
      console.error('Error registering', error);
      console.error('Error response:', error.response?.data); // Debug log
      alert(`Registration failed: ${error.response?.data?.error || error.message}`);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        backgroundColor: '#f0f4f8',
      }}
    >
      <Container maxWidth="xs">
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            padding: 4,
            borderRadius: 2,
            boxShadow: 3,
            backgroundColor: '#ffffff',
          }}
        >
          <Typography 
            component="h1" 
            variant="h4" 
            gutterBottom 
            sx={{ color: '#1976d2', fontWeight: 'bold' }}
          >
            Create Account
          </Typography>
          <form onSubmit={handleSubmit(onSubmit)} style={{ width: '100%' }}>
            <TextField
              label="Username"
              {...register("username", { required: "Username is required" })}
              error={!!errors.username}
              helperText={errors.username?.message}
              variant="outlined"
              margin="normal"
              fullWidth
              sx={{ backgroundColor: '#f9f9f9' }}
            />
            <TextField
              label="Email"
              type="email"
              {...register("email", { 
                required: "Email is required",
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: "Invalid email address"
                }
              })}
              error={!!errors.email}
              helperText={errors.email?.message}
              variant="outlined"
              margin="normal"
              fullWidth
              sx={{ backgroundColor: '#f9f9f9' }}
            />
            <TextField
              label="Password"
              type="password"
              {...register("password", { 
                required: "Password is required",
                minLength: {
                  value: 6,
                  message: "Password must be at least 6 characters"
                }
              })}
              error={!!errors.password}
              helperText={errors.password?.message}
              variant="outlined"
              margin="normal"
              fullWidth
              sx={{ backgroundColor: '#f9f9f9' }}
            />
            
            {/* Buttons Section */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', marginTop: 2 }}>
              <Button 
                type="submit" 
                variant="contained" 
                color="primary" 
                sx={{ fontWeight: 'bold' }}
              >
                Register
              </Button>
              
              {/* Login Button */}
              <Button 
                variant="outlined" 
                color="secondary" 
                onClick={() => navigate('/')}
                sx={{ fontWeight: 'bold' }}
              >
                Login
              </Button>
            </Box>
          </form>
        </Box>
      </Container>
    </Box>
  );
};

export default RegistrationPage;
