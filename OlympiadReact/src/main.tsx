//import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
//import './index.css'
import React from 'react'
import App from './App.tsx'
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';


createRoot(document.getElementById('root')!).render(
  <App />
)
