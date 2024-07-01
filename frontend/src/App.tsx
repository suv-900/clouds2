import React from 'react';
import {createBrowserRouter,RouterProvider} from 'react-router-dom';
import {PostViewer} from './components/posts/PostViewer';
import Login from './components/user/Login';
import Register from './components/user/Register';
import Home from './components/Home';
import CreatePost from './components/user/CreatePost';
import ErrorPage from './components/errors/ErrorPage';

const router = createBrowserRouter([
  {
    path:"/",
    element:<h2>Main Page</h2>
  },
  {
    path:"/login",
    element:<Login/>
  }, 
  {
    path:"/register",
    element:<Register/>
  },
  {
    path:"/home",
    element:<Home/>
  },
  {
    path:"/post/view",
    element:<PostViewer />
  },
  {
    path:"/createpost",
    element:<CreatePost/>
  },
  {
    path:"/error",
    element:<ErrorPage/>
  }
])

function App() {
  return (
  <div>
    <RouterProvider router={router} />
  </div> 
  );
}

export default App;
