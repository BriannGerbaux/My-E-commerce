import React from 'react'
import { Sun } from 'lucide-react';
import { Moon } from 'lucide-react';
import { ShoppingCart } from 'lucide-react';
import { Button } from '../ui/button';
import { useTheme } from '../theme-provider';

import { Link } from "react-router";

const Navbar = () => {
  const { theme, setTheme } = useTheme();

  return (
    <div className='w-full px-8 border-b-1 border-border h-1/15 bg-background text-secondary-foreground'>
      <div className='w-full h-full flex items-center'>
        <h1 className="text-xl mr-auto underline underline-offset-4 decoration-amber-300">Yeon Shop</h1>
        <Button className='relative mr-3' variant='outline' size='icon'>
          <span className="absolute -left-1 -top-1 w-4 h-4 bg-amber-300 text-black text-xs text-center rounded-full pointer-events-none">4</span>
          <Link className='flex w-full h-full' to="/cart">
            <ShoppingCart className='m-auto'/>
          </Link>
        </Button>
        <Button onClick={() => setTheme(theme == "dark" ? "light" : "dark")} variant="outline" size="icon">
          { theme === "dark" ?
            <Sun />
          : <Moon />
          }
        </Button>
      </div>
    </div>
  )
}

export default Navbar