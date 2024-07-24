import { BrowserRouter, Routes, Route } from "react-router-dom";
import ErrorPage from "./rc/components/errors/ErrorPage";
import { PostViewer } from "./rc/components/posts/PostViewer";
import { ProtectRoutes } from "./rc/components/ProtectRoutes";
import CreatePost from "./rc/components/user/CreatePost";
import Login from "./rc/components/user/Login";
import Register from "./rc/components/user/Register";
import UserInfo from "./rc/components/user/UserInfo";
import UserProfile from "./rc/components/user/UserProfile";
import Home from "./rc/components/Home";

function App() {
  return (
  <div>
    <BrowserRouter>
      <Routes>
        <Route path="/v" element={<ProtectRoutes/>}>
          <Route path="/v/createpost" element={<CreatePost/>}/>
          <Route path="/v/user-profile/:username" element={<UserProfile/>}/>
        </Route>

        <Route path="/" element={<div>
           <div className="text-3xl font-bold underline">
            Hello world!
          </div>
          </div>}/>

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