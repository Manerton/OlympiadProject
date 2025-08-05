import { useContext } from 'react';
import { ThemeContext } from './ThemeContext';
import { Button } from 'react-bootstrap';
import { Sun, Moon } from 'react-bootstrap-icons';

const ThemeToggleButton = () => {
    const { theme, toggleTheme } = useContext(ThemeContext);

    return (
        <Button
            onClick={toggleTheme}
            variant={theme === 'light' ? 'dark' : 'light'}
            className="m-2"
            aria-label={theme === 'light' ? 'Switch to Dark Mode' : 'Switch to Light Mode'}
        >
            {theme === 'light' ? <Moon size={24} /> : <Sun size={24} />}
        </Button>
    );
};

export default ThemeToggleButton;