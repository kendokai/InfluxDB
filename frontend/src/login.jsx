import { useState } from 'react';

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'


const sendPostRequest = async (username, password) => {
  try {
    const response = await fetch('/attempt-login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json', // specify the content type
      },
      body: JSON.stringify({
        'username': username,
        'password': password
      }), // the body must be serialized into JSON format
    });
    console.log('Response:', response.status);
    if (response.ok) {
      window.location.href = '/'
    }
  } catch (error) {
    console.error('Error:', error);
  }
  
};

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  //const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
  
    if (username === "" || password === "") {
      console.log("Need A Username and password")
      return
    }

    sendPostRequest(username, password)
    
  };

  return (
    <div className="wrapper">
      <h1>Sign In</h1>
      <form onSubmit={handleSubmit}>
        <label htmlFor="fname"></label>
        <input 
          type="text" 
          id="fname" 
          name="fname"
          placeholder="USERNAME" required 
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <label htmlFor="lname"></label>
        <input 
          type="password" 
          id="lname" 
          name="lname"
          placeholder="PASSWORD" required 
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <input type="submit" value="SUBMIT" />
      </form>
    </div>
  );
};

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <Login />
  </StrictMode>,
)

