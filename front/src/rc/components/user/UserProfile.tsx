import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import User from "../../types/User";
import Post from "../../types/Post";
import { getTime, getTimeForPosts } from "../../utils/utils";
import Loading from "../Loading";
import PostCule from "../posts/PostCule";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs";
import { FaRegCalendarAlt } from "react-icons/fa";

export default function UserProfile(){
    const[username,setUsername] = useState("");
    const[token,setToken] = useState("") 
    const[fetching,setFetching] = useState(false);
    const[notFound,setNotFound] = useState(false);
    const[user,setUser] = useState<User>();
    const[posts,setPosts] = useState<Post[]>([]); 
    const[userdetailsSuccess,setUserDetailsSuccess] = useState(false);
    const[offset,setOffset] = useState(0);
    const[hasMore,setHasmore] = useState(true);
    const[authorid,setAuthorid] = useState<number>()
    const[loading,setLoading] = useState(false);
    const[errorMessage,setErrorMessage] = useState("");
    const[open,setOpen] = useState(false);
    const[deleteMultiplePosts,setDeleteMultiplePosts] = useState(false);    
    const[userAbout,setUserAbout] = useState("");
    const[userAboutInput,setUserAboutInput] = useState("");

    const navigator = useNavigate()
    useEffect(()=>{
        setLoading(true)
        const location = window.location
        const parts = location.pathname.split("/")
        const username = parts[parts.length - 1];

        if(username.length === 0){
            navigator("/error");
        }else{
            setUsername(username)
            const userLoggedIn = localStorage.getItem("viewer-loggedin");
            if(userLoggedIn != null && userLoggedIn === "true"){
                const token = localStorage.getItem("token");
                if(token != null){
                    setToken(token)
                }else{
                    navigator("/login")
                    return;
                }
            }else{
                navigator("/login")
                return;
            }
        }
        getUserData(username)
        window.addEventListener("scroll",handleScroll)
        return ()=>window.removeEventListener("scroll",handleScroll)
    },[])
    function handleScroll(){
        if(window.innerHeight + document.documentElement.scrollTop === document.documentElement.offsetHeight && !fetching){
            setFetching(true);
        }
    }

    useEffect(()=>{
        if(!fetching) return;
        setOffset(prevValue=>prevValue+5)
        getPosts()
    },[fetching])


    useEffect(()=>{
        if(!userdetailsSuccess) return;
        if(user){
            setTimeout(()=>{
                getPosts()
            },2000)
            setOffset(prevValue=>prevValue+5)
        }
    },[userdetailsSuccess])
    

    async function getUserData(uname:string){
        const response = await fetch(`http://localhost:8000/user?username=${uname}`)
       
        if(response.status === 200){
            const res = await response.json();
            const user = new User(
                res.UserID,
                res.Username,
                res.About,
                res.JoinDate
            )
            user.createdat = getTime(user.createdat)
            setUser(user);
            setAuthorid(user.userID)
            setUserAbout(user.about)
            setUserDetailsSuccess(true);
            setLoading(false)
        }else if(response.status === 404){
            setNotFound(true);
        }

    }
    async function getPosts(){
        if(!hasMore){
            setFetching(false);
            return;
        }
        const response = await fetch(`http://localhost:8000/posts/get-author-posts?authorid=${authorid}&limit=5&offset=${offset}`)
        
        if(response.status === 200){
            const res = await response.json();
        
            if(res === null){
                setHasmore(false);
                setFetching(false);
                return;
            }
            
            let postsarr = posts
            for(let i=0;i<res.length;i++){
            const k = res[i];
            const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    "",
                    k.Author_name,
                    0,
                    k.Post_likes,
                    false,
                    false,
                    k.Createdat_str
            )
            post.createdat = getTimeForPosts(post.createdat)
            postsarr.push(post)
            }

            setTimeout(()=>{
                setFetching(false);
                setPosts(postsarr);
            },3000)
        }else if(response.status === 500){

        }
        
    }
    
    function displayError(s:string){
        setErrorMessage(s);
        setTimeout(()=>{
            setErrorMessage("");
        },2000)
    }
    function closePopUp(){
        setTimeout(()=>{
            setUserAbout(userAboutInput)
            setOpen(false)
        },1000)
    }
    async function update(){
        const response = await fetch(`http://localhost:8000/users/update-about`,{
            method:"PUT",
            headers:{"Authorization":token},
            body:userAboutInput
        })
        if(response.status === 200){
            closePopUp()
        }else if(response.status === 400){
            displayError("please login")
        }else if(response.status === 500){
            displayError("server error.try again later")
        }
    }

    return(
        <div>
            <Loading enable={loading}/>
            {notFound && <h1>User not found</h1>}
            {!loading && user?
            <div>
               <div className="user-info-container">
                    <div className="user-name">{user.username}</div>
                    <div className="user-about">{userAbout}</div>
                    <Tabs defaultValue="posts" className="w-[800px]" >
                        <TabsList>
                            <TabsTrigger value="posts">Posts</TabsTrigger>
                        </TabsList>
                        <TabsContent value="posts">Posts tab</TabsContent>
                    </Tabs>
                    

                {fetching && 
                <div className="scroll-loader-div">
                    <span className="scroll-loader"></span>
                </div>}
                </div> 
            </div>:<></>}
        </div>
    )
}