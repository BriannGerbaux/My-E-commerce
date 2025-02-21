import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from "react-router";

import './index.css'

import Home from './pages/Home'
import Cart from './pages/Cart'
import Product from './pages/Product'
import { ThemeProvider } from './components/theme-provider';


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <BrowserRouter>
        <Routes>
          <Route path='/home' element={ <Home /> }></Route>
          <Route path='/cart' element={ <Cart /> }></Route>
          <Route path='/product' element={ <Product /> }></Route>
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  </StrictMode>,
)
