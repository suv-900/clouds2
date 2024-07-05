import React, { useEffect, useState } from "react";
import User from "../../types/User";
import { getTime, getTimeForPosts } from "../../utils/utils";
import Post from "../../types/Post";
import Loading from "../Loading";
import PostCule from "../posts/PostCule";

export default function UserInfo(){
    const[username,setUsername] = useState("");
    const[notFound,setNotFound] = useState(false);
    const[user,setUser] = useState<User>();
    const[posts,setPosts] = useState<Post[]>([]); 
    const[userdetailsSuccess,setUserDetailsSuccess] = useState(false);
    const[offset,setOffset] = useState(0);
    const[hasMore,setHasmore] = useState(true);
    const[fetching,setFetching] = useState(false);
    const[authorid,setAuthorid] = useState<number>()
    const[loading,setLoading] = useState(true);
    
    useEffect(()=>{
        const location = window.location
        const parts = location.pathname.split("/");
        const username = parts[parts.length - 1];
        window.addEventListener("scroll",handleScroll)

        setTimeout(()=>{
            setLoading(false)
        },3000)
        getUserData(username)
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
            setUserDetailsSuccess(true);
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
    return(
        <div>
            <Loading enable={loading}/>
            {notFound && <h1>User not found</h1>}
            {!loading && user?
            <div>
               <div className="user-info-container">
                    <div className="user-name">{user.username}</div>
                    <div className="userinfo-post-separator"></div> 
                </div> 
               <div>
                <div className="userinfo-posts">
                    {posts.map(post=><PostCule post={post}/>)}
                </div>

                {fetching && 
                <div className="scroll-loader-div">
                    <span className="scroll-loader"></span>
                </div>}
                {!hasMore && <div className="scroll-loader-div">End.</div>} 
                </div> 
            </div>:<></>}
        </div>
    )
}