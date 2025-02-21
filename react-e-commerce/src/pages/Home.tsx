import Navbar from '@/components/common/Navbar'
import ProductList from '@/components/common/ProductList'
import React from 'react'

const Home = () => {
  return (
    <div className='w-screen h-screen bg-background'>
      <Navbar></Navbar>
      <ProductList></ProductList>
    </div>
  )
}

export default Home