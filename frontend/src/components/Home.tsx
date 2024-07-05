import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import PostCule from "./posts/PostCule";
import Loading from "./Loading";
import { useNavigate } from "react-router-dom";

export default function Home(){
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[loading,setLoading] = useState(false);
    const[fetching,setFetching] = useState(false);
    const[hasMore,setHasmore] = useState(true)
    const[username,setUsername] = useState<string>();
    const[loggedIn,setLoggedIn] = useState(false);

    const navigator = useNavigate();
    useEffect(()=>{
        const loggedIn = localStorage.getItem("viewer-loggedin")
        if(loggedIn != null && loggedIn === "true"){
            setLoggedIn(true)
            const username = localStorage.getItem("username");
            if(username != null){
                setUsername(username)
            }else{
                //better
                navigator("/login")
            }
        }

        window.addEventListener("scroll",handleScroll)    
        startLoading()
        getPosts()
        setOffset(prevValue=>prevValue+5)
        
        return ()=>window.removeEventListener("scroll",handleScroll)
    },[])

    function handleScroll(){
        console.log("called")
        if((window.innerHeight + document.documentElement.scrollTop) === document.documentElement.offsetHeight){
            console.log("reached.")
            setFetching(true);
        }
    }

    useEffect(()=>{
        if(!fetching) return;
        setOffset(prevValue=>prevValue+5)
        getPosts()
    },[fetching])

    function startLoading(){
        setLoading(true);

        setTimeout(()=>{
            setLoading(false)
        },3000)
    }
    async function getPosts(){
        console.log("fetching")
        if(!hasMore){
            setFetching(false);
            return;
        }

        const response = await fetch(`http://localhost:8000/posts/getposts?offset=${offset}&limit=5`)
        const res = await response.json();
        if(res === null){
            setHasmore(false)
            setFetching(false)
            return
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
                    ""
            )
            postsarr.push(post)
        }
        setTimeout(()=>{
            setFetching(false)
            setPosts(postsarr)
        },3000)
    }

    function goToCreatePost(){
        setLoading(true);
        setTimeout(()=>{
            setLoading(false);
            navigator("/v/createpost")
        },2000)
    }
    function goToUserProfile(){
        setLoading(true);
        setTimeout(()=>{
            setLoading(false);
            navigator(`/user/${username}`)
        },2000)
    }
    function logOut(){
        setLoading(true);
        localStorage.removeItem("token");
        localStorage.setItem("viewer-loggedin","false");
        localStorage.removeItem("username")
        setTimeout(()=>{
            setLoading(false);
            navigator("/")
        },2000)
    }

    return(
        <div>
            <Loading enable={loading}/>
            {!loading?
            <div>
                <h3>Home</h3>

                {
                    loggedIn 
                    &&
                    <div>
                        <div onClick={goToCreatePost} className="i-links">Create Post</div>
                        <div onClick={goToUserProfile} className="i-links">Your Profile</div>
                        <div onClick={logOut} className="i-links">Logout</div>
                    </div>
                }
                {posts.map(post=><PostCule post={post}/>)}

                {fetching && 
                <div className="scroll-loader-div">
                    <span className="scroll-loader"></span>
                </div>}
                {!hasMore && <div className="scroll-loader-div">End.</div>} 
            </div>:<></>}
        </div>
    )
}