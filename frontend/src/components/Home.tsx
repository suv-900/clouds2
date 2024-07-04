import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import PostCule from "./posts/PostCule";
import Loading from "./Loading";

export default function Home(){
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[loading,setLoading] = useState(false);
    const[fetching,setFetching] = useState(false);
    const[hasMore,setHasmore] = useState(true)

    useEffect(()=>{
        window.addEventListener("scroll",handleScroll)    
        startLoading()
        getPosts()
        setOffset(prevValue=>prevValue+5)
        
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

    function startLoading(){
        setLoading(true);

        setTimeout(()=>{
            setLoading(false)
        },3000)
    }
    async function getPosts(){
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

    return(
        <div>
            <Loading enable={loading}/>
            {!loading?
            <div>
                <h3>Home</h3>
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