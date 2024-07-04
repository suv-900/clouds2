import React from 'react';
import {BrowserRouter, Route,Routes} from 'react-router-dom';
import {PostViewer} from './components/posts/PostViewer';
import Login from './components/user/Login';
import Register from './components/user/Register';
import Home from './components/Home';
import CreatePost from './components/user/CreatePost';
import ErrorPage from './components/errors/ErrorPage';
import { ProtectRoutes } from './components/ProtectRoutes';
import UserInfo from './components/user/UserInfo';

function App() {
  return (
  <div>
    <BrowserRouter>
      <Routes>
        <Route path="/v" element={<ProtectRoutes/>}>
          <Route path="/v/createpost" element={<CreatePost/>}/>
        </Route>

        <Route path="/" element={<h1>Main Page</h1>}/>

        <Route path="/login" element={<Login/>} />
        <Route path="/home" element={<Home/>} />
        <Route path="/register" element={<Register/>} />

        <Route path="/post/view" element={<PostViewer/>} />

        <Route path="/user/:username" element={<UserInfo/>}/>
        <Route path="/error" element={<ErrorPage/>}/>
      </Routes>
    </BrowserRouter>
  </div> 
  );
}

export default App;
